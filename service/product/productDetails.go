package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kart-challenge/model"
	"github.com/kart-challenge/utils"
)

const defaultProductDetailsTimeoutSeconds = 30

func GetProductDetailsByID(ctx *gin.Context, productID string) (model.Product, error) {
	if utils.ShouldReadFromServer() {
		productDetails, err := getProductDetailsByIDFromServer(ctx, productID)
		if err != nil {
			return model.Product{}, err
		}
		return productDetails, nil
	} else {
		//read from DB
		return model.Product{}, nil
	}
}

func getProductDetailsByIDFromServer(ctx *gin.Context, productID string) (model.Product, error) {
	// Check if context is cancelled before proceeding
	select {
	case <-ctx.Done():
		return model.Product{}, ctx.Err()
	default:
		// Continue with server call
	}

	baseEndpoint, timeout, err := getProductDetailsServerConfig()
	if err != nil {
		return model.Product{}, fmt.Errorf("failed to get server configuration: %w", err)
	}

	endpoint := fmt.Sprintf("%s/%s", baseEndpoint, productID)

	// Create HTTP request with context for proper cancellation
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return model.Product{}, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{
		Timeout: timeout,
	}

	response, err := client.Do(req)
	if err != nil {
		return model.Product{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return model.Product{}, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return model.Product{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var product model.Product
	if err := json.Unmarshal(body, &product); err != nil {
		return model.Product{}, fmt.Errorf("failed to parse product JSON: %w", err)
	}

	return product, nil
}

// getProductDetailsServerConfig retrieves server endpoint and timeout from environment variables
func getProductDetailsServerConfig() (string, time.Duration, error) {
	baseEndpoint := os.Getenv("PRODUCT_DETAILS_ENDPOINT")
	if baseEndpoint == "" {
		return "", 0, errors.New("PRODUCT_DETAILS_ENDPOINT environment variable is not set")
	}

	timeoutStr := os.Getenv("PRODUCT_DETAILS_ENDPOINT_TIMEOUT")
	if timeoutStr == "" {
		return baseEndpoint, defaultProductDetailsTimeoutSeconds * time.Second, nil
	}

	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid timeout value '%s': %w", timeoutStr, err)
	}

	return baseEndpoint, time.Duration(timeout) * time.Second, nil
}
