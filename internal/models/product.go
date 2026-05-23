package models

import "time"

type Product struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	SKU        string    `json:"sku"`
	ImageCount int       `json:"image_count"`
	VideoCount int       `json:"video_count"`
	CreatedAt  time.Time `json:"created_at"`
}

type ProductMedia struct {
	ImageURLs []string `json:"image_urls"`
	VideoURLs []string `json:"video_urls"`
}

type CreateProductRequest struct {
	Name      string   `json:"name"`
	SKU       string   `json:"sku"`
	ImageURLs []string `json:"image_urls"`
	VideoURLs []string `json:"video_urls"`
}

type ProductDetail struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	SKU        string    `json:"sku"`
	ImageURLs  []string  `json:"image_urls"`
	VideoURLs  []string  `json:"video_urls"`
	ImageCount int       `json:"image_count"`
	VideoCount int       `json:"video_count"`
	CreatedAt  time.Time `json:"created_at"`
}
type AddMediaRequest struct {
	ImageURLs []string `json:"image_urls"`
	VideoURLs []string `json:"video_urls"`
}
