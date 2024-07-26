package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const baseURL = "https://vault.immudb.io/ics/api/v1"

type Document struct {
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	IBAN          string `json:"iban"`
	Address       string `json:"address"`
	Amount       string `json:"amount"`
	Type          string `json:"type"`
}

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	router.Use(cors.New(config))

	router.PUT("/document", handlePutDocument)
	router.GET("/documents", handleListDocuments)

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handlePutDocument(c *gin.Context) {
	var doc Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonData, err := json.Marshal(doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal document"})
		return
	}

	resp, err := makeRequest("PUT", "/ledger/default/collection/default/document", jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	c.Data(resp.StatusCode, "application/json", body)
}

type SearchRequest struct {
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
}

func handleListDocuments(c *gin.Context) {
	// Get page and perPage from query parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(c.DefaultQuery("perPage", "100"))
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 100
	}

	// Construct the search request
	searchReq := SearchRequest{
		Page:    page,
		PerPage: perPage,
	}

	// Marshal the search request to JSON
	jsonData, err := json.Marshal(searchReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to construct search request"})
		return
	}

	resp, err := makeRequest("POST", "/ledger/default/collection/default/documents/search", jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	// Parse the response
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response body"})
		return
	}

	c.JSON(resp.StatusCode, result)
}

func makeRequest(method, path string, body []byte) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, baseURL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	apiKey := "default.4c10jOO-Pj0mc5_Bjl4kHA.C77o6i3NY78x3xBQ15RO3UaM2-IQuPHGAYjDU7VMSoH6TXG6"
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return client.Do(req)
}
