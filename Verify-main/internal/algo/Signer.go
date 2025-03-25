package crypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"io/ioutil"
	"os"
)

func HashFile(filepath string) ([]byte, error) {

	fileData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(fileData)
	//it will return the hash of the contents of the file
	//it returns an array of 32 bytes representing the sha-256 hash
	//we will use this hash to verify the integrity of our files
	return hash[:], nil
}

// this is used to genrate the digital signature for the hash using the privateKey that got generated
func SignHash(privateKey *rsa.PrivateKey, hash []byte) ([]byte, error) {
	signature, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hash)
	//here we are signing the file
	//1.nil -> default is cryptographically random no generator
	//2.privateKey taken as paramter
	//3.The hash algo that was used to create the hash of the file
	//4.hash - The hash that got generated for the file
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func SaveSignature(signature []byte, path string) error {
	return ioutil.WriteFile(path, signature, 0644)
	//0644 - only owner will have read and write permission rest all will have read permission only
}
