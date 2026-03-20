package obfs

import (
	"testing"
	"time"
)

func TestTLSTransportRoundtrip(t *testing.T) {
	server, err := NewTLSListener("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	addr := server.Addr()

	received := make(chan []byte, 1)
	go func() {
		buf, _ := server.ReadPacket()
		received <- buf
	}()

	client, err := NewTLSClient(addr)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	payload := []byte("hello cosvpn tls test")
	err = client.WritePacket(payload)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case data := <-received:
		if string(data) != string(payload) {
			t.Errorf("got %q, want %q", data, payload)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timeout")
	}
}

func TestTLSMultiplePackets(t *testing.T) {
	server, err := NewTLSListener("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	done := make(chan struct{})
	var packets [][]byte
	go func() {
		for i := 0; i < 5; i++ {
			buf, err := server.ReadPacket()
			if err != nil {
				break
			}
			packets = append(packets, buf)
		}
		close(done)
	}()

	client, err := NewTLSClient(server.Addr())
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	for i := 0; i < 5; i++ {
		client.WritePacket([]byte{byte(i), byte(i + 1), byte(i + 2)})
	}

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("timeout")
	}

	if len(packets) != 5 {
		t.Fatalf("got %d packets, want 5", len(packets))
	}
}
