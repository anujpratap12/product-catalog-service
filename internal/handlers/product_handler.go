package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"sourceasia-backend-assignment/internal/models"
	"sourceasia-backend-assignment/internal/store"
	"sourceasia-backend-assignment/internal/utils"
)

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var request models.CreateProductRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	request.Name = strings.TrimSpace(request.Name)
	request.SKU = strings.TrimSpace(request.SKU)

	if request.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}

	if request.SKU == "" {
		utils.WriteError(w, http.StatusBadRequest, "sku is required")
		return
	}

	if len(request.ImageURLs) > 20 {
		utils.WriteError(w, http.StatusBadRequest, "maximum 20 image URLs allowed")
		return
	}

	if len(request.VideoURLs) > 20 {
		utils.WriteError(w, http.StatusBadRequest, "maximum 20 video URLs allowed")
		return
	}

	for _, url := range request.ImageURLs {

		if !strings.HasPrefix(url, "http://") &&
			!strings.HasPrefix(url, "https://") {

			utils.WriteError(w, http.StatusBadRequest, "invalid image URL")
			return
		}

		if len(url) > 2048 {
			utils.WriteError(w, http.StatusBadRequest, "image URL too long")
			return
		}
	}

	for _, url := range request.VideoURLs {

		if !strings.HasPrefix(url, "http://") &&
			!strings.HasPrefix(url, "https://") {

			utils.WriteError(w, http.StatusBadRequest, "invalid video URL")
			return
		}

		if len(url) > 2048 {
			utils.WriteError(w, http.StatusBadRequest, "video URL too long")
			return
		}
	}

	store.ProductMutex.Lock()
	defer store.ProductMutex.Unlock()

	_, exists := store.SKUIndex[request.SKU]

	if exists {
		utils.WriteError(w, http.StatusConflict, "sku already exists")
		return
	}

	productID := store.NextProductID
	store.NextProductID++

	product := &models.Product{
		ID:         productID,
		Name:       request.Name,
		SKU:        request.SKU,
		ImageCount: len(request.ImageURLs),
		VideoCount: len(request.VideoURLs),
		CreatedAt:  time.Now(),
	}

	store.Products[productID] = product

	imageURLs := request.ImageURLs
	videoURLs := request.VideoURLs

	if imageURLs == nil {
		imageURLs = []string{}
	}

	if videoURLs == nil {
		videoURLs = []string{}
	}

	store.ProductMediaStore[productID] = &models.ProductMedia{
		ImageURLs: imageURLs,
		VideoURLs: videoURLs,
	}

	store.SKUIndex[request.SKU] = productID

	response := models.ProductDetail{
		ID:         product.ID,
		Name:       product.Name,
		SKU:        product.SKU,
		ImageURLs:  request.ImageURLs,
		VideoURLs:  request.VideoURLs,
		ImageCount: product.ImageCount,
		VideoCount: product.VideoCount,
		CreatedAt:  product.CreatedAt,
	}

	utils.WriteJSON(w, http.StatusCreated, response)
}

func ListProductsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	limit := 10
	offset := 0

	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")

	if limitParam != "" {

		parsedLimit, err := strconv.Atoi(limitParam)

		if err != nil || parsedLimit <= 0 {
			utils.WriteError(w, http.StatusBadRequest, "invalid limit")
			return
		}

		if parsedLimit > 100 {
			parsedLimit = 100
		}

		limit = parsedLimit
	}

	if offsetParam != "" {

		parsedOffset, err := strconv.Atoi(offsetParam)

		if err != nil || parsedOffset < 0 {
			utils.WriteError(w, http.StatusBadRequest, "invalid offset")
			return
		}

		offset = parsedOffset
	}

	store.ProductMutex.RLock()
	defer store.ProductMutex.RUnlock()

	allProducts := make([]*models.Product, 0)

	for _, product := range store.Products {
		allProducts = append(allProducts, product)
	}

	total := len(allProducts)

	if offset >= total {
		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"products": []models.Product{},
			"count":    0,
			"total":    total,
		})

		return
	}

	end := offset + limit

	if end > total {
		end = total
	}

	paginatedProducts := allProducts[offset:end]

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"products": paginatedProducts,
		"count":    len(paginatedProducts),
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

func GetProductHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	idParam := strings.TrimPrefix(r.URL.Path, "/products/")

	productID, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	store.ProductMutex.RLock()
	defer store.ProductMutex.RUnlock()

	product, exists := store.Products[productID]

	if !exists {
		utils.WriteError(w, http.StatusNotFound, "product not found")
		return
	}

	media := store.ProductMediaStore[productID]

	response := models.ProductDetail{
		ID:         product.ID,
		Name:       product.Name,
		SKU:        product.SKU,
		ImageURLs:  media.ImageURLs,
		VideoURLs:  media.VideoURLs,
		ImageCount: product.ImageCount,
		VideoCount: product.VideoCount,
		CreatedAt:  product.CreatedAt,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func AddMediaHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimSuffix(r.URL.Path, "/media")
	idParam := strings.TrimPrefix(path, "/products/")

	productID, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	var request models.AddMediaRequest

	err = json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if len(request.ImageURLs) == 0 && len(request.VideoURLs) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "at least one media field is required")
		return
	}

	if len(request.ImageURLs) > 20 {
		utils.WriteError(w, http.StatusBadRequest, "maximum 20 image URLs allowed")
		return
	}

	if len(request.VideoURLs) > 20 {
		utils.WriteError(w, http.StatusBadRequest, "maximum 20 video URLs allowed")
		return
	}

	for _, url := range request.ImageURLs {

		if !strings.HasPrefix(url, "http://") &&
			!strings.HasPrefix(url, "https://") {

			utils.WriteError(w, http.StatusBadRequest, "invalid image URL")
			return
		}
	}

	for _, url := range request.VideoURLs {

		if !strings.HasPrefix(url, "http://") &&
			!strings.HasPrefix(url, "https://") {

			utils.WriteError(w, http.StatusBadRequest, "invalid video URL")
			return
		}
	}

	store.ProductMutex.Lock()
	defer store.ProductMutex.Unlock()

	product, exists := store.Products[productID]

	if !exists {
		utils.WriteError(w, http.StatusNotFound, "product not found")
		return
	}

	media := store.ProductMediaStore[productID]

	media.ImageURLs = append(media.ImageURLs, request.ImageURLs...)
	media.VideoURLs = append(media.VideoURLs, request.VideoURLs...)

	product.ImageCount = len(media.ImageURLs)
	product.VideoCount = len(media.VideoURLs)

	response := models.ProductDetail{
		ID:         product.ID,
		Name:       product.Name,
		SKU:        product.SKU,
		ImageURLs:  media.ImageURLs,
		VideoURLs:  media.VideoURLs,
		ImageCount: product.ImageCount,
		VideoCount: product.VideoCount,
		CreatedAt:  product.CreatedAt,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
