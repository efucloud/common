/*
Copyright 2022 The efucloud.com Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"context"
	"embed"
	"github.com/emicklei/go-restful/v3"
	"github.com/ghodss/yaml"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"k8s.io/klog/v2"
	"path"
)

func GetLanguageFromCtx(ctx context.Context, reqAttributeKey string) (lang string) {
	lan := ctx.Value(reqAttributeKey)
	if lan != nil {
		return lan.(string)
	}
	return I18nZH
}
func GetLanguageFromReq(req *restful.Request, reqAttributeKey string) (lang string) {
	langAttr := req.Attribute(reqAttributeKey)
	lang = I18nZH
	if langAttr != nil {
		lang = langAttr.(string)
	}
	return lang
}
func I18nInit(i18nFiles embed.FS) (bundle *i18n.Bundle, universalTranslator *ut.UniversalTranslator) {
	// todo add other language
	universalTranslator = ut.New(en.New(), zh.New())
	bundle = i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	_, err := bundle.LoadMessageFileFS(i18nFiles, path.Join("locales", "en.yaml"))
	if err != nil {
		klog.Fatalf("load i18n message file failed, err: %s", err.Error())
	}
	_, err = bundle.LoadMessageFileFS(i18nFiles, path.Join("locales", "zh.yaml"))
	if err != nil {
		klog.Fatalf("load i18n message file failed, err: %s", err.Error())
	}
	return
}
func GetLocaleMessage(bundle *i18n.Bundle, templateData map[string]interface{}, lang string, id string) (msg string, err error) {
	localizer := i18n.NewLocalizer(bundle, lang)
	msg, err = localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: id}, TemplateData: templateData})
	if len(msg) == 0 {
		msg = id
	}
	return msg, err
}

func ValidateTrans(unTrans *ut.UniversalTranslator, validate *validator.Validate, lang string, err error) FiledValidFailed {
	trans, _ := unTrans.GetTranslator(lang)
	switch lang {
	case I18nZH:
		_ = zh_translations.RegisterDefaultTranslations(validate, trans)
	case I18nEN:
		_ = en_translations.RegisterDefaultTranslations(validate, trans)
	default:
		_ = zh_translations.RegisterDefaultTranslations(validate, trans)
	}
	errs := err.(validator.ValidationErrors)
	return removeTopStruct(errs.Translate(trans))
}

func ValidateTransCtx(unTrans *ut.UniversalTranslator, ctxLangKey string, validate *validator.Validate, ctx context.Context, err error) FiledValidFailed {
	lang := ctx.Value(ctxLangKey)
	lan := I18nZH
	if lang != nil {
		lan = lang.(string)
	}
	trans, _ := unTrans.GetTranslator(lan)
	switch lang {
	case I18nZH:
		_ = zh_translations.RegisterDefaultTranslations(validate, trans)
	case I18nEN:
		_ = en_translations.RegisterDefaultTranslations(validate, trans)
	default:
		_ = zh_translations.RegisterDefaultTranslations(validate, trans)
	}
	errs := err.(validator.ValidationErrors)

	return removeTopStruct(errs.Translate(trans))
}
