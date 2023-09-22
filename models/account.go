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
	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
	"time"
)

type Account struct {
	ID           uint      `gorm:"primarykey,omitempty" json:"id" description:"主键"` // 记录ID
	AuthId       string    `gorm:"type:varchar(50)" json:"authId" validate:"required"  description:"第三方认证的ID"`
	CreatedAt    time.Time `json:"createdAt,omitempty" description:"创建时间"` // 创建时间
	UpdatedAt    time.Time `json:"updatedAt,omitempty" description:"更新时间"` // 更新时间
	Organization string    `gorm:"type:varchar(50);default:efu-cloud" json:"organization" validate:"code" description:"所属组织(组织编码)"`
	Username     string    `gorm:"type:varchar(100)" json:"username" description:"用户名"` // 用户名
	Nickname     string    `gorm:"type:varchar(255)" json:"nickname" description:"别名"`  // 别名
	Email        string    `gorm:"type:varchar(255)" json:"email" description:"邮箱" unique:"true"`
	Phone        string    `gorm:"type:varchar(255)" json:"phone" description:"手机号码" unique:"true"`
	OrgRole      string    `gorm:"type:varchar(50);default:none" json:"orgRole" validate:"oneof=admin view edit none" enum:"admin|view|edit|none" description:"组织角色"` // 组织角色
	Enable       uint      `gorm:"type:uint;size:8" json:"enable" description:"是否启用" validate:"oneof=0 1" enum:"1|0"`
}
type AccountList struct {
	Data  []*Account `json:"data"`  //
	Total int64      `json:"total"` //
}

func (t *Account) Default() {

}
func (t *Account) TableName() string {
	return "account"
}
func (t *Account) Indexes() (results map[string][]string) {
	results = make(map[string][]string)
	results["idx_enable"] = []string{"enable"}
	return
}
func (t *Account) UniqueIndexes() (results map[string][]string) {
	results = make(map[string][]string)
	results["uniq_idx_username"] = []string{"organization", "username"}
	results["uniq_idx_username"] = []string{"organization", "auth_id"}
	return
}
func (t *Account) Validate() (err error) {
	validate := validator.New()
	validate.RegisterTagNameFunc(common.TagNameFunc)
	err = validate.Struct(t)
	return
}
func (t *Account) ValueMap() (values map[string]interface{}) {
	values = structs.Map(t)
	return
}
