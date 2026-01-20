package services

import (
	"app-noti/common"
	"app-noti/config"
	"app-noti/internal/models"
	"app-noti/services/digital_ocean_storage"
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

// CreateBill - TASK-010: Create bill from table
// Creates a bill from a table with all orders and items
func (s *Service) CreateBill(ctx context.Context, req *models.CreateBillRequest) (*models.BillDetailResponse, error) {
	// 1. Validate table exists
	table, err := s.tableRepo.GetByID(ctx, req.TableID)
	if err != nil {
		return nil, common.ErrTableNotFound
	}

	// 2. Get all orders for this table (not completed/cancelled only)
	allOrders, err := s.orderRepo.ListByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("table_id = ?", req.TableID)
	})
	if err != nil || len(allOrders) == 0 {
		return nil, common.ErrOrderNotFound
	}

	// 3. Extract order IDs for batch query
	orderIDs := make([]int, 0, len(allOrders))
	for _, order := range allOrders {
		orderIDs = append(orderIDs, order.ID)
	}

	// 4. Get ALL order items in ONE query
	allOrderItems, err := s.orderItemRepo.ListByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("order_id IN ?", orderIDs)
	})
	if err != nil || len(allOrderItems) == 0 {
		return nil, common.ErrOrderNotFound
	}

	// 5. Calculate totals and group items by order
	var subtotal float64
	var modifiersTotal float64
	itemsByOrderID := make(map[int][]*models.OrderItem)

	// Group items by order_id and calculate totals
	for _, item := range allOrderItems {
		itemsByOrderID[item.OrderID] = append(itemsByOrderID[item.OrderID], item)
		subtotal += item.Subtotal
		if item.ModifiersTotal > 0 {
			modifiersTotal += item.ModifiersTotal
		}
	}

	// 6. Apply discount if provided
	var discountAmount float64 = 0
	var discountCode *string
	if req.DiscountCode != nil {
		validatedDiscount, err := s.ValidateDiscountCode(ctx, *req.DiscountCode, 0, subtotal+modifiersTotal)
		if err == nil && validatedDiscount != nil {
			discountCode = req.DiscountCode
			discountAmount = validatedDiscount.DiscountAmount
		}
	}

	// 7. Calculate tax
	taxRate := 8.0
	subtotalAfterDiscount := subtotal + modifiersTotal - discountAmount
	taxAmount := (subtotalAfterDiscount * taxRate) / 100

	// 8. Calculate final total
	totalAmount := subtotalAfterDiscount + taxAmount

	// 9. Generate bill number
	billNumber := common.GenerateBillNumber()

	// 10. Create bill at table level
	firstOrderID := allOrders[0].ID
	bill := &models.Bill{
		BillNumber:     billNumber,
		OrderID:        firstOrderID,
		RestaurantID:   &table.RestaurantId,
		TableID:        &req.TableID,
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
	}

	createdBill, err := s.billRepo.Create(ctx, bill)
	if err != nil {
		return nil, err
	}

	// Mark table as ready to bill
	_, err = s.tableRepo.UpdateColumns(ctx, req.TableID, map[string]interface{}{
		"is_ready_to_bill": true,
		"updated_at":       time.Now(),
	})
	if err != nil {
		// Log error but don't fail the bill creation
		fmt.Printf("Warning: failed to update table is_ready_to_bill: %v\n", err)
	}

	// 11. Build detailed response with all orders and items
	return s.buildBillDetailResponse(ctx, createdBill, allOrders, itemsByOrderID, table)
}

// GetBill - TASK-011: Get bill details with all orders and items
// Retrieves bill details including all orders and items on the table
func (s *Service) GetBill(ctx context.Context, billID int, format string) (interface{}, error) {
	// Get bill
	bill, err := s.billRepo.GetByID(ctx, billID)
	if err != nil {
		return nil, common.ErrBillNotFound
	}

	// Get table
	table, err := s.tableRepo.GetByID(ctx, *bill.TableID)
	if err != nil {
		return nil, err
	}

	// Get all orders for this table
	allOrders, err := s.orderRepo.ListByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("table_id = ?", *bill.TableID)
	})
	if err != nil {
		return nil, err
	}

	// Extract order IDs
	orderIDs := make([]int, 0, len(allOrders))
	for _, order := range allOrders {
		orderIDs = append(orderIDs, order.ID)
	}

	// Get all order items
	var allOrderItems []*models.OrderItem
	if len(orderIDs) > 0 {
		allOrderItems, err = s.orderItemRepo.ListByConditions(ctx, func(tx *gorm.DB) {
			tx.Where("order_id IN ?", orderIDs)
		})
		if err != nil {
			return nil, err
		}
	}

	// Group items by order_id
	itemsByOrderID := make(map[int][]*models.OrderItem)
	for _, item := range allOrderItems {
		itemsByOrderID[item.OrderID] = append(itemsByOrderID[item.OrderID], item)
	}

	// For PDF format, generate PDF and return URL
	if format == "pdf" {
		// Generate PDF
		pdfBytes, err := s.generateBillPDF(ctx, bill, allOrders, itemsByOrderID, table)
		if err != nil {
			return nil, fmt.Errorf("failed to generate PDF: %w", err)
		}

		// Upload PDF to Digital Ocean
		pdfURL, err := s.uploadBillPDF(bill.BillNumber, pdfBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to upload PDF: %w", err)
		}

		expiresAt := time.Now().Add(24 * time.Hour)
		return &models.BillPDFResponse{
			PdfURL:    pdfURL,
			ExpiresAt: expiresAt,
		}, nil
	}

	// For JSON format, return full bill detail with all orders and items
	return s.buildBillDetailResponse(ctx, bill, allOrders, itemsByOrderID, table)
}

// UpdateBill - TASK-012: Update bill (add discount, mark paid)
// Updates bill with new discount or payment information
func (s *Service) UpdateBill(ctx context.Context, billID int, req *models.UpdateBillRequest) (*models.BillDetailResponse, error) {
	// Get current bill
	bill, err := s.billRepo.GetByID(ctx, billID)
	if err != nil {
		return nil, common.ErrBillNotFound
	}

	// Get table
	table, err := s.tableRepo.GetByID(ctx, *bill.TableID)
	if err != nil {
		return nil, err
	}

	// Update discount if provided
	if req.DiscountCode != nil {
		validatedDiscount, err := s.ValidateDiscountCode(ctx, *req.DiscountCode, 0, bill.Subtotal)
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

	// Get all orders for this table
	allOrders, err := s.orderRepo.ListByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("table_id = ?", *bill.TableID)
	})
	if err != nil {
		return nil, err
	}

	// Extract order IDs
	orderIDs := make([]int, 0, len(allOrders))
	for _, order := range allOrders {
		orderIDs = append(orderIDs, order.ID)
	}

	// Get all order items
	var allOrderItems []*models.OrderItem
	if len(orderIDs) > 0 {
		allOrderItems, err = s.orderItemRepo.ListByConditions(ctx, func(tx *gorm.DB) {
			tx.Where("order_id IN ?", orderIDs)
		})
		if err != nil {
			return nil, err
		}
	}

	// Group items by order_id
	itemsByOrderID := make(map[int][]*models.OrderItem)
	for _, item := range allOrderItems {
		itemsByOrderID[item.OrderID] = append(itemsByOrderID[item.OrderID], item)
	}

	return s.buildBillDetailResponse(ctx, updatedBill, allOrders, itemsByOrderID, table)
}

// Helper function to build detailed bill response with all orders and items
func (s *Service) buildBillDetailResponse(ctx context.Context, bill *models.Bill, allOrders []*models.Order, itemsByOrderID map[int][]*models.OrderItem, table *models.Table) (*models.BillDetailResponse, error) {
	// Build orders with items
	ordersForBill := make([]models.OrderForBill, 0)
	totalItemsCount := 0

	for _, order := range allOrders {
		orderItems := itemsByOrderID[order.ID]

		// Build order items
		itemsForBill := make([]models.OrderItemForBill, 0)
		for _, item := range orderItems {
			itemsForBill = append(itemsForBill, models.OrderItemForBill{
				ID:        item.ID,
				OrderID:   item.OrderID,
				ItemName:  item.ItemName,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
				Status:    item.Status,
			})
			totalItemsCount++
		}

		ordersForBill = append(ordersForBill, models.OrderForBill{
			ID:          order.ID,
			OrderNumber: order.OrderNumber,
			Status:      order.Status,
			TotalAmount: order.Total,
			Items:       itemsForBill,
		})
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

	return &models.BillDetailResponse{
		ID:             bill.ID,
		BillNumber:     bill.BillNumber,
		TableID:        bill.TableID,
		TableNumber:    &table.TableNumber,
		RestaurantID:   bill.RestaurantID,
		CustomerID:     bill.CustomerID,
		CustomerName:   nil, // Can be fetched if needed
		CustomerPhone:  nil,
		Orders:         ordersForBill,
		OrdersCount:    len(allOrders),
		ItemsCount:     totalItemsCount,
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
		Breakdown:      breakdown,
	}, nil
}

// Helper function to generate bill PDF
func (s *Service) generateBillPDF(ctx context.Context, bill *models.Bill, allOrders []*models.Order, itemsByOrderID map[int][]*models.OrderItem, table *models.Table) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "BILL")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, fmt.Sprintf("Bill #: %s", bill.BillNumber))
	pdf.Ln(5)
	pdf.Cell(0, 5, fmt.Sprintf("Date: %s", bill.CreatedAt.Format("2006-01-02 15:04:05")))
	pdf.Ln(5)

	pdf.Ln(5)

	// Table Info
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 5, fmt.Sprintf("Table: %d", *bill.TableID))
	pdf.Cell(0, 5, fmt.Sprintf("Restaurant: %d", *bill.RestaurantID))
	pdf.Ln(5)

	pdf.Ln(5)

	// Items Table Header
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(80, 5, "Item")
	pdf.Cell(20, 5, "Qty")
	pdf.Cell(30, 5, "Price")
	pdf.Cell(30, 5, "Subtotal")
	pdf.Ln(5)

	// Items Table Body
	pdf.SetFont("Arial", "", 9)
	for _, order := range allOrders {
		orderItems := itemsByOrderID[order.ID]
		for _, item := range orderItems {
			pdf.Cell(80, 5, fmt.Sprintf("%s (%s)", item.ItemName, order.OrderNumber))
			pdf.Cell(20, 5, fmt.Sprintf("%d", item.Quantity))
			pdf.Cell(30, 5, fmt.Sprintf("%.2f", item.UnitPrice))
			pdf.Cell(30, 5, fmt.Sprintf("%.2f", item.Subtotal))
			pdf.Ln(5)
		}
	}

	pdf.Ln(5)

	// Summary Section
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(130, 5, "Subtotal:")
	pdf.Cell(30, 5, fmt.Sprintf("%.2f", bill.Subtotal))
	pdf.Ln(5)

	if bill.DiscountAmount > 0 {
		pdf.Cell(130, 5, fmt.Sprintf("Discount (%s):", *bill.DiscountCode))
		pdf.Cell(30, 5, fmt.Sprintf("-%.2f", bill.DiscountAmount))
		pdf.Ln(5)
	}

	pdf.Cell(130, 5, fmt.Sprintf("Tax (%.1f%%):", bill.TaxRate))
	pdf.Cell(30, 5, fmt.Sprintf("%.2f", bill.TaxAmount))
	pdf.Ln(5)

	if bill.ServiceCharge > 0 {
		pdf.Cell(130, 5, "Service Charge:")
		pdf.Cell(30, 5, fmt.Sprintf("%.2f", bill.ServiceCharge))
		pdf.Ln(5)
	}

	// Grand Total
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(130, 5, "TOTAL:")
	pdf.Cell(30, 5, fmt.Sprintf("%.2f", bill.TotalAmount))
	pdf.Ln(5)

	// Status
	pdf.SetFont("Arial", "", 9)
	pdf.Ln(5)
	pdf.Cell(0, 5, fmt.Sprintf("Status: %s", bill.Status))
	pdf.Ln(5)
	if bill.PaidAt != nil {
		pdf.Cell(0, 5, fmt.Sprintf("Paid At: %s", bill.PaidAt.Format("2006-01-02 15:04:05")))
		pdf.Ln(5)
	}

	// Convert to bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Helper function to upload PDF to Digital Ocean
func (s *Service) uploadBillPDF(billNumber string, pdfBytes []byte) (string, error) {
	// Initialize DO storage
	err, doStorage := digital_ocean_storage.NewDOStorage("bills")
	if err != nil {
		return "", err
	}

	if err := doStorage.Run(); err != nil {
		return "", err
	}

	// Get S3 client
	s3Client, ok := doStorage.Get().(*s3.S3)
	if !ok {
		return "", fmt.Errorf("failed to get S3 client")
	}

	// Upload file
	bucketName := config.Config.DigitalOcean.StorageBucket
	fileName := fmt.Sprintf("bills/%s.pdf", billNumber)

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(pdfBytes),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("application/pdf"),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload PDF: %w", err)
	}

	// Build URL
	storageEndpoint := config.Config.DigitalOcean.StorageEndPoint
	fileURL := fmt.Sprintf("https://%s.%s/%s", bucketName, storageEndpoint, fileName)

	return fileURL, nil
}
