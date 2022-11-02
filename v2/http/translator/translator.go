package translator

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhs "github.com/go-playground/validator/v10/translations/zh"
)

var (
	validate = validator.New()          // 实例化验证器
	chinese  = zh.New()                 // 获取中文翻译器
	uni      = ut.New(chinese, chinese) // 设置成中文翻译器
	trans, _ = uni.GetTranslator("zh")  // 获取翻译字典
)

// InitTrans 初始化翻译器
func InitTrans(p interface{}) interface{} {
	_ = zhs.RegisterDefaultTranslations(validate, trans)
	err := validate.Struct(p)
	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			fmt.Println(errors.Error(), errors.Translate(trans))
			return errors.Translate(trans)
		}
	}
	return nil
}
