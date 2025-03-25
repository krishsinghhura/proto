package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	algo "verification/internal/algo"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
	fmt.Println("Upload handler triggered")

	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed", "err": err.Error()})
		return
	}

	storageDir := "./storage"
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		err := os.MkdirAll(storageDir, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create storage directory"})
			return
		}
	}

	savePath := fmt.Sprintf("%s/%s", storageDir, file.Filename)
	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File save failed", "details": err.Error()})
		return
	}

	privateKey, err := algo.GenerateRSAKeys(2048)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Key generation failed", "details": err.Error()})
		return
	}

	publicKeyPath := filepath.Join(storageDir, file.Filename+"_public_key.pem")
	err = algo.SavePublicKey(&privateKey.PublicKey, publicKeyPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Public key save failed", "details": err.Error()})
		return
	}

	hash, err := algo.HashFile(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File hashing failed", "details": err.Error()})
		return
	}

	signature, err := algo.SignHash(privateKey, hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Hash signing failed", "details": err.Error()})
		return
	}

	signaturePath := filepath.Join(storageDir, file.Filename+"_signature.txt")
	err = algo.SaveSignature(signature, signaturePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Signature save failed", "details": err.Error()})
		return
	}

	jsonData := map[string]interface{}{
		"fileName":         file.Filename,
		"fileHash":         hash,
		"SignatureContent": signature,
		"publickey":        privateKey.PublicKey,
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON marshaling failed", "details": err.Error()})
		return
	}

	req, err := http.NewRequest("POST", "https://api.pinata.cloud/pinning/pinJSONToIPFS", bytes.NewBuffer(jsonBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request creation failed", "details": err.Error()})
		return
	}

	apiKey := os.Getenv("PINATA_API_KEY")
	secretKey := os.Getenv("PINATA_SECRET_KEY")
	if apiKey == "" || secretKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing Pinata API credentials"})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("pinata_api_key", apiKey)
	req.Header.Set("pinata_secret_api_key", secretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "HTTP request failed", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body", "details": err.Error()})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON response", "details": err.Error()})
		return
	}

	fileData, err := os.Open(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file", "details": err.Error()})
		return
	}
	defer fileData.Close()

	fileBody := &bytes.Buffer{}
	writer := multipart.NewWriter(fileBody)
	formFile, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create form file", "details": err.Error()})
		return
	}

	_, err = io.Copy(formFile, fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy file data", "details": err.Error()})
		return
	}

	writer.Close()

	req, err = http.NewRequest("POST", "https://api.pinata.cloud/pinning/pinFileToIPFS", fileBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request creation failed", "details": err.Error()})
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("pinata_api_key", apiKey)
	req.Header.Set("pinata_secret_api_key", secretKey)

	fileResp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File upload failed", "details": err.Error()})
		return
	}
	defer fileResp.Body.Close()

	fileRespBody, err := io.ReadAll(fileResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file response body", "details": err.Error()})
		return
	}

	var fileResult map[string]interface{}
	if err := json.Unmarshal(fileRespBody, &fileResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse file response JSON", "details": err.Error()})
		return
	}

	if ipfsHash, ok := result["IpfsHash"].(string); ok {
		if fileIpfsHash, ok := fileResult["IpfsHash"].(string); ok {
			c.JSON(http.StatusOK, gin.H{
				"message":        "File and JSON uploaded successfully",
				"json_ipfs_hash": ipfsHash,
				"file_ipfs_hash": fileIpfsHash,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "File IPFS hash not found in response"})
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON IPFS hash not found in response"})
	}
}
