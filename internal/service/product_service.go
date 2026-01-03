package service

import (
	"github.com/programmerjide/ecommerce/internal/dto"
	"github.com/programmerjide/ecommerce/internal/models"
	"github.com/programmerjide/ecommerce/internal/utils"
	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{
		db: db,
	}
}

func (s *ProductService) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	// Implementation for creating a product category
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.db.Create(&category).Error; err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
	}, nil
}

func (s *ProductService) GetCategories() ([]dto.CategoryResponse, error) {
	var categories []models.Category
	if err := s.db.Find(&categories).Error; err != nil {
		return nil, err
	}

	if err := s.db.Where("is_active = ?", true).Find(&categories).Error; err != nil {
		return nil, err
	}

	response := make([]dto.CategoryResponse, len(categories))
	for i, category := range categories {
		response[i] = dto.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			IsActive:    category.IsActive,
		}
	}
	return response, nil
}

func (s *ProductService) UpdateCategory(categoryID uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	// Implementation for updating a product category
	var category models.Category
	if err := s.db.First(&category, categoryID).Error; err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.db.Save(&category).Error; err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
	}, nil
}

func (s *ProductService) DeleteCategory(categoryID uint) error {
	// Implementation for deleting a product category
	if err := s.db.Delete(&models.Category{}, categoryID).Error; err != nil {
		return err
	}
	return nil
}

func (s *ProductService) CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	// Implementation for creating a product
	product := &models.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.SKU,
	}

	if err := s.db.Create(&product).Error; err != nil {
		return nil, err
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		SKU:         product.SKU,
		IsActive:    product.IsActive,
	}, nil
}

func (s *ProductService) GetProducts(page, limit int) ([]dto.ProductResponse, *utils.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	var products []models.Product
	var total int64

	s.db.Model(&models.Product{}).Where("is_active = ?", true).Count(&total)

	if err := s.db.Preload("Category").Preload("Images").
		Where("is_active = ?", true).
		Offset(offset).Limit(limit).
		Find(&products).Error; err != nil {
		return nil, nil, err
	}

	response := make([]dto.ProductResponse, len(products))
	for i := range products {
		response[i] = s.convertToProductResponse(&products[i])
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
	}

	return response, meta, nil
}

func (s *ProductService) GetProduct(id uint) (*dto.ProductResponse, error) {
	var product models.Product
	if err := s.db.Preload("Category").Preload("Images").First(&product, id).Error; err != nil {
		return nil, err
	}

	response := s.convertToProductResponse(&product)
	return &response, nil
}

func (s *ProductService) UpdateProduct(productID uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	// Implementation for updating a product
	var product models.Product
	if err := s.db.First(&product, productID).Error; err != nil {
		return nil, err
	}

	product.CategoryID = req.CategoryID
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	product.IsActive = *req.IsActive

	if err := s.db.Save(&product).Error; err != nil {
		return nil, err
	}

	return s.GetProduct(productID)
}

func (s *ProductService) DeleteProduct(productID uint) error {
	// Implementation for deleting a product
	if err := s.db.Delete(&models.Product{}, productID).Error; err != nil {
		return err
	}
	return nil
}

func (s *ProductService) GetProductsByCategory(categoryID uint) ([]dto.ProductResponse, error) {
	var products []models.Product
	if err := s.db.Where("category_id = ? AND is_active = ?", categoryID, true).Preload("Category").Preload("Images").Find(&products).Error; err != nil {
		return nil, err
	}

	response := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		images := make([]dto.ProductImageResponse, len(product.Images))
		for j, img := range product.Images {
			images[j] = dto.ProductImageResponse{
				ID:        img.ID,
				URL:       img.URL,
				AltText:   img.AltText,
				IsPrimary: img.IsPrimary,
			}
		}

		response[i] = dto.ProductResponse{
			ID:          product.ID,
			CategoryID:  product.CategoryID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			SKU:         product.SKU,
			IsActive:    product.IsActive,
			Category: dto.CategoryResponse{
				ID:          product.Category.ID,
				Name:        product.Category.Name,
				Description: product.Category.Description,
				IsActive:    product.Category.IsActive,
			},
			Images: images,
		}
	}
	return response, nil
}

func (s *ProductService) GetProductByID(productID uint) (*dto.ProductResponse, error) {
	var product models.Product
	if err := s.db.Preload("Category").Preload("Images").First(&product, productID).Error; err != nil {
		return nil, err
	}

	images := make([]dto.ProductImageResponse, len(product.Images))
	for j, img := range product.Images {
		images[j] = dto.ProductImageResponse{
			ID:        img.ID,
			URL:       img.URL,
			AltText:   img.AltText,
			IsPrimary: img.IsPrimary,
		}
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		SKU:         product.SKU,
		IsActive:    product.IsActive,
		Category: dto.CategoryResponse{
			ID:          product.Category.ID,
			Name:        product.Category.Name,
			Description: product.Category.Description,
			IsActive:    product.Category.IsActive,
		},
		Images: images,
	}, nil
}

func (s *ProductService) SearchProducts(req *dto.SearchProductsRequest) ([]dto.ProductSearchResult, *utils.PaginationMeta, error) {

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 1 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit

	// build query
	query := s.db.Model(&models.Product{}).
		Select("products.*, ts_rank(search_vector, plainto_tsquery('english', ?)) as rank", req.Query).
		Where("search_vector @@ plainto_tsquery('english', ?)", req.Query).
		Where("is_active = ?", true)

	if req.CategoryID != nil {
		query = query.Where("category_id = ?", *req.CategoryID)
	}

	if req.MinPrice != nil {
		query = query.Where("price >= ?", *req.MinPrice)
	}

	if req.MaxPrice != nil {
		query = query.Where("price <= ?", *req.MaxPrice)
	}

	// Count total results
	var total int64
	query.Count(&total)

	// Execute query with ranking and create product slices
	type productsWithRank struct {
		models.Product
		Rank float32 `gorm:"column:rank"`
	}
	var rows []productsWithRank
	if err := query.
		Order("rank DESC, created_at DESC"). // order by relevance
		Preload("Category").
		Preload("Images").
		Offset(offset).
		Limit(req.Limit).
		Find(&rows).Error; err != nil {
		return nil, nil, err
	}

	// Build output response
	results := make([]dto.ProductSearchResult, len(rows))
	for i := range rows {
		results[i] = dto.ProductSearchResult{
			ProductResponse: s.convertToProductResponse(&rows[i].Product),
			Rank:            rows[i].Rank,
		}
	}

	// build pagination meta
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))
	meta := &utils.PaginationMeta{
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      int(total),
		TotalPages: totalPages,
	}

	return results, meta, nil
}

func (s *ProductService) convertToProductResponse(product *models.Product) dto.ProductResponse {
	images := make([]dto.ProductImageResponse, len(product.Images))
	for i := range product.Images {
		images[i] = dto.ProductImageResponse{
			ID:        product.Images[i].ID,
			URL:       product.Images[i].URL,
			AltText:   product.Images[i].AltText,
			IsPrimary: product.Images[i].IsPrimary,
			CreatedAt: product.Images[i].CreatedAt,
		}
	}

	return dto.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		SKU:         product.SKU,
		IsActive:    product.IsActive,
		Category: dto.CategoryResponse{
			ID:          product.Category.ID,
			Name:        product.Category.Name,
			Description: product.Category.Description,
			IsActive:    product.Category.IsActive,
			CreatedAt:   product.Category.CreatedAt,
			UpdatedAt:   product.Category.UpdatedAt,
		},
		Images:    images,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}
