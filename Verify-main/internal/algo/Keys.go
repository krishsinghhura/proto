package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func GenerateRSAKeys(bits int) (*rsa.PrivateKey, error) {
	//it returns a pointer to a struct that defines the structure of the private key and error
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	//this is the core of the function it takes two arguments , first -> a cryptographically secure random no genrato , second -> size of the private key
	//it returns two things 1->Private key , 2->error
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func SavePublicKey(publicKey *rsa.PublicKey, path string) error {
	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	//converting the public key to ASN1 format which is the standard format to represent the data structures and is commonly used in cryptography
	//above function will return two value
	//1.Public key in the format of ASN1
	//2.error
	if err != nil {
		return err
	}

	//next we convert this public key from ASN1 format to PEM format
	//Privacy enhanced mail is a text-format used to represent the cryptographic keys and certificates
	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubASN1,
		},
	)
	return os.WriteFile(path, pubPEM, 0644)
}
