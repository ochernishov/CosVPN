package obfs

import (
	"bytes"
	"crypto/rand"
	"testing"
)

// makeKey генерирует случайный 16-байтный ключ для тестов.
func makeKey(t *testing.T) [16]byte {
	t.Helper()
	var key [16]byte
	if _, err := rand.Read(key[:]); err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	return key
}

func TestObfuscateDeobfuscateRoundtrip(t *testing.T) {
	key := makeKey(t)

	// Создаём пакет 148 байт с type=1
	original := make([]byte, 148)
	original[0] = 1 // type=1
	if _, err := rand.Read(original[1:]); err != nil {
		t.Fatalf("failed to fill packet: %v", err)
	}

	obfuscated, err := Obfuscate(original, key)
	if err != nil {
		t.Fatalf("Obfuscate failed: %v", err)
	}

	// Обфусцированный пакет должен отличаться от оригинала
	if bytes.Equal(obfuscated, original) {
		t.Error("obfuscated packet should differ from original")
	}

	deobfuscated, err := Deobfuscate(obfuscated, key)
	if err != nil {
		t.Fatalf("Deobfuscate failed: %v", err)
	}

	if !bytes.Equal(deobfuscated, original) {
		t.Errorf("roundtrip failed: got %d bytes, want %d bytes", len(deobfuscated), len(original))
	}
}

func TestJunkPacketDetection(t *testing.T) {
	key := makeKey(t)

	// Junk-пакет должен распознаваться
	junk := MakeJunkPacket(key)
	if !IsJunkPacket(junk, key) {
		t.Error("MakeJunkPacket result should be detected as junk")
	}

	// Обычный пакет не должен быть junk
	normal := make([]byte, 100)
	normal[0] = 1 // type=1, гарантированно не 0xFF ^ key[0]
	// Убедимся что первый байт не совпадёт с маркером
	if (normal[0] ^ key[0]) == 0xFF {
		normal[0] = 2
	}
	if IsJunkPacket(normal, key) {
		t.Error("normal packet should not be detected as junk")
	}
}

func TestObfuscateDisabledWithZeroKey(t *testing.T) {
	var zeroKey [16]byte

	original := make([]byte, 64)
	if _, err := rand.Read(original); err != nil {
		t.Fatalf("failed to fill packet: %v", err)
	}

	obfuscated, err := Obfuscate(original, zeroKey)
	if err != nil {
		t.Fatalf("Obfuscate with zero key failed: %v", err)
	}

	if !bytes.Equal(obfuscated, original) {
		t.Error("with zero key, packet should not be modified")
	}

	deobfuscated, err := Deobfuscate(original, zeroKey)
	if err != nil {
		t.Fatalf("Deobfuscate with zero key failed: %v", err)
	}

	if !bytes.Equal(deobfuscated, original) {
		t.Error("with zero key, deobfuscated packet should equal original")
	}
}

func TestMultipleRoundtrips(t *testing.T) {
	key := makeKey(t)

	for i := 0; i < 100; i++ {
		size := 50 + i*15 // от 50 до ~1535 байт
		if size > 1500 {
			size = 1500
		}

		original := make([]byte, size)
		if _, err := rand.Read(original); err != nil {
			t.Fatalf("iteration %d: failed to fill packet: %v", i, err)
		}

		obfuscated, err := Obfuscate(original, key)
		if err != nil {
			t.Fatalf("iteration %d: Obfuscate failed: %v", i, err)
		}

		deobfuscated, err := Deobfuscate(obfuscated, key)
		if err != nil {
			t.Fatalf("iteration %d: Deobfuscate failed: %v", i, err)
		}

		if !bytes.Equal(deobfuscated, original) {
			t.Errorf("iteration %d: roundtrip failed for %d-byte packet", i, size)
		}
	}
}

func TestPaddingVariesSize(t *testing.T) {
	key := makeKey(t)

	original := make([]byte, 100)
	if _, err := rand.Read(original); err != nil {
		t.Fatalf("failed to fill packet: %v", err)
	}

	sizes := make(map[int]bool)
	for i := 0; i < 10; i++ {
		obfuscated, err := Obfuscate(original, key)
		if err != nil {
			t.Fatalf("Obfuscate failed on iteration %d: %v", i, err)
		}
		sizes[len(obfuscated)] = true
	}

	if len(sizes) < 2 {
		t.Errorf("expected varying sizes from random padding, got %d unique size(s): %v", len(sizes), sizes)
	}
}
