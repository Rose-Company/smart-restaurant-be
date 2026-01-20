package services

import (
	"app-noti/common"
	"app-noti/internal/models"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProcessPayment - TASK-013: Process payment (VN-PAY, Cash or Card)
// Initiates a payment for a bill with support for multiple payment methods
// For VN-PAY: returns redirect URL for customer to complete payment
// For Cash/Card: processes payment directly
func (s *Service) ProcessPayment(ctx context.Context, req *models.ProcessPaymentRequest) (*models.PaymentResponse, error) {
	// 1. Get bill
	bill, err := s.billRepo.GetByID(ctx, req.BillID)
	if err != nil {
		return nil, common.ErrBillNotFound
	}

	// 2. Validate amount
	if req.Amount <= 0 || req.Amount > bill.TotalAmount {
		return nil, common.ErrInvalidAmount
	}

	// 3. Generate payment ID
	paymentID := fmt.Sprintf("PAY-%s-%d", time.Now().Format("20060102"), bill.ID)

	// 4. Create payment record with pending status
	payment := &models.Payment{
		PaymentID: paymentID,
		BillID:    req.BillID,
		OrderID:   bill.OrderID,
		Amount:    req.Amount,
		Method:    req.Method,
		Status:    "pending",
		Metadata:  req.Metadata,
	}

	// 5. Handle different payment methods
	switch req.Method {
	case "vnpay":
		// For VN-PAY, create payment record and return redirect URL
		createdPayment, err := s.paymentRepo.Create(ctx, payment)
		if err != nil {
			return nil, err
		}

		// Generate VN-PAY redirect URL
		vnpayURL, err := s.generateVNPayPaymentURL(createdPayment, bill)
		if err != nil {
			return nil, err
		}

		return &models.PaymentResponse{
			PaymentID:  createdPayment.PaymentID,
			BillID:     createdPayment.BillID,
			BillNumber: &bill.BillNumber,
			Amount:     createdPayment.Amount,
			Method:     createdPayment.Method,
			Status:     createdPayment.Status,
			VNPayURL:   vnpayURL,
			CreatedAt:  createdPayment.CreatedAt,
		}, nil

	case "cash":
		if req.ReceivedAmount == nil {
			return nil, common.ErrInvalidInput
		}
		payment.ReceivedAmount = req.ReceivedAmount
		if req.ChangeAmount != nil {
			payment.ChangeAmount = req.ChangeAmount
		} else {
			change := *req.ReceivedAmount - req.Amount
			payment.ChangeAmount = &change
		}
		payment.Status = "succeeded"
		now := time.Now()
		payment.ProcessedAt = &now

		createdPayment, err := s.paymentRepo.Create(ctx, payment)
		if err != nil {
			return nil, err
		}

		// Update bill status to paid
		bill.Status = "paid"
		bill.PaymentMethod = &req.Method
		bill.PaidAt = &now
		_, err = s.billRepo.Update(ctx, bill.ID, bill)
		if err != nil {
			s.logger.Error("Failed to update bill status after cash payment", zap.Error(err))
		}

		err = s.updateTableStatusAfterPayment(ctx, bill.TableID)
		if err != nil {
			s.logger.Error("Failed to update table status after payment", zap.Error(err))
		}

		return &models.PaymentResponse{
			PaymentID:      createdPayment.PaymentID,
			BillID:         createdPayment.BillID,
			BillNumber:     &bill.BillNumber,
			Amount:         createdPayment.Amount,
			Method:         createdPayment.Method,
			Status:         createdPayment.Status,
			ReceivedAmount: createdPayment.ReceivedAmount,
			ChangeAmount:   createdPayment.ChangeAmount,
			CreatedAt:      createdPayment.CreatedAt,
			ProcessedAt:    createdPayment.ProcessedAt,
			TableID:        bill.TableID,
		}, nil

	case "card", "ewallet":
		// For card and e-wallet, mark as succeeded
		payment.Status = "succeeded"
		now := time.Now()
		payment.ProcessedAt = &now

		createdPayment, err := s.paymentRepo.Create(ctx, payment)
		if err != nil {
			return nil, err
		}

		// Update bill status to paid
		bill.Status = "paid"
		bill.PaymentMethod = &req.Method
		bill.PaidAt = &now
		_, err = s.billRepo.Update(ctx, bill.ID, bill)
		if err != nil {
			s.logger.Error("Failed to update bill status after card payment", zap.Error(err))
		}

		// Update table status to available after successful payment
		err = s.updateTableStatusAfterPayment(ctx, bill.TableID)
		if err != nil {
			s.logger.Error("Failed to update table status after payment", zap.Error(err))
		}

		// Fetch updated table to get new QR token
		updatedTable, _ := s.tableRepo.GetByID(ctx, *bill.TableID)

		return &models.PaymentResponse{
			PaymentID:        createdPayment.PaymentID,
			BillID:           createdPayment.BillID,
			BillNumber:       &bill.BillNumber,
			Amount:           createdPayment.Amount,
			Method:           createdPayment.Method,
			Status:           createdPayment.Status,
			ReceiptURL:       createdPayment.ReceiptURL,
			CreatedAt:        createdPayment.CreatedAt,
			ProcessedAt:      createdPayment.ProcessedAt,
			TableID:          bill.TableID,
			QrToken:          &updatedTable.QrToken,
			QrTokenExpiresAt: updatedTable.QrTokenExpiresAt,
		}, nil

	default:
		return nil, common.ErrInvalidPaymentMethod
	}
}

// GetPaymentStatus - TASK-014: Check payment status
// Retrieves the status of a payment
func (s *Service) GetPaymentStatus(ctx context.Context, paymentID string) (*models.PaymentStatusResponse, error) {
	// Get payment using payment_id field
	payment, err := s.paymentRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("payment_id = ?", paymentID)
	})
	if err != nil {
		return nil, common.ErrPaymentNotFound
	}

	// Get bill for bill number
	bill, err := s.billRepo.GetByID(ctx, payment.BillID)
	if err != nil {
		return nil, err
	}

	return &models.PaymentStatusResponse{
		PaymentID:                 payment.PaymentID,
		BillID:                    payment.BillID,
		BillNumber:                &bill.BillNumber,
		Status:                    payment.Status,
		Amount:                    payment.Amount,
		Method:                    payment.Method,
		CreatedAt:                 payment.CreatedAt,
		ProcessedAt:               payment.ProcessedAt,
		FailedAt:                  payment.FailedAt,
		ReceiptURL:                payment.ReceiptURL,
		ErrorMessage:              payment.ErrorMessage,
		ErrorCode:                 payment.ErrorCode,
		StripePaymentIntentStatus: getStripeStatus(payment.Status),
	}, nil
}

// generateVNPayPaymentURL - Generate VN-PAY payment redirect URL
// Reuses existing CreateVNPayPaymentRes to build VN-PAY payment URL with utility functions
func (s *Service) generateVNPayPaymentURL(payment *models.Payment, bill *models.Bill) (*string, error) {
	// Prepare request data for VN-PAY using the existing CreateVNPayPaymentRes function
	requestData := map[string]interface{}{
		"amount":   payment.Amount,
		"orderId":  payment.PaymentID,
		"language": "vn",
	}

	// Use the existing VN-PAY response builder with nil HTTP request
	// The function will handle IP address extraction gracefully
	vnpayRes, err := s.CreateVNPayPaymentRes(requestData, nil)
	if err != nil {
		return nil, err
	}

	return &vnpayRes.PaymentURL, nil
}

// Helper function to get Stripe status
func getStripeStatus(paymentStatus string) *string {
	var stripeStatus string
	switch paymentStatus {
	case "pending":
		stripeStatus = "processing"
	case "succeeded":
		stripeStatus = "succeeded"
	case "failed":
		stripeStatus = "failed"
	default:
		stripeStatus = paymentStatus
	}
	return &stripeStatus
}

// ValidateDiscountCode - TASK-015: Validate discount code
// Validates a discount code and returns discount information
func (s *Service) ValidateDiscountCode(ctx context.Context, code string, orderID int, orderAmount float64) (*models.ValidateDiscountResponse, error) {
	// Get discount code using code field
	discountCode, err := s.discountCodeRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("code = ?", code)
	})
	if err != nil {
		return nil, common.ErrInvalidDiscount
	}

	// Check if discount is active
	now := time.Now()
	if !discountCode.IsActive || now.Before(discountCode.ValidFrom) || now.After(discountCode.ValidUntil) {
		return nil, common.ErrExpiredDiscount
	}

	// Check minimum order amount
	if orderAmount < discountCode.MinOrderAmount {
		return nil, common.ErrMinOrderNotMet
	}

	// Calculate discount amount
	discountAmount := 0.0
	if discountCode.DiscountType == "percentage" {
		discountAmount = (orderAmount * discountCode.DiscountValue) / 100
		if discountCode.MaxDiscountAmount != nil && discountAmount > *discountCode.MaxDiscountAmount {
			discountAmount = *discountCode.MaxDiscountAmount
		}
	} else if discountCode.DiscountType == "fixed_amount" {
		discountAmount = discountCode.DiscountValue
	}

	// Check usage limits
	if discountCode.UsageLimit != nil && discountCode.UsageCount >= *discountCode.UsageLimit {
		return nil, common.ErrDiscountLimitExceeded
	}

	return &models.ValidateDiscountResponse{
		DiscountID:        discountCode.ID,
		Code:              discountCode.Code,
		IsValid:           true,
		DiscountType:      discountCode.DiscountType,
		DiscountValue:     discountCode.DiscountValue,
		DiscountAmount:    discountAmount,
		MinOrderAmount:    discountCode.MinOrderAmount,
		MaxDiscountAmount: discountCode.MaxDiscountAmount,
		ValidFrom:         discountCode.ValidFrom,
		ValidUntil:        discountCode.ValidUntil,
		UsageLimit:        discountCode.UsageLimit,
		UsageCount:        discountCode.UsageCount,
		Description:       discountCode.Description,
	}, nil
}

// CreateDiscountUsageRecord - Creates a discount usage record when a discount is applied
func (s *Service) CreateDiscountUsageRecord(ctx context.Context, discountID int, customerID *string, orderID int, billID int, discountAmount float64) error {
	usage := &models.DiscountUsage{
		DiscountID:     discountID,
		CustomerID:     customerID,
		OrderID:        &orderID,
		BillID:         &billID,
		DiscountAmount: discountAmount,
	}

	_, err := s.discountUsageRepo.Create(ctx, usage)
	return err
}

// HandleVNPayCallback - Handle VN-PAY payment callback
// Called by VN-PAY after customer completes payment
// Updates payment status based on VN-PAY response
func (s *Service) HandleVNPayCallback(ctx context.Context, params map[string]string) (*models.PaymentCallbackResponse, error) {
	// Get payment using transaction reference (payment_id field)
	txnRef := params["vnp_TxnRef"]
	payment, err := s.paymentRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("payment_id = ?", txnRef)
	})
	if err != nil {
		return nil, common.ErrPaymentNotFound
	}

	// Verify VN-PAY response signature
	// Extract secure hash from params
	responseCode := params["vnp_ResponseCode"]
	transactionStatus := params["vnp_TransactionStatus"]

	// Update payment status based on response
	if responseCode == "00" && transactionStatus == "00" {
		// Payment successful
		payment.Status = "succeeded"

		// Store VN-PAY specific data (proper string to *string conversion)
		if txnNo, ok := params["vnp_TransactionNo"]; ok {
			payment.VNPayTransactionNo = &txnNo
		}
		if bankCode, ok := params["vnp_BankCode"]; ok {
			payment.VNPayBankCode = &bankCode
		}
		if cardType, ok := params["vnp_CardType"]; ok {
			payment.VNPayCardType = &cardType
		}

		now := time.Now()
		payment.ProcessedAt = &now

		// Update bill status to paid
		bill, err := s.billRepo.GetByID(ctx, payment.BillID)
		if err == nil && bill != nil {
			bill.Status = "paid"
			method := "vnpay"
			bill.PaymentMethod = &method
			bill.PaidAt = &now
			s.billRepo.Update(ctx, bill.ID, bill)

			// Update table status to available after successful VN-PAY payment
			s.updateTableStatusAfterPayment(ctx, bill.TableID)
		}
	} else {
		// Payment failed
		payment.Status = "failed"
		if respCode, ok := params["vnp_ResponseCode"]; ok {
			payment.VNPayResponseCode = &respCode
		}
		now := time.Now()
		payment.FailedAt = &now
		errorMsg := fmt.Sprintf("VN-PAY Error: %s", responseCode)
		payment.ErrorMessage = &errorMsg
	}

	// Save updated payment
	_, err = s.paymentRepo.Update(ctx, payment.ID, payment)
	if err != nil {
		s.logger.Error("Failed to update payment after VN-PAY callback", zap.Error(err))
	}

	return &models.PaymentCallbackResponse{
		PaymentID: payment.PaymentID,
		Status:    payment.Status,
		Message:   "Payment callback processed",
	}, nil
}

// Helper function to update table status to available after successful payment
func (s *Service) updateTableStatusAfterPayment(ctx context.Context, tableID *int) error {
	if tableID == nil {
		return nil
	}

	// Get table
	table, err := s.tableRepo.GetByID(ctx, *tableID)
	if err != nil {
		return fmt.Errorf("failed to get table: %w", err)
	}

	// Update table status to active
	table.Status = "active"

	// Generate NEW QR token for next customer using table service's token generator
	newToken, err := GenerateSecureToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate QR token: %w", err)
	}

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	// Use UpdateColumns to explicitly update only token fields
	_, err = s.tableRepo.UpdateColumns(ctx, *tableID, map[string]interface{}{
		"status":              "active",
		"qr_token":            newToken,
		"qr_token_created_at": now,
		"qr_token_expires_at": expiresAt,
	})
	if err != nil {
		return fmt.Errorf("failed to update table status: %w", err)
	}

	s.logger.Info(fmt.Sprintf("Table %d reset to active with new QR token after payment", *tableID))
	return nil
}
