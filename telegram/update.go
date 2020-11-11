package telegram

// pack

/*Getting updates
有两种相互排斥的方式来接收您的机器人的更新：一方面是getUpdates方法，另一方面是Webhooks。传入更新存储在服务器上，直到僵尸程序以任何一种方式接收到它们为止，但它们的保存时间不会超过24小时。
无论选择哪个选项，都会收到JSON序列化的Update对象。
https://core.telegram.org/bots/api#update*/

import (
	"io"

	"github.com/elissa2333/httpc"

	"github.com/elissa2333/tgbot/utils"
)

// Update 该对象表示传入的更新。
//最多一个可选参数可以出现在任何给定的更新。
// https://core.telegram.org/bots/api#update
type Update struct {
	UpdateID           int64               `json:"update_id,omitempty"`            // 更新的唯一标识符。更新标识符从某个正数开始，并依次增加。如果您使用的是Webhooks，则此ID变得特别方便，因为它允许您忽略重复的更新或恢复正确的更新顺序（如果它们不正常的话）。如果至少有一个星期没有新更新，则将随机选择下一个更新的标识符，而不是顺序选择。
	Message            *Message            `json:"message,omitempty"`              // 可选的。任何形式的新传入消息-文本，照片，标签等。
	EditedMessage      *Message            `json:"edited_message,omitempty"`       // 可选的。机器人已知并已编辑的消息的新版本
	ChannelPost        *Message            `json:"channel_post,omitempty"`         // 可选的。任何形式的新传入频道帖子-文字，照片，贴纸等。
	EditedChannelPost  *Message            `json:"edited_channel_post,omitempty"`  // 可选的。机器人已知并已编辑的频道发布的新版本
	InlineQuery        *InlineQuery        `json:"inline_query,omitempty"`         // 可选的。新的传入内联查询
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"` // 可选的。用户选择并发送给他们的聊天伙伴的内联查询的结果。请参阅收集反馈的文档，以获取有关如何为您的机器人启用这些更新的详细信息。
	CallbackQuery      *CallbackQuery      `json:"callback_query,omitempty"`       // 可选的。新传入的回调查询
	ShippingQuery      *ShippingQuery      `json:"shipping_query,omitempty"`       // 可选的。新的收货查询。仅适用于价格灵活的发票
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query,omitempty"`   // 可选的。新的传入预结帐查询。包含有关结帐的完整信息
	Poll               *Poll               `json:"poll,omitempty"`                 // 可选的。新的投票状态。漫游器仅接收有关僵尸程序发送的有关已停止的轮询和轮询的更新
	PollAnswer         *PollAnswer         `json:"poll_answer,omitempty"`          // 可选的。用户在非匿名调查中更改了答案。僵尸程序仅在由僵尸程序本身发送的民意调查中才能获得新的选票。
}

// GetUpdates 使用此方法可以使用长轮询（wiki）接收传入的更新。返回更新对象数组
// https://core.telegram.org/bots/api#getupdates
func (a *API) GetUpdates(offset int64, limit uint, timeout uint, allowedUpdates ...string) ([]Update, error) { // TODO 改为GET
	res, err := a.HTTPClient.SetBody(map[string]interface{}{"offset": offset, "limit": limit, "timeout": timeout, "allowed_updates": allowedUpdates}).Post("/getUpdates")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var m []Update
	err = HandleResp(res, &m)
	return m, err
}

// WebhookOptional SetWebhook 可选参数
type WebhookOptional struct {
	Certificate        InputFile `json:"certificate,omitempty"`     // 上传您的公共密钥证书，以便可以检查正在使用的根证书。有关详细信息，请参见我们的自签名指南。
	IPAddress          string    `json:"ip_address"`                // 固定IP地址将用于发送Webhook请求，而不是通过DNS解析的IP地址
	MaxConnections     int       `json:"max_connections,omitempty"` // 与Webhook进行更新交付的同时HTTPS连接的最大允许数量为1-100。默认为40。使用较低的值可以限制bot服务器的负载，使用较高的值可以增加bot的吞吐量。
	AllowedUpdates     []string  `json:"allowed_updates,omitempty"` // 您希望机器人接收的更新类型的JSON序列化列表。例如，指定[“ message”，“ edited_channel_post”，“ callback_query”]仅接收这些类型的更新。请参阅更新以获取可用更新类型的完整列表。指定一个空列表以接收所有更新，无论类型如何（默认）。如果未指定，将使用以前的设置。
	DropPendingUpdates bool      `json:"drop_pending_updates"`      // 传递True以删除所有待处理的更新
}

// SetWebhook 指定URL并通过传出的Webhook接收传入的更新。只要该漫游器有更新，我们就会向指定的URL发送一个HTTPS POST请求，其中包含JSON序列化的Update。如果请求失败，我们将在合理的尝试后放弃。成功返回True。
//如果您想确保Webhook请求来自Telegram，建议您在URL中使用秘密路径，例如 `https://www.example.com/<token>`。由于没有其他人知道您的漫游器令牌，因此您可以确定它是我们。
// https://core.telegram.org/bots/api#setwebhook
func (a *API) SetWebhook(url string, optional *WebhookOptional) error { // TODO 仅允许 telegram 的 ip 进行访问 我直接监听指定的URI
	var rows []httpc.FromDataRow
	rows = append(rows, httpc.FromDataRow{
		Key:   "url",
		Value: url,
	})

	if optional != nil { // TODO 检查一下是否可以这样干
		m, err := utils.StructToMap(optional)
		if err != nil {
			return err
		}

		for k, v := range m {
			reader, ok := v.(io.Reader)
			if ok {
				rows = append(rows, httpc.FromDataRow{
					Key:  k,
					Data: reader,
				})
			} else {
				rows = append(rows, httpc.FromDataRow{
					Key:   k,
					Value: utils.ToString(v),
				})
			}
		}
	}

	res, err := a.HTTPClient.SetFromData(rows...).Post("/setWebhook")
	if err != nil {
		return err
	}

	var result bool
	return HandleResp(res, &result)
}

// DeleteWebhookOptional DeleteWebhook 可选参数
type DeleteWebhookOptional struct {
	DropPendingUpdates string `json:"drop_pending_updates"` // 传递True以删除所有待处理的更新
}

// DeleteWebhook 如果您决定切换回 getUpdates ，请使用此方法删除webhook集成
// https://core.telegram.org/bots/api#deletewebhook
func (a API) DeleteWebhook(optional *DeleteWebhookOptional) error {
	res, err := a.HTTPClient.SetBody(optional).Post("/deleteWebhook")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var result bool
	return HandleResp(res, &result)
}

// WebhookInfo  Webhook当前状态的信息
// https://core.telegram.org/bots/api#webhookinfo
type WebhookInfo struct {
	URL                  string   `json:"url,omitempty"`                    // Webhook URL，如果未设置webhook，则可能为空
	HasCustomCertificate bool     `json:"has_custom_certificate,omitempty"` // 如果为webhook证书检查提供了自定义证书
	PendingUpdateCount   int64    `json:"pending_update_count,omitempty"`   // 等待交付的更新数量
	IPAddress            string   `json:"ip_address"`                       // 可选的。当前使用的Webhook IP地址
	LastErrorDate        int64    `json:"last_error_date,omitempty"`        // 可选的。尝试通过Webhook传递更新时发生的最新错误的Unix时间
	LastErrorMessage     string   `json:"last_error_message,omitempty"`     // 可选的。尝试通过Webhook传递更新时发生的最新错误的人类可读格式的错误消息
	MaxConnections       int      `json:"max_connections,omitempty"`        // 可选的。与Webhook进行更新交付的同时HTTPS连接的最大允许数量
	AllowedUpdates       []string `json:"allowed_updates,omitempty"`        // 可选的。机器人已订阅的更新类型的列表。默认为所有更新类型
}

// GetWebhookInfo 使用此方法获取当前的Webhook状态。不需要参数。成功时，返回 WebhookInfo 对象。如果机器人正在使用 GetUpdates，将返回一个url字段为空的对象。
// https://core.telegram.org/bots/api#getwebhookinfo
func (a API) GetWebhookInfo() (*WebhookInfo, error) {
	res, err := a.HTTPClient.Get("/GetWebhookInfo")
	if err != nil {
		return nil, err
	}

	m := &WebhookInfo{}
	err = HandleResp(res, m)
	return m, err
}
