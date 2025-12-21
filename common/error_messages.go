package common

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	TokenNotFound     = "TokenNotFound"
	TokenUnAuthorized = "TokenUnAuthorized"
	PerDenied         = "PerDenied"
	ResultFailed      = "ResultFailed"
	InternalError     = "InternalError"
)

const (
	DuplicatedDataErr = "Lỗi! Trùng lặp dữ liệu"
	DefaultError      = "Lỗi! Xảy ra lỗi không xác định, vui lòng liên hệ quản trị viên"
)

var DataIsNullErr = func(obj string) string {
	return fmt.Sprintf("%v cannot use nil", obj)
}

var DataIsExisted = func(obj string) string {
	return fmt.Sprintf("%v is existed", obj)
}

var DataIsSmallerZero = func(obj string) string {
	return fmt.Sprintf("%v is not smaller zero", obj)
}

var DataIsBeforeNow = func(obj string) string {
	return fmt.Sprintf("%v is not before now", obj)
}

var ErrorWrapper = func(prefix string, err error) error {
	return fmt.Errorf("%v: %v", prefix, err.Error())
}

var PgErrorTransform = func(err error) error {
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "duplicate key value") {
		return fmt.Errorf(DuplicatedDataErr)
	}

	return err
}

var (
	ErrCodeNotAuthorized = errors.New("not_authorized")
	ErrCodeSystemErr     = errors.New("system_error")
	ErrActionNotAllowed  = errors.New("action_not_allowed")
	ErrTokenNotFound     = errors.New("token_not_found")
	ErrNotAuthorized     = errors.New("not_authorized")
)

var (
	ErrCodeInvalidTimeRange = errors.New("invalid_time_range")
)

var listErrorData = []errData{
	{
		Code:        "cart_not_found",
		HTTPCode:    404,
		MessageViVn: "Giỏ hàng không tồn tại",
		MessageEnUs: "Cart not found",
	},
	{
		Code:        "not_authorized",
		HTTPCode:    401,
		MessageViVn: "Không có quyền truy cập",
		MessageEnUs: "Not authorized",
	},
	{
		Code:        "invalid_cart_status",
		HTTPCode:    400,
		MessageViVn: "Trạng thái giỏ hàng không hợp lệ",
		MessageEnUs: "Invalid cart status",
	},
	{
		Code:        "item_not_in_cart",
		HTTPCode:    400,
		MessageViVn: "Sản phẩm không có trong giỏ hàng",
		MessageEnUs: "Item not in cart",
	},
	{
		Code:        "invalid_promo_code",
		HTTPCode:    400,
		MessageViVn: "Mã khuyến mại không hợp lệ",
		MessageEnUs: "Invalid promo code",
	},
	{
		Code:        "exist_item_in_cart",
		HTTPCode:    400,
		MessageViVn: "Sản phẩm đã có trong giỏ hàng",
		MessageEnUs: "Exist item in cart",
	},
	{
		Code:        "invalid_item_qty",
		HTTPCode:    400,
		MessageViVn: "Số lượng sản phẩm không hợp lệ",
		MessageEnUs: "Invalid item qty",
	},
	{
		Code:        "has_invalid_item_in_cart",
		HTTPCode:    400,
		MessageViVn: "Giỏ hàng có sản phẩm không hợp lệ",
		MessageEnUs: "Cart has invalid item",
	},
	{
		Code:        "you_already_in_class",
		HTTPCode:    400,
		MessageViVn: "Bạn đã tham gia lớp học này rồi",
		MessageEnUs: "You already in class",
	},
	{
		Code:        "promo_code_is_used_for_specific_user",
		HTTPCode:    400,
		MessageViVn: "Mã khuyến mại chỉ dành cho người dùng cụ thể",
		MessageEnUs: "Promo code is used for specific user",
	},
	{
		Code:        "promo_code_is_exceed_max_use",
		HTTPCode:    400,
		MessageViVn: "Mã khuyến mại đã vượt quá số lần sử dụng",
		MessageEnUs: "Promo code is exceed max use",
	},
	{
		Code:        "promo_code_is_used_for_min_amount",
		HTTPCode:    400,
		MessageViVn: "Mã khuyến mại chỉ dành cho đơn hàng có giá trị tối thiểu",
		MessageEnUs: "Promo code is used for min amount",
	},
	{
		Code:        "promo_code_is_used_for_specific_item",
		HTTPCode:    400,
		MessageViVn: "Mã khuyến mại chỉ dành cho sản phẩm cụ thể",
		MessageEnUs: "Promo code is used for specific item",
	},
	{
		Code:        "promo_code_is_expired",
		HTTPCode:    400,
		MessageViVn: "Mã khuyến mại đã hết hạn",
		MessageEnUs: "Promo code is expired",
	},
	{
		Code:        "invalid_time_range",
		HTTPCode:    400,
		MessageViVn: "Thời gian không hợp lệ",
		MessageEnUs: "Invalid time range",
	},
	{
		Code:        "invalid_payment_method",
		HTTPCode:    400,
		MessageViVn: "Phương thức thanh toán không hợp lệ",
		MessageEnUs: "Invalid payment method",
	},
	{
		Code:        "invalid_order_status",
		HTTPCode:    400,
		MessageViVn: "Trạng thái đơn hàng không hợp lệ",
		MessageEnUs: "Invalid order status",
	},
	{
		Code:        "invalid_data",
		HTTPCode:    400,
		MessageViVn: "Dữ liệu không hợp lệ",
		MessageEnUs: "Invalid data",
	},
	{
		Code:        "order_not_found",
		HTTPCode:    404,
		MessageViVn: "Đơn hàng không tồn tại",
		MessageEnUs: "Order not found",
	},
	{
		Code:        "invalid_total_price",
		HTTPCode:    400,
		MessageViVn: "Giá trị đơn hàng đã thay đổi, vui lòng kiểm tra lại!",
		MessageEnUs: "Total price is invalid",
	},
	{
		Code:        "system_error",
		HTTPCode:    500,
		MessageViVn: "Đã có lỗi xảy ra, vui lòng thử lại!",
		MessageEnUs: "System error, please try again!",
	},
	{
		Code:        "item_not_found",
		HTTPCode:    404,
		MessageViVn: "Sản phẩm không tồn tại",
		MessageEnUs: "Item not found",
	},
	{
		Code:        "class_closed",
		HTTPCode:    400,
		MessageViVn: "Lớp học đã hết hạn",
		MessageEnUs: "Class closed",
	},
	{
		Code:        "class_not_found",
		HTTPCode:    404,
		MessageViVn: "Khóa học không tồn tại hoặc không khả dụng",
		MessageEnUs: "Class not found",
	},
	{
		Code:        "user_deactivated",
		HTTPCode:    403,
		MessageViVn: "Tài khoản của bạn đã bị vô hiệu hóa",
		MessageEnUs: "Your account has been deactivated",
	},
	{
		Code:        "promo_code_is_used_for_new_user",
		HTTPCode:    400,
		MessageViVn: "Mã khuyến mại chỉ dành cho người dùng mới",
		MessageEnUs: "Promo code is used for new user",
	},
	{
		Code:        "token_expired",
		HTTPCode:    403,
		MessageViVn: "Phiên làm việc của bạn đã hết hạn",
		MessageEnUs: "Your session has expired",
	},
	{
		Code:        "item_already_exists_in_other_order",
		HTTPCode:    400,
		MessageViVn: "Đã tồn tại sản phẩm này trong 1 đơn hàng khác, vui lòng kiểm tra lại!",
		MessageEnUs: "This item is already in another order, please check again!",
	},
	{
		Code:        "not_allow_trial_mode",
		HTTPCode:    400,
		MessageViVn: "Lớp học không hỗ trợ dạng học thử!",
		MessageEnUs: "Class not support trial mode!",
	},
	{
		Code:        "not_found_registration_request",
		HTTPCode:    400,
		MessageViVn: "Không tìm thấy đơn phù hợp",
		MessageEnUs: "Not found valid request",
	},
	{
		Code:        "invalid_registration_request_payment_status",
		HTTPCode:    400,
		MessageViVn: "Trạng thái thanh toán đã hoàn thành hoặc không phù hợp",
		MessageEnUs: "Payment status completed or invalid",
	},
	{
		Code:        "invalid_registration_request_entrance_test_status",
		HTTPCode:    400,
		MessageViVn: "Trạng thái bài test đã hoàn thành hoặc không phù hợp",
		MessageEnUs: "Entrance test status completed or invalid",
	},
	{
		Code:        "invalid_registration_request_enroll_status",
		HTTPCode:    400,
		MessageViVn: "Đã tham gia vào lớp học hoặc trạng thái không phù hợp",
		MessageEnUs: "Enroll status completed or invalid",
	},
	{
		Code:        "registration_request_not_found",
		HTTPCode:    404,
		MessageViVn: "Yêu cầu đăng ký không tồn tại",
		MessageEnUs: "Registration request not found",
	},
	{
		Code:        "user_already_exists",
		HTTPCode:    400,
		MessageViVn: "Tài khoản đã tồn tại",
		MessageEnUs: "User already exists",
	},
	{
		Code:        "user_not_exists",
		HTTPCode:    404,
		MessageViVn: "Tài khoản không tồn tại",
		MessageEnUs: "User not exists",
	},
	{
		Code:        "email_already_exists",
		HTTPCode:    400,
		MessageViVn: "Email đã tồn tại",
		MessageEnUs: "Email already exists",
	},
	{
		Code:        "not_found_registration_request",
		HTTPCode:    400,
		MessageViVn: "Không tìm thấy đơn phù hợp",
		MessageEnUs: "Not found valid request",
	},
	{
		Code:        "invalid_registration_request_payment_status",
		HTTPCode:    400,
		MessageViVn: "Trạng thái thanh toán đã hoàn thành hoặc không phù hợp",
		MessageEnUs: "Payment status completed or invalid",
	},
	{
		Code:        "invalid_registration_request_entrance_test_status",
		HTTPCode:    400,
		MessageViVn: "Trạng thái bài test đã hoàn thành hoặc không phù hợp",
		MessageEnUs: "Entrance test status completed or invalid",
	},
	{
		Code:        "invalid_registration_request_enroll_status",
		HTTPCode:    400,
		MessageViVn: "Đã tham gia vào lớp học hoặc trạng thái không phù hợp",
		MessageEnUs: "Enroll status completed or invalid",
	},
	{
		Code:        "missing_email",
		HTTPCode:    400,
		MessageViVn: "Email rỗng",
		MessageEnUs: "Empty email",
	},
	{
		Code:        "action_not_allowed",
		HTTPCode:    403,
		MessageViVn: "Bạn không có quyền thực hiện",
		MessageEnUs: "You're not allow",
	},
	{
		Code:        "token_not_found",
		HTTPCode:    401,
		MessageViVn: "Không tìm thấy token",
		MessageEnUs: "Token not found",
	},
	{
		Code:        ErrActionNotAllowed.Error(),
		HTTPCode:    403,
		MessageViVn: "Hành động không được phép",
		MessageEnUs: "Action not allowed",
	},
}

var (
	AllErrors *MasterErrData
)

func FetchMasterErrData() {
	AllErrors = NewMasterErrData()
	AllErrors.fetchAll()
}

type errData struct {
	Code        string `json:"code" gorm:"column:code"`
	HTTPCode    int    `json:"http_code" gorm:"column:http_code"`
	MessageViVn string `json:"message_vi_vn" gorm:"column:message_vi_vn"`
	MessageEnUs string `json:"message_en_us" gorm:"column:message_en_us"`
}

type ExtraData struct {
	OrderID int64 `json:"order_id,omitempty"`
}

type LocalizeErrRes struct {
	Code      string     `json:"code,omitempty"`
	Message   string     `json:"message,omitempty"`
	HTTPCode  int        `json:"-"`
	Internal  string     `json:"internal,omitempty"`
	ExtraData *ExtraData `json:"extra_data,omitempty"`
}

func (a *LocalizeErrRes) Error() string {
	return a.Code
}

type MasterErrData struct {
	mutex sync.Mutex
	data  map[string]errData
}

// Error data

func NewMasterErrData() *MasterErrData {
	return &MasterErrData{}
}

func (a *MasterErrData) fetchAll() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	for _, errMessage := range listErrorData {
		if a.data == nil {
			a.data = make(map[string]errData)
		}
		a.data[errMessage.Code] = errMessage
	}
}

func (a *MasterErrData) New(err error, language string, internal ...string) *LocalizeErrRes {
	errRes := new(LocalizeErrRes)
	ok := errors.As(err, &errRes)
	if !ok {
		errRes = &LocalizeErrRes{
			Code:    "bad_request",
			Message: "Đã có lỗi xảy ra, vui lòng thử lại!",
		}
		if len(internal) > 0 {
			errRes.Internal = internal[0]
		}
		errFromDB, exists := a.data[err.Error()]
		if exists {
			errRes.Code = errFromDB.Code
			errRes.HTTPCode = errFromDB.HTTPCode
			switch language {
			case "vi":
				errRes.Message = errFromDB.MessageViVn
			default:
				errRes.Message = errFromDB.MessageEnUs
			}
		} else {
			errRes.HTTPCode = 400
		}
	}

	if len(internal) > 0 {
		errRes.Internal = internal[0]
	}
	return errRes
}

// Error res

func (a *LocalizeErrRes) SetMessage(message string) *LocalizeErrRes {
	a.Message = message
	return a
}

func (a *LocalizeErrRes) ReplaceDescByVars(args ...interface{}) *LocalizeErrRes {
	for _, arg := range args {
		a.Message = fmt.Sprintf(a.Message, arg)
	}
	return a
}

func (a *LocalizeErrRes) SetOrderIDToExtraData(orderID int64) *LocalizeErrRes {
	if a.ExtraData == nil {
		a.ExtraData = new(ExtraData)
	}
	a.ExtraData.OrderID = orderID
	return a
}

func (a *LocalizeErrRes) ConvertToBaseError() Response {
	res := BaseResponse(REQUEST_FAILED, a.Message, a.Internal, a.ExtraData)
	res.SetErrorCode(a.Code)
	return res
}

func AbortWithError(c *gin.Context, err error) {
	errJSON := AllErrors.New(err, "vi", err.Error())
	c.AbortWithStatusJSON(errJSON.HTTPCode, errJSON.ConvertToBaseError())
}
