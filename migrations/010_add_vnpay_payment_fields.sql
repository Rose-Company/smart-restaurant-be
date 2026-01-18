-- Migration: Add VN-PAY payment fields
-- Description: Adds VN-PAY specific fields to payments table for VN-PAY integration support
-- Created: 2026-01-18

-- Add VN-PAY specific columns to payments table
ALTER TABLE payments
ADD COLUMN IF NOT EXISTS vnpay_transaction_no VARCHAR(100),
ADD COLUMN IF NOT EXISTS vnpay_bank_code VARCHAR(50),
ADD COLUMN IF NOT EXISTS vnpay_card_type VARCHAR(50),
ADD COLUMN IF NOT EXISTS vnpay_response_code VARCHAR(10);

-- Create index for VN-PAY transaction number for quick lookup
CREATE INDEX IF NOT EXISTS idx_payments_vnpay_transaction_no ON payments(vnpay_transaction_no);

-- Update method enum to include vnpay
-- Note: If method is a VARCHAR, this is already flexible
-- If it's an ENUM, you'd need to: ALTER TYPE payment_method ADD VALUE 'vnpay';

-- Verify payments table has all necessary fields
-- SELECT column_name, data_type FROM information_schema.columns WHERE table_name='payments';
