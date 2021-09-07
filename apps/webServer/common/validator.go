package common

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"strings"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

//定义字段名称
var fieldName = map[string]string{
	"Username": "用户名",
	"Password": "用户密码",
	"Avatar":   "用户头像",
}

//定义自定义
var myTags = map[string]string{
	"myvalidate": "必须通过自定义方法验证",
}

//自定义验证方法
var myvalidate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(string)
	if ok {
		if date != "gangan" {
			return false
		}
	}

	return true
}

func init() {
	translator := zh.New()
	uni = ut.New(translator, translator)
	trans, _ = uni.GetTranslator("zh")

	var (
		validate *validator.Validate
		ok       bool
	)

	//注册自定义验证方法
	if validate, ok = binding.Validator.Engine().(*validator.Validate); ok {
		validate.RegisterValidation("myvalidate", myvalidate)
	}

	_ = zhTranslations.RegisterDefaultTranslations(validate, trans)
}

func Translate(err error) string {
	var result string

	errors := err.(validator.ValidationErrors)

	for _, err := range errors {

		//判断是自定义验证方法
		var (
			tag   = err.Tag()   //绑定的验证方法
			field = err.Field() //绑定的验证字段
			msg   string
		)
		if val, exist := myTags[tag]; exist {
			msg = field + val + ";"
		} else {
			msg = err.Translate(trans) + ";"
		}

		if val, exist := fieldName[field]; exist {
			msg = strings.Replace(msg, field, val, 1)
		}

		result += msg
	}
	return result
}
