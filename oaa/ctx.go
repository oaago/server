package oaa

import (
	"encoding/json"
	"github.com/duke-git/lancet/convertor"
	"github.com/gin-gonic/gin"
	"github.com/oaago/component/op"
	"strconv"
)

type Ctx struct {
	*gin.Context
}

func (c *Ctx) Return(code int, arg ...interface{}) {
	// arg 第一位代表 data 字段 第二位代表 message
	var message string
	if len(arg) > 1 {
		message = convertor.ToString(arg[1])
		var msg map[string]interface{}
		json.Unmarshal([]byte(message), &msg)
		for _, i := range msg {
			message = convertor.ToString(i)
		}
	}
	if len(op.ConfigData.CodeMap) == 0 {
		panic("没有配置好 codeMap 请配置好 在启动项目......")
	}
	CodeMap := op.ConfigData.CodeMap //nolint:typecheck
	// 200 无异常
	// 10000 自己以下定义异常
	// 其他 consts.CodeMap 定义异常
	if code != 200 {
		message = CodeMap[code]
		if message == "" {
			panic("暂未定义错误码: " + strconv.Itoa(code))
		}
	}

	c.JSON(200, Result{
		Code:    code,
		Message: message,
		Data:    arg[0],
	})
}
