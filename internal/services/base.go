package services

import (
	"app-noti/common"
	"app-noti/internal/repositories"
	l "app-noti/pkg/logger"
	redisClient "app-noti/pkg/redis"
	"app-noti/server"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	logger                    *zap.Logger
	redisClient               redisClient.ClientI
	tableRepo                 *repositories.TableRepo
	restaurantRepo            *repositories.RestaurantRepo
	menuCategoryRepo          *repositories.MenuCategoryRepo
	menuItemRepo              *repositories.MenuItemRepo
	menuItemPhotoRepo         *repositories.MenuItemPhotoRepo
	modifierGroupRepo         *repositories.ModifierGroupRepo
	modifierOptionRepo        *repositories.ModifierOptionRepo
	menuItemModifierGroupRepo *repositories.MenuItemModifierGroupRepo
	reviewItemRepo            repositories.ReviewItemRepository
	userRepo                  *repositories.UserRepo
	otpRepo                   *repositories.OTPRepo
	otpAttemptRepo            *repositories.OTPAttemptRepo
	roleRepo                  *repositories.RoleRepo
	orderRepo                 repositories.OrderRepository
	orderItemRepo             repositories.OrderItemRepository
	orderModifierRepo         repositories.OrderModifierRepository
	orderTimelineRepo         repositories.OrderTimelineRepository
	kitchenAlertRepo          repositories.KitchenAlertRepository
	billRepo                  repositories.BillRepository
	billItemRepo              repositories.BillItemRepository
	paymentRepo               repositories.PaymentRepository
	discountCodeRepo          repositories.DiscountCodeRepository
	discountUsageRepo         repositories.DiscountUsageRepository
	reviewRepo                repositories.ReviewRepository
	preferencesRepo           repositories.UserPreferencesRepository
	staffProfileRepo          repositories.StaffProfileRepository
	waiterTableAssignmentRepo repositories.WaiterTableAssignmentRepository
	staffInvitationRepo       repositories.StaffInvitationRepository
}

func NewService(sc server.ServerContext) *Service {
	db := sc.GetService(common.PREFIX_MAIN_POSTGRES).(*gorm.DB)

	return &Service{
		logger:                    l.New(),
		redisClient:               redisClient.NewRedisClient(),
		reviewItemRepo:            repositories.NewReviewItemRepository(db),
		tableRepo:                 repositories.NewTableRepository(db),
		restaurantRepo:            repositories.NewRestaurantRepository(db),
		menuCategoryRepo:          repositories.NewMenuCategoryRepository(db),
		menuItemRepo:              repositories.NewMenuItemRepository(db),
		menuItemPhotoRepo:         repositories.NewMenuItemPhotoRepository(db),
		modifierGroupRepo:         repositories.NewModifierGroupRepository(db),
		modifierOptionRepo:        repositories.NewModifierOptionRepository(db),
		menuItemModifierGroupRepo: repositories.NewMenuItemModifierGroupRepository(db),
		userRepo:                  repositories.NewUserRepository(db),
		otpRepo:                   repositories.NewOTPRepository(db),
		otpAttemptRepo:            repositories.NewOTPAttemptRepository(db),
		roleRepo:                  repositories.NewRoleRepository(db),
		orderRepo:                 repositories.NewOrderRepository(db),
		orderItemRepo:             repositories.NewOrderItemRepository(db),
		orderModifierRepo:         repositories.NewOrderModifierRepository(db),
		orderTimelineRepo:         repositories.NewOrderTimelineRepository(db),
		kitchenAlertRepo:          repositories.NewKitchenAlertRepository(db),
		billRepo:                  repositories.NewBillRepository(db),
		billItemRepo:              repositories.NewBillItemRepository(db),
		paymentRepo:               repositories.NewPaymentRepository(db),
		discountCodeRepo:          repositories.NewDiscountCodeRepository(db),
		discountUsageRepo:         repositories.NewDiscountUsageRepository(db),
		reviewRepo:                repositories.NewReviewRepository(db),
		preferencesRepo:           repositories.NewUserPreferencesRepository(db),
		staffProfileRepo:          repositories.NewStaffProfileRepository(db),
		waiterTableAssignmentRepo: repositories.NewWaiterTableAssignmentRepository(db),
		staffInvitationRepo:       repositories.NewStaffInvitationRepository(db),
	}
}
