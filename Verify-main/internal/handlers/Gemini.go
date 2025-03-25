package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type LLMRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type LLMResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func FetchGeminiData(systemPrompt, userPrompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API key not set")
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=" + apiKey
	combinedPrompt := systemPrompt + " " + userPrompt

	// Create the correct JSON structure
	requestBody, err := json.Marshal(LLMRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: combinedPrompt},
				},
			},
		},
	})
	if err != nil {
		return "", err
	}

	log.Printf("[DEBUG] Final Request Body: %s", string(requestBody))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Printf("[DEBUG] Raw Response: %s", string(body))

	var llmResponse LLMResponse
	if err := json.Unmarshal(body, &llmResponse); err != nil {
		return "", err
	}

	if len(llmResponse.Candidates) == 0 ||
		len(llmResponse.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("No predictions received")
	}

	return llmResponse.Candidates[0].Content.Parts[0].Text, nil
}

func FetchData(c *gin.Context) {
	var requestData struct {
		UserPrompt string `json:"userPrompt"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	systemPrompt := "You are an expert assistant designed to help with a digital document signing and verification project. Your role is to understand the user's requests related to this project and provide accurate, helpful, and concise responses. The project involves an administrator digitally signing a document and a user verifying that signature. The process should result in a successful verification if the signature is valid, and a failure if the signature is invalid. Your tasks include: 1.  Understanding the user's input, which may involve requests for information, code snippets, explanations, or troubleshooting assistance related to the digital signing and verification process. 2.  Providing clear and accurate information about digital signatures, verification methods, relevant libraries or tools, and potential issues. 3.  Generating code examples (in any requested language) that demonstrate the signing and verification process. 4.  Explaining the steps involved in the signing and verification process, including cryptographic concepts if necessary. 5.  Troubleshooting issues that the user may encounter during the implementation. 6.  Maintaining a focus on the project's core functionality: administrator signing, user verification, and success/failure determination. 7.  If the user asks for code, provide secure and efficient code snippets. 8.  If the user asks for concepts explain them in simple terms. Act as a helpful and knowledgeable assistant, providing the user with the necessary information and support to successfully implement their digital document signing and verification system. When responding, always consider the user's technical proficiency and adjust the level of detail accordingly. Now, await the user's input `the response should be in 3 lines. Donth ask questions give accurate answers`"

	data, err := FetchGeminiData(systemPrompt, requestData.UserPrompt)
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from Gemini LLM API"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": data})
}
