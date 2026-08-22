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
