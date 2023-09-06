package validator

import (
    "fmt"
    "reflect"

    "github.com/go-playground/locales/zh_Hans_CN"
    unTrans "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
    zhTrans "github.com/go-playground/validator/v10/translations/zh"
    "github.com/hd2yao/ginblog/utils/err_msg"
)

func Validate(data any) (string, int) {
    validate := validator.New()

    // zh_Hans_CN.New() 创建一个中文简体的语言包
    // 创建一个该语言包相关联的翻译器实例
    uni := unTrans.New(zh_Hans_CN.New())
    // 获取与特定语言包相关联的翻译器（中文简体 翻译器）
    trans, _ := uni.GetTranslator("zh_Hans_CN")

    // 注册默认的翻译信息，这一步将会让验证器知道如何将错误信息翻译成指定语言
    err := zhTrans.RegisterDefaultTranslations(validate, trans)
    if err != nil {
        fmt.Println("err:", err)
    }
    validate.RegisterTagNameFunc(func(field reflect.StructField) string {
        label := field.Tag.Get("label")
        return label
    })

    err = validate.Struct(data)
    if err != nil {
        for _, v := range err.(validator.ValidationErrors) {
            return v.Translate(trans), err_msg.ERROR
        }
    }
    return "", err_msg.SUCCESS
}
