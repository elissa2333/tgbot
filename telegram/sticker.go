package telegram

// pack

/*Stickers
机器人可以通过以下方法和对象来处理贴纸和贴纸集。
https://core.telegram.org/bots/api#stickers*/

// Sticker 贴纸和贴纸集
// https://core.telegram.org/bots/api#sticker
type Sticker struct {
	FileID       string        `json:"file_id,omitempty"`        // 该文件的标识符，可用于下载或重复使用该文件
	FileUniqueID string        `json:"file_unique_id,omitempty"` // 此文件的唯一标识符，随着时间的推移，对于不同的漫游器，应该是相同的。不能用于下载或重复使用文件。
	Width        int64         `json:"width,omitempty"`          // 贴纸宽度
	Height       int64         `json:"height,omitempty"`         // 贴纸高度
	IsAnimated   bool          `json:"is_animated,omitempty"`    // 如果贴纸是动画的，则为True
	Thumb        *PhotoSize    `json:"thumb,omitempty"`          // 可选的。.WEBP 或 .JPG 格式的贴纸缩略图
	Emoji        string        `json:"emoji,omitempty"`          // 可选的。与贴纸相关的表情符号
	SetName      string        `json:"set_name,omitempty"`       // 可选的。贴纸所属的贴纸集的名称
	MaskPosition *MaskPosition `json:"mask_position,omitempty"`  // 可选的。对于遮罩贴纸，应放置遮罩的位置
	FileSize     int64         `json:"file_size,omitempty"`      // 可选的。文件大小
}

// StickerSet 该对象代表一个贴纸集。
// https://core.telegram.org/bots/api#stickerset
type StickerSet struct {
	Name          string     `json:"name,omitempty"`           // 贴纸集名称
	Title         string     `json:"title,omitempty"`          // 贴纸集标题
	IsAnimated    bool       `json:"is_animated,omitempty"`    // 如果贴纸集包含动画贴纸，则为True
	ContainsMasks bool       `json:"contains_masks,omitempty"` // 如果贴纸组包含遮罩，则为True
	Stickers      []Sticker  `json:",omitemptystickers"`       // 全部套装一览
	Thumb         *PhotoSize `json:"thumb,omitempty"`          // 可选的。.WEBP或.TGS格式的贴纸集缩略图
}

// MaskPosition 默认情况下应放置遮罩的面上的位置
// https://core.telegram.org/bots/api#maskposition
type MaskPosition struct {
	Point  string  `json:"point,omitempty"`   // 面罩应相对放置的面部部分。“额头”，“眼睛”，“嘴”或“下巴”之一。
	XShift float64 `json:"x_shift,omitempty"` // X轴移动量（从遮罩的宽度开始缩放），缩放比例与面部尺寸成比例，从左到右。例如，选择-1.0将把遮罩放置在默认遮罩位置的左侧。
	YShift float64 `json:"y_shift,omitempty"` // 在蒙版的高度中测量的Y轴移位从顶部到底部，缩放比例与面部大小成比例。例如，1.0将遮罩放置在默认遮罩位置的正下方。
	Scale  float64 `json:"scale,omitempty"`   // 蒙版缩放系数。例如，2.0表示双倍大小。
}

// SendStickerOptional SendSticker 可选参数
type SendStickerOptional struct {
	DisableNotification      bool        `json:"disable_notification,omitempty"`        // 静默发送消息。用户将收到没有声音的通知。
	ReplyToMessageID         int64       `json:"reply_to_message_id,omitempty"`         // 如果消息是答复，则为原始消息的ID
	AllowSendingWithoutReply bool        `json:"allow_sending_without_reply,omitempty"` // 如果未发送指定的回复消息也应发送消息，则传递True
	ReplyMarkup              interface{} `json:"reply_markup,omitempty"`                // 其他界面选项。内联键盘，自定义回复键盘，删除回复键盘或强制用户回复的说明的JSON序列化对象。
}

// SendSticker 使用此方法发送静态.WEBP或动画的.TGS贴纸。成功后，将返回发送的 Message。
// https://core.telegram.org/bots/api#sendsticker
func (a API) SendSticker(chatID string, sticker InputFile, optional *SendStickerOptional) (*Message, error) {
	result := &Message{}
	err := a.handleOptional("/sendSticker", map[string]interface{}{"chat_id": chatID, "sticker": sticker}, optional, result)
	return result, err
}

// GetStickerSet 使用此方法获取贴纸集。成功后，将返回 StickerSet 对象。
// https://core.telegram.org/bots/api#getstickerset
func (a API) GetStickerSet(name string) (*StickerSet, error) {
	result := &StickerSet{}
	err := a.handleOptional("/getStickerSet", map[string]interface{}{"name": name}, nil, result)
	return result, err
}

// UploadStickerFile 使用此方法可以上传带有标签的.PNG文件，以供以后在createNewStickerSet和addStickerToSet方法中使用（可以多次使用）。成功返回上载的 File。
// https://core.telegram.org/bots/api#uploadstickerfile
func (a API) UploadStickerFile(userID int64, pngSticker InputFile) (*File, error) {
	result := &File{}
	err := a.handleOptional("/uploadStickerFile", map[string]interface{}{"user_id": userID, "png_sticker": pngSticker}, nil, result)
	return result, err

}

// CreateNewStickerSetOptional CreateNewStickerSet 可选参数
type CreateNewStickerSetOptional struct {
	PngSticker    InputFile     `json:"png_sticker,omitempty"`    // 带标签的PNG图片，最大不能超过512 KB，尺寸不能超过512px，宽度或高度必须恰好是512px。传递file_id作为字符串以发送Telegram服务器上已经存在的文件，传递HTTP URL作为Telegram以字符串形式从Internet获取文件，或者使用multipart / form-data上载新文件。有关发送文件的更多信息»
	TgsSticker    InputFile     `json:"tgs_sticker,omitempty"`    // 带标签的TGS动画，使用多部分/表单数据上传。看到https://core.telegram.org/animated_stickers#technical-技术要求
	ContainsMasks bool          `json:"contains_masks,omitempty"` // 如果需要创建一组遮罩贴纸，则传递True
	MaskPosition  *MaskPosition `json:"mask_position,omitempty"`  // JSON序列化的对象，用于将遮罩放置在脸上的位置
}

// CreateNewStickerSet 使用此方法可以创建用户拥有的新贴纸集。机器人将能够编辑由此创建的贴纸集。您必须完全使用png_sticker或tgs_sticker字段之一。成功返回True。
// https://core.telegram.org/bots/api#createnewstickerset
func (a API) CreateNewStickerSet(userID int64, name string, title string, emojis string, optional *CreateNewStickerSetOptional) (bool, error) {
	var result bool
	err := a.handleOptional("/createNewStickerSet", map[string]interface{}{"user_id": userID, "name": name, "title": title}, optional, &result)
	return result, err
}

// AddStickerToSetOptional AddStickerToSet 可选参数
type AddStickerToSetOptional struct {
	PngSticker   InputFile     `json:"png_sticker,omitempty"`   // 带标签的PNG图片，最大不能超过512 KB，尺寸不能超过512px，宽度或高度必须恰好是512px。传递file_id作为字符串以发送Telegram服务器上已经存在的文件，传递HTTP URL作为Telegram以字符串形式从Internet获取文件，或者使用multipart / form-data上载新文件。有关发送文件的更多信息»
	TgsSticker   InputFile     `json:"tgs_sticker,omitempty"`   // 带标签的TGS动画，使用多部分/表单数据上传。看到https://core.telegram.org/animated_stickers#technical-技术要求
	MaskPosition *MaskPosition `json:"mask_position,omitempty"` // JSON序列化的对象，用于将遮罩放置在脸上的位置
}

// AddStickerToSet 使用此方法可将新标签添加到由机器人创建的集合中。您必须完全使用png_sticker或tgs_sticker字段之一。可以将动画贴纸添加到动画贴纸集中，并且只能添加到它们。动画贴纸集最多可以包含50个贴纸。静态贴纸集最多可包含120个贴纸。成功返回
// https://core.telegram.org/bots/api#addstickertoset
func (a API) AddStickerToSet(userID string, name string, emojis string, optional *AddStickerToSetOptional) (bool, error) {
	var result bool
	err := a.handleOptional("/addStickerToSet", map[string]interface{}{"user_id": userID, "name": name, "emojis": emojis}, optional, result)
	return result, err
}

// SetStickerPositionInSet 使用此方法将机器人创建的集合中的贴纸移动到特定位置。成功返回True。
// https://core.telegram.org/bots/api#setstickerpositioninset
func (a API) SetStickerPositionInSet(sticker string, position int64) (bool, error) {
	var result bool
	err := a.handleOptional("/setStickerPositionInSet", map[string]interface{}{"sticker": sticker, "position": position}, nil, &result)
	return result, err
}

// DeleteStickerFromSet 使用此方法从机器人创建的集合中删除标签。成功返回True。
// https://core.telegram.org/bots/api#deletestickerfromset
func (a API) DeleteStickerFromSet(sticker string) (bool, error) {
	var result bool
	err := a.handleOptional("/deleteStickerFromSet", map[string]interface{}{"sticker": sticker}, nil, result)
	return result, err
}

// SetStickerSetThumbOptional SetStickerSetThumb 可选参数
type SetStickerSetThumbOptional struct {
	Thumb InputFile `json:"thumb,omitempty"` // PNG与缩略图图像，必须大于128kb的大小，并且具有的宽度和高度准确100像素，或TGS动画与缩略图大小高达32千字节; 看到有关动画贴纸技术要求的https://core.telegram.org/animated_stickers#technical-requirements。传递file_id作为字符串以发送Telegram服务器上已经存在的文件，传递HTTP URL作为Telegram以字符串形式从Internet获取文件，或者使用multipart / form-data上载新文件。有关发送文件的详细信息»。动画贴纸集缩略图无法通过HTTP URL上传。
}

// SetStickerSetThumb 使用此方法设置贴纸集的缩略图。只能为动画贴纸集设置动画缩略图。成功返回True。
// https://core.telegram.org/bots/api#setstickersetthumb
func (a API) SetStickerSetThumb(name string, userID int64, optional *SetStickerSetThumbOptional) (bool, error) {
	var result bool
	err := a.handleOptional("/setStickerSetThumb", map[string]interface{}{"name": name, "user_id": userID}, optional, &result)
	return result, err
}
