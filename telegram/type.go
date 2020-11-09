package telegram

import "io"

// pack

/*Available types
API API响应中使用的所有类型均表示为JSON对象。
除非另有说明，否则使用32位带符号整数存储所有 Integer 字段是安全的。
https://core.telegram.org/bots/api#available-types*/

// User 电报用户或bot
// https://core.telegram.org/bots/api#user
type User struct {
	ID                      int64  `json:"id,omitempty"`                          // 该用户或bot的唯一标识符
	IsBot                   bool   `json:"is_bot,omitempty"`                      // 如果该用户是机器人
	FirstName               string `json:"first_name,omitempty"`                  // 用户或bot的名字
	LastName                string `json:"last_name,omitempty"`                   // 可选的。用户或bot的姓氏
	Username                string `json:"username,omitempty"`                    // 可选的。用户或bot的用户名
	LanguageCode            string `json:"language_code,omitempty"`               // 可选的。用户语言的IETF语言标签
	CanJoinGroups           bool   `json:"can_join_groups,omitempty"`             // 可选的。是的，如果该bot可以被邀请加入小组。仅在 GetMe 中返回。
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"` // 可选的。是的，如果该bot禁用了隐私模式。仅在 GetMe 中返回。
	SupportsInlineQueries   bool   `json:"supports_inline_queries,omitempty"`     // 可选的。如果机器人支持内联查询，则为 true。仅在 GetMe中返回。
}

const (
	// ChatTypeAtPrivate 私有
	ChatTypeAtPrivate = "private"
	// ChatTypeAtGroup 群组
	ChatTypeAtGroup = "group"
	// ChatTypeAtSuperGroup 超级群组
	ChatTypeAtSuperGroup = "supergroup"
	// ChatTypeAtChannel 频道
	ChatTypeAtChannel = "channel"
)

// Chat 聊天
// https://core.telegram.org/bots/api#chat
type Chat struct {
	ID               int64            `json:"id,omitempty"`                  // 此聊天的唯一标识符。此数字可能大于32位，并且某些编程语言在解释它时可能会有缺陷。但是它小于52位，因此带符号的64位整数或双精度浮点类型对于存储此标识符是安全的。
	Type             string           `json:"type,omitempty"`                // 聊天类型，可以是 “private”, “group”, “supergroup” or “channel”
	Title            string           `json:"title,omitempty"`               // 可选的。标题，用于超级群组，频道和群组聊天
	Username         string           `json:"username,omitempty"`            // 可选的。用户名，用于私人聊天，超组和频道（如果有）
	FirstName        string           `json:"first_name,omitempty"`          // 可选的。私人聊天中对方的名字
	LastName         string           `json:"last_name,omitempty"`           // 可选的。私人聊天中对方的姓氏
	Photo            *ChatPhoto       `json:"photo,omitempty"`               // 可选的。聊天照片。仅在 GetChat 中返回。
	Bio              string           `json:"bio"`                           // 可选的。私人聊天中对方的个人简介。仅在getChat中返回。
	Description      string           `json:"description,omitempty"`         // 可选的。说明，用于群组，超群组和频道聊天。仅在 GetChat 中返回。
	InviteLink       string           `json:"invite_link,omitempty"`         // 可选的。聊天邀请链接，用于群组，超级群组和频道聊天。聊天中的每个管理员都会生成自己的邀请链接，因此bot必须首先使用exportChatInviteLink生成链接。仅在 GetChat 中返回。
	PinnedMessage    *Message         `json:"pinned_message,omitempty"`      // 可选的。固定消息，用于组，超组和通道。仅在 GetChat 中返回。
	Permissions      *ChatPermissions `json:"permissions,omitempty"`         // 可选的。组和超级组的默认聊天成员权限。仅在 GetChat 中返回。
	SlowModeDelay    int64            `json:"slow_mode_delay,omitempty"`     // 可选的。对于超组，每个非特权用户发送的连续消息之间允许的最小延迟。仅在 GetChat 中返回。
	StickerSetName   string           `json:"sticker_set_name,omitempty"`    // 可选的。对于超组，请使用组贴纸集的名称。仅在 GetChat 中返回。
	CanSetStickerSet bool             `json:"can_set_sticker_set,omitempty"` // 可选的。是的，如果漫游器可以更改组标签集。仅在 GetChat 中返回。
	LinkedChatID     int64            `json:"linked_chat_id"`                // 可选的。链接聊天的唯一标识符，即频道的讨论组标识符，反之亦然；用于超级群组和频道聊天。该标识符可能大于32位，并且某些编程语言在解释它时可能会有困难/无声的缺陷。但是它小于52位，因此带符号的64位整数或双精度浮点类型对于存储此标识符是安全的。仅在getChat中返回。
	Location         *ChatLocation    `json:"location"`                      // 可选的。对于超组，是超组连接到的位置。仅在getChat中返回。
}

// Message 一条消息
// https://core.telegram.org/bots/api#message
type Message struct {
	MessageID               int64                    `json:"message_id,omitempty"`              // 此聊天中的唯一消息标识符
	From                    *User                    `json:"from,omitempty"`                    // 可选的。发件人，对于发送到渠道的消息为空
	SenderChat              *Chat                    `json:"sender_chat"`                       // 可选的。消息发送方，代表聊天室发送。频道本身用于频道消息。超组本身用于接收来自匿名组管理员的消息。消息的链接通道自动转发到讨论组
	Data                    int64                    `json:"data,omitempty"`                    // 消息在Unix时间中发送的日期
	Chat                    *Chat                    `json:"chat,omitempty"`                    // 消息所属的会话
	ForwardFrom             *User                    `json:"forward_from,omitempty"`            // 可选的。对于转发的邮件，原始邮件的发件人
	ForwardFromChat         *Chat                    `json:"forward_from_chat,omitempty"`       // 可选的。对于从频道转发的消息，有关原始频道的信息
	ForwardFromMessageID    int64                    `json:"forward_from_message_id,omitempty"` // 可选的。对于从通道转发的消息，是通道中原始消息的标识符
	ForwardSignature        string                   `json:"forward_signature,omitempty"`       // 可选的。对于从频道转发的消息，请提供帖子作者的签名（如果有）
	ForwardSenderName       string                   `json:"forward_sender_name,omitempty"`     // 可选的。从用户转发的邮件的发件人名称，这些用户不允许在转发的邮件中添加指向其帐户的链接
	ForwardDate             int64                    `json:"forward_date,omitempty"`            // 可选的。对于转发的消息，原始消息的发送日期为Unix时间
	ReplyToMessage          *Message                 `json:"reply_to_message,omitempty"`        // 可选的。对于答复，原始消息。请注意，即使此字段本身是答复，该字段中的 Message 对象也不会包含其他的reply_to_message字段。
	ViaBot                  *User                    `json:"via_bot,omitempty"`                 // 可选的。发送消息的机器人
	EditDate                int64                    `json:"edit_date,omitempty"`               // 可选的。消息最后一次在Unix时间中编辑的日期
	MediaGroupID            string                   `json:"media_group_id,omitempty"`          // 可选的。该消息所属的媒体消息组的唯一标识符
	AuthorSignature         string                   `json:"author_signature,omitempty"`        // 可选的。在频道中留言的帖子作者的签名
	Text                    string                   `json:"text,omitempty"`                    // 可选的。对于文本消息，消息的实际UTF-8文本，0-4096个字符
	Entities                []MessageEntity          `json:"entities,omitempty"`                // 可选的。对于文本消息，出现在文本中的特殊实体，例如用户名，URL，机器人命令等。
	Animation               *Animation               `json:"animation,omitempty"`               // 可选的。消息是动画，有关动画的信息。为了向后兼容，设置此字段时，还将设置文档字段
	Audio                   *Audio                   `json:"audio,omitempty"`                   // 可选的。消息是音频文件，有关该文件的信息
	Document                *Document                `json:"document,omitempty"`                // 可选的。消息是常规文件，有关文件的信息
	Photo                   []PhotoSize              `json:"photo,omitempty"`                   // 可选的。邮件是照片，照片的可用尺寸
	Sticker                 *Sticker                 `json:"sticker,omitempty"`                 // 可选的。消息是贴纸，有关贴纸的信息
	Video                   *Video                   `json:"video,omitempty"`                   // 可选的。消息是视频，有关视频的信息
	VideoNote               *VideoNote               `json:"video_note,omitempty"`              // 可选的。留言是视频笔记，有关视频留言的信息
	Voice                   *Voice                   `json:"voice,omitempty"`                   // 可选的。消息是语音消息，有关文件的信息
	Caption                 string                   `json:"caption,omitempty"`                 // 可选的。动画，音频，文档，照片，视频或语音的标题，0-1024个字符
	CaptionEntities         []MessageEntity          `json:"caption_entities,omitempty"`        // 可选的。对于带标题的邮件，出现在标题中的特殊实体，例如用户名，URL，漫游器命令等。
	Contact                 *Contact                 `json:"contact,omitempty"`                 // 可选的。消息是共享的联系人，有关该联系人的信息
	Dice                    *Dice                    `json:"dice,omitempty"`                    // 可选的。消息是一个骰子，具有从1到6的随机值
	Game                    *Game                    `json:"game,omitempty"`                    // 可选的。消息是一个游戏，有关游戏的信息
	Poll                    *Poll                    `json:"poll,omitempty"`                    // 可选的。消息是本机民意测验，有关民意测验的信息
	Venue                   *Venue                   `json:"venue,omitempty"`                   // 可选的。消息是一个场地，有关该场地的信息。为了向后兼容，设置此字段时，还将设置位置字段
	Location                *Location                `json:"location,omitempty"`                // 可选的。消息是共享位置，有关位置的信息
	NewChatMembers          []User                   `json:"new_chat_members,omitempty"`        // 可选的。添加到组或超组中的新成员以及有关它们的信息（机器人本身可能是这些成员之一）
	LeftChatMember          *User                    `json:"left_chat_member,omitempty"`        // 可选的。成员已从群组中删除，有关他们的信息（该成员可能是漫游器本身）
	EwChatTitle             string                   `json:"ew_chat_title,omitempty"`           // 可选的。聊天标题已更改为此值
	NewChatPhoto            []PhotoSize              `json:"new_chat_photo,omitempty"`          // 可选的。聊天照片已更改为此值
	DeleteChatPhoto         bool                     `json:"delete_chat_photo,omitempty"`       // 可选的。服务消息：聊天照片已删除
	GroupChatCreated        bool                     `json:"group_chat_created,omitempty"`      // 可选的。服务信息：组已创建
	SupergroupChatCreated   bool                     `json:"supergroup_chat_created,omitempty"` // 可选的。服务消息：超组已创建。在通过更新发送的消息中无法接收到该字段，因为bot在创建时不能成为超组的成员。仅当有人回复直接创建的超组中的第一条消息时，才可以在reply_to_message中找到该消息。
	ChannelChatCreated      bool                     `json:"channel_chat_created,omitempty"`    // 可选的。服务信息：频道已创建。在通过更新发送的消息中无法接收到该字段，因为bot在创建时不能成为频道的成员。如果有人回复频道中的第一条消息，则只能在reply_to_message中找到它。
	MigrateToChatID         int64                    `json:"migrate_to_chat_id,omitempty"`      // 可选的。该组已迁移到具有指定标识符的超组。该数字可能大于32位，并且某些编程语言在解释它时可能会有困难/无声的缺陷。但是它小于52位，因此带符号的64位整数或双精度浮点类型对于存储此标识符是安全的。
	MigrateFromChatID       int64                    `json:"migrate_from_chat_id,omitempty"`    // 可选的。超级组已从具有指定标识的组中迁移。该数字可能大于32位，并且某些编程语言在解释它时可能会有困难/无声的缺陷。但是它小于52位，因此带符号的64位整数或双精度浮点类型对于存储此标识符是安全的。
	PinnedMessage           *Message                 `json:"pinned_message,omitempty"`          // 可选的。指定的消息已固定。请注意，即使该字段本身是答复，该字段中的Message对象也不会包含其他的reply_to_message字段。
	Invoice                 *Invoice                 `json:"invoice,omitempty"`                 // 可选的。消息是付款的发票，有关发票的信息
	SuccessfulPayment       *SuccessfulPayment       `json:"successful_payment,omitempty"`      // 可选的。消息是有关成功付款的服务消息，有关付款的信息
	ConnectedWebsite        string                   `json:"connected_website,omitempty"`       // 可选的。用户登录的网站的域名
	PassportData            *PassportData            `json:"passport_data,omitempty"`           // 可选的。电报护照数据
	ProximityAlertTriggered *ProximityAlertTriggered `json:"proximity_alert_triggered"`         // 可选的。服务消息。聊天中的用户在共享实时位置时触发了另一个用户的接近警报。
	ReplyMarkup             *InlineKeyboardButton    `json:"reply_markup,omitempty"`            // 可选的。消息附带的嵌入式键盘。login_url按钮表示为普通url按钮
}

const (
	// MessageEntityAtMention 提及
	MessageEntityAtMention = "mention"
	// MessageEntityAtHashtag has和标签
	MessageEntityAtHashtag = "hashtag"
	// MessageEntityAtCashtag 现金标签
	MessageEntityAtCashtag = "cashtag"
	// MessageEntityAtBotCommand bot 命令
	MessageEntityAtBotCommand = "bot_command"
	// MessageEntityAtURL URL
	MessageEntityAtURL = "url"
	// MessageEntityAtEmail 电子邮件
	MessageEntityAtEmail = "email"
	// MessageEntityAtPhoneNumber // 手机号
	MessageEntityAtPhoneNumber = "phone_number"
	// MessageEntityAtBold 粗体
	MessageEntityAtBold = "bold"
	// MessageEntityAtItalic 斜体
	MessageEntityAtItalic = "italic"
	// MessageEntityAtUnderline 下划线
	MessageEntityAtUnderline = "underline"
	// MessageEntityAtStrikethrough 删除线
	MessageEntityAtStrikethrough = "strikethrough"
	//MessageEntityAtCode 代码
	MessageEntityAtCode = "code"
	// MessageEntityAtPre 等宽块
	MessageEntityAtPre = "pre"
	// MessageEntityAtTextLink 用于可点击的文本网址
	MessageEntityAtTextLink = "text_link"
	// MessageEntityAtTextMention 适用于没有用户名的用户
	MessageEntityAtTextMention = "text_mention"
)

// MessageEntity 消息中的一个特殊实体。例如，标签，用户名，URL等
// https://core.telegram.org/bots/api#messageentity
type MessageEntity struct {
	Type     string `json:"type,omitempty"`     // 实体的类型。可以是 “mention” (@username), “hashtag” (#hashtag), “cashtag” ($USD), “bot_command” (/start@jobs_bot), “url” (https://telegram.org), “email” (do-not-reply@telegram.org), “phone_number” (+1-212-555-0123), “bold” (bold text), “italic” (italic text), “underline” (underlined text), “strikethrough” (strikethrough text), “code” (monowidth string), “pre” (monowidth block), “text_link” (for clickable text URLs), “text_mention” (for users without usernames)
	Offset   int64  `json:"offset,omitempty"`   // 以UTF-16代码单位向实体开始的偏移量
	Length   int64  `json:"length,omitempty"`   // 实体的长度（以UTF-16代码单元为单位）
	URL      string `json:"url,omitempty"`      // 可选的。仅对于“text_link”，用户点击文本后将打开的URL
	User     *User  `json:"user,omitempty"`     // 可选的。仅针对“text_mention”，提到的用户
	Language string `json:"language,omitempty"` // 可选的。仅对于“ pre”，实体文本的编程语言
}

// PhotoSize 照片或文件/标签缩略图的一种尺寸
// https://core.telegram.org/bots/api#photosize
type PhotoSize struct {
	FileID       string `json:"file_id,omitempty"`        // 该文件的标识符，可用于下载或重复使用该文件
	FileUniqueID string `json:"file_unique_id,omitempty"` // 此文件的唯一标识符，随着时间的推移，对于不同的漫游器，应该是相同的。不能用于下载或重复使用文件。
	Width        int64  `json:"width,omitempty"`          // 照片宽度
	Height       int64  `json:"height,omitempty"`         // 照片高度
	FileSize     int64  `json:"file_size,omitempty"`      // 可选的。文件大小
}

// Animation 动画文件（无声音的GIF或H.264 / MPEG-4 AVC视频）。
// https://core.telegram.org/bots/api#animation
type Animation struct {
	FileID       string     `json:"file_id,omitempty"`        // 该文件的标识符，可用于下载或重复使用该文件
	FileUniqueID string     `json:"file_unique_id,omitempty"` // 此文件的唯一标识符，随着时间的推移，对于不同的漫游器，应该是相同的。不能用于下载或重复使用文件。
	Width        int64      `json:"width,omitempty"`          // 发件人定义的视频宽度
	Height       int64      `json:"height,omitempty"`         // 发件人定义的视频高度
	Duration     int64      `json:"duration,omitempty"`       // 发件人定义的视频时长（以秒为单位）
	Thumb        *PhotoSize `json:"thumb,omitempty"`          // 可选的。发件人定义的动画缩略图
	FileName     string     `json:"file_name,omitempty"`      // 可选的。由发送方定义的原始动画文件名
	MimeType     string     `json:"mime_type,omitempty"`      // 可选的。发件人定义的文件的MIME类型
	FileSize     int64      `json:"file_size,omitempty"`      // 可选的。文件大小
}

// Audio 由电报客户端视为音频的音频文件
// https://core.telegram.org/bots/api#audio
type Audio struct {
	FileID       string     `json:"file_id,omitempty"`        // 该文件的标识符，可用于下载或重复使用该文件
	FileUniqueID string     `json:"file_unique_id,omitempty"` // 此文件的唯一标识符，随着时间的推移，对于不同的漫游器，应该是相同的。不能用于下载或重复使用文件。
	Duration     int64      `json:"duration,omitempty"`       // 发件人定义的音频持续时间（以秒为单位）
	Performer    string     `json:"performer,omitempty"`      // 可选的。由发送者或音频标签定义的音频执行者
	Title        string     `json:"title,omitempty"`          // 可选的。由发送者或音频标签定义的音频标题
	FileName     string     `json:"file_name"`                // 可选的。发件人定义的原始文件名
	MimeType     string     `json:"mime_type,omitempty"`      // 可选的。发件人定义的文件的MIME类型
	FileSize     int64      `json:"file_size,omitempty"`      // 可选的。文件大小
	Thumb        *PhotoSize `json:"thumb,omitempty"`          // 可选的。音乐文件所属专辑封面的缩略图
}

// Document 常规文件（与照片，语音消息和音频文件相对）
// https://core.telegram.org/bots/api#document
type Document struct {
	FileID       string     `json:"file_id,omitempty"`        // 该文件的标识符，可用于下载或重复使用该文件
	FileUniqueID string     `json:"file_unique_id,omitempty"` // 此文件的唯一标识符，随着时间的推移，对于不同的漫游器，应该是相同的。不能用于下载或重复使用文件。
	Thumb        *PhotoSize `json:"thumb,omitempty"`          // 可选的。发件人定义的文档缩略图
	FileName     string     `json:"file_name,omitempty"`      // 可选的。发件人定义的原始文件名
	MimeType     string     `json:"mime_type,omitempty"`      // 可选的。发件人定义的文件的MIME类型
	FileSize     int64      `json:"file_size,omitempty"`      // 可选的。文件大小
}

// Video 视频文件
// https://core.telegram.org/bots/api#video
type Video struct {
	FileID       string     `json:"file_id,omitempty"`        // 该文件的标识符，可用于下载或重复使用该文件
	FileUniqueID string     `json:"file_unique_id,omitempty"` // 此文件的唯一标识符，随着时间的推移，对于不同的漫游器，应该是相同的。不能用于下载或重复使用文件。
	Width        int64      `json:"width,omitempty"`          // 发件人定义的视频宽度
	Height       int64      `json:"height,omitempty"`         // 发件人定义的视频高度
	Duration     int64      `json:"duration,omitempty"`       // 发件人定义的视频高度
	Thumb        *PhotoSize `json:"thumb,omitempty"`          // 可选的。影片缩图
	FileName     string     `json:"file_name"`                // 可选的。发件人定义的原始文件名
	MimeType     string     `json:"mime_type,omitempty"`      // 可选的。发件人定义的文件的MIME类型
	FileSize     int64      `json:"file_size,omitempty"`      // 可选的。文件大小
}

// VideoNote 视频消息（自v.4.0起在Telegram应用中可用）
// https://core.telegram.org/bots/api#videonote
type VideoNote struct {
	FileID       string     `json:"file_id,omitempty"`        // 该文件的标识符，可用于下载或重复使用该文件
	FileUniqueID string     `json:"file_unique_id,omitempty"` // 此文件的唯一标识符，随着时间的推移，对于不同的漫游器，应该是相同的。不能用于下载或重复使用文件。
	Length       int64      `json:"length,omitempty"`         // 发件人定义的视频宽度和高度（视频消息的直径）
	Duration     int64      `json:"duration,omitempty"`       // 发件人定义的视频时长（以秒为单位）
	Thumb        *PhotoSize `json:"thumb,omitempty"`          // 可选的。影片缩图
	FileSize     int64      `json:"file_size,omitempty"`      // 可选的。文件大小
}

// Voice 语音笔记
// https://core.telegram.org/bots/api#voice
type Voice struct {
	FileID       string `json:"file_id,omitempty"`        // 该文件的标识符，可用于下载或重复使用该文件
	FileUniqueID string `json:"file_unique_id,omitempty"` // 此文件的唯一标识符，随着时间的推移，对于不同的漫游器，应该是相同的。不能用于下载或重复使用文件。
	Duration     int64  `json:"duration,omitempty"`       // 发件人定义的音频持续时间（以秒为单位）
	MimeType     string `json:"mime_type,omitempty"`      // 可选的。发件人定义的文件的MIME类型
	FileSize     int64  `json:"file_size,omitempty"`      // 可选的。文件大小
}

// Contact 电话联系人
// https://core.telegram.org/bots/api#contact
type Contact struct {
	PhoneNumber string `json:"phone_number,omitempty"` // 联系人的电话号码
	FirstName   string `json:"first_name,omitempty"`   // 联系人的名字
	LastName    string `json:"last_name,omitempty"`    // 可选的。联系人的姓氏
	UserID      int64  `json:"user_id,omitempty"`      // 可选的。电报中联系人的用户标识符
	Vcard       string `json:"vcard,omitempty"`        // 可选的。以vCard形式的有关联系人的其他数据
}

// Dice 显示随机值的动画表情符号
// https://core.telegram.org/bots/api#dice
type Dice struct {
	Emoji string `json:"emoji,omitempty"` // 掷骰子动画所基于的表情符号
	Value int    `json:"value,omitempty"` // 骰子的值，“🎲”和“🎯”基本表情符号为1-6，“🏀”基本表情符号为1-5
}

// PollOption 民意测验中一个答案选项的信息
// https://core.telegram.org/bots/api#polloption
type PollOption struct {
	Text       string `json:"text,omitempty"`        // 选项文字，1-100个字符
	VoterCount int64  `json:"voter_count,omitempty"` // 对该选项投票的用户数
}

// PollAnswer 非匿名调查中的回答
// https://core.telegram.org/bots/api#poll_answer
type PollAnswer struct {
	PollID    string  `json:"poll_id,omitempty"`    // 唯一的投票标识符
	User      *User   `json:"user,omitempty"`       // 更改调查答案的用户
	OptionIds []int64 `json:"option_ids,omitempty"` // 用户选择的答案选项的基于0的标识符。如果用户撤回其投票，则可能为空。
}

const (
	// PollTypeAtQuiz 测验
	PollTypeAtQuiz = "quiz"
	// PollTypeAtRegular 常规
	PollTypeAtRegular = "regular"
)

// Poll 轮询的信息
// https://core.telegram.org/bots/api#poll
type Poll struct {
	ID                    string          `json:"id,omitempty"`                      // 唯一的投票标识符
	Question              string          `json:"question,omitempty"`                // 投票问题，1-255个字符
	Options               []PollOption    `json:"options,omitempty"`                 // 投票选项清单
	TotalVoterCount       int64           `json:"total_voter_count,omitempty"`       // 在民意调查中投票的用户总数
	IsClosed              bool            `json:"is_closed,omitempty"`               // True，如果民意调查已关闭
	IsAnonymous           bool            `json:"is_anonymous,omitempty"`            // True，如果民意调查是匿名的
	Type                  string          `json:"type,omitempty"`                    // 投票类型，当前可以是  “regular” 或 “quiz”
	AllowsMultipleAnswers bool            `json:"allows_multiple_answers,omitempty"` // True，如果民意测验允许多个答案
	CorrectOptionID       int64           `json:"correct_option_id,omitempty"`       // 可选的。正确答案选项的从0开始的标识符。仅适用于处于测验模式，已关闭或已由漫游器发送（不转发）或与漫游器进行私人聊天的测验模式。
	Explanation           string          `json:"explanation,omitempty"`             // 可选的。当用户选择不正确的答案或轻按测验样式的民意测验中的灯泡图标时显示的文本，0-200个字符
	ExplanationEntities   []MessageEntity `json:"explanation_entities,omitempty"`    // 可选的。解释中出现的特殊实体，例如用户名，URL，机器人命令等
	OpenPeriod            int64           `json:"open_period,omitempty"`             // 可选的。创建轮询后将激活活动的时间（以秒为单位）
	CloseDate             int64           `json:"close_date,omitempty"`              // 可选的。轮询将自动关闭的时间点（Unix时间戳记）
}

// Location 地图上的一个点
// https://core.telegram.org/bots/api#location
type Location struct {
	Longitude            float64 `json:"longitude,omitempty"`    // 经度
	Latitude             float64 `json:"latitude,omitempty"`     // 纬度
	HorizontalAccuracy   float64 `json:"horizontal_accuracy"`    // 可选的。位置的不确定性半径，以米为单位；0-1500
	LivePeriod           int64   `json:"live_period"`            // 可选的。相对于消息发送日期的时间，在此期间可以更新位置，以秒为单位。仅适用于活动的实时位置。
	Heading              int64   `json:"heading"`                // 可选的。用户移动的方向，以度为单位；1-360。仅适用于活动的实时位置。
	ProximityAlertRadius int64   `json:"proximity_alert_radius"` // 可选的。有关接近另一个聊天成员的接近警报的最大距离（以米为单位）。仅适用于已发送的实时位置。
}

// Venue 表示场地
// https://core.telegram.org/bots/api#venue
type Venue struct {
	Location        *Location `json:"location,omitempty"`        // 场地位置
	Title           string    `json:"title,omitempty"`           // 会场名称
	Address         string    `json:"address,omitempty"`         // 会场地址
	FoursquareID    string    `json:"foursquare_id,omitempty"`   // 可选的。场地的Foursquare标识符
	FoursquareType  string    `json:"foursquare_type,omitempty"` // 可选的。场地的Foursquare类型。（例如，“arts_entertainment/default”, “arts_entertainment/aquarium” or “food/icecream”。）
	GooglePlaceID   string    `json:"google_place_id"`           // 可选的。场地的Google地方信息标识符
	GooglePlaceType string    `json:"google_place_type"`         // 可选的。场所的Google地方信息类型。（请参阅支持的类型 https://developers.google.com/places/web-service/supported_types 。）
}

// ProximityAlertTriggered 该对象表示服务消息的内容，每当聊天中的某个用户触发另一个用户设置的接近警报时，就会发送该消息。
// https://core.telegram.org/bots/api#proximityalerttriggered
type ProximityAlertTriggered struct {
	Traveler *User `json:"traveler"` // 触发警报的用户
	Watcher  *User `json:"watcher"`  // 设置警报的用户
	Distance int64 `json:"distance"` // 用户之间的距离
}

// UserProfilePhotos 用户的个人资料图片
// https://core.telegram.org/bots/api#userprofilephotos
type UserProfilePhotos struct {
	TotalCount int64      `json:"total_count,omitempty"` // 目标用户拥有的个人资料图片总数
	Photos     [][]string `json:"photos,omitempty"`      // 个人资料图片（每张最多4个尺寸）
}

// File 准备下载的文件。可以通过链接下载文件 `https://api.telegram.org/file/bot<token>/<file_path>`。可以确保链接至少有效1个小时。当链接过期时，可以通过调用 GetFile 来请求一个新的链接。
// https://core.telegram.org/bots/api#file
type File struct {
	FileID       string `json:"file_id,omitempty"`        // 该文件的标识符，可用于下载或重复使用该文件
	FileUniqueID string `json:"file_unique_id,omitempty"` // 此文件的唯一标识符，随着时间的推移，对于不同的漫游器，应该是相同的。不能用于下载或重复使用文件。
	FileSize     int64  `json:"file_size,omitempty"`      // 可选的。文件大小（如果已知）
	FilePath     string `json:"file_path,omitempty"`      // 可选的。文件路径。使用 `https://api.telegram.org/file/bot<token>/<file_path>` 来获取文件。
}

// ReplyKeyboardMarkup 带有回复选项的自定义键盘（有关详细信息和示例，请参阅机器人简介）
// https://core.telegram.org/bots/api#replykeyboardmarkup
type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard,omitempty"`          // 按钮行数组，每个行由一个KeyboardButton对象数组表示
	ResizeKeyboard  bool               `json:"resize_keyboard,omitempty"`   // 可选的。请求客户垂直调整键盘大小以达到最佳配合（例如，如果只有两行按钮，则使键盘变小）。默认值为false，在这种情况下，自定义键盘的高度始终与应用程序的标准键盘相同。
	OneTimeKeyboard bool               `json:"one_time_keyboard,omitempty"` // 可选的。要求客户在使用键盘后立即隐藏它。键盘仍然可用，但是客户端将在聊天中自动显示常用的字母键盘-用户可以在输入字段中按特殊按钮以再次查看自定义键盘。默认为false。
	Selective       bool               `json:"selective,omitempty"`         // 可选的。如果只想向特定用户显示键盘，请使用此参数。目标：1）在Message对象的文本中@提及的用户；2）如果漫游器的消息是回复（具有reply_to_message_id），则为原始消息的发送者。示例：用户请求更改漫游器的语言，漫游器用键盘答复选择新语言的请求。群组中的其他用户看不到键盘。
}

// KeyboardButton 回复键盘的一个按钮。对于简单的文本按钮，可以使用String代替此对象来指定按钮的文本。可选字段request_contact，request_location和request_poll是互斥的。
// https://core.telegram.org/bots/api#keyboardbutton
type KeyboardButton struct {
	Text            string                 `json:"text,omitempty"`             // 按钮的文字。如果未使用任何可选字段，则在按下按钮时它将作为消息发送
	RequestContact  bool                   `json:"request_contact,omitempty"`  // 可选的。如果为True，则当按下按钮时，用户的电话号码将作为联系人发送。仅在私人聊天中可用
	RequestLocation bool                   `json:"request_location,omitempty"` // 可选的。如果为True，则在按下按钮时将发送用户的当前位置。仅在私人聊天中可用
	RequestPoll     KeyboardButtonPollType `json:"request_poll,omitempty"`     // 可选的。如果指定，则将要求用户创建一个民意调查，并在按下按钮时将其发送给机器人。仅在私人聊天中可用
}

// KeyboardButtonPollType 民意调查的类型，可以在按下相应按钮时创建和发送该民意调查。
// https://core.telegram.org/bots/api#keyboardbuttonpolltype
type KeyboardButtonPollType struct {
	Type string `json:"type,omitempty"` // 可选的。如果通过测验，将仅允许用户以测验模式创建民意测验。如果通过常规，则仅允许常规轮询。否则，将允许用户创建任何类型的民意测验。
}

// ReplyKeyboardRemove // 收到带有该对象的消息后，Telegram客户端将删除当前的自定义键盘并显示默认的字母键盘。默认情况下，显示自定义键盘，直到机器人发送新键盘为止。一次性键盘的例外情况是用户按下按钮后立即将其隐藏（请参见ReplyKeyboardMarkup）。
// https://core.telegram.org/bots/api#replykeyboardremove
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard,omitempty"` // 请求客户端删除自定义键盘（用户将无法召唤此键盘；如果您希望隐藏键盘，但仍可访问，请在ReplyKeyboardMarkup中使用one_time_keyboard）
	Selective      bool `json:"selective,omitempty"`       // 可选的。如果仅要为特定用户卸下键盘，请使用此参数。目标：1）在Message对象的文本中@提及的用户；2）如果漫游器的消息是回复（具有reply_to_message_id），则为原始消息的发送者。示例：用户在投票中投票，机器人返回确认消息以回应投票，并删除该用户的键盘，同时仍向尚未投票的用户显示带有投票选项的键盘。
}

// InlineKeyboardMarkup 嵌入式键盘，出现在其所属消息的旁边
// https://core.telegram.org/bots/api#inlinekeyboardmarkup
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard,omitempty"` // 按钮行数组，每个行由一个InlineKeyboardButton对象数组表示
}

// InlineKeyboardButton 嵌入式键盘的一个按钮。您必须恰好使用可选字段之一
// https://core.telegram.org/bots/api#inlinekeyboardbutton
type InlineKeyboardButton struct {
	Text                         string        `json:"text,omitempty"`                             // 在按钮上标记文本
	URL                          string        `json:"url,omitempty"`                              // 可选的。按下按钮时将打开HTTP或 tg:// url
	LoginURL                     *LoginURL     `json:"login_url,omitempty"`                        // 可选的。用于自动授权用户的HTTP URL。可以替代电报登录小部件。
	CallbackData                 string        `json:"callback_data,omitempty"`                    // 可选的。按下按钮时要在回调查询中发送到bot的数据，1-64个字节
	SwitchInlineQuery            string        `json:"switch_inline_query,omitempty"`              // 可选的。如果已设置，则按下按钮将提示用户选择其聊天之一，打开该聊天并将bot的用户名和指定的内联查询插入输入字段。可以为空，在这种情况下，只会插入机器人的用户名。
	SwitchInlineQueryCurrentChat string        `json:"switch_inline_query_current_chat,omitempty"` // 可选的。如果已设置，则按下按钮会将bot的用户名和指定的嵌入式查询插入当前聊天的输入字段中。可以为空，在这种情况下，只会插入机器人的用户名。
	CallbackGame                 *CallbackGame `json:"callback_game,omitempty"`                    // 可选的。用户按下按钮时将启动的游戏的描述。
	Pay                          bool          `json:"pay,omitempty"`                              // 可选的。指定True，发送“付款”按钮。
}

// LoginURL 用于自动授权用户的嵌入式键盘按钮的参数。当用户来自Telegram时，可以用作Telegram Login Widget的替代品。用户所需要做的就是点击/单击按钮并确认他们要登录
// https://core.telegram.org/bots/api#loginurl
type LoginURL struct {
	URL                string `json:"url,omitempty"`                  // 按下按钮时，将打开一个HTTP URL，并将用户授权数据添加到查询字符串中。如果用户拒绝提供授权数据，则将打开不包含有关用户信息的原始URL。添加的数据与接收授权数据中所述的相同。
	ForwardText        string `json:"forward_text,omitempty"`         // 可选的。转发邮件中按钮的新文本。
	BotUsername        string `json:"bot_username,omitempty"`         // 可选的。机器人的用户名，将用于用户授权。有关更多详细信息，请参见设置机器人。如果未指定，将使用当前机器人的用户名。该网址的域名必须是相同的与机器人相关联的域。有关更多详细信息，请参见将您的域链接到机器人。
	RequestWriteAccess bool   `json:"request_write_access,omitempty"` // 可选的。传递True以请求您的漫游器向用户发送消息的权限。
}

// CallbackQuery 从嵌入式键盘中的回调按钮传入的回调查询。如果发起查询的按钮已附加到机器人发送的消息中，则将显示字段消息。如果该按钮已附加到通过漫游器发送的邮件（以内联模式），则将显示inline_message_id字段。会出现字段data或game_short_name之一。
// https://core.telegram.org/bots/api#callbackquery
type CallbackQuery struct {
	ID              string   `json:"id,omitempty"`                // 此查询的唯一标识符
	From            *User    `json:"from,omitempty"`              // 发件人
	Message         *Message `json:"message,omitempty"`           // 可选的。带有发起查询的回调按钮的消息。请注意，如果消息太旧，则消息内容和消息日期将不可用
	InlineMessageID string   `json:"inline_message_id,omitempty"` // 可选的。通过机器人以串联模式发送的消息的标识符，该消息是发起查询的
	ChatInstance    string   `json:"chat_instance,omitempty"`     // 全局标识符，唯一地与带有回调按钮的消息发送到的聊天相对应。对于游戏中的高分有用。
	Data            string   `json:"data,omitempty"`              // 可选的。与回调按钮关联的数据。请注意，错误的客户端可以在此字段中发送任意数据。
	GameShortName   string   `json:"game_short_name,omitempty"`   // 可选的。要返回的游戏的简称，用作游戏的唯一标识符
}

// ForceReply // 收到带有该对象的消息后，Telegram客户端将向用户显示一个答复界面（就像用户选择了漫游器的消息并点按“答复”一样）。如果您想创建用户友好的逐步界面而不必牺牲隐私模式，这将非常有用。
// https://core.telegram.org/bots/api#forcereply
type ForceReply struct {
	ForceReply bool `json:"force_reply,omitempty"` // 向用户显示回复界面，就像他们手动选择了机器人的消息并点按“回复”一样
	Selective  bool `json:"selective,omitempty"`   // 可选的。如果只想强制特定用户答复，请使用此参数。目标：1）在Message对象的文本中@提及的用户；2）如果漫游器的消息是回复（具有reply_to_message_id），则为原始消息的发送者。
}

// ChatPhoto 聊天照片
// https://core.telegram.org/bots/api#chatphoto
type ChatPhoto struct {
	SmallFileID       string `json:"small_file_id,omitempty"`        // 小型（160x160）聊天照片的文件标识符。该file_id仅可用于照片下载，并且仅在不改变照片的情况下使用。
	SmallFileUniqueID string `json:"small_file_unique_id,omitempty"` // 小型（160x160）聊天照片的唯一文件标识符，随着时间的推移，对于不同的漫游器，应该是相同的。不能用于下载或重复使用文件。
	BigFileID         string `json:"big_file_id,omitempty"`          // 大（640x640）聊天照片的文件标识符。该file_id仅可用于照片下载，并且仅在不改变照片的情况下使用。
	BigFileUniqueID   string `json:"big_file_unique_id,omitempty"`   // 大（640x640）大聊天照片的唯一文件标识符，随着时间的推移，对于不同的漫游器，该标识符应该是相同的。不能用于下载或重复使用文件。
}

const (
	// ChatMemberAtCreator 创建者
	ChatMemberAtCreator = "creator"
	// ChatMemberAtAdministrator 管理员
	ChatMemberAtAdministrator = "administrator”"
	// ChatMemberAtMember 成员
	ChatMemberAtMember = "member"
	// ChatMemberAtRestricted 受限制
	ChatMemberAtRestricted = "restricted"
	// ChatMemberAtLeft 左
	ChatMemberAtLeft = "left"
	// ChatMemberAtKicked 被踢
	ChatMemberAtKicked = "kicked"
)

// ChatMember 聊天成员的信息
// https://core.telegram.org/bots/api#chatmember
type ChatMember struct {
	User                  *User  `json:"user,omitempty"`                      // 有关用户的信息
	Status                string `json:"status,omitempty"`                    // 成员在聊天中的状态。可以是 “creator”, “administrator”, “member”, “restricted”, “left” or “kicked”
	CustomTitle           string `json:"custom_title,omitempty"`              // 可选的。仅所有者和管理员。该用户的自定义标题
	IsAnonymous           bool   `json:"is_anonymous"`                        // 可选的。仅所有者和管理员。可以，如果隐藏了用户在聊天中的状态
	CanBeEdited           bool   `json:"can_be_edited,omitempty"`             // 可选的。仅管理员。是的，如果允许漫游器编辑该用户的管理员权限
	CanPostMessages       bool   `json:"can_post_messages,omitempty"`         // 可选的。仅管理员。是的，如果管理员可以在频道中发布；仅频道
	CanEditMessages       bool   `json:"can_edit_messages,omitempty"`         // 可选的。仅管理员。是的，如果管理员可以编辑其他用户的消息并可以固定消息；仅频道
	CanDeleteMessages     bool   `json:"can_delete_messages,omitempty"`       // 可选的。仅管理员。是的，如果管理员可以删除其他用户的邮件
	CanRestrictMembers    bool   `json:"can_restrict_members,omitempty"`      // 可选的。仅管理员。 是的，如果管理员可以限制，禁止或取消禁止聊天成员
	CanPromoteMembers     bool   `json:"can_promote_members,omitempty"`       // 可选的。仅管理员。的确，如果管理员可以添加具有自己特权子集的新管理员，或者将直接或间接晋升的管理员降级（由用户任命的管理员晋升）
	CanChangeInfo         bool   `json:"can_change_info,omitempty"`           // 选的。仅限管理员和受限。是的，如果允许用户更改聊天标题，照片和其他设置
	CanInviteUsers        bool   `json:"can_invite_users,omitempty"`          // 可选的。仅限管理员和受限。是的，如果允许用户邀请新用户加入聊天
	CanPinMessages        bool   `json:"can_pin_messages,omitempty"`          // 可选的。仅限管理员和受限。是的，如果允许用户固定消息；仅组和超组
	IsMember              bool   `json:"is_member,omitempty"`                 // 可选的。仅受限制。是的，如果用户是请求时聊天的成员
	CanSendMessages       bool   `json:"can_send_messages,omitempty"`         // 可选的。仅受限制。是的，如果允许用户发送短信，联系人，位置和地点
	CanSendMediaMessages  bool   `json:"can_send_media_messages,omitempty"`   // 可选的。仅受限制。是的，如果允许用户发送音频，文档，照片，视频，视频注释和语音注释
	CanSendPolls          bool   `json:"can_send_polls,omitempty"`            // 可选的。仅受限制。是的，如果允许用户发送民意调查
	CanSendOtherMessages  bool   `json:"can_send_other_messages,omitempty"`   // 可选的。仅受限制。是的，如果允许用户发送动画，游戏，贴纸并使用嵌入式机器人
	CanAddWebPagePreviews bool   `json:"can_add_web_page_previews,omitempty"` // 可选的。仅受限制。是的，如果允许用户将网页预览添加到他们的消息中
	UntilDate             int64  `json:"until_date,omitempty"`                // 可选的。限制和踢。对该用户取消限制的日期；Unix时间
}

// ChatPermissions 允许非管理员用户进行聊天的操作
// https://core.telegram.org/bots/api#chatpermissions
type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages,omitempty"`         // 可选的。True，如果允许用户发送短信，联系人，位置和地点
	CanSendMediaMessages  bool `json:"can_send_media_messages,omitempty"`   // 可选的。True，如果允许用户发送音频，文档，照片，视频，视频注释和语音注释，则意味着 can_send_messages
	CanSendPolls          bool `json:"can_send_polls,omitempty"`            // 可选的。True，如果允许用户发送民意调查，则意味着 can_send_messages
	CanSendOtherMessages  bool `json:"can_send_other_messages,omitempty"`   // 可选的。True，如果允许用户发送动画，游戏，贴纸和使用嵌入式机器人，则意味着 can_send_media_messages
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews,omitempty"` // 可选的。True，如果允许用户向其消息添加网页预览，则意味着 can_send_media_messages
	CanChangeInfo         bool `json:"can_change_info,omitempty"`           // 可选的。True，如果允许用户更改聊天标题，照片和其他设置。在公共超级组中被忽略
	CanInviteUsers        bool `json:"can_invite_users,omitempty"`          // 可选的。True，如果允许用户邀请新用户加入聊天
	CanPinMessages        bool `json:"can_pin_messages,omitempty"`          // 可选的。True，如果允许用户固定消息。在公共超级组中被忽略
}

// ChatLocation 代表聊天连接的位置
// https://core.telegram.org/bots/api#chatlocation
type ChatLocation struct {
	Location *Location `json:"location"` // 超组连接到的位置。不能是居住地点。
	Address  string    `json:"address"`  // 位置地址；1-64个字符，由聊天所有者定义
}

// BotCommand 机器人命令
// https://core.telegram.org/bots/api#botcommand
type BotCommand struct {
	Command     string `json:"command,omitempty"`     // 命令文本，1-32个字符。只能包含小写英文字母，数字和下划线。
	Description string `json:"description,omitempty"` // 命令说明，3-256个字符。
}

// InputMedia 要发送的媒体消息的内容（注：应该是下面几种类型中的一种）
// https://core.telegram.org/bots/api#inputmedia
type InputMedia struct {
	*InputMediaAnimation
	*InputMediaDocument
	*InputMediaAudio
	*InputMediaPhoto
	*InputMediaVideo
}

// InputMediaPhotoType 照片类型
const InputMediaPhotoType = "photo"

// InputMediaPhoto 要发送的照片
// https://core.telegram.org/bots/api#inputmediaphoto
type InputMediaPhoto struct {
	Type            string          `json:"type,omitempty"`       // 结果类型，必须是 photo
	Media           string          `json:"media,omitempty"`      // 文件发送。传递file_id以发送电报服务器上存在的文件（推荐），传递电报的HTTP URL以从Internet获取文件，或传递“attach://<file_attach_name>”以使用multipart/<file_attach_name>名称下的form-data。
	Caption         string          `json:"caption,omitempty"`    // 可选的。要发送的照片的标题，实体解析后0-1024个字符
	ParseMode       string          `json:"parse_mode,omitempty"` // 可选的。解析照片标题中的实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities []MessageEntity `json:"caption_entities"`     // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
}

// InputMediaVideoType 视频类型
const InputMediaVideoType = "video"

// InputMediaVideo 要发送的视频
// https://core.telegram.org/bots/api#inputmediavideo
type InputMediaVideo struct {
	Type              string          `json:"type,omitempty"`               // 结果类型，必须是 video
	Media             string          `json:"media,omitempty"`              // 文件发送。传递file_id以发送电报服务器上存在的文件（推荐），传递电报的HTTP URL以从Internet获取文件，或传递“ attach：// <file_attach_name>”以使用multipart / <file_attach_name>名称下的form-data。有关发送文件的更多信息»
	Thumb             InputFile       `json:"thumb,omitempty"`              // 可选的。已发送文件的缩略图；如果在服务器端支持为文件生成缩略图，则可以忽略。缩略图应为JPEG格式，并且大小应小于200 kB。缩略图的宽度和高度不应超过320。如果未使用multipart / form-data上传文件，则忽略该缩略图。缩略图不能重复使用，只能作为新文件上传，因此如果缩略图是使用<file_attach_name>下的multipart / form-data上传的，则可以传递“ attach：// <file_attach_name>”。
	Caption           string          `json:"caption,omitempty"`            // 可选的。要发送的视频的标题，实体解析后0-1024个字符
	ParseMode         string          `json:"parse_mode,omitempty"`         // 可选的。视频字幕中的实体解析模式。有关更多详细信息，请参见格式化选项
	CaptionEntities   []MessageEntity `json:"caption_entities"`             // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	Width             int64           `json:"width,omitempty"`              // 可选的。影片宽度
	Height            int64           `json:"height,omitempty"`             // 可选的。影片高度
	Duration          int64           `json:"duration,omitempty"`           // 可选的。影片时长
	SupportsStreaming bool            `json:"supports_streaming,omitempty"` // 可选的。如果上传的视频适合流式传输，则通过True
}

// InputMediaAnimationType 动画文件 类型
const InputMediaAnimationType = "animation"

// InputMediaAnimation 要发送的动画文件（GIF或H.264 / MPEG-4 AVC视频，无声音）。
// https://core.telegram.org/bots/api#inputmediaanimation
type InputMediaAnimation struct {
	Type            string          `json:"type,omitempty"`       // 结果类型，必须是 animation
	Media           string          `json:"media,omitempty"`      // 文件发送。传递file_id以发送电报服务器上存在的文件（推荐），传递电报的HTTP URL以从Internet获取文件，或传递“ attach：// <file_attach_name>”以使用multipart / <file_attach_name>名称下的form-data。
	Thumb           InputFile       `json:"thumb,omitempty"`      // 可选的。已发送文件的缩略图；如果在服务器端支持为文件生成缩略图，则可以忽略。缩略图应为JPEG格式，并且大小应小于200 kB。缩略图的宽度和高度不应超过320。如果未使用multipart / form-data上传文件，则忽略该缩略图。缩略图不能重复使用，只能作为新文件上传，因此如果缩略图是使用<file_attach_name>下的multipart / form-data上传的，则可以传递“ attach：// <file_attach_name>”
	Caption         string          `json:"caption,omitempty"`    // 可选的。要发送的动画的标题，实体解析后为0-1024个字符
	ParseMode       string          `json:"parse_mode,omitempty"` // 可选的。解析动画标题中实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities []MessageEntity `json:"caption_entities"`     // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	Width           int64           `json:"width,omitempty"`      // 可选的。动画宽度
	Height          int64           `json:"height,omitempty"`     // 可选的。动画高度
	Duration        int64           `json:"duration,omitempty"`   // 可选的。动画时长
}

// InputMediaAudioType 音乐类型
const InputMediaAudioType = "audio"

// InputMediaAudio 发送的音乐
// https://core.telegram.org/bots/api#inputmediaaudio
type InputMediaAudio struct {
	Type            string          `json:"type,omitempty"`       // 结果类型，必须为 audio
	Media           string          `json:"media,omitempty"`      // 文件发送。传递file_id以发送电报服务器上存在的文件（推荐），传递电报的HTTP URL以从Internet获取文件，或传递“ attach：// <file_attach_name>”以使用multipart / <file_attach_name>名称下的form-data
	Thumb           InputFile       `json:"thumb,omitempty"`      // 可选的。已发送文件的缩略图；如果在服务器端支持为文件生成缩略图，则可以忽略。缩略图应为JPEG格式，并且大小应小于200 kB。缩略图的宽度和高度不应超过320。如果未使用multipart / form-data上传文件，则忽略该缩略图。缩略图不能重复使用，只能作为新文件上传，因此如果缩略图是使用<file_attach_name>下的multipart / form-data上传的，则可以传递“ attach：// <file_attach_name>”。
	Caption         string          `json:"caption,omitempty"`    // 可选的。要发送的音频的标题，实体解析后为0-1024个字符
	ParseMode       string          `json:"parse_mode,omitempty"` // 可选的。解析音频字幕中实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities []MessageEntity `json:"caption_entities"`     // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	Duration        int64           `json:"duration,omitempty"`   // 可选的。音频持续时间（以秒为单位）
	Performer       string          `json:"performer,omitempty"`  // 可选的。音频表演者
	Title           string          `json:"title,omitempty"`      // 可选的。音频标题
}

// InputMediaDocument 要发送的常规文件
// https://core.telegram.org/bots/api#inputmediadocument
type InputMediaDocument struct {
	Type                        string          `json:"type,omitempty"`                 // 结果类型，必须为文件
	Media                       string          `json:"media,omitempty"`                // 文件发送。传递file_id以发送电报服务器上存在的文件（推荐），传递电报的HTTP URL以从Internet获取文件，或传递“ attach：// <file_attach_name>”以使用multipart / <file_attach_name>名称下的form-data。
	Thumb                       InputFile       `json:"thumb,omitempty"`                // 可选的。已发送文件的缩略图；如果在服务器端支持为文件生成缩略图，则可以忽略。缩略图应为JPEG格式，并且大小应小于200 kB。缩略图的宽度和高度不应超过320。如果未使用multipart / form-data上传文件，则忽略该缩略图。缩略图不能重复使用，只能作为新文件上传，因此如果缩略图是使用<file_attach_name>下的multipart / form-data上传的，则可以传递“ attach：// <file_attach_name>”。
	Caption                     string          `json:"caption,omitempty"`              // 可选的。待发送文档的标题，实体解析后0-1024个字符
	ParseMode                   string          `json:"parse_mode,omitempty"`           // 可选的。解析文档标题中的实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities             []MessageEntity `json:"caption_entities"`               // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	DisableContentTypeDetection bool            `json:"disable_content_type_detection"` // 可选的。对使用multipart / form-data上传的文件禁用服务器端内容类型自动检测。如果文档是作为相册的一部分发送的，则始终为true。
}

// InputFile 上传文件的内容。必须使用 multipart/form-data 以通过浏览器上传文件的通常方式进行发布。
// 库作者特别注释: InputFile 就是 io.Reader 建议在 Golang 直接使用 io.Reader
// https://core.telegram.org/bots/api#inputfile
type InputFile = io.Reader

// Sending files
// https://core.telegram.org/bots/api#sending-files
/*有三种发送文件的方法（照片，贴纸，音频，媒体等）：

	1. 如果文件已经存储在Telegram服务器上的某个位置，则无需重新上传它：每个文件对象都有一个file_id字段，只需将此file_id作为参数传递而不是上传。有没有限制为发送这样的文件。
	2. 为Telegram提供要发送文件的HTTP URL。电报将下载并发送文件。照片的最大大小为5 MB，其他类型的内容的最大大小为20 MB。
	3. 使用multipart/form-data发布文件的方式与通过浏览器上传文件的通常方式相同。照片最大大小为10 MB，其他文件最大为50 MB。
	通过file_id发送

通过file_id重新发送时，无法更改文件类型。即视频不能作为照片发送，照片不能作为文档发送，等等。
	无法重新发送缩略图。
	通过file_id重新发送照片将发送其所有尺寸。
	file_id对于每个单独的漫游器都是唯一的，并且不能从一个漫游器转移到另一个漫游器。
	file_id唯一地标识一个文件，但是即使对于同一漫游器，文件也可以具有不同的有效file_id。
	通过URL发送

通过URL发送时，目标文件必须具有正确的MIME类型（例如，sendAudio的音频/ mpeg 等）。
	在sendDocument中，按URL发送当前仅适用于gif，pdf和zip文件。
	要使用sendVoice，文件必须具有audio / ogg类型，并且大小不得超过1MB。1-20MB语音便笺将作为文件发送。
	其他配置可能会起作用，但我们不能保证一定会。*/
