package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/codenotary/immudb/pkg/api/schema"
	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "backend/docs" // This is where Swag will generate its docs
)

var client immudb.ImmuClient

// @title Accounting Information API
// @version 1.0
// @description This is a simple accounting information API.
// @host localhost:8080
// @BasePath /

func main() {
	immudbAddress := getEnv("IMMUDB_ADDRESS", "immudb")
	immudbPort := getEnv("IMMUDB_PORT", "3322")
	immudbUsername := getEnv("IMMUDB_USERNAME", "immudb")
	immudbPassword := getEnv("IMMUDB_PASSWORD", "immudb")
	immudbDatabase := getEnv("IMMUDB_DATABASE", "defaultdb")

	log.Printf("Connecting to immudb at %s:%s", immudbAddress, immudbPort)

	opts := immudb.DefaultOptions().
		WithAddress(immudbAddress).
		WithPort(3322)

	var err error
	client = immudb.NewClient().WithOptions(opts)

	// Retry logic
	for i := 0; i < 5; i++ {
		log.Printf("Attempt %d to connect to immudb", i+1)
		err = client.OpenSession(
			context.Background(),
			[]byte(immudbUsername),
			[]byte(immudbPassword),
			immudbDatabase,
		)
		if err == nil {
			break
		}
		log.Printf("Failed to open session: %v", err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to open session after 5 attempts: %v", err)
	}

	log.Println("Successfully connected to immudb")

	defer client.CloseSession(context.Background())

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173"} // Add your frontend URL here
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	r.Use(cors.New(config))

	r.POST("/accounts", createAccount)
	r.GET("/accounts", getAccounts)

		// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}


// ErrorResponse represents the structure of an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// AccountsResponse represents the structure of the response for multiple accounts
type AccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

// @Summary Create a new account
// @Description Create a new account with the provided information
// @Accept  json
// @Produce  json
// @Param account body Account true "Account information"
// @Success 201 {object} Account
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts [post]
func createAccount(c *gin.Context) {
	var account Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Store account in immudb
	key := []byte("account:" + account.AccountNumber)
	value, err := json.Marshal(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to marshal account"})
		return
	}

	_, err = client.Set(context.Background(), key, value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to store account"})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// @Summary Get all accounts
// @Description Retrieve all accounts stored in the system
// @Produce  json
// @Success 200 {object} AccountsResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts [get]
func getAccounts(c *gin.Context) {
	// Retrieve all accounts from immudb
	prefix := []byte("account:")
	scanRequest := &schema.ScanRequest{
		Prefix:  prefix,
		SinceTx: 0,
		Limit:   1000, // Adjust this value based on your needs
		Desc:    false,
	}

	entries, err := client.Scan(context.Background(), scanRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve accounts"})
		return
	}

	var accounts []Account
	for _, entry := range entries.Entries {
		var account Account
		err := json.Unmarshal(entry.Value, &account)
		if err != nil {
			log.Printf("Failed to unmarshal account: %v", err)
			continue
		}
		accounts = append(accounts, account)
	}

	c.JSON(http.StatusOK, AccountsResponse{Accounts: accounts})
}

type Account struct {
	AccountNumber string  `json:"account_number"`
	AccountName   string  `json:"account_name"`
	IBAN          string  `json:"iban"`
	Address       string  `json:"address"`
	Amount        string `json:"amount"`
	Type          string  `json:"type"`
}