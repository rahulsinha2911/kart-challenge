package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	productService "github.com/kart-challenge/service/product"
)

func ListProducts(ctx *gin.Context) {
	// Generate a random transaction ID and insert into context
	transactionID := uuid.New().String()
	ctx.Set("transaction_id", transactionID)

	products, err := productService.ListProducts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Products listed successfully",
		"data":    products,
	})
}

func GetProductDetails(ctx *gin.Context) {
	// Generate a random transaction ID and insert into context
	transactionID := uuid.New().String()
	ctx.Set("transaction_id", transactionID)

	productID := ctx.Param("productId")
	if productID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Product ID is required",
		})
		return
	}

	productDetails, err := productService.GetProductDetailsByID(ctx, productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Product details fetched successfully",
		"data":    productDetails,
	})
}
