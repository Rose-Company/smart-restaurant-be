package services

import (
	"app-noti/config"
	"app-noti/internal/models"
	"app-noti/internal/util"
	"net/http"
	"strconv"
)

func (s *Service) CreateVNPayPaymentRes(
	requestData map[string]interface{},
	r *http.Request,
) (*models.VNPayResponse, error) {
	vnpCfg := config.Config.Payment.VNPay
	// amount * 100 (Follow VNPAY doc)
	amountRaw := requestData["amount"]
	amount := int64(amountRaw.(float64)) * 100

	orderID := requestData["orderId"].(string)
	language := requestData["language"].(string)

	bankCode := ""
	if v, ok := requestData["bankCode"]; ok {
		bankCode = v.(string)
	}

	vnpParams := vnpCfg.BuildVNPayParams()

	vnpParams["vnp_Amount"] = strconv.FormatInt(amount, 10)
	vnpParams["vnp_OrderInfo"] = orderID
	vnpParams["vnp_TxnRef"] = orderID
	vnpParams["vnp_Locale"] = language
	vnpParams["vnp_IpAddr"] = util.GetIPAddress(r)

	if bankCode != "" {
		vnpParams["vnp_BankCode"] = bankCode
	}

	queryURL := util.GetPaymentURL(vnpParams, true)
	hashData := util.GetPaymentURL(vnpParams, false)

	secureHash := util.HmacSHA512(vnpCfg.SecretKey, hashData)

	paymentURL := vnpCfg.PayURL + "?" + queryURL + "&vnp_SecureHash=" + secureHash

	return &models.VNPayResponse{
		Code:       "OK",
		Message:    "success",
		Amount:     amount,
		PaymentURL: paymentURL,
	}, nil
}
