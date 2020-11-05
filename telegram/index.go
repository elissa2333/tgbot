package telegram

import (
	"fmt"
	"net/http"

	"github.com/elissa2333/httpc"
)

// API 实体
type API struct {
	ID    int    // ID
	Token string // 令牌

	HTTPClient *httpc.Client // http 客户端
}

// New 新建 API 调用器
func New(httpClient *http.Client, id int, token string) *API {
	b := &API{
		ID:    id,
		Token: token,
	}

	if httpClient == nil {
		b.HTTPClient = httpc.New()
	} else {
		b.HTTPClient = httpc.UseClient(*httpClient)
	}

	b.HTTPClient = b.HTTPClient.SetBaseURL(fmt.Sprintf("https://api.telegram.org/bot%d:%s", id, token))

	return b
}
