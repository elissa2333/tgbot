package telegram

// pack

/*Payments
您的机器人可以接受来自Telegram用户的付款。请参阅付款简介，以详细了解该流程以及如何为您的机器人设置付款。请注意，用户需要使用Telegram v.4.0或更高版本才能使用付款（于2017年5月18日发布）。
https://core.telegram.org/bots/api#payments*/

// Invoice 有关发票的基本信息
// https://core.telegram.org/bots/api#invoice
type Invoice struct {
	Title          string `json:"title"`           // 产品名称
	Description    string `json:"description"`     // 产品描述
	StartParameter string `json:"start_parameter"` // 唯一的漫游器深层链接参数可用于生成此发票
	Currency       string `json:"currency"`        // 三字母的ISO 4217 货币代码
	TotalAmount    int64  `json:"total_amount"`    // 以货币的最小单位表示的总价格（整数，不是浮动/双精度）。例如，对于US$ 1.45 pass 的价格amount = 145。看到EXP在参数currencies.json，它显示的数字过去小数点每种货币（2为广大货币）的数目。
}

// ShippingAddress 收货地址
// https://core.telegram.org/bots/api#shippingaddress
type ShippingAddress struct {
	CountryCode string `json:"country_code"`  // ISO 3166-1 alpha-2国家/地区代码
	State       string `json:"state"`         // 州（如果适用）
	City        string `json:"city"`          // 市
	StreetLine1 string `json:"street_line_1"` // 地址的第一行
	StreetLine2 string `json:"street_line_2"` // 地址的第二行
	PostCode    string `json:"post_code"`     // 邮政邮编
}

// OrderInfo 订单的信息
// https://core.telegram.org/bots/api#orderinfo
type OrderInfo struct {
	Name            string           `json:"name"`             // 可选的。用户名
	PhoneNumber     string           `json:"phone_number"`     // 可选的。用户电话
	Email           string           `json:"email"`            // 可选的。用户电子邮件
	ShippingAddress *ShippingAddress `json:"shipping_address"` // 可选的。用户送货地址
}

// SuccessfulPayment 包含有关成功付款的基本信息
// https://core.telegram.org/bots/api#successfulpayment
type SuccessfulPayment struct {
	Currency                string     `json:"currency"`                   // 三字母的 ISO4217 货币代码
	TotalAmount             int64      `json:"total_amount"`               // 以货币的最小单位表示的总价格（整数，不是浮动/双精度）。例如，对于US$ 1.45pass 的价格amount = 145。看到EXP在参数currencies.json，它显示的数字过去小数点每种货币（2为广大货币）的数目。
	InvoicePayload          string     `json:"invoice_payload"`            // Bot指定的发票有效载荷
	ShippingOptionID        string     `json:"shipping_option_id"`         // 可选的。用户选择的运输选项的标识符
	OrderInfo               *OrderInfo `json:"order_info"`                 // 可选的。用户提供的订单信息
	TelegramPaymentChargeID string     `json:"telegram_payment_charge_id"` // 电报付款标识符
	ProviderPaymentChargeID string     `json:"provider_payment_charge_id"` // 提供商付款标识符
}

// ShippingQuery 有关入库运输查询的信息
// https://core.telegram.org/bots/api#shippingquery
type ShippingQuery struct {
	ID              string           `json:"id"`               // 唯一查询标识符
	From            *User            `json:"from"`             // 发送查询的用户
	InvoicePayload  string           `json:"invoice_payload"`  // Bot指定的发票有效载荷
	ShippingAddress *ShippingAddress `json:"shipping_address"` // 用户指定的送货地址
}

// PreCheckoutQuery 有关传入的检出前查询的信息
// https://core.telegram.org/bots/api#precheckoutquery
type PreCheckoutQuery struct {
	ID               string     `json:"id"`                 // 唯一查询标识符
	From             *User      `json:"from"`               // 发送查询的用户
	Currency         string     `json:"currency"`           // 三字母的ISO 4217 货币代码
	TotalAmount      int64      `json:"total_amount"`       // 以货币的最小单位表示的总价格（整数，不是浮动/双精度）。例如，对于US$ 1.45pass 的价格amount = 145。看到EXP在参数currencies.json，它显示的数字过去小数点每种货币（2为广大货币）的数目。
	InvoicePayload   string     `json:"invoice_payload"`    // Bot指定的发票有效载荷
	ShippingOptionID string     `json:"shipping_option_id"` // 可选的。用户选择的运输选项的标识符
	OrderInfo        *OrderInfo `json:"order_info"`         // 可选的。用户提供的订单信息
}
