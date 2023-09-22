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
	"fmt"
	"github.com/efucloud/common"
	"github.com/go-playground/validator/v10"
	"time"
)

type OrganizationApply struct {
	ID          uint      `gorm:"primarykey" json:"id" description:"主键"`                                                                                             // 记录ID
	CreatedAt   time.Time `json:"createdAt" description:"创建时间"`                                                                                                      // 创建时间
	UpdatedAt   time.Time `json:"updatedAt" description:"更新时间"`                                                                                                      // 更新时间
	Name        string    `gorm:"type:varchar(50)" json:"name" description:"组织名称"`                                                                                   // 组织名称
	Code        string    `gorm:"type:varchar(50)" json:"code" validate:"code" description:"组织编码不可修改,不能为:company public top global default org organization system"` // 组织编码不可修改
	Description string    `gorm:"type:varchar(255)" json:"description" description:"应用描述"`                                                                           // 应用描述
	OwnerID     uint      `json:"ownerId" validate:"required" description:"管理员ID"`                                                                                   // 管理员ID
	Owner       Account   `gorm:"-" validate:"-" json:"owner" description:"拥有者"`
	State       string    `gorm:"type:varchar(50);default:draft" json:"state" validate:"oneof=draft confirm approval reject" enum:"draft|confirm|approval|reject" description:"申请状态"` // 申请状态
}
type OrganizationApplyList struct {
	Data  []*OrganizationApply `json:"data" description:"数据列表"`   //
	Total int64                `json:"total" description:"记录总数量"` //
}
type StateChange struct {
	ID    uint   `json:"id" description:""`                                                   //  记录ID
	State string `json:"state" validate:"oneof=draft confirm approval reject" description:""` //
	Note  string `json:"note" description:""`                                                 //  说明
}

func (org *OrganizationApply) Indexes() (results map[string][]string) {
	results = make(map[string][]string)
	return
}
func (org *OrganizationApply) UniqueIndexes() (results map[string][]string) {
	results = make(map[string][]string)
	results["uniq_idx_code"] = []string{"code"}
	return
}
func (org OrganizationApply) TableName() string {
	return "organization_apply"
}
func (org *OrganizationApply) Default() {
	if len(org.State) == 0 {
		org.State = "draft"
	}
}
func (org *OrganizationApply) Validate() (err error) {
	if org.Code == "efu-cloud" {
		err = fmt.Errorf("org code: %s is invalid", org.Code)
		return
	}
	validate := validator.New()
	_ = validate.RegisterValidation("code", common.CodeValidate)
	validate.RegisterTagNameFunc(common.TagNameFunc)
	err = validate.Struct(org)
	return
}
