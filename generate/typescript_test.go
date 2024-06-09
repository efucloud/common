package generate

import (
	"github.com/efucloud/common/datatypes"
	"os"
	"reflect"
	"testing"
	"time"
)

type EmbedStruct struct {
	Doc     string  `gorm:"-" json:"-" description:"内嵌表"`
	Field1  string  `json:"field1" description:""`
	Field2  float64 `json:"field2" description:""`
	Boolean bool    `json:"boolean" description:""`
}
type Account struct {
	Doc           string                 `gorm:"-" json:"-" description:"系统账户表"`
	ID            uint                   `gorm:"primarykey" json:"id" description:"主键"`                                           // 记录ID
	CreatedAt     *time.Time             `json:"createdAt" description:"创建时间"`                                                    // 创建时间
	UpdatedAt     time.Time              `json:"updatedAt" description:"更新时间"`                                                    // 更新时间
	Username      string                 `gorm:"type:varchar(255)" json:"username" validate:"required,max=255" description:"用户名"` // 用户名 组织内唯一
	Nickname      string                 `gorm:"type:varchar(255)" json:"nickname" validate:"max=255" description:"昵称"`           // 昵称，如中文名
	JobNumber     string                 `gorm:"type:varchar(255)" json:"jobNumber" validate:"max=255" description:"工号"`
	Password      string                 `gorm:"-" json:"password" example:"admin" description:"密码"`                                                                             // 密码
	PasswordStore string                 `gorm:"type:varchar(255)" json:"-"`                                                                                                     // 数据库保存的加密密码
	Role          string                 `gorm:"type:varchar(50);default:none" json:"role" validate:"oneof=admin view edit none" enum:"admin|view|edit|none" description:"组织角色"` // 组织角色
	Enable        uint                   `gorm:"type:uint;size:8;default:1" json:"enable" validate:"oneof=0 1" enum:"1|0" description:"是否有效"`                                    // 是否有效，组织管理员不能设置为无效
	Email         string                 `gorm:"type:varchar(255)" json:"email" validate:"max=255" description:"邮箱"`                                                             // 邮箱
	Phone         string                 `gorm:"type:varchar(50)" json:"phone" validate:"max=50" description:"电话"`                                                               // 手机号码
	Language      string                 `gorm:"type:varchar(50)" json:"language" validate:"max=50" description:"语言"`
	ArrayStrings  *datatypes.ArrayString `json:"arrayStrings" description:""`
	ArrayUint     datatypes.ArrayUint    `json:"arrayUint" description:""`
	ParentPtr     *EmbedStruct           `json:"parentPtr" description:""`
	Parent        EmbedStruct            `json:"parent" description:""`
}

func TestAccount(t *testing.T) {
	typ := NewTypeScript()
	typ.AddStruct(reflect.TypeOf(Account{}))
	dir, _ := os.Getwd()
	typ.GenerateToDir(dir)
}
