package util

import (
	"crypto/x509"
	"encoding/pem"
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
	//t.Logf("private key validated: %s", private)

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
	//t.Logf("public key validated: %s", public)
}
