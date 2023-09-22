package models

import (
	"github.com/efucloud/common"
	"github.com/go-playground/validator/v10"
	"time"
)

type Workspace struct {
	ID           uint      `gorm:"primarykey,omitempty" json:"id" description:"主键"` // 记录ID
	CreatedAt    time.Time `json:"createdAt,omitempty" description:"创建时间"`          // 创建时间
	UpdatedAt    time.Time `json:"updatedAt,omitempty" description:"更新时间"`          // 更新时间
	Organization string    `gorm:"type:varchar(50)" json:"organization" validate:"code" description:"所属组织(组织编码)"`
	Code         string    `gorm:"type:varchar(50)" json:"code" validate:"code" description:"编码"` // 编码
	Name         string    `gorm:"type:varchar(50)" json:"name"  description:"名称"`                // 名称
	Description  string    `gorm:"type:varchar(255)" json:"description"  description:"描述"`        // 说明
	OwnerID      uint      `json:"ownerId" validate:"required,gt=0" description:"工作空间管理员ID"`      // 管理员ID
	Owner        Account   `gorm:"-" validate:"-" json:"owner" description:"管理员"`
}
type WorkspaceList struct {
	Data  []*Workspace `json:"data"`  //
	Total int64        `json:"total"` //
}

func (t *Workspace) Default() {

}
func (t Workspace) TableName() string {
	return "workspace"
}
func (t *Workspace) Indexes() (results map[string][]string) {
	results = make(map[string][]string)
	return
}
func (t *Workspace) UniqueIndexes() (results map[string][]string) {
	results = make(map[string][]string)
	results["uniq_idx_code"] = []string{"organization", "code"}
	results["uniq_idx_name"] = []string{"organization", "name"}
	return
}

func (t *Workspace) Validate() (err error) {
	validate := validator.New()
	_ = validate.RegisterValidation("code", common.CodeValidate)
	validate.RegisterTagNameFunc(common.TagNameFunc)
	err = validate.Struct(t)
	return
}
