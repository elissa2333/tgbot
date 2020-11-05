package telegram

import (
	"github.com/elissa2333/httpc"
)

// Response 响应实体
type Response struct {
	Ok     bool        `json:"ok"` // 是否请求成功
	Result interface{} // 响应结果

	ErrorCode   int    `json:"error_code"`  // 错误状态码
	Description string `json:"description"` // 错误说明
}

// Error 实现 error 接口包裹错误
func (r Response) Error() string {
	return r.Description
}

// HandleResp 处理响应
func HandleResp(res *httpc.Response, resultByPtr interface{}) error {
	defer res.Body.Close()

	m := &Response{Result: resultByPtr}
	if err := res.ToJSON(m); err != nil {
		return err
	}
	if !m.Ok {
		return m
	}

	return nil
}
