package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/kart-challenge/globals"
	"github.com/kart-challenge/model"
	orderService "github.com/kart-challenge/service/order"
)

const couponCodesFilePath = "assets/couponCodes.json"

var (
	couponCodes     map[string]float64
	couponCodesOnce sync.Once
	couponCodesErr  error
)

func PlaceOrder(ctx *gin.Context) {
	orderRequestData := model.OrderRequestData{}
	if err := ctx.ShouldBindJSON(&orderRequestData); err != nil {
		ctx.JSON(http.StatusBadRequest, model.GenericErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Please provide valid order request data",
		})
		return
	}

	err := validateOrderRequestData(orderRequestData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.GenericErrorResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: "Please provide valid order request data values: " + err.Error(),
			Ecode:   globals.ERROR_ORDER_PLACE_REQUEST_DATA_INVALID,
			Edesc:   globals.ERROR_ORDER_PLACE_REQUEST_DATA_INVALID_DESC,
		})
	}

	orderResponseData, err := orderService.PlaceOrder(ctx, orderRequestData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.GenericSuccessResponse{
		Status:  http.StatusOK,
		Message: "Order placed successfully",
		Data:    orderResponseData,
	})
}

func validateOrderRequestData(orderRequestData model.OrderRequestData) error {
	if orderRequestData.CouponCode != "" && !validateCouponCode(orderRequestData.CouponCode) {
		return errors.New("invalid coupon code")
	}
	if len(orderRequestData.Items) == 0 {
		return errors.New("items are required")
	}

	for _, item := range orderRequestData.Items {
		if item.ProductID == "" {
			return errors.New("product id is required")
		}
		if item.Quantity <= 0 {
			return errors.New("quantity must be greater than 0")
		}
	}
	return nil
}

// loadCouponCodes loads coupon codes from JSON file (thread-safe, loads once)
func loadCouponCodes() error {
	couponCodesOnce.Do(func() {
		wd, err := os.Getwd()
		if err != nil {
			couponCodesErr = err
			return
		}

		// Construct the full path to the coupon codes file
		filePath := filepath.Join(wd, couponCodesFilePath)

		data, err := os.ReadFile(filePath)
		if err != nil {
			couponCodesErr = err
			return
		}

		codes := make(map[string]float64)
		if err := json.Unmarshal(data, &codes); err != nil {
			couponCodesErr = err
			return
		}

		couponCodes = codes
	})
	return couponCodesErr
}

// validateCouponCode checks if the provided coupon code exists in the coupon codes file
func validateCouponCode(couponCode string) bool {
	// Load coupon codes if not already loaded
	if err := loadCouponCodes(); err != nil {
		return false
	}

	_, exists := couponCodes[couponCode]
	return exists
}
