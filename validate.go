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
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var dns1123Reg *regexp.Regexp

func init() {
	dns1123Reg = regexp.MustCompile(`[a-z0-9]([-a-z0-9]*[a-z0-9])?`)
}
func ValidateTime(fl validator.FieldLevel) bool {
	_, err := time.Parse(TimeFormat, fl.Field().String())
	return err == nil
}
func ValidateDNS1123(fl validator.FieldLevel) bool {
	return dns1123Reg.MatchString(fl.Field().String())
}
func NotBlank(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(field.String())) > 0
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}
func TagNameFunc(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return fld.Name
	}
	return name
}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

type FiledValidFailed map[string]string

func (f FiledValidFailed) String() string {
	var infos []string
	for k, v := range f {
		infos = append(infos, fmt.Sprintf("%s:%s", k, v))
	}
	return strings.Join(infos, ";")
}
func (f FiledValidFailed) LocaleString(localeMap map[string]interface{}) string {
	var infos []string
	for key, value := range f {
		if v, exist := localeMap[key]; exist {
			infos = append(infos, fmt.Sprintf("%s:%s", v, value))
		} else {
			infos = append(infos, fmt.Sprintf("%s:%s", key, value))
		}
	}
	return strings.Join(infos, ";")
}
func (f FiledValidFailed) LocaleMap(localeMap map[string]string) (result map[string]string) {
	result = make(map[string]string)
	for key, value := range f {
		if v, exist := localeMap[key]; exist {
			result[v] = value
		} else {
			result[key] = value
		}
	}
	return result
}
func CodeValidate(fl validator.FieldLevel) bool {
	alphaRegex := regexp.MustCompile("^[a-z-]+$")
	code := fl.Field().String()
	if len(code) == 0 {
		return false
	}
	allows := []string{"company", "public", "top", "global", "default", "org", "organization", "system"}
	if common.StringKeyInArray(fl.Field().String(), allows) {
		return false
	}
	return alphaRegex.MatchString(code)
}
