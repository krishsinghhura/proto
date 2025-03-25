package handlers

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func VerifyHandler(c *gin.Context) {
	fmt.Println("Verification triggered")

	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file upload failed", "details": err.Error()})
	}

	//save the uploaded file temporarily as u have to upload the file again

	tempFilePath := filepath.Join("./storage", "temp_"+file.Filename)
	err = c.SaveUploadedFile(file, tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save temp file", "detail": err.Error()})
		return
	}

	defer os.Remove(tempFilePath)

	//read the file document
	fileData, err := ioutil.ReadFile(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Reading of file failed", "detail": err.Error()})
		return
	}

	//compute the hash of the uploaded file
	hash := sha256.Sum256(fileData)

	//load the stored public key
	publicKeyPath := filepath.Join("./storage", file.Filename+"_public_key.pem")
	pubKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Public key not found", "details": err.Error()})
		return
	}

	block, _ := pem.Decode(pubKeyBytes)
	if block == nil || (block.Type != "RSA PUBLIC KEY" && block.Type != "PUBLIC KEY") {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Public key failed to decode", "details": err.Error()})
		return
	}
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid public key", "details": err.Error()})
		return
	}

	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not an RSA public key"})
		return
	}

	//now get the signature file
	signaturePath := filepath.Join("./storage", file.Filename+"_signature.txt")

	signature, err := ioutil.ReadFile(signaturePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Signature not found", "details": err.Error()})
		return
	}

	//verify the signature
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"verification": "FAILED", "details": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"verification": "SUCCESS"})
	}

}
