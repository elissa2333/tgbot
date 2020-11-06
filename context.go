package tgbot

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/elissa2333/tgbot/telegram"
	"github.com/elissa2333/tgbot/utils"
)

const (
	// ContextTypeAtText 文本
	ContextTypeAtText = "text"
	// ContextTypeAtPhoto 照片
	ContextTypeAtPhoto = "photo"
	// ContextTypeAtVoice 语音
	ContextTypeAtVoice = "voice"
	// ContextTypeAtAudio 音频
	ContextTypeAtAudio = "audio"
	// ContextTypeAtVideo 视频
	ContextTypeAtVideo = "video"
	// ContextTypeAtAnimation 动画
	ContextTypeAtAnimation = "animation"
	// ContextTypeAtDocument 文件
	ContextTypeAtDocument = "document"
	// ContextTypeAtSticker 贴纸
	ContextTypeAtSticker = "sticker"
	// ContextTypeAtVideoNote 视频笔记
	ContextTypeAtVideoNote = "videoNote"
	// ContextTypeAtContact 联系人
	ContextTypeAtContact = "contact"
	// ContextTypeAtDice 色子
	ContextTypeAtDice = "dice"
	// ContextTypeAtGame 游戏
	ContextTypeAtGame = "game"
	// ContextTypeAtPoll 调查
	ContextTypeAtPoll = "poll"
	// ContextTypeAtVenue 场地
	ContextTypeAtVenue = "venue"
	// ContextTypeAtLocation 共享位置
	ContextTypeAtLocation = "location"
)

// Context 上下文
type Context struct {
	*telegram.API                   // 所有 api 方法
	MessageType   string            // 消息类型
	Message       *telegram.Message // 接收到的消息
}

// GetChatID 获取会话 ID
func (c Context) GetChatID() string {
	if c.Message != nil {
		if c.Message.Chat != nil {
			return utils.ToString(c.Message.Chat.ID)
		}
	}
	return ""
}

// GetDownloadURL 获取文件下载地址下载地址
func (c Context) GetDownloadURL(filePath string) string {
	return fmt.Sprintf("https://api.telegram.org/file/bot%d:%s/%s", c.ID, c.Token, filePath)
}

//DownloadFile 下载文件
func (c Context) DownloadFile(filePath string) (io.ReadCloser, error) {
	res, err := c.HTTPClient.DeleteBaseURL().Get(c.GetDownloadURL(filePath))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("http response code is not a 200")
	}

	return res.Body, err
}
