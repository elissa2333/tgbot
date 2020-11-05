package telegram

// pack

/*Stickers
机器人可以通过以下方法和对象来处理贴纸和贴纸集。
https://core.telegram.org/bots/api#stickers*/

// Sticker 贴纸和贴纸集
// https://core.telegram.org/bots/api#sticker
type Sticker struct {
	FileID       string        `json:"file_id"`        // 该文件的标识符，可用于下载或重复使用该文件
	FileUniqueID string        `json:"file_unique_id"` // 此文件的唯一标识符，随着时间的推移，对于不同的漫游器，应该是相同的。不能用于下载或重复使用文件。
	Width        int64         `json:"width"`          // 贴纸宽度
	Height       int64         `json:"height"`         // 贴纸高度
	IsAnimated   bool          `json:"is_animated"`    // 如果贴纸是动画的，则为True
	Thumb        *PhotoSize    `json:"thumb"`          // 可选的。.WEBP 或 .JPG 格式的贴纸缩略图
	Emoji        string        `json:"emoji"`          // 可选的。与贴纸相关的表情符号
	SetName      string        `json:"set_name"`       // 可选的。贴纸所属的贴纸集的名称
	MaskPosition *MaskPosition `json:"mask_position"`  // 可选的。对于遮罩贴纸，应放置遮罩的位置
	FileSize     int64         `json:"file_size"`      // 可选的。文件大小
}

// MaskPosition 默认情况下应放置遮罩的面上的位置
// https://core.telegram.org/bots/api#maskposition
type MaskPosition struct {
	Point  string  `json:"point"`   // 面罩应相对放置的面部部分。“额头”，“眼睛”，“嘴”或“下巴”之一。
	XShift float64 `json:"x_shift"` // X轴移动量（从遮罩的宽度开始缩放），缩放比例与面部尺寸成比例，从左到右。例如，选择-1.0将把遮罩放置在默认遮罩位置的左侧。
	YShift float64 `json:"y_shift"` // 在蒙版的高度中测量的Y轴移位从顶部到底部，缩放比例与面部大小成比例。例如，1.0将遮罩放置在默认遮罩位置的正下方。
	Scale  float64 `json:"scale"`   // 蒙版缩放系数。例如，2.0表示双倍大小。
}
