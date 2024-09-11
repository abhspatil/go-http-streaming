package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/abhspatil/go-http-streaming/domain"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func RandomStatus() string {
	statuses := []string{"PENDING", "COMPLETED", "FAILED"}
	return statuses[rand.Intn(len(statuses))]
}

func RandomNetwork() string {
	return "mainnet"
}

func GenerateRandomTransactions(count int) []domain.TransactionHistoryData {
	rand.Seed(time.Now().UnixNano())
	var transactions []domain.TransactionHistoryData
	for i := 0; i < count; i++ {
		tx := domain.TransactionHistoryData{
			TransactionId: strconv.Itoa(rand.Intn(100000000)),
			Amount:        strconv.Itoa(rand.Intn(100000000)),
			Network:       RandomNetwork(),
			CreatedOn:     time.Now(),
		}
		transactions = append(transactions, tx)
	}
	return transactions
}

func GetTransactionsInChunks(c *gin.Context) {
	totalRecords := 10000
	chunkSize := 1000

	for i := 0; i < totalRecords; i += chunkSize {
		end := i + chunkSize
		if end > totalRecords {
			end = totalRecords
		}

		// Generate a chunk of random transactions
		transactions := GenerateRandomTransactions(end - i)
		c.JSON(http.StatusOK, transactions)

		// Flush the response to send the current chunk
		c.Writer.Flush()
		time.Sleep(1 * time.Second) // Simulate processing delay
	}
}

func ProcessChunkData(c *gin.Context) {
	resp, err := http.Get("http://localhost:8080/transactions")
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Create a buffered reader to handle chunked data
	reader := bufio.NewReader(resp.Body)

	// Loop through the chunks and process them as they come
	decoder := json.NewDecoder(reader)
	for decoder.More() {
		var transactions []domain.TransactionHistoryData
		err := decoder.Decode(&transactions)
		if err != nil {
			fmt.Printf("Error decoding JSON: %v\n", err)
			return
		}

		// Process each chunk of transactions
		for _, tx := range transactions {
			fmt.Printf("TxID: %s, Network: %s, Price: %s\n", tx.TransactionId, tx.Network, tx.Amount)
		}

		fmt.Printf("/n/n/n Processed chuck of length %d\n", len(transactions))
	}

	c.JSON(http.StatusOK, gin.H{"message": "All chunks processed"})
}

func main() {

	r := gin.Default()

	r.GET("/transactions", GetTransactionsInChunks)
	r.GET("/transactions/process", ProcessChunkData)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
