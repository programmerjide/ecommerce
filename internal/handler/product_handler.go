package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/programmerjide/ecommerce/internal/dto"
	"github.com/programmerjide/ecommerce/internal/service"
	"github.com/programmerjide/ecommerce/internal/utils"
	"github.com/rs/zerolog"
	"strconv"
)

type ProductHandler struct {
	// Placeholder for product handler dependencies
	productService *service.ProductService
	logger         zerolog.Logger
}

func NewProductHandler(productService *service.ProductService, logger zerolog.Logger) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		logger:         logger,
	}
}

func (h *ProductHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid request payload for creating category")
		utils.BadRequestResponse(c, "Invalid request payload", err)
		return
	}

	response, err := h.productService.CreateCategory(&req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create category")
		utils.InternalServerErrorResponse(c, "Failed to create category", err)
		return
	}

	utils.CreatedResponse(c, "Category created successfully", response)
}

func (h *ProductHandler) GetCategories(c *gin.Context) {
	categories, err := h.productService.GetCategories()
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list categories")
		utils.InternalServerErrorResponse(c, "Failed to list categories", err)
		return
	}

	utils.SuccessResponse(c, "Categories retrieved successfully", categories)
}

func (h *ProductHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid category ID")
		utils.BadRequestResponse(c, "Invalid category ID", err)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid request payload for updating category")
		utils.BadRequestResponse(c, "Invalid request payload", err)
		return
	}

	response, err := h.productService.UpdateCategory(uint(id), &req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update category")
		utils.InternalServerErrorResponse(c, "Failed to update category", err)
		return
	}

	utils.SuccessResponse(c, "Category updated successfully", response)
}

func (h *ProductHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid category ID")
		utils.BadRequestResponse(c, "Invalid category ID", err)
		return
	}

	err = h.productService.DeleteCategory(uint(id))
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to delete category")
		utils.InternalServerErrorResponse(c, "Failed to delete category", err)
		return
	}

	utils.SuccessResponse(c, "Category deleted successfully", nil)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid request payload for creating product")
		utils.BadRequestResponse(c, "Invalid request payload", err)
		return
	}

	response, err := h.productService.CreateProduct(&req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create product")
		utils.InternalServerErrorResponse(c, "Failed to create product", err)
		return
	}

	utils.CreatedResponse(c, "Product created successfully", response)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	products, meta, err := h.productService.GetProducts(page, limit)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get products")
		utils.InternalServerErrorResponse(c, "Failed to get products", err)
		return
	}
	utils.PaginatedSuccessResponse(c, "Products retrieved successfully", products, meta)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid product ID")
		utils.BadRequestResponse(c, "Invalid product ID", err)
		return
	}

	product, err := h.productService.GetProduct(uint(id))
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get product")
		utils.InternalServerErrorResponse(c, "Failed to get product", err)
		return
	}

	utils.SuccessResponse(c, "Product retrieved successfully", product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid product ID")
		utils.BadRequestResponse(c, "Invalid product ID", err)
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid request payload for updating product")
		utils.BadRequestResponse(c, "Invalid request payload", err)
		return
	}

	response, err := h.productService.UpdateProduct(uint(id), &req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update product")
		utils.InternalServerErrorResponse(c, "Failed to update product", err)
		return
	}

	utils.SuccessResponse(c, "Product updated successfully", response)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid product ID")
		utils.BadRequestResponse(c, "Invalid product ID", err)
		return
	}

	err = h.productService.DeleteProduct(uint(id))
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to delete product")
		utils.InternalServerErrorResponse(c, "Failed to delete product", err)
		return
	}

	utils.SuccessResponse(c, "Product deleted successfully", nil)
}

func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("category_id"), 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid category ID")
		utils.BadRequestResponse(c, "Invalid category ID", err)
		return
	}

	products, err := h.productService.GetProductsByCategory(uint(categoryID))
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get products by category")
		utils.InternalServerErrorResponse(c, "Failed to get products by category", err)
		return
	}

	utils.SuccessResponse(c, "Products retrieved successfully", products)
}

func (h *ProductHandler) SearchProducts(c *gin.Context) {
	var req dto.SearchProductsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid search query parameters")
		utils.BadRequestResponse(c, "Invalid search query parameters", err)
		return
	}

	results, meta, err := h.productService.SearchProducts(&req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to search products")
		utils.InternalServerErrorResponse(c, "Failed to search products", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "Products search results", results, meta)
}
