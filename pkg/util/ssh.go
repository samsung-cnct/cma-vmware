package util

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

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

	// generate public key fingerprint
	sshPubKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		fmt.Println("Error creating ssh public key")
		return "", "", err
	}
	pubKeyBytes := ssh.MarshalAuthorizedKey(sshPubKey)

	return string(privatePEMBytes), string(pubKeyBytes), nil
}

// AddPublicKeyToRemoteNode will add the publicKey to the username@host:port's authorized_keys file w/password
func AddPublicKeyToRemoteNode(host string, port string, username string, password string, publicKey string) error {
	var remoteAuthorizedKeysFile = filepath.Join("${HOME}", ".ssh", "authorized_keys")

	remoteCmd := fmt.Sprintf("echo %s >> %s && chmod 600 %s",
		strings.TrimSuffix(publicKey, "\n"),
		remoteAuthorizedKeysFile,
		remoteAuthorizedKeysFile)

	err := ExecuteCommandOnRemoteNode(host, port, username, ssh.Password(password), remoteCmd)
	if err != nil {
		fmt.Printf("ERROR: Failed to add public key to remote node (%s) via password: %s\n", host, err)
		return err
	}

	return nil
}

// ExecuteCommandOnRemoteNode executes the commmand on username@host:port using the authMethed
func ExecuteCommandOnRemoteNode(host string, port string, username string, authMethod ssh.AuthMethod, command string) error {
	config := sshClientConfig(username, authMethod)

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), config)
	if err != nil {
		fmt.Printf("ERROR: Failed to ssh into remote node (%s): %s\n", host, err)
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("ERROR: Failed to creae ssh session: %s\n", err)
		return err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		fmt.Printf("ERROR: Failed to run command on remote node (%s): %s", host, err)
		return err
	}
	fmt.Println(b.String())

	return nil
}

// SSHAuthMethPublicKey generates a ssh public key authentication method based on privateKey
func SSHAuthMethPublicKey(privateKey string) (ssh.AuthMethod, error) {
	buffer := []byte(privateKey)

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		fmt.Printf("ERROR: could not parse private key")
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

func sshClientConfig(username string, authMethod ssh.AuthMethod) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //TODO: implement known_hosts
	}
}
