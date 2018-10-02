package util

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestGenerateSSHKeyPair(t *testing.T) {
	private, public, err := GenerateSSHKeyPair()
	if err != nil {
		t.Errorf("Error generating ssh key pair: %s", err)
		return
	}

	// validate private key
	pemBlock, _ := pem.Decode([]byte(private))
	ecdsaPrivateKey, err := x509.ParseECPrivateKey(pemBlock.Bytes)
	ecdsaPublicKey := &ecdsaPrivateKey.PublicKey

	if !ecdsaPublicKey.Curve.IsOnCurve(ecdsaPublicKey.X, ecdsaPublicKey.Y) {
		t.Errorf("Invalid private key generated")
		return
	}
	//t.Logf("private key validated:\n%s\n", private)

	// validate public key
	sshPubKey, err := ssh.NewPublicKey(ecdsaPublicKey)
	if err != nil {
		t.Errorf("Invalid ssh public key generation: %s", err)
		return
	}
	testPublicKey := ssh.MarshalAuthorizedKey(sshPubKey)
	if public != string(testPublicKey) {
		t.Errorf("Invalid public key generated")
		return
	}
	//t.Logf("public key validated:\n%s\n", public)
}

// This test requires:
//   - SSHD be running locally (mac: sudo systemsetup -setremotelogin on)
//   - SSH_TEST_PASSWORD environment set and exported
func TestAddPublicKeyToRemoteNode(t *testing.T) {
	_, public, err := GenerateSSHKeyPair()
	if err != nil {
		t.Errorf("Error generating ssh key pair: %s", err)
		return
	}

	username := os.Getenv("USER")
	password := os.Getenv("SSH_TEST_PASSWORD")
	if password == "" {
		t.Skipf("Skipping because SSH_TEST_PASSWORD is not set and/or exported")
		return
	}
	AddPublicKeyToRemoteNode("localhost", "22", username, password, public)

	authKeysFile := filepath.Join(os.Getenv("HOME"), ".ssh", "authorized_keys")
	authorizedKeysBytes, err := ioutil.ReadFile(authKeysFile)
	if err != nil {
		t.Errorf("Failed to load authorized_keys (%s): %v", authKeysFile, err)
		return
	}

	foundKey := false
	for len(authorizedKeysBytes) > 0 {
		sshPublicKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			t.Fatal(err)
		}
		parsedKey := string(ssh.MarshalAuthorizedKey(sshPublicKey))
		//t.Logf("parsed key: %s\n", parsedKey)

		if public == parsedKey {
			foundKey = true
		}
		authorizedKeysBytes = rest
	}

	if !foundKey {
		t.Errorf("Did not find the key: %s", public)
	}
}
