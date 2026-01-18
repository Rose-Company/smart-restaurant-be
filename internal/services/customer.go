package services

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"time"

	"app-noti/common"
	"app-noti/internal/models"
	"app-noti/pkg/argon2id"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GetProfile - TASK-016: Get customer profile
// Retrieves customer profile information
func (s *Service) GetProfile(ctx context.Context, userID string) (*models.ProfileResponse, error) {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get user", zap.Error(err))
		return nil, common.ErrOrderNotFound // Use a generic not found error or create new one
	}

	resp := &models.ProfileResponse{
		ID:            int(0), // Convert string ID if needed
		Email:         user.Email,
		Name:          s.getFullName(user),
		Phone:         &user.PhoneNumber,
		AvatarURL:     user.AvatarURL,
		DateOfBirth:   user.DateOfBirth,
		Gender:        user.Gender,
		CreatedAt:     user.DateCreated,
		UpdatedAt:     user.UpdatedAt,
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
	}

	// Build address from inline fields
	if user.StreetAddress != nil || user.City != nil || user.State != nil || user.PostalCode != nil || user.Country != nil {
		resp.Address = &models.AddressResponse{
			Street:     user.StreetAddress,
			City:       user.City,
			State:      user.State,
			PostalCode: user.PostalCode,
			Country:    user.Country,
		}
	}

	// Get preferences
	if user.Preferences != nil {
		prefs := &models.PreferencesResponse{
			NotificationEnabled: user.Preferences.NotificationEnabled,
			EmailNotifications:  user.Preferences.EmailNotifications,
			SMSNotifications:    user.Preferences.SMSNotifications,
		}

		if err := json.Unmarshal([]byte(user.Preferences.DietaryRestrictions), &prefs.DietaryRestrictions); err != nil {
			s.logger.Debug("failed to unmarshal dietary restrictions", zap.Error(err))
		}

		if err := json.Unmarshal([]byte(user.Preferences.FavoriteCuisines), &prefs.FavoriteCuisines); err != nil {
			s.logger.Debug("failed to unmarshal favorite cuisines", zap.Error(err))
		}

		resp.Preferences = prefs
	}

	// Get loyalty info
	resp.Loyalty = s.calculateLoyalty(user.LoyaltyPoints)

	// Get statistics
	resp.Statistics = &models.StatisticsResponse{
		TotalOrders: user.TotalOrders,
		TotalSpent:  user.TotalSpent,
	}

	return resp, nil
}

// UpdateProfile - TASK-017: Update customer profile
// Updates customer profile information (name, email, phone, address, preferences)
func (s *Service) UpdateProfile(ctx context.Context, userID string, req *models.UpdateProfileRequest) (*models.ProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get user", zap.Error(err))
		return nil, common.ErrOrderNotFound
	}

	// Update basic fields
	if req.Name != nil {
		parts := s.splitName(*req.Name)
		user.FirstName = &parts[0]
		if len(parts) > 1 {
			user.LastName = &parts[1]
		}
	}

	if req.Phone != nil {
		// Check if phone already exists
		existing, err := s.userRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
			tx.Where("phone_number = ? AND id != ?", req.Phone, userID)
		})
		if err == nil && existing != nil {
			return nil, errors.New("phone number already in use")
		}
		user.PhoneNumber = *req.Phone
	}

	if req.DateOfBirth != nil {
		user.DateOfBirth = req.DateOfBirth
	}

	if req.Gender != nil {
		user.Gender = req.Gender
	}

	// Update inline address fields
	if req.Address != nil {
		if req.Address.Street != nil {
			user.StreetAddress = req.Address.Street
		}
		if req.Address.City != nil {
			user.City = req.Address.City
		}
		if req.Address.State != nil {
			user.State = req.Address.State
		}
		if req.Address.PostalCode != nil {
			user.PostalCode = req.Address.PostalCode
		}
		if req.Address.Country != nil {
			user.Country = req.Address.Country
		}
	}

	// Update user (includes address fields)
	_, err = s.userRepo.Update(ctx, userID, user)
	if err != nil {
		s.logger.Error("failed to update user", zap.Error(err))
		return nil, err
	}

	// Update preferences
	if req.Preferences != nil {
		prefs := user.Preferences
		if prefs == nil {
			prefs = &models.UserPreferences{CustomerID: userID}
		}

		if len(req.Preferences.DietaryRestrictions) > 0 {
			data, _ := json.Marshal(req.Preferences.DietaryRestrictions)
			prefs.DietaryRestrictions = string(data)
		}

		if len(req.Preferences.FavoriteCuisines) > 0 {
			data, _ := json.Marshal(req.Preferences.FavoriteCuisines)
			prefs.FavoriteCuisines = string(data)
		}

		if req.Preferences.NotificationEnabled != nil {
			prefs.NotificationEnabled = *req.Preferences.NotificationEnabled
		}

		if req.Preferences.EmailNotifications != nil {
			prefs.EmailNotifications = *req.Preferences.EmailNotifications
		}

		if req.Preferences.SMSNotifications != nil {
			prefs.SMSNotifications = *req.Preferences.SMSNotifications
		}

		if prefs.ID == 0 {
			_, err := s.preferencesRepo.Create(ctx, prefs)
			if err != nil {
				s.logger.Error("failed to create preferences", zap.Error(err))
			}
		} else {
			_, err := s.preferencesRepo.Update(ctx, prefs.ID, prefs)
			if err != nil {
				s.logger.Error("failed to update preferences", zap.Error(err))
			}
		}
	}

	// Return updated profile
	return s.GetProfile(ctx, userID)
}

// ChangePassword - TASK-019: Change customer password
// Changes customer account password
func (s *Service) ChangePassword(ctx context.Context, userID string, req *models.ChangePasswordRequest) (*models.ChangePasswordResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get user", zap.Error(err))
		return nil, common.ErrOrderNotFound
	}

	// Verify old password
	match, err := argon2id.ComparePasswordAndHash(req.OldPassword, user.Password)
	if err != nil || !match {
		return nil, errors.New("current password is incorrect")
	}

	// Validate new password strength
	if !s.validatePasswordStrength(req.NewPassword) {
		return nil, errors.New("new password does not meet security requirements")
	}

	// Hash new password
	hashedPassword, err := argon2id.CreateHash(req.NewPassword, argon2id.DefaultParams)
	if err != nil {
		s.logger.Error("failed to hash password", zap.Error(err))
		return nil, err
	}

	user.Password = hashedPassword
	_, err = s.userRepo.Update(ctx, userID, user)
	if err != nil {
		s.logger.Error("failed to update user password", zap.Error(err))
		return nil, err
	}

	return &models.ChangePasswordResponse{
		ChangedAt: time.Now(),
	}, nil
}

// GetCustomerReviews - TASK-020: Get customer reviews
// Retrieves customer reviews with pagination
func (s *Service) GetCustomerReviews(ctx context.Context, userID string, page, pageSize int, sort string) (*models.ReviewListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get total count
	var total int64
	total, err := s.reviewRepo.Count(ctx, models.QueryParams{}, func(tx *gorm.DB) {
		tx.Where("user_id = ?", userID)
	})
	if err != nil {
		s.logger.Error("failed to count reviews", zap.Error(err))
	}

	// Get reviews with sorting
	reviews, err := s.reviewRepo.ListByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("user_id = ?", userID)

		switch sort {
		case "created_at_asc":
			tx.Order("created_at ASC")
		case "rating_desc":
			tx.Order("rating DESC")
		case "rating_asc":
			tx.Order("rating ASC")
		default:
			tx.Order("created_at DESC")
		}

		tx.Offset(offset).Limit(pageSize)
	})

	if err != nil {
		s.logger.Error("failed to get reviews", zap.Error(err))
		return nil, err
	}

	// Convert to response
	items := make([]models.ReviewDetailResponse, len(reviews))
	averageRating := 0.0
	ratingBreakdown := models.RatingBreakdown{}

	for i, review := range reviews {
		items[i] = s.buildReviewResponse(ctx, review)
		averageRating += float64(review.Rating)

		// Count rating breakdown
		switch review.Rating {
		case 5:
			ratingBreakdown.FiveStar++
		case 4:
			ratingBreakdown.FourStar++
		case 3:
			ratingBreakdown.ThreeStar++
		case 2:
			ratingBreakdown.TwoStar++
		case 1:
			ratingBreakdown.OneStar++
		}
	}

	if len(reviews) > 0 {
		averageRating = averageRating / float64(len(reviews))
	}

	return &models.ReviewListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    items,
		Extra: models.ReviewExtraInfo{
			AverageRating:   math.Round(averageRating*10) / 10,
			TotalReviews:    total,
			RatingBreakdown: ratingBreakdown,
		},
	}, nil
}

// Helper functions

func (s *Service) getFullName(user *models.User) *string {
	if user.FirstName == nil && user.LastName == nil {
		return nil
	}

	name := ""
	if user.FirstName != nil {
		name = *user.FirstName
	}
	if user.LastName != nil {
		if name != "" {
			name += " "
		}
		name += *user.LastName
	}

	return &name
}

func (s *Service) splitName(fullName string) []string {
	// Simple split on last space
	var parts []string
	if len(fullName) > 0 {
		parts = append(parts, fullName)
	}
	return parts
}

func (s *Service) validatePasswordStrength(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper, hasLower, hasNumber, hasSpecial := false, false, false, false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		default:
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func (s *Service) calculateLoyalty(points int) *models.LoyaltyResponse {
	tier := "bronze"
	nextTier := "silver"
	pointsToNext := 100

	if points >= 500 {
		tier = "gold"
		nextTier = "platinum"
		pointsToNext = 1000 - (points % 1000)
	} else if points >= 100 {
		tier = "silver"
		nextTier = "gold"
		pointsToNext = 500 - (points % 500)
	}

	return &models.LoyaltyResponse{
		Points:       points,
		Tier:         tier,
		NextTier:     nextTier,
		PointsToNext: pointsToNext,
	}
}

func (s *Service) buildReviewResponse(ctx context.Context, review *models.Review) models.ReviewDetailResponse {
	resp := models.ReviewDetailResponse{
		ID:        review.ID,
		OrderID:   review.OrderID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		Status:    review.Status,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}

	// Add order number if order exists
	if review.Order != nil {
		resp.OrderNumber = review.Order.OrderNumber
	}

	// Add photos
	if len(review.Photos) > 0 {
		resp.Photos = make([]models.ReviewPhotoResponse, len(review.Photos))
		for i, photo := range review.Photos {
			resp.Photos[i] = models.ReviewPhotoResponse{
				ID:           photo.ID,
				URL:          photo.URL,
				ThumbnailURL: photo.ThumbnailURL,
			}
		}
	}

	// Add items reviewed
	if len(review.ItemsReviewed) > 0 {
		resp.ItemsReviewed = make([]models.ReviewItemResponse, len(review.ItemsReviewed))
		for i, item := range review.ItemsReviewed {
			itemResp := models.ReviewItemResponse{
				MenuItemID:  item.MenuItemID,
				ItemRating:  item.ItemRating,
				ItemComment: item.ItemComment,
			}

			if item.MenuItem != nil {
				itemResp.MenuItemName = item.MenuItem.Name
				// Get menu item image from database if needed
				itemResp.MenuItemImage = ""
			}

			resp.ItemsReviewed[i] = itemResp
		}
	}

	// Add restaurant response
	if review.RestaurantResponse != nil {
		resp.RestaurantResponse = &models.RestaurantResponseView{
			Message:     review.RestaurantResponse.Message,
			RespondedBy: review.RestaurantResponse.RespondedBy,
			RespondedAt: review.RestaurantResponse.RespondedAt,
		}
	}

	return resp
}
