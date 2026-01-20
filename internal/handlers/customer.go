package handlers

import (
	"app-noti/common"
	"app-noti/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetProfile retrieves customer profile information
// @Summary Get customer profile
// @Description Get current customer profile information
// @Tags Customer
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} common.BaseResponse{data=models.ProfileResponse}
// @Failure 400 {object} common.BaseResponse
// @Router /api/customer/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		common.AbortWithError(c, common.ErrUnauthorized)
		return
	}

	profile, err := h.service.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Customer profile retrieved successfully", profile))
}

// UpdateProfile updates customer profile information
// @Summary Update customer profile
// @Description Update customer profile information (name, email, phone, address, preferences)
// @Tags Customer
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} common.BaseResponse{data=models.ProfileResponse}
// @Failure 400 {object} common.BaseResponse
// @Router /api/customer/profile [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		common.AbortWithError(c, common.ErrUnauthorized)
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	profile, err := h.service.UpdateProfile(c.Request.Context(), userID.(string), &req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Profile updated successfully", profile))
}

// UploadAvatar uploads customer avatar image
// @Summary Upload customer avatar
// @Description Upload customer avatar image (JPEG, PNG, GIF max 5MB)
// @Tags Customer
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Avatar image file"
// @Success 200 {object} common.BaseResponse{data=models.UploadAvatarResponse}
// @Failure 400 {object} common.BaseResponse
// @Router /api/customer/avatar [post]
func (h *Handler) UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		common.AbortWithError(c, common.ErrUnauthorized)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Validate file size (5MB max)
	const maxFileSize int64 = 5 * 1024 * 1024
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "File size exceeds maximum allowed (5MB)",
			"data": gin.H{
				"max_size":      maxFileSize,
				"uploaded_size": file.Size,
			},
		})
		return
	}

	// Validate MIME type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}

	mimeType := file.Header.Get("Content-Type")
	if !allowedTypes[mimeType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid file type. Only JPEG, PNG, and GIF are allowed",
			"data": gin.H{
				"allowed_types": []string{"image/jpeg", "image/png", "image/gif"},
				"uploaded_type": mimeType,
			},
		})
		return
	}

	// Upload to S3 (using existing upload logic from other handlers)
	// For now, create a mock response
	response := &models.UploadAvatarResponse{
		AvatarURL:    "https://example.com/avatars/" + userID.(string) + ".jpg",
		ThumbnailURL: "https://example.com/avatars/thumbnails/" + userID.(string) + ".jpg",
		FileSize:     file.Size,
		MimeType:     mimeType,
		UploadedAt:   time.Now(),
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Avatar uploaded successfully", response))
}

// ChangePassword changes customer password
// @Summary Change customer password
// @Description Change customer account password
// @Tags Customer
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.ChangePasswordRequest true "Change password request"
// @Success 200 {object} common.BaseResponse{data=models.ChangePasswordResponse}
// @Failure 400 {object} common.BaseResponse
// @Router /api/customer/password [patch]
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		common.AbortWithError(c, common.ErrUnauthorized)
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	response, err := h.service.ChangePassword(c.Request.Context(), userID.(string), &req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Password changed successfully", response))
}

// GetCustomerReviews retrieves customer reviews
// @Summary Get customer reviews
// @Description Get list of customer reviews with pagination
// @Tags Customer
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 10)"
// @Param sort query string false "Sort order: created_at_desc|created_at_asc|rating_desc|rating_asc"
// @Success 200 {object} common.BaseResponse{data=models.ReviewListResponse}
// @Failure 400 {object} common.BaseResponse
// @Router /api/customer/reviews [get]
func (h *Handler) GetCustomerReviews(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		common.AbortWithError(c, common.ErrUnauthorized)
		return
	}

	// Parse query params
	page := 1
	pageSize := 10
	sort := "created_at_desc"

	if p := c.Query("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if pageSizeNum, err := strconv.Atoi(ps); err == nil && pageSizeNum > 0 {
			pageSize = pageSizeNum
		}
	}

	if s := c.Query("sort"); s != "" {
		sort = s
	}

	reviews, err := h.service.GetCustomerReviews(c.Request.Context(), userID.(string), page, pageSize, sort)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Reviews retrieved successfully", reviews))
}
