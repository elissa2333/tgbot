package telegram

// pack

/*Games
您的机器人可以为用户提供HTML5游戏，使他们可以独自玩耍或在小组和一对一聊天中相互竞争。使用 /newgame 命令通过 @BotFather 创建游戏。请注意，这种功能需要承担责任：您将需要接受机器人将提供的每款游戏的条款。

游戏是电报上一种新型的内容，由Game和InlineQueryResultGame对象表示。
通过BotFather创建游戏后，您可以使用sendGame方法将游戏作为常规消息发送到聊天室，或者通过InlineQueryResultGame使用内联模式。
如果您发送的游戏消息没有任何按钮，它将自动具有一个“Play GameName”按钮。当按下此按钮时，您的机器人会收到一个CallbackQuery，其中包含所请求游戏的 game_short_name。您为此特定用户提供正确的URL，然后该应用会在应用内浏览器中打开游戏。
您可以在游戏消息中手动添加多个按钮。请注意，第一行中的第一个按钮必须始终运行游戏，使用领域callback_game在InlineKeyboardButton。您可以根据喜好添加其他按钮：例如，有关规则的描述或打开游戏的官方社区。
为了使您的游戏更具吸引力，您可以通过BotFather向用户上传GIF动画，向用户演示游戏内容（例如，参见Lumberjack）。
游戏消息还将显示当前聊天的高分。使用setGameScore将高分发布到与游戏的聊天中，添加edit_message参数以使用当前记分板自动更新消息。
使用getGameHighScores获取游戏内高分表的数据。
您还可以添加一个额外的共享按钮，供用户将其最佳分数分享到不同的聊天中。
有关使用此新内容可以完成的操作的示例，请检查@gamebot和@gamee机器人。
https://core.telegram.org/bots/api#games*/

// SendGameOptionl SendGame 可选参数
type SendGameOptionl struct {
	DisableNotification      bool                  `json:"disable_notification"`
	ReplyToMessageID         int64                 `json:"reply_to_message_id"`
	AllowSendingWithoutReply bool                  `json:"allow_sending_without_reply"`
	ReplyMarkup              *InlineKeyboardMarkup `json:"reply_markup *InlineKeyboardMarkup"`
}

// SendGame 使用此方法发送游戏。成功后，将返回发送的 Message。
// https://core.telegram.org/bots/api#sendgame
func (a API) SendGame(chatID int64, gameShortName string, optionl *SendGameOptionl) (*Message, error) {
	result := &Message{}
	err := a.handleOptional("/sendGame", map[string]interface{}{"chat_id": chatID, "game_short_name": gameShortName}, optionl, result)
	return result, err
}

//Game 代表一个游戏。使用BotFather创建和编辑游戏，它们的简称将用作唯一的标识符
//https://core.telegram.org/bots/api#game
type Game struct {
	Title        string          `json:"title"`         // 游戏标题
	Description  string          `json:"description"`   // 游戏说明
	Photo        []PhotoSize     `json:"photo"`         // 将在聊天中的游戏消息中显示的照片。
	Text         string          `json:"text"`          // 可选的。对游戏的简要说明或游戏消息中包含的高分。当漫游器调用setGameScore时，可以自动编辑以包括当前游戏的高分，或使用editMessageText进行手动编辑。0-4096个字符。
	TextEntities []MessageEntity `json:"text_entities"` // 可选的。以文本形式出现的特殊实体，例如用户名，URL，机器人命令等。
	Animation    *Animation      `json:"animation"`     // 可选的。将在聊天中显示在游戏消息中的动画。通过BotFather上传
}

// CallbackGame 占位符，当前不保存任何信息。使用BotFather来设置您的游戏
// https://core.telegram.org/bots/api#callbackgame
type CallbackGame struct{}

// SetGameScoreOptionl SetGameScore可选参数
type SetGameScoreOptionl struct {
	Force              bool   `json:"force"`                // 如果允许降低高分，则通过True。在纠正错误或禁止作弊者时，这很有用
	DisableDditMessage bool   `json:"disable_edit_message"` // 如果游戏消息不应该自动编辑为包括当前记分板，则通过True
	ChatID             int64  `json:"chat_id"`              // 如果未指定inline_message_id，则为必需。目标聊天的唯一标识符
	MessageID          int64  `json:"message_id"`           // 如果未指定inline_message_id，则为必需。发送邮件的标识符
	InlineMessageID    string `json:"inline_message_id"`    // 如果未指定chat_id和message_id，则为必需。内联消息的标识符
}

// SetGameScore 使用此方法可以设置游戏中指定用户的分数。成功后，如果该消息是由漫游器发送的，则返回已编辑的Message，否则返回True。如果新分数不超过用户在聊天中的当前分数并且强制为False，则返回错误。
// https://core.telegram.org/bots/api#setgamescore
func (a API) SetGameScore(userID int64, score int64, optionl *SetGameScoreOptionl) (*Message, error) {
	result := &Message{}
	err := a.handleOptional("/setGameScore", map[string]interface{}{"user_id": userID, "score": score}, optionl, result)

	return result, err
}

// GetGameHighScoresOptionl GetGameHighScores可选参数
type GetGameHighScoresOptionl struct {
	ChatID          int64  `json:"chat_id"`           // 可选的	如果未指定inline_message_id，则为必需。目标聊天的唯一标识符
	MessageID       int64  `json:"message_id"`        // 如果未指定inline_message_id，则为必需。发送邮件的标识符
	InlineMessageID string `json:"inline_message_id"` // 如果未指定chat_id和message_id，则为必需。内联消息的标识符
}

// GetGameHighScores 使用此方法获取高分表的数据。将返回游戏中指定用户及其几个邻居的分数。如果成功，则返回数组的GameHighScore对象。
// https://core.telegram.org/bots/api#getgamehighscores
func (a API) GetGameHighScores(userID int64, optionl *GetGameHighScoresOptionl) ([]GameHighScore, error) {
	var result []GameHighScore
	err := a.handleOptional("/getGameHighScores", map[string]interface{}{"user_id": userID}, optionl, &result)

	return result, err

}

// GameHighScore 该对象表示游戏的高分表的一行。
// https://core.telegram.org/bots/api#gamehighscore
type GameHighScore struct {
	Position int64 `json:"position"` // 在游戏高分表中的位置
	User     *User `json:"user"`     // 用户
	Score    int64 `json:"score"`    // 得分了
}
