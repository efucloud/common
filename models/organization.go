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

package models

import (
	"github.com/efucloud/common"
	"time"
)

type Organization struct {
	ID          uint      `gorm:"primarykey" json:"id" description:"主键"`                                                                                             // 记录ID
	CreatedAt   time.Time `json:"createdAt" description:"创建时间"`                                                                                                      // 创建时间
	UpdatedAt   time.Time `json:"updatedAt" description:"更新时间"`                                                                                                      // 更新时间
	Name        string    `gorm:"type:varchar(50)" json:"name" description:"组织名称"`                                                                                   // 组织名称
	Code        string    `gorm:"type:varchar(50)" json:"code" validate:"code" description:"组织编码不可修改,不能为:company public top global default org organization system"` // 组织编码不可修改
	Description string    `gorm:"type:varchar(255)" json:"description" description:"组织描述"`                                                                           // 应用描述
	OwnerID     uint      `json:"ownerId" validate:"required,gt=0" description:"组织管理员ID"`                                                                            // 管理员ID
	Owner       Account   `gorm:"-" validate:"-" json:"owner" description:"管理员"`
	Enable      uint      `gorm:"type:uint;size:8;default:1" json:"enable" validate:"oneof=0 1" enum:"1|0" description:"是有有效"` // 组织是否有效
}
type OrganizationList struct {
	Data  []*Organization `json:"data" description:"数据列表"`   //
	Total int64           `json:"total" description:"记录总数量"` //
}

func (t *Organization) Indexes() (results map[string][]string) {
	results = make(map[string][]string)
	return
}
func (t *Organization) UniqueIndexes() (results map[string][]string) {
	results = make(map[string][]string)
	results["uniq_idx_org_name"] = []string{"name"}
	results["uniq_idx_org_code"] = []string{"code"}
	return
}

func (t Organization) TableName() string {
	return "organization"
}
func (t *Organization) Default() {

}
func (t *Organization) Validate() (err error) {
	validate := validator.New()
	_ = validate.RegisterValidation("supportLogins", supportLoginsValidate)
	_ = validate.RegisterValidation("category", categoryValidate)
	_ = validate.RegisterValidation("code", common.CodeValidate)
	_ = validate.RegisterValidation("mfa", mfaValidate)
	validate.RegisterTagNameFunc(common.TagNameFunc)
	err = validate.Struct(t)
	return
}
func mfaValidate(fl validator.FieldLevel) bool {
	allows := []string{"totp", "sms", "email"}
	if fl.Field().Interface() != nil {
		items := fl.Field().Interface().
		(dt.
		[]string)
		for _, item := range items {
			if common.StringKeyInArray(item, allows) {
				return true
			}
		}

	}
	return false
}
func supportLoginsValidate(fl validator.FieldLevel) bool {
	allows := []string{"username", "phone", "email", "phoneCode", "emailCode", "oidc", "ldap"}
	if fl.Field().Interface() != nil {
		items := fl.Field().Interface().
		(dt.
		[]string)
		for _, item := range items {
			if common.StringKeyInArray(item, allows) {
				return true
			}
		}

	}
	return false
}
func categoryValidate(fl validator.FieldLevel) bool {
	allows := []string{"enterprise", "customer", "provider", "customer_and_provider", "consumer"}
	return common.StringKeyInArray(fl.Field().String(), allows)
}
