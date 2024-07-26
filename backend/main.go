package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "backend/docs"
)

const baseURL = "https://vault.immudb.io/ics/api/v1"

type Document struct {
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	IBAN          string `json:"iban"`
	Address       string `json:"address"`
	Amount        string `json:"amount"`
	Type          string `json:"type"`
}

type SearchRequest struct {
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
}

// @title ImmuDB Vault API
// @version 1.0
// @description This is an API server for interacting with ImmuDB Vault.
// @host localhost:8080
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	router.Use(cors.New(config))

	router.PUT("/document", handlePutDocument)
	router.GET("/documents", handleListDocuments)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// @Summary Create or update a document
// @Description Create a new document or update an existing one
// @Tags documents
// @Accept json
// @Produce json
// @Param document body Document true "Document object"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /document [put]
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

// @Summary List documents
// @Description Get a paginated list of documents
// @Tags documents
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param perPage query int false "Items per page" default(100)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /documents [get]
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
	apiKey := os.Getenv("IMMUDB_API_KEY")

	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return client.Do(req)
}
