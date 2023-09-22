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
	"github.com/go-playground/validator/v10"
	"time"
)

type DbHistory struct {
	ID        uint       `gorm:"primarykey,omitempty" json:"id" description:"主键"` // 记录ID
	CreatedAt *time.Time `json:"createdAt,omitempty" description:"创建时间"`          // 创建时间
	Table     string     `json:"table" validate:"required" `                      // 表名称
	Name      string     `json:"name" validate:"required"`                        // 索引或unique的Name
	Category  string     `gorm:"type:varchar(100)" json:"category" validate:"oneof=index unique"`
	Content   string     `gorm:"type:longtext" json:"content" validate:"required"`
	Hash      string     `gorm:"type:varchar(100)" json:"hash" validate:"required"`
}

type DbHistoryList struct {
	Data  []*DbHistory `json:"data"`  //
	Total int64        `json:"total"` //
}

func (h *DbHistory) TableName() string {
	return "db_history"
}
func (t *DbHistory) Indexes() (results map[string][]string) {
	results = make(map[string][]string)
	return
}
func (t *DbHistory) UniqueIndexes() (results map[string][]string) {
	results = make(map[string][]string)
	results["uniq_idx_history"] = []string{"hash"}
	return
}
func (t *DbHistory) Default() {
	t.Hash = common.MD5V(t.Content)
}
func (h *DbHistory) Validate() (err error) {
	validate := validator.New()
	validate.RegisterTagNameFunc(common.TagNameFunc)
	err = validate.Struct(h)
	return
}
