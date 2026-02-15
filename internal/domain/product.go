package domain

type Product struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	SKU       string `json:"sku"`
	Price     float64 `json:"price"`
	Stock     int    `json:"stock"`
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type SellProductRequest struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type BulkPriceUpdateItem struct {
	ProductID int     `json:"productId"`
	NewPrice  float64 `json:"newPrice"`
}

