package services

import (
	"app-noti/common"
	"app-noti/internal/models"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CreateBill - TASK-010: Create bill from order
// Creates a bill from an order with all line items and calculates totals
func (s *Service) CreateBill(ctx context.Context, req *models.CreateBillRequest) (*models.BillResponse, error) {
	// 1. Validate order exists and is in a valid state
	order, err := s.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, common.ErrOrderNotFound
	}

	// 2. Get order items using repository method
	orderItems, err := s.orderItemRepo.ListByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("order_id = ?", req.OrderID)
	})
	if err != nil || len(orderItems) == 0 {
		return nil, common.ErrOrderNotFound
	}

	// 3. Calculate totals
	var subtotal float64
	var modifiersTotal float64
	billItems := make([]*models.BillItem, 0)

	for _, item := range orderItems {
		subtotal += item.Subtotal
		if item.ModifiersTotal > 0 {
			modifiersTotal += item.ModifiersTotal
		}

		billItem := &models.BillItem{
			MenuItemID:     item.MenuItemID,
			MenuItemName:   item.ItemName,
			Quantity:       item.Quantity,
			UnitPrice:      item.UnitPrice,
			ModifiersTotal: item.ModifiersTotal,
			Subtotal:       item.Subtotal,
		}
		billItems = append(billItems, billItem)
	}

	// 4. Apply discount if provided
	var discountAmount float64 = 0
	var discountCode *string
	if req.DiscountCode != nil {
		validatedDiscount, err := s.ValidateDiscountCode(ctx, *req.DiscountCode, req.OrderID, subtotal+modifiersTotal)
		if err == nil && validatedDiscount != nil {
			discountCode = req.DiscountCode
			discountAmount = validatedDiscount.DiscountAmount
		}
	}

	// 5. Calculate tax (assuming 10% tax rate from database default)
	taxRate := 10.0
	subtotalAfterDiscount := subtotal + modifiersTotal - discountAmount
	taxAmount := (subtotalAfterDiscount * taxRate) / 100

	// 6. Calculate final total
	totalAmount := subtotalAfterDiscount + taxAmount

	// 7. Generate bill number
	billNumber := common.GenerateBillNumber()

	// 8. Create bill
	bill := &models.Bill{
		BillNumber:     billNumber,
		OrderID:        req.OrderID,
		RestaurantID:   order.RestaurantID,
		TableID:        &order.TableID,
		CustomerID:     order.CustomerUserID,
		Subtotal:       subtotal,
		TaxAmount:      taxAmount,
		TaxRate:        taxRate,
		DiscountAmount: discountAmount,
		DiscountCode:   discountCode,
		ServiceCharge:  0,
		TotalAmount:    totalAmount,
		Status:         "pending",
		BillType:       req.Type,
		RequestedBy:    &req.RequestedBy,
		BillItems:      billItems,
	}

	createdBill, err := s.billRepo.Create(ctx, bill)
	if err != nil {
		return nil, err
	}

	// 9. Create bill items
	for _, item := range billItems {
		item.BillID = createdBill.ID
	}
	err = s.billItemRepo.CreatesMultiple(ctx, billItems)
	if err != nil {
		return nil, err
	}

	// 10. Build response
	return s.buildBillResponse(ctx, createdBill)
}

// GetBill - TASK-011: Get bill details
// Retrieves bill details with optional format (json or pdf)
func (s *Service) GetBill(ctx context.Context, billID int, format string) (interface{}, error) {
	// Get bill
	bill, err := s.billRepo.GetByID(ctx, billID)
	if err != nil {
		return nil, common.ErrBillNotFound
	}

	// For JSON format, return full bill details
	if format != "pdf" {
		return s.buildBillResponse(ctx, bill)
	}

	// For PDF format, generate and return PDF URL
	// In a real application, this would generate a PDF and upload it to storage
	pdfURL := fmt.Sprintf("https://storage.example.com/bills/%s.pdf", bill.BillNumber)
	expiresAt := time.Now().Add(24 * time.Hour)

	return &models.BillPDFResponse{
		PdfURL:    pdfURL,
		ExpiresAt: expiresAt,
	}, nil
}

// UpdateBill - TASK-012: Update bill (add discount, mark paid)
// Updates bill with new discount or payment information
func (s *Service) UpdateBill(ctx context.Context, billID int, req *models.UpdateBillRequest) (*models.BillResponse, error) {
	// Get current bill
	bill, err := s.billRepo.GetByID(ctx, billID)
	if err != nil {
		return nil, common.ErrBillNotFound
	}

	// Update discount if provided
	if req.DiscountCode != nil {
		_, err := s.orderRepo.GetByID(ctx, bill.OrderID)
		if err != nil {
			return nil, err
		}

		validatedDiscount, err := s.ValidateDiscountCode(ctx, *req.DiscountCode, bill.OrderID, bill.Subtotal)
		if err == nil && validatedDiscount != nil {
			bill.DiscountCode = req.DiscountCode
			bill.DiscountAmount = validatedDiscount.DiscountAmount

			// Recalculate totals
			subtotalAfterDiscount := bill.Subtotal + bill.ServiceCharge - bill.DiscountAmount
			bill.TaxAmount = (subtotalAfterDiscount * bill.TaxRate) / 100
			bill.TotalAmount = subtotalAfterDiscount + bill.TaxAmount
		}
	} else if req.DiscountAmount != nil {
		bill.DiscountAmount = *req.DiscountAmount

		// Recalculate totals
		subtotalAfterDiscount := bill.Subtotal + bill.ServiceCharge - bill.DiscountAmount
		bill.TaxAmount = (subtotalAfterDiscount * bill.TaxRate) / 100
		bill.TotalAmount = subtotalAfterDiscount + bill.TaxAmount
	}

	// Update status if provided
	if req.Status != nil {
		bill.Status = *req.Status
		if *req.Status == "paid" {
			now := time.Now()
			bill.PaidAt = &now
		}
	}

	// Update payment method if provided
	if req.PaymentMethod != nil {
		bill.PaymentMethod = req.PaymentMethod
	}

	// Save updated bill
	updatedBill, err := s.billRepo.Update(ctx, billID, bill)
	if err != nil {
		return nil, err
	}

	return s.buildBillResponse(ctx, updatedBill)
}

// Helper function to build bill response
func (s *Service) buildBillResponse(ctx context.Context, bill *models.Bill) (*models.BillResponse, error) {
	// Get order for order number and table info
	order, err := s.orderRepo.GetByID(ctx, bill.OrderID)
	if err != nil {
		return nil, err
	}

	// Get table info
	var tableName *string
	if bill.TableID != nil {
		table, err := s.tableRepo.GetByID(ctx, *bill.TableID)
		if err == nil && table != nil {
			tableName = &table.TableNumber
		}
	}

	// Get customer info
	var customerName *string
	var customerPhone *string
	if bill.CustomerID != nil {
		user, err := s.userRepo.GetByID(ctx, *bill.CustomerID)
		if err == nil && user != nil {
			fullName := fmt.Sprintf("%s %s", *user.FirstName, *user.LastName)
			customerName = &fullName
			customerPhone = &user.PhoneNumber
		}
	}

	// Build bill items response
	billItems := make([]*models.BillItemResponse, 0)
	if bill.BillItems != nil {
		for _, item := range bill.BillItems {
			billItems = append(billItems, &models.BillItemResponse{
				ID:             item.ID,
				MenuItemName:   item.MenuItemName,
				Quantity:       item.Quantity,
				UnitPrice:      item.UnitPrice,
				ModifiersTotal: item.ModifiersTotal,
				Subtotal:       item.Subtotal,
			})
		}
	}

	// Build breakdown
	subtotalBeforeDiscount := bill.Subtotal
	subtotalAfterDiscount := bill.Subtotal - bill.DiscountAmount
	breakdown := &models.BillBreakdown{
		ItemsTotal:             bill.Subtotal,
		ModifiersTotal:         0,
		SubtotalBeforeDiscount: subtotalBeforeDiscount,
		Discount:               -bill.DiscountAmount,
		SubtotalAfterDiscount:  subtotalAfterDiscount,
		Tax:                    bill.TaxAmount,
		ServiceCharge:          bill.ServiceCharge,
		GrandTotal:             bill.TotalAmount,
	}

	return &models.BillResponse{
		ID:             bill.ID,
		BillNumber:     bill.BillNumber,
		OrderID:        bill.OrderID,
		OrderNumber:    &order.OrderNumber,
		TableID:        bill.TableID,
		TableName:      tableName,
		CustomerID:     bill.CustomerID,
		CustomerName:   customerName,
		CustomerPhone:  customerPhone,
		Subtotal:       bill.Subtotal,
		TaxAmount:      bill.TaxAmount,
		TaxRate:        &bill.TaxRate,
		DiscountAmount: bill.DiscountAmount,
		DiscountCode:   bill.DiscountCode,
		ServiceCharge:  bill.ServiceCharge,
		TotalAmount:    bill.TotalAmount,
		Status:         bill.Status,
		Type:           bill.BillType,
		PaymentMethod:  bill.PaymentMethod,
		RequestedBy:    bill.RequestedBy,
		CreatedAt:      bill.CreatedAt,
		UpdatedAt:      &bill.UpdatedAt,
		PaidAt:         bill.PaidAt,
		Items:          billItems,
		Breakdown:      breakdown,
	}, nil
}
