package orderService

import (
	"bytes"
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

func PlaceOrder(ctx *gin.Context, orderRequestData model.OrderRequestData) (model.OrderResponseData, error) {
	if utils.ShouldReadFromServer() {
		orderResponseData, err := getOrderResponseDataFromServer(ctx, orderRequestData)
		if err != nil {
			return model.OrderResponseData{}, err
		}
		return orderResponseData, nil
	} else {
		//read from DB
		return model.OrderResponseData{}, nil
	}
}

func getOrderResponseDataFromServer(ctx *gin.Context, orderRequestData model.OrderRequestData) (model.OrderResponseData, error) {
	// Check if context is cancelled before proceeding
	select {
	case <-ctx.Done():
		return model.OrderResponseData{}, ctx.Err()
	default:
		// Continue with server call
	}

	baseEndpoint := os.Getenv("ORDER_ENDPOINT")
	if baseEndpoint == "" {
		return model.OrderResponseData{}, errors.New("ORDER_ENDPOINT environment variable is not set")
	}

	timeoutStr := os.Getenv("ORDER_ENDPOINT_TIMEOUT")
	if timeoutStr == "" {
		return model.OrderResponseData{}, errors.New("ORDER_ENDPOINT_TIMEOUT environment variable is not set")
	}
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return model.OrderResponseData{}, err
	}

	requestBody, err := json.Marshal(orderRequestData)
	if err != nil {
		return model.OrderResponseData{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create HTTP request with context for proper cancellation
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return model.OrderResponseData{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if authToken := os.Getenv("ORDER_SERVICE_TOKEN"); authToken != "" {
		req.Header.Set("Authorization", authToken)
	}

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	response, err := client.Do(req)
	if err != nil {
		return model.OrderResponseData{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return model.OrderResponseData{}, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return model.OrderResponseData{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var orderResponseData model.OrderResponseData
	if err := json.Unmarshal(body, &orderResponseData); err != nil {
		return model.OrderResponseData{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return orderResponseData, nil
}
