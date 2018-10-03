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
}

// This test requires:
//   - SSHD be running locally (mac: sudo systemsetup -setremotelogin on)
//   - SSH_TEST_PASSWORD environment set and exported
func TestAddPublicKeyToRemoteNode(t *testing.T) {
	username := os.Getenv("USER")
	password := os.Getenv("SSH_TEST_PASSWORD")
	if password == "" {
		t.Skipf("Skipping because SSH_TEST_PASSWORD is not set and/or exported")
		return
	}
	_, public, err := generateKeyPairAndAddToRemote(t, "localhost", "22", username, password)

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

		if public == parsedKey {
			foundKey = true
		}
		authorizedKeysBytes = rest
	}

	if !foundKey {
		t.Errorf("Did not find the key: %s", public)
	}
}

// This test requires:
//   - SSHD be running locally (mac: sudo systemsetup -setremotelogin on)
//   - SSH_TEST_PASSWORD environment set and exported
func TestPublicKeyAccess(t *testing.T) {
	username := os.Getenv("USER")
	password := os.Getenv("SSH_TEST_PASSWORD")
	if password == "" {
		t.Skipf("Skipping because SSH_TEST_PASSWORD is not set and/or exported")
		return
	}
	private, _, err := generateKeyPairAndAddToRemote(t, "localhost", "22", username, password)

	// Test private key
	testCmd := "echo cma-vmware: $(date) >> ~/.ssh/test-pvka"

	authMethod, err := SSHAuthMethPublicKey(private)
	if err != nil {
		t.Errorf("Failed to generate public key access for ssh")
		return
	}

	err = ExecuteCommandOnRemoteNode("localhost", "22", username, authMethod, testCmd)
	if err != nil {
		t.Errorf("Failed to execute test command via private key")
		return
	}
}

func generateKeyPairAndAddToRemote(t *testing.T, host string, port string, username string, password string) (private string, public string, err error) {
	private, public, err = GenerateSSHKeyPair()
	if err != nil {
		t.Errorf("Error generating ssh key pair: %s", err)
		return "", "", err
	}

	err = AddPublicKeyToRemoteNode(host, port, username, password, public)
	if err != nil {
		t.Errorf("Error adding public key to remote node")
		return "", "", err
	}

	return private, public, nil
}
