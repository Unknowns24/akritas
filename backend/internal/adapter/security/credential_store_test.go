package security

import (
	"bytes"
	"context"
	"errors"
	"testing"
)

func TestCredentialStoreEncrypt(t *testing.T) {
	t.Parallel()

	store, err := NewCredentialStore("a-high-entropy-master-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	plaintext := []byte("JBSWY3DPEHPK3PXP")
	first, err := store.Encrypt(context.Background(), plaintext)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := store.Encrypt(context.Background(), plaintext)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if bytes.Equal(first, second) {
		t.Fatal("two encryptions of the same plaintext must differ (random nonce)")
	}
	if bytes.Contains(first, plaintext) || bytes.Contains(second, plaintext) {
		t.Fatal("ciphertext must not contain the plaintext")
	}
}

func TestNewCredentialStoreRejectsEmptyMasterKey(t *testing.T) {
	t.Parallel()

	if _, err := NewCredentialStore(""); !errors.Is(err, ErrEmptyMasterKey) {
		t.Fatalf("expected ErrEmptyMasterKey, got %v", err)
	}
}

func TestCredentialStoreDecryptRoundTrip(t *testing.T) {
	t.Parallel()

	store, err := NewCredentialStore("a-high-entropy-master-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	plaintexts := [][]byte{
		[]byte("JBSWY3DPEHPK3PXP"),
		[]byte(""),
		[]byte("a longer secret value with spaces and symbols !@#"),
	}

	for _, plaintext := range plaintexts {
		ciphertext, err := store.Encrypt(context.Background(), plaintext)
		if err != nil {
			t.Fatalf("encrypt: %v", err)
		}
		decrypted, err := store.Decrypt(context.Background(), ciphertext)
		if err != nil {
			t.Fatalf("decrypt: %v", err)
		}
		if !bytes.Equal(decrypted, plaintext) {
			t.Fatalf("round trip mismatch: got %q, want %q", decrypted, plaintext)
		}
	}
}

func TestCredentialStoreDecryptRejectsTooShortCiphertext(t *testing.T) {
	t.Parallel()

	store, err := NewCredentialStore("a-high-entropy-master-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := store.Decrypt(context.Background(), []byte("short")); !errors.Is(err, ErrCiphertextTooShort) {
		t.Fatalf("expected ErrCiphertextTooShort, got %v", err)
	}
}

func TestCredentialStoreDecryptRejectsCorruptedCiphertext(t *testing.T) {
	t.Parallel()

	store, err := NewCredentialStore("a-high-entropy-master-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ciphertext, err := store.Encrypt(context.Background(), []byte("JBSWY3DPEHPK3PXP"))
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	corrupted := append([]byte(nil), ciphertext...)
	corrupted[len(corrupted)-1] ^= 0xFF

	if _, err := store.Decrypt(context.Background(), corrupted); err == nil {
		t.Fatal("expected an error decrypting a corrupted ciphertext")
	}
}
