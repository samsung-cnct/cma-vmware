package util

import (
	"fmt"
	// "net"
	// "os"

	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

// GenerateSSHKeyPair creates a ECDSA a x509 ASN.1-DER format-PEM encoded private
// key string and a SHA256 encoded public key string
func GenerateSSHKeyPair() (private string, public string, err error) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println("Error generating private key")
		return "", "", err
	}

	// validate private key
	publicKey := &privateKey.PublicKey
	if !curve.IsOnCurve(publicKey.X, publicKey.Y) {
		fmt.Println("Error validating private key")
		return "", "", err
	}

	// convert to x509 ASN.1, DER format
	privateDERBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		fmt.Println("Error encoding private key")
		return "", "", err
	}

	// generate pem encoded private key
	privatePEMBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateDERBytes,
	})

	// generate public key fingerprint (problem)
	sshPubKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		fmt.Println("Error creating ssh public key")
		return "", "", err
	}
	pubKeyBytes := ssh.MarshalAuthorizedKey(sshPubKey)

	return string(privatePEMBytes), string(pubKeyBytes), nil
}
