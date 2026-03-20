package obfs

import (
	crand "crypto/rand"
	"errors"
	"math/rand/v2"
)

// maxPadding — максимальный размер случайного padding (0-64 байт).
const maxPadding = 64

// Obfuscate применяет XOR к первым 16 байтам заголовка пакета,
// добавляет случайный padding (0-64 байт) и записывает длину padding
// в последний байт результата.
// С нулевым ключом пакет возвращается без изменений.
func Obfuscate(packet []byte, key [16]byte) ([]byte, error) {
	if len(packet) == 0 {
		return nil, errors.New("obfs: empty packet")
	}

	// Нулевой ключ — обфускация отключена
	var zero [16]byte
	if key == zero {
		dst := make([]byte, len(packet))
		copy(dst, packet)
		return dst, nil
	}

	// XOR первых min(16, len(packet)) байт заголовка
	xored := make([]byte, len(packet))
	copy(xored, packet)
	xorLen := len(xored)
	if xorLen > 16 {
		xorLen = 16
	}
	for i := 0; i < xorLen; i++ {
		xored[i] ^= key[i]
	}

	// Случайный padding: 0-64 байт
	padLen := rand.IntN(maxPadding + 1) // [0, 64]
	padding := make([]byte, padLen)
	if padLen > 0 {
		if _, err := crand.Read(padding); err != nil {
			return nil, err
		}
	}

	// Результат: xored_packet + padding + padLen_byte
	result := make([]byte, 0, len(xored)+padLen+1)
	result = append(result, xored...)
	result = append(result, padding...)
	result = append(result, byte(padLen))

	return result, nil
}

// Deobfuscate — обратная операция: читает длину padding из последнего байта,
// отрезает padding и применяет XOR к заголовку.
// С нулевым ключом пакет возвращается без изменений.
func Deobfuscate(packet []byte, key [16]byte) ([]byte, error) {
	if len(packet) == 0 {
		return nil, errors.New("obfs: empty packet")
	}

	// Нулевой ключ — обфускация отключена
	var zero [16]byte
	if key == zero {
		dst := make([]byte, len(packet))
		copy(dst, packet)
		return dst, nil
	}

	// Последний байт — длина padding
	padLen := int(packet[len(packet)-1])

	// Валидация: padLen + 1 (сам байт длины) не должны превышать длину пакета
	if padLen+1 >= len(packet) {
		return nil, errors.New("obfs: invalid padding length")
	}

	// Вырезаем оригинальные данные (без padding и байта длины)
	dataLen := len(packet) - padLen - 1
	result := make([]byte, dataLen)
	copy(result, packet[:dataLen])

	// XOR первых min(16, dataLen) байт
	xorLen := dataLen
	if xorLen > 16 {
		xorLen = 16
	}
	for i := 0; i < xorLen; i++ {
		result[i] ^= key[i]
	}

	return result, nil
}

// MakeJunkPacket создаёт junk-пакет размером 64-1264 байт.
// Первый байт после XOR с key[0] равен 0xFF — маркер junk-пакета.
func MakeJunkPacket(key [16]byte) []byte {
	size := 64 + rand.IntN(1201) // [64, 1264]
	junk := make([]byte, size)

	// Заполняем криптографически случайными данными
	if _, err := crand.Read(junk); err != nil {
		// fallback: оставить нули (крайне маловероятно)
		_ = err
	}

	// Устанавливаем маркер: первый байт XOR key[0] = 0xFF
	junk[0] = 0xFF ^ key[0]

	return junk
}

// IsJunkPacket проверяет, является ли пакет junk-пакетом.
// Проверяет маркер: первый байт XOR key[0] == 0xFF.
func IsJunkPacket(packet []byte, key [16]byte) bool {
	if len(packet) == 0 {
		return false
	}
	return (packet[0] ^ key[0]) == 0xFF
}

// ShouldSendJunk возвращает true с вероятностью 30%.
func ShouldSendJunk() bool {
	return rand.IntN(100) < 30
}
