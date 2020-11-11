package tgbot

import (
	"fmt"
	"os"

	"github.com/elissa2333/tgbot/telegram"
	"github.com/elissa2333/tgbot/utils"
)

func ExampleBot_SetInlineQueryProcessor() {
	tg := New(utils.ToInt(os.Getenv("id")), os.Getenv("token"), nil)
	tg.SetInlineQueryProcessor(func(c *InlineQueryContext) error {
		content := "content" // 必须有内容
		if c.Query != "" {
			content = c.Query
		}
		var result []telegram.InlineQueryResult
		result = append(result, telegram.InlineQueryResult{
			InlineQueryResultArticle: &telegram.InlineQueryResultArticle{
				Type:  telegram.InlineQueryResultArticleType,
				ID:    "123456",
				Title: content,
				InputMessageContent: telegram.InputMessageContent{
					InputTextMessageContent: &telegram.InputTextMessageContent{
						MessageText: content,
					}}}})
		ok, err := c.AnswerInlineQuery(c.InlineQuery.ID, result, nil)
		if err != nil {
			return err
		}
		fmt.Println(ok) // 每次成功调用会返回 true
		return nil
	})

	if err := tg.Run(); err != nil {
		panic(err)
	}
}
