package credentialcipher

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestCipherRoundTripAndNonceUniqueness(t *testing.T) {
	t.Parallel()

	key := base64.StdEncoding.EncodeToString([]byte("0123456789abcdef0123456789abcdef"))
	cipher, err := New(key)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	aad := Metadata{OwnerType: "github_account", OwnerID: "acc-1", SecretKind: "github_pat", Version: 1}
	first, err := cipher.Encrypt([]byte("secret-value"), aad)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	second, err := cipher.Encrypt([]byte("secret-value"), aad)
	if err != nil {
		t.Fatalf("Encrypt() second error = %v", err)
	}
	if bytes.Equal(first.Nonce, second.Nonce) || bytes.Equal(first.Ciphertext, second.Ciphertext) {
		t.Fatal("Encrypt() reused nonce or ciphertext")
	}
	if bytes.Contains(first.Ciphertext, []byte("secret-value")) {
		t.Fatal("ciphertext contains plaintext")
	}

	plaintext, err := cipher.Decrypt(first, aad)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if string(plaintext) != "secret-value" {
		t.Fatalf("Decrypt() = %q", plaintext)
	}
}

func TestCipherBindsAADAndRejectsInvalidKeys(t *testing.T) {
	t.Parallel()

	if _, err := New("invalid"); err == nil {
		t.Fatal("New() expected invalid key error")
	}

	key := base64.StdEncoding.EncodeToString([]byte("0123456789abcdef0123456789abcdef"))
	cipher, err := New(key)
	if err != nil {
		t.Fatal(err)
	}
	aad := Metadata{OwnerType: "dokploy_server", OwnerID: "server-1", SecretKind: "dokploy_api_key", Version: 1}
	sealed, err := cipher.Encrypt([]byte("api-key"), aad)
	if err != nil {
		t.Fatal(err)
	}

	wrong := aad
	wrong.OwnerID = "server-2"
	if _, err := cipher.Decrypt(sealed, wrong); err == nil {
		t.Fatal("Decrypt() accepted different AAD")
	}
	if bytes.Contains([]byte(errString(cipher, sealed, wrong)), []byte("api-key")) {
		t.Fatal("Decrypt() error leaked plaintext")
	}
}

func errString(cipher *Cipher, sealed SealedValue, metadata Metadata) string {
	_, err := cipher.Decrypt(sealed, metadata)
	if err == nil {
		return ""
	}
	return err.Error()
}
