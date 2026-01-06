package handlers

import (
	"app-noti/common"
	"app-noti/config"
	"app-noti/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SignUp(c *gin.Context) {
	var req models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	err := h.service.SignupUser(c.Request.Context(), req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.ResponseOk("User created successfully"))
}

func (h *Handler) RequestSignupOTP(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	otp, err := h.service.GenerateOTP(c.Request.Context(), req.Email, common.TypeVerifyEmail)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	err = common.SendOTPEmail(config.Config.Mail.FromEmail, req.Email, otp)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.BaseResponse(http.StatusOK, "OTP email sent successfully", "", ""))
}

func (h *Handler) LogIn(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	if req.Email == nil && req.IdToken == nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	token, err := h.service.LoginUser(c.Request.Context(), req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.BaseResponseMess(common.SUCCESS_STATUS, "Login successfully", token))
}

func (h *Handler) RequestResetPassword(c *gin.Context) {
	var req models.RequestResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	otp, err := h.service.GenerateOTP(c.Request.Context(), req.Email, common.TypeResetPassword)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	err = common.SendOTPEmail(config.Config.Mail.FromEmail, req.Email, otp)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.BaseResponse(http.StatusOK, "OTP email sent successfully", "", ""))
}

func (h *Handler) ValidateOTP(c *gin.Context) {
	var req models.OTPValidateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	verifyToken, err := h.service.ValidateOTP(c.Request.Context(), req.Email, req.OTP, common.TypeResetPassword)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.BaseResponseMess(common.SUCCESS_STATUS, "OTP validated successfully", verifyToken))
}

func (h *Handler) ValidateSignupOTP(c *gin.Context) {
	var req models.OTPValidateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	verifyToken, err := h.service.ValidateOTP(c.Request.Context(), req.Email, req.OTP, common.TypeVerifyEmail)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.BaseResponseMess(common.SUCCESS_STATUS, "OTP validated successfully", verifyToken))
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	err := h.service.ResetPassword(c.Request.Context(), req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.BaseResponse(http.StatusOK, "Password reset successfully", "", ""))
}
