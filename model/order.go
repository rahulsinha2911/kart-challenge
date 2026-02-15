package model

// {
// 	"couponCode": "",
// 	"items": [
// 	  {
// 		"productId": "",
// 		"quantity": 1
// 	  }
// 	]
//   }

type OrderRequestData struct {
	CouponCode string `json:"couponCode"`
	Items      []Item `json:"items"`
}

type Item struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

////////
// {
// 	"id": "0000-0000-0000-0000",
// 	"items": [
// 	  {
// 		"productId": "string",
// 		"quantity": 1
// 	  }
// 	],
// 	"products": [
// 	  {
// 		"id": "10",
// 		"name": "Chicken Waffle",
// 		"price": 1,
// 		"category": "Waffle"
// 	  }
// 	]
//   }

type OrderResponseData struct {
	ID       string    `json:"id"`
	Items    []Item    `json:"items"`
	Products []Product `json:"products"`
}
