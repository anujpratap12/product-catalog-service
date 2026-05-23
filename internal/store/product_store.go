package store

import (
	"sync"

	"sourceasia-backend-assignment/internal/models"
)

var (
	Products = make(map[int64]*models.Product)

	ProductMediaStore = make(map[int64]*models.ProductMedia)

	SKUIndex = make(map[string]int64)

	ProductMutex sync.RWMutex

	NextProductID int64 = 1
)
