package services

import (
	"app-noti/internal/models"
	"app-noti/internal/repositories"
	"app-noti/pkg/argon2id"
	"app-noti/pkg/utils"
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ListStaff - TASK-021: List staff with filters and pagination
func (s *Service) ListStaff(ctx context.Context, req models.ListStaffRequest) (*models.StaffListResponse, error) {
	page, pageSize := utils.GetPageAndPageSize(req.Page, req.PageSize)

	filters := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Joins("JOIN roles ON users.role = roles.id").
				Where("roles.is_staff = TRUE")
		},
	}

	// Filter by role
	if req.Role != nil && *req.Role != "" {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("roles.name = ?", *req.Role)
		})
	}

	// Filter by status
	if req.Status != nil && *req.Status != "" {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("users.status = ?", *req.Status)
		})
	}

	// Search by name or email
	if req.Search != nil && *req.Search != "" {
		searchTerm := "%" + *req.Search + "%"
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("users.email ILIKE ? OR users.first_name ILIKE ? OR users.last_name ILIKE ?",
				searchTerm, searchTerm, searchTerm)
		})
	}

	// Get total count
	totalCount, err := s.userRepo.Count(ctx, models.QueryParams{}, filters...)
	if err != nil {
		return nil, err
	}

	if totalCount == 0 {
		return &models.StaffListResponse{
			Total:    0,
			Page:     page,
			PageSize: pageSize,
			Items:    []models.StaffResponse{},
		}, nil
	}

	queryParams := models.QueryParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
		QuerySort: models.QuerySort{
			Origin: "created_at.desc",
		},
	}

	users, err := s.userRepo.List(ctx, queryParams, filters...)
	if err != nil {
		return nil, err
	}

	// Build staff responses
	items := make([]models.StaffResponse, 0, len(users))
	for _, user := range users {
		staffResp, err := s.buildStaffResponse(ctx, user, false)
		if err != nil {
			s.logger.Error("Failed to build staff response", zap.String("user_id", user.ID), zap.Error(err))
			continue
		}
		items = append(items, *staffResp)
	}

	return &models.StaffListResponse{
		Total:    totalCount,
		Page:     page,
		PageSize: pageSize,
		Items:    items,
	}, nil
}

// CreateStaff - TASK-022: Create staff account
func (s *Service) CreateStaff(ctx context.Context, req models.CreateStaffRequest, adminID string) (*models.StaffResponse, error) {
	// Check if email already exists
	existing, _ := s.userRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("email = ?", req.Email)
	})
	if existing != nil {
		return nil, fmt.Errorf("email already in use")
	}

	// Get role by name
	role, err := s.roleRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("name = ?", req.Role)
	})
	if err != nil || role == nil {
		return nil, fmt.Errorf("invalid role")
	}

	// Hash password
	hashedPassword, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Split name into first and last name
	firstName, lastName := splitName(req.Name)

	// Create user
	user := &models.User{
		Email:       req.Email,
		FirstName:   &firstName,
		LastName:    &lastName,
		Password:    hashedPassword,
		RoleID:      role.ID,
		PhoneNumber: req.Phone,
		Status:      "active",
		IsActive:    true,
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// If role is waiter and tables assigned, create table assignments
	if req.Role == "waiter" && len(req.AssignedTables) > 0 {
		for _, tableID := range req.AssignedTables {
			assignment := &models.WaiterTableAssignment{
				WaiterID:   createdUser.ID,
				TableID:    tableID,
				AssignedBy: &adminID,
				IsActive:   true,
			}
			_, err := s.waiterTableAssignmentRepo.Create(ctx, assignment)
			if err != nil {
				s.logger.Error("Failed to create table assignment",
					zap.String("waiter_id", createdUser.ID),
					zap.Int("table_id", tableID),
					zap.Error(err))
			}
		}
	}

	// Build response
	return s.buildStaffResponse(ctx, createdUser, false)
}

// GetStaffByID - TASK-023: Get staff details
func (s *Service) GetStaffByID(ctx context.Context, staffID string) (*models.StaffResponse, error) {
	user, err := s.userRepo.GetByID(ctx, staffID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return s.buildStaffResponse(ctx, user, true)
}

// UpdateStaff - TASK-024: Update staff account
func (s *Service) UpdateStaff(ctx context.Context, staffID string, req models.UpdateStaffRequest, adminID string) (*models.StaffResponse, error) {
	_, err := s.userRepo.GetByID(ctx, staffID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	columns := make(map[string]interface{})

	// Update fields if provided
	if req.Name != nil {
		firstName, lastName := splitName(*req.Name)
		columns["first_name"] = firstName
		columns["last_name"] = lastName
	}
	if req.Email != nil {
		columns["email"] = *req.Email
	}
	if req.Phone != nil {
		columns["phone_number"] = *req.Phone
	}
	if req.Status != nil {
		columns["status"] = *req.Status
		columns["is_active"] = (*req.Status == "active")
	}

	// Update role if provided
	if req.Role != nil {
		role, err := s.roleRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
			tx.Where("name = ?", *req.Role)
		})
		if err != nil || role == nil {
			return nil, fmt.Errorf("invalid role")
		}
		columns["role"] = role.ID
	}

	// Save user
	if len(columns) > 0 {
		_, err = s.userRepo.UpdateColumns(ctx, staffID, columns)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	// Handle table assignments for waiters
	if req.AssignedTables != nil {
		// Deactivate existing assignments
		err = s.waiterTableAssignmentRepo.UpdatesColumnsByConditions(ctx,
			map[string]interface{}{"is_active": false},
			func(tx *gorm.DB) {
				tx.Where("waiter_id = ?", staffID)
			},
		)
		if err != nil {
			s.logger.Error("Failed to deactivate old table assignments", zap.Error(err))
		}

		// Create new assignments
		for _, tableID := range req.AssignedTables {
			assignment := &models.WaiterTableAssignment{
				WaiterID:   staffID,
				TableID:    tableID,
				AssignedBy: &adminID,
				IsActive:   true,
			}
			_, err := s.waiterTableAssignmentRepo.Create(ctx, assignment)
			if err != nil {
				s.logger.Error("Failed to create table assignment",
					zap.String("waiter_id", staffID),
					zap.Int("table_id", tableID),
					zap.Error(err))
			}
		}
	}

	// Refetch user to get updated data
	user, err := s.userRepo.GetByID(ctx, staffID)
	if err != nil {
		return nil, err
	}

	return s.buildStaffResponse(ctx, user, true)
}

// DeleteStaff - TASK-025: Delete staff account
func (s *Service) DeleteStaff(ctx context.Context, staffID string) error {
	_, err := s.userRepo.GetByID(ctx, staffID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Soft delete by setting status to inactive
	columns := map[string]interface{}{
		"status":    "inactive",
		"is_active": false,
	}

	_, err = s.userRepo.UpdateColumns(ctx, staffID, columns)
	if err != nil {
		return fmt.Errorf("failed to delete staff: %w", err)
	}

	return nil
}

// SendStaffInvite - TASK-026: Send staff invitation email
func (s *Service) SendStaffInvite(ctx context.Context, staffID string, req models.SendInviteRequest, adminID string) error {
	// Get role
	role, err := s.roleRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("name = ?", req.Role)
	})
	if err != nil || role == nil {
		return fmt.Errorf("invalid role")
	}

	// Generate unique token
	token := fmt.Sprintf("invite_%s_%d", staffID, time.Now().Unix())

	// Create invitation
	invitation := &models.StaffInvitation{
		Email:     req.Email,
		RoleID:    role.ID,
		InvitedBy: adminID,
		Token:     token,
		Status:    "pending",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	_, err = s.staffInvitationRepo.Create(ctx, invitation)
	if err != nil {
		return fmt.Errorf("failed to create invitation: %w", err)
	}

	// TODO: Send email using AWS SES or email service
	s.logger.Info("Staff invitation sent", zap.String("email", req.Email), zap.String("token", token))

	return nil
}

// AssignTablesToWaiter - TASK-027: Assign tables to waiter
func (s *Service) AssignTablesToWaiter(ctx context.Context, waiterID string, req models.AssignTablesRequest, adminID string) error {
	// Verify waiter exists and has waiter role
	user, err := s.userRepo.GetByID(ctx, waiterID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Get role
	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil || role.Name != "waiter" {
		return fmt.Errorf("user is not a waiter")
	}

	// Deactivate existing assignments
	err = s.waiterTableAssignmentRepo.UpdatesColumnsByConditions(ctx,
		map[string]interface{}{"is_active": false},
		func(tx *gorm.DB) {
			tx.Where("waiter_id = ?", waiterID)
		},
	)
	if err != nil {
		s.logger.Error("Failed to deactivate old table assignments", zap.Error(err))
	}

	// Create new assignments
	for _, tableID := range req.TableIDs {
		assignment := &models.WaiterTableAssignment{
			WaiterID:   waiterID,
			TableID:    tableID,
			AssignedBy: &adminID,
			IsActive:   true,
		}
		_, err := s.waiterTableAssignmentRepo.Create(ctx, assignment)
		if err != nil {
			s.logger.Error("Failed to create table assignment",
				zap.String("waiter_id", waiterID),
				zap.Int("table_id", tableID),
				zap.Error(err))
		}
	}

	return nil
}

// GetStaffProfile - TASK-028: Get logged-in staff profile (role-specific)
func (s *Service) GetStaffProfile(ctx context.Context, staffID string) (*models.StaffResponse, error) {
	user, err := s.userRepo.GetByID(ctx, staffID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return s.buildStaffResponse(ctx, user, true)
}

// UpdateStaffProfile - TASK-029: Update own profile
func (s *Service) UpdateStaffProfile(ctx context.Context, staffID string, req models.UpdateStaffProfileRequest) (*models.StaffResponse, error) {
	_, err := s.userRepo.GetByID(ctx, staffID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	columns := make(map[string]interface{})

	// Update fields
	if req.Name != nil {
		firstName, lastName := splitName(*req.Name)
		columns["first_name"] = firstName
		columns["last_name"] = lastName
	}
	if req.Email != nil {
		columns["email"] = *req.Email
	}
	if req.Phone != nil {
		columns["phone_number"] = *req.Phone
	}
	if req.Avatar != nil {
		columns["avatar_url"] = *req.Avatar
	}

	if len(columns) > 0 {
		_, err = s.userRepo.UpdateColumns(ctx, staffID, columns)
		if err != nil {
			return nil, fmt.Errorf("failed to update profile: %w", err)
		}
	}

	// Refetch user
	user, err := s.userRepo.GetByID(ctx, staffID)
	if err != nil {
		return nil, err
	}

	return s.buildStaffResponse(ctx, user, true)
}

// ChangeStaffPassword - TASK-030: Change password for staff
func (s *Service) ChangeStaffPassword(ctx context.Context, staffID string, req models.ChangeStaffPasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, staffID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify old password
	match, err := argon2id.ComparePasswordAndHash(req.OldPassword, user.Password)
	if err != nil {
		return fmt.Errorf("failed to verify password: %w", err)
	}
	if !match {
		return fmt.Errorf("incorrect old password")
	}

	// Validate new password
	if err := validatePasswordStrength(req.NewPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := argon2id.CreateHash(req.NewPassword, argon2id.DefaultParams)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	columns := map[string]interface{}{
		"password": hashedPassword,
	}

	_, err = s.userRepo.UpdateColumns(ctx, staffID, columns)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// Helper function to build staff response
func (s *Service) buildStaffResponse(ctx context.Context, user *models.User, includeDetails bool) (*models.StaffResponse, error) {
	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	name := getFullName(user)

	// Get role name
	roleName := ""
	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err == nil && role != nil {
		roleName = role.Name
	}

	resp := &models.StaffResponse{
		Name:      name,
		Email:     user.Email,
		Phone:     &user.PhoneNumber,
		Role:      roleName,
		Status:    user.Status,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.DateCreated,
		LastLogin: user.LastLoginAt,
	}

	// Get assigned tables for waiters
	if roleName == "waiter" {
		assignments, err := s.waiterTableAssignmentRepo.ListByConditions(ctx, func(tx *gorm.DB) {
			tx.Preload("Table").Where("waiter_id = ? AND is_active = ?", user.ID, true)
		})
		if err == nil && len(assignments) > 0 {
			tableIDs := make([]int, 0)
			tableNames := make([]string, 0)
			tableDetails := make([]models.TableDetailResponse, 0)

			for _, assignment := range assignments {
				tableIDs = append(tableIDs, assignment.TableID)
				if assignment.Table != nil {
					tableNames = append(tableNames, assignment.Table.TableNumber)
					if includeDetails {
						tableDetails = append(tableDetails, models.TableDetailResponse{
							ID:       assignment.Table.ID,
							Name:     assignment.Table.TableNumber,
							Capacity: assignment.Table.Capacity,
							Status:   assignment.Table.Status,
						})
					}
				}
			}

			resp.AssignedTables = tableIDs
			resp.AssignedTablesNames = tableNames
			if includeDetails {
				resp.AssignedTablesDetails = tableDetails
			}
		}
	}

	// Get statistics if detailed view
	if includeDetails {
		stats, err := s.getStaffStatistics(ctx, user.ID, roleName)
		if err == nil {
			resp.Statistics = stats
		}
	}

	// Get permissions for admin
	if includeDetails && roleName == "admin" {
		resp.Permissions = []string{
			"manage_staff",
			"manage_menu",
			"manage_tables",
			"manage_orders",
			"view_reports",
		}
	}

	return resp, nil
}

// Helper function to get staff statistics
func (s *Service) getStaffStatistics(ctx context.Context, staffID string, role string) (*models.StaffStatistics, error) {
	stats := &models.StaffStatistics{}

	// Get staff profile
	profile, err := s.staffProfileRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("user_id = ?", staffID)
	})
	if err == nil && profile != nil {
		if role == "waiter" {
			stats.TotalOrdersServed = &profile.TotalOrdersServed
			stats.AverageRating = profile.AverageRating
		} else if role == "kitchen_staff" {
			stats.TotalOrdersPrepared = &profile.TotalOrdersPrepared
			stats.AveragePrepTime = profile.AverageOrderTime
		}
	}

	// Get today's orders count
	today := time.Now().Truncate(24 * time.Hour)
	var ordersToday int64
	if role == "waiter" {
		ordersToday, _ = s.orderRepo.Count(ctx, models.QueryParams{}, func(tx *gorm.DB) {
			tx.Where("waiter_id = ? AND created_at >= ?", staffID, today)
		})
	} else if role == "kitchen_staff" {
		ordersToday, _ = s.orderRepo.Count(ctx, models.QueryParams{}, func(tx *gorm.DB) {
			tx.Where("kitchen_staff_id = ? AND created_at >= ?", staffID, today)
		})
	}
	if ordersToday > 0 {
		ordersTodayInt := int(ordersToday)
		stats.OrdersToday = &ordersTodayInt
	}

	// Get active orders count
	var activeOrders int64
	if role == "waiter" {
		activeOrders, _ = s.orderRepo.Count(ctx, models.QueryParams{}, func(tx *gorm.DB) {
			tx.Where("waiter_id = ? AND status IN ?", staffID, []string{"pending", "confirmed", "preparing"})
		})
	} else if role == "kitchen_staff" {
		activeOrders, _ = s.orderRepo.Count(ctx, models.QueryParams{}, func(tx *gorm.DB) {
			tx.Where("kitchen_staff_id = ? AND status IN ?", staffID, []string{"confirmed", "preparing"})
		})
		pendingOrders, _ := s.orderRepo.Count(ctx, models.QueryParams{}, func(tx *gorm.DB) {
			tx.Where("kitchen_staff_id = ? AND status = ?", staffID, "confirmed")
		})
		if pendingOrders > 0 {
			pendingOrdersInt := int(pendingOrders)
			stats.PendingOrders = &pendingOrdersInt
		}
	}
	if activeOrders > 0 {
		activeOrdersInt := int(activeOrders)
		stats.ActiveOrders = &activeOrdersInt
	}

	return stats, nil
}

// Helper functions
func splitName(fullName string) (string, string) {
	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

func getFullName(user *models.User) string {
	if user.FirstName == nil && user.LastName == nil {
		return ""
	}
	firstName := ""
	if user.FirstName != nil {
		firstName = *user.FirstName
	}
	lastName := ""
	if user.LastName != nil {
		lastName = *user.LastName
	}
	return strings.TrimSpace(firstName + " " + lastName)
}

func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return fmt.Errorf("password must contain uppercase, lowercase, number, and special character")
	}

	return nil
}
