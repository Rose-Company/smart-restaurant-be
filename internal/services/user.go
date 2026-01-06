package services

import (
	"app-noti/common"
	"app-noti/internal/models"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var JWTSecret = []byte("your_secret_key")

func (s *Service) SignupUser(ctx context.Context, req models.SignupRequest) error {
	// Step 1: Verify email via OTP
	storedOTP, err := s.otpRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("target = ? AND type = ?", req.Email, common.TypeVerifyEmail)
		tx.Order("created_at desc")
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrOTPNotFound
		}
		return err
	}

	// Check if OTP is verified and verify token matches
	if !storedOTP.IsVerified || storedOTP.VerifyToken != req.VerifyToken {
		return common.ErrInvalidVerifyToken
	}

	// Step 2: Check if user already exists
	_, err = s.userRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("email = ?", req.Email)
	})

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	} else {
		return common.ErrEmailAlreadyExists
	}

	// Step 3: Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Step 4: Get role from database
	role, err := s.roleRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("name = ?", req.Role)
	})
	if err != nil {
		return common.ErrInvalidInput
	}

	// Step 5: Create new user
	newUser := models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		RoleID:    role.ID,
		FirstName: &req.FirstName,
		LastName:  &req.LastName,
		IsActive:  true,
		Provider:  common.USER_PROVIDER_LOCAL,
		Status:    common.USER_STATUS_ACTIVE,
	}

	_, err = s.userRepo.Create(ctx, &newUser)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) LoginUser(ctx context.Context, req models.LoginRequest) (*string, error) {
	var user *models.User
	var err error

	// Google OAuth login
	if req.IdToken != nil {
		googleUser, err := verifyGoogleOAuthToken(*req.IdToken)
		if err != nil {
			return nil, err
		}

		user, err = s.userRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
			tx.Where("email = ? AND provider = ?", googleUser.Email, common.USER_PROVIDER_GOOGLE)
		})

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				roleStr := common.ROLE_END_USER
				if req.Role != nil && *req.Role != "" {
					roleStr = *req.Role
				}

				// Get role from database
				role, err := s.roleRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
					tx.Where("name = ?", roleStr)
				})
				if err != nil {
					return nil, common.ErrInvalidInput
				}

				newUser := models.User{
					FirstName: &googleUser.GivenName,
					LastName:  &googleUser.FamilyName,
					Email:     googleUser.Email,
					RoleID:    role.ID,
					Provider:  common.USER_PROVIDER_GOOGLE,
					IsActive:  true,
					Status:    common.USER_STATUS_ACTIVE,
				}

				user, err = s.userRepo.Create(ctx, &newUser)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	} else {
		// Email/password login
		if req.Email == nil || req.Password == nil {
			return nil, common.ErrInvalidInput
		}

		user, err = s.userRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
			tx.Where("email = ?", *req.Email)
		})

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, common.ErrInvalidEmailOrPassword
			}
			return nil, err
		}

		// Verify password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*req.Password)); err != nil {
			return nil, common.ErrInvalidEmailOrPassword
		}
	}

	// Check if user is active
	if !user.IsActive {
		return nil, common.ErrUserInactive
	}

	return generateJWTToken(user)
}

func generateJWTToken(user *models.User) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.RoleID,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func verifyGoogleOAuthToken(idToken string) (*models.GoogleUser, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + idToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, common.ErrInvalidGoogleAuthenToken
	}

	var googleUser models.GoogleUser
	if err := json.Unmarshal(bodyBytes, &googleUser); err != nil {
		return nil, err
	}

	return &googleUser, nil
}

func (s *Service) GenerateOTP(ctx context.Context, email string, otpType string) (string, error) {
	otp := common.GenerateRandomOTP()
	expiry := time.Now().UTC().Add(5 * time.Minute) // 5 minutes expiry

	// Invalidate existing OTP if any
	existingOTP, err := s.otpRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("target = ? AND type = ?", email, otpType)
		tx.Order("created_at desc")
	})

	if err == nil {
		existingOTP.IsVerified = true
		_, err = s.otpRepo.Update(ctx, existingOTP.ID, existingOTP)
		if err != nil {
			return "", common.ErrFailedToInValidateExistingOTP
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	// Create new OTP
	newOTP := models.OTP{
		Target:     email,
		Type:       otpType,
		OTPCode:    otp,
		ExpiredAt:  expiry,
		IsVerified: false,
	}

	_, err = s.otpRepo.Create(ctx, &newOTP)
	if err != nil {
		return "", err
	}

	return otp, nil
}

func (s *Service) ValidateOTP(ctx context.Context, email, otp, otpType string) (string, error) {
	storedOTP, err := s.otpRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("target = ? AND type = ? AND is_verified = ?", email, otpType, false)
		tx.Order("created_at desc")
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", common.ErrOTPNotFound
		}
		return "", err
	}

	currentTime := time.Now().UTC()

	// Create attempt record
	newAttempt := models.OTPAttempt{
		OTPID:     storedOTP.ID,
		Value:     otp,
		IsSuccess: false,
		CreatedAt: currentTime,
	}

	// Check if OTP is expired
	if storedOTP.ExpiredAt.Before(currentTime) {
		_, _ = s.otpAttemptRepo.Create(ctx, &newAttempt)
		return "", common.ErrOTPExpired
	}

	// Check if OTP matches
	if storedOTP.OTPCode != otp {
		_, _ = s.otpAttemptRepo.Create(ctx, &newAttempt)
		return "", common.ErrInvalidOTP
	}

	// Mark OTP as verified and generate verify token
	storedOTP.IsVerified = true
	verifyToken, err := common.GenerateShortUUID()
	if err != nil {
		return "", common.ErrFailedToUpdateOTPStatus
	}
	storedOTP.VerifyToken = verifyToken
	_, err = s.otpRepo.Update(ctx, storedOTP.ID, storedOTP)
	if err != nil {
		return "", common.ErrFailedToUpdateOTPStatus
	}

	// Record successful attempt
	newAttempt.IsSuccess = true
	_, _ = s.otpAttemptRepo.Create(ctx, &newAttempt)

	return verifyToken, nil
}

func (s *Service) ResetPassword(ctx context.Context, req models.ResetPasswordRequest) error {
	// Step 1: Check if user exists
	user, err := s.userRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("email = ?", req.Email)
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrEmailNotFound
		}
		return err
	}

	// Step 2: Verify OTP and token match
	storedOTP, err := s.otpRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("target = ? AND type = ? AND is_verified = ?", req.Email, common.TypeResetPassword, true)
		tx.Order("created_at desc")
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrOTPNotVerified
		}
		return err
	}

	// Check verify token matches
	if storedOTP.VerifyToken != req.VerifyToken {
		return common.ErrInvalidVerifyToken
	}

	// Step 3: Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Step 4: Update user password
	updatedUser := models.User{
		Password: string(hashedPassword),
	}

	return s.userRepo.UpdatesByConditions(ctx, &updatedUser, func(tx *gorm.DB) {
		tx.Where("id = ?", user.ID)
	})
}
