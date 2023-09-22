package models

import (
	"github.com/efucloud/common"
	"github.com/go-playground/validator/v10"
	"time"
)

type WorkspaceMember struct {
	ID           uint      `gorm:"primarykey,omitempty" json:"id" description:"主键"` // 记录ID
	CreatedAt    time.Time `json:"createdAt,omitempty" description:"创建时间"`          // 创建时间
	UpdatedAt    time.Time `json:"updatedAt,omitempty" description:"更新时间"`          // 更新时间
	Organization string    `gorm:"type:varchar(50)" json:"organization" validate:"code" description:"所属组织(组织编码)"`
	AccountID    uint      `json:"accountId" description:"用户ID"`                                                  //
	Workspace    string    `json:"workspace" description:"工作空间编码"`                                                //
	Role         string    `json:"role" validate:"oneof=admin edit view" enum:"admin|edit|view" description:"角色"` //
}

type WorkspaceMemberList struct {
	Data  []*WorkspaceMember `json:"data"`  //
	Total int64              `json:"total"` //
}

func (t WorkspaceMember) TableName() string {
	return "workspace_member"
}
func (t *WorkspaceMember) Indexes() (results map[string][]string) {
	results = make(map[string][]string)
	return
}
func (t *WorkspaceMember) UniqueIndexes() (results map[string][]string) {
	results = make(map[string][]string)
	results["uniq_idx_ids"] = []string{"organization", "workspace", "account_id"}
	return
}
func (t *WorkspaceMember) Default() {
}
func (t *WorkspaceMember) Validate() (err error) {
	validate := validator.New()
	validate.RegisterTagNameFunc(common.TagNameFunc)
	err = validate.Struct(t)
	return
}
