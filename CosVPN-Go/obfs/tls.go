package obfs

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"sync"
	"time"
)

// TLSListener — TLS-сервер, принимает соединения и читает/пишет CosVPN-пакеты
// поверх TLS 1.3 с самоподписанным сертификатом.
type TLSListener struct {
	listener net.Listener
	conn     *tls.Conn
	mu       sync.Mutex
	accepted chan *tls.Conn
	closed   chan struct{}
}

// TLSClient — TLS-клиент, подключается к серверу и читает/пишет пакеты.
type TLSClient struct {
	conn *tls.Conn
	mu   sync.Mutex
}

// NewTLSListener создаёт TLS-сервер с самоподписанным сертификатом на указанном адресе.
// Адрес "127.0.0.1:0" выберет случайный свободный порт.
func NewTLSListener(addr string) (*TLSListener, error) {
	cert, err := generateSelfSignedCert()
	if err != nil {
		return nil, fmt.Errorf("obfs/tls: failed to generate certificate: %w", err)
	}

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
	}

	ln, err := tls.Listen("tcp", addr, tlsCfg)
	if err != nil {
		return nil, fmt.Errorf("obfs/tls: failed to listen: %w", err)
	}

	s := &TLSListener{
		listener: ln,
		accepted: make(chan *tls.Conn, 1),
		closed:   make(chan struct{}),
	}

	// Запускаем accept loop в горутине
	go s.acceptLoop()

	return s, nil
}

// acceptLoop принимает входящие соединения.
func (s *TLSListener) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.closed:
				return
			default:
				continue
			}
		}
		tlsConn, ok := conn.(*tls.Conn)
		if !ok {
			conn.Close()
			continue
		}

		// Устанавливаем первое принятое соединение
		s.mu.Lock()
		if s.conn == nil {
			s.conn = tlsConn
		}
		s.mu.Unlock()

		select {
		case s.accepted <- tlsConn:
		default:
			// Если канал полон, сохраняем только в s.conn
		}
	}
}

// Addr возвращает адрес, на котором слушает сервер.
func (s *TLSListener) Addr() string {
	return s.listener.Addr().String()
}

// Close закрывает listener и активное соединение.
func (s *TLSListener) Close() error {
	close(s.closed)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn != nil {
		s.conn.Close()
	}
	return s.listener.Close()
}

// waitConn ожидает принятия соединения, если его ещё нет.
func (s *TLSListener) waitConn() (*tls.Conn, error) {
	s.mu.Lock()
	c := s.conn
	s.mu.Unlock()
	if c != nil {
		return c, nil
	}

	select {
	case c := <-s.accepted:
		return c, nil
	case <-s.closed:
		return nil, errors.New("obfs/tls: listener closed")
	}
}

// WritePacket отправляет пакет через TLS-соединение сервера.
// Framing: [2 bytes BE length][payload].
func (s *TLSListener) WritePacket(data []byte) error {
	conn, err := s.waitConn()
	if err != nil {
		return err
	}
	return writePacket(conn, data)
}

// ReadPacket читает пакет из TLS-соединения сервера.
func (s *TLSListener) ReadPacket() ([]byte, error) {
	conn, err := s.waitConn()
	if err != nil {
		return nil, err
	}
	return readPacket(conn)
}

// NewTLSClient подключается к TLS-серверу по указанному адресу.
// Использует InsecureSkipVerify=true для работы с самоподписанными сертификатами.
func NewTLSClient(addr string) (*TLSClient, error) {
	tlsCfg := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS13,
	}

	conn, err := tls.Dial("tcp", addr, tlsCfg)
	if err != nil {
		return nil, fmt.Errorf("obfs/tls: failed to connect: %w", err)
	}

	return &TLSClient{conn: conn}, nil
}

// WritePacket отправляет пакет через TLS-соединение клиента.
// Framing: [2 bytes BE length][payload].
func (c *TLSClient) WritePacket(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return writePacket(c.conn, data)
}

// ReadPacket читает пакет из TLS-соединения клиента.
func (c *TLSClient) ReadPacket() ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return readPacket(c.conn)
}

// Close закрывает TLS-соединение клиента.
func (c *TLSClient) Close() error {
	return c.conn.Close()
}

// writePacket записывает пакет с framing: [2 bytes BE length][payload].
func writePacket(conn net.Conn, data []byte) error {
	if len(data) > 65535 {
		return errors.New("obfs/tls: packet too large (max 65535 bytes)")
	}

	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, uint16(len(data)))

	if _, err := conn.Write(header); err != nil {
		return fmt.Errorf("obfs/tls: write header: %w", err)
	}
	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("obfs/tls: write payload: %w", err)
	}
	return nil
}

// readPacket читает пакет: сначала 2 байта длины (BE), затем payload.
func readPacket(conn net.Conn) ([]byte, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, fmt.Errorf("obfs/tls: read header: %w", err)
	}

	length := binary.BigEndian.Uint16(header)
	if length == 0 {
		return nil, errors.New("obfs/tls: zero-length packet")
	}

	payload := make([]byte, length)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return nil, fmt.Errorf("obfs/tls: read payload: %w", err)
	}

	return payload, nil
}

// generateSelfSignedCert генерирует самоподписанный ECDSA P-256 сертификат.
func generateSelfSignedCert() (tls.Certificate, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("generate ECDSA key: %w", err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("generate serial number: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("create certificate: %w", err)
	}

	return tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  key,
	}, nil
}
