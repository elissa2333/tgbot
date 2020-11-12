package telegram

// pack

/*Payments
您的机器人可以接受来自Telegram用户的付款。请参阅付款简介，以详细了解该流程以及如何为您的机器人设置付款。请注意，用户需要使用Telegram v.4.0或更高版本才能使用付款（于2017年5月18日发布）。
https://core.telegram.org/bots/api#payments*/

// SendInvoiceOptional SendInvoice 可选参数
type SendInvoiceOptional struct {
	ProviderData              string                `json:"provider_data"`                 // 有关发票的JSON序列化数据，将与付款提供商共享。付款人应提供必填字段的详细说明。
	PhotoURL                  string                `json:"photo_url"`                     // 发票的产品照片网址。可以是商品的照片或服务的营销图片。人们看到付款后会更喜欢它。
	PhotoSize                 int64                 `json:"photo_size"`                    // 照片尺寸
	PhotoWidth                int64                 `json:"photo_width"`                   // 照片宽度
	PhotoHeight               int64                 `json:"photo_height"`                  // 照片高度
	NeedName                  bool                  `json:"need_name"`                     // 如果您需要用户的全名来完成订单，请传递True。
	NeedPhoneNumber           bool                  `json:"need_phone_number"`             // 如果您需要用户的电话号码来完成订单，请传递True。
	NeedEmail                 bool                  `json:"need_email"`                    // 如果您需要用户的电子邮件地址来完成订单，请传递True。
	NeedShippingAddress       bool                  `json:"need_shipping_address"`         // 如果您需要用户的送货地址来完成订单，请传递True。
	SendPhoneNumberToProvider bool                  `json:"send_phone_number_to_provider"` // 如果将用户的电话号码发送给提供商，则通过True
	SendEmailToProvider       bool                  `json:"send_email_to_provider"`        // 如果将用户的电子邮件地址发送给提供商，则传递True
	IsFlexible                bool                  `json:"is_flexible"`                   // 如果最终价格取决于运输方式，请通过True
	DisableNotification       bool                  `json:"disable_notification"`          // 静默发送消息。用户将收到没有声音的通知。
	ReplyToMessageID          int64                 `json:"reply_to_message_id"`           // 如果消息是答复，则为原始消息的ID
	AllowSendingWithoutReply  bool                  `json:"allow_sending_without_reply"`   // 如果未发送指定的回复消息也应发送消息，则传递True
	ReplyMarkup               *InlineKeyboardMarkup `json:"reply_markup"`                  // 嵌入式键盘的JSON序列化对象。如果为空，total price将显示一个“付款”按钮。如果不为空，则第一个按钮必须是“付款”按钮。
}

// SendInvoice 使用此方法发送发票。成功后，将返回发送的 Message
// https://core.telegram.org/bots/api#sendinvoice
func (a API) SendInvoice(chatID string, title string, description string, payload string, providerToken string, startParameter string, currency string, prices []LabeledPrice, optional *SendInvoiceOptional) (*Message, error) {
	result := &Message{}
	err := a.handleOptional("/sendInvoice", map[string]interface{}{"chat_id": chatID, "title": title, "description": description, "payload": payload, "provider_token": providerToken, "start_parameter": startParameter, "currency": currency, "prices": prices}, optional, result)
	return result, err
}

// AnswerShippingQueryOptional AnswerShippingQuery 可选参数
type AnswerShippingQueryOptional struct {
	ShippingOptions []ShippingOption `json:"shipping_options"` // 如果ok为True，则为必需。可用的送货选项的JSON序列化数组。
	ErrorMessage    string           `json:"error_message"`    // 如果ok为False，则为必需。错误消息以易于阅读的形式解释了为什么无法完成订单（例如，“抱歉，无法发送到您想要的地址”。）电报将向用户显示此消息。
}

// AnswerShippingQuery  如果您发送了要求提供送货地址的发票，并且指定了参数is_flexible，则Bot API会将带有Shipping_query字段的 Update 发送给机器人。使用此方法可以答复运输查询。成功时，返回True。
// https://core.telegram.org/bots/api#answershippingquery
func (a API) AnswerShippingQuery(shippingQueryID string, ok bool, optional *AnswerShippingQueryOptional) (bool, error) {
	var result bool
	err := a.handleOptional("/answerShippingQuery", map[string]interface{}{"shipping_query_id": shippingQueryID, "ok": ok}, optional, &result)
	return result, err
}

// AnswerPreCheckoutQueryOptional AnswerPreCheckoutQuery 可选参数
type AnswerPreCheckoutQueryOptional struct {
	ErrorMessage string `json:"error_message"` // 如果ok为False，则为必需。易读格式的错误消息，说明无法进行结帐的原因（例如，“抱歉，有人在您忙于填写付款详细信息时，刚买了我们最后一件令人惊异的黑色T恤。请选择其他颜色或服装！”）。电报将向用户显示此消息。
}

// AnswerPreCheckoutQuery 一旦用户确认了付款和运输明细，Bot API就会以带有字段pre_checkout_query的Update形式发送最终确认。使用此方法来响应此类预结帐查询。成功时，返回True。注意： Bot API必须在发送预结帐查询后的10秒钟内收到答案。
// https://core.telegram.org/bots/api#answerprecheckoutquery
func (a API) AnswerPreCheckoutQuery(preCheckoutQueryID string, ok bool, optional *AnswerPreCheckoutQueryOptional) (bool, error) {
	var result bool
	err := a.handleOptional("/answerPreCheckoutQuery", map[string]interface{}{"pre_checkout_query_id": preCheckoutQueryID, "ok": ok}, optional, &result)
	return result, err
}

// LabeledPrice 此对象代表商品或服务价格的一部分。
// https://core.telegram.org/bots/api#labeledprice
type LabeledPrice struct {
	Label  string `json:"label"`  // 部分标签
	Amount int64  `json:"amount"` // 价在该产品的最小单位的的货币（整数，不浮动/双）。例如，对于US$ 1.45pass的价格amount = 145。看到EXP在参数currencies.json，它显示的数字过去小数点每种货币（2为广大货币）的数目。
}

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

// ShippingOption 该对象表示一种运输选项。
// https://core.telegram.org/bots/api#shippingoption
type ShippingOption struct {
	ID     string         `json:"id"`     // 运输选项标识符
	Title  string         `json:"title"`  // 选项标题
	Prices []LabeledPrice `json:"prices"` // 价格部分清单
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
