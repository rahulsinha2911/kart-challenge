package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kart-challenge/model"
	"github.com/kart-challenge/utils"
)

const defaultTimeoutSeconds = 30

func ListProducts(ctx *gin.Context) ([]model.Product, error) {
	if utils.ShouldReadFromServer() {
		//read from server
		productListFromServer, err := getProductsFromServer(ctx)
		if err != nil {
			log.Println("Error getting products list from server:", err)
			return nil, err
		}
		return productListFromServer, nil
	} else {
		//read from DB
		return []model.Product{}, nil
	}
}

func getProductsFromServer(ctx *gin.Context) ([]model.Product, error) {
	// Check if context is cancelled before proceeding
	select {
	case <-ctx.Done():
		log.Println("Context cancelled")
		return nil, ctx.Err()
	default:
		// Continue with server call
	}

	endpoint, timeout, err := getServerConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get server configuration: %w", err)
	}

	// Create HTTP request with context for proper cancellation
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{
		Timeout: timeout,
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var products []model.Product
	if err := json.Unmarshal(body, &products); err != nil {
		return nil, fmt.Errorf("failed to parse products JSON: %w", err)
	}
	return products, nil
}

// getServerConfig retrieves server endpoint and timeout from environment variables
func getServerConfig() (string, time.Duration, error) {
	endpoint := os.Getenv("PRODUCT_LIST_ENDPOINT")
	if endpoint == "" {
		return "", 0, errors.New("PRODUCT_LIST_ENDPOINT environment variable is not set")
	}

	timeoutStr := os.Getenv("PRODUCT_LIST_ENDPOINT_TIMEOUT")
	if timeoutStr == "" {
		return endpoint, defaultTimeoutSeconds * time.Second, nil
	}

	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid timeout value '%s': %w", timeoutStr, err)
	}

	return endpoint, time.Duration(timeout) * time.Second, nil
}
