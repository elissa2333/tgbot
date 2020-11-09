package telegram

import (
	"os"
	"testing"

	"github.com/elissa2333/httpc"

	"github.com/elissa2333/tgbot/utils"
)

var (
	tAPI *API

	// 在进行测试前你需要先设置bot的 `id` `token` 变量以及 `chat_id`
	tID     = utils.ToInt(os.Getenv("id"))
	tToken  = os.Getenv("token")
	tChatID = os.Getenv("chat_id")

	tProxyS = os.Getenv("http_proxy") // 如果你需要使用代理的话请设置代理 `http_proxy` 变量
)

func TestMain(m *testing.M) {
	if tProxyS != "" {
		tAPI = New(&httpc.SetProxy(tProxyS).Client, tID, tToken)
	} else {
		tAPI = New(nil, tID, tToken)
	}

	os.Exit(m.Run())
}
