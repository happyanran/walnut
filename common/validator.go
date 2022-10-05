package common

import (
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

type Vlidate struct {
	Val   *validator.Validate
	Trans ut.Translator
}

func NewValidate() (*Vlidate, error) {
	// 校验器
	val := validator.New()

	// 中文翻译器
	uni := ut.New(zh.New())
	trans, _ := uni.GetTranslator("zh")

	// 注册翻译器到校验器
	err := zh_translations.RegisterDefaultTranslations(val, trans)
	if err != nil {
		return nil, err
	}

	return &Vlidate{
		Val:   val,
		Trans: trans,
	}, nil
}

func (s Vlidate) Struct(i interface{}) *map[string]string {
	err := s.Val.Struct(i)
	if err != nil {
		msg := make(map[string]string)
		for _, m := range err.(validator.ValidationErrors) {
			msg[m.Field()] = m.Translate(s.Trans)
		}
		return &msg
	}
	return nil
}
