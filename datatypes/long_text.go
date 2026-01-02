package datatypes

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// LongText 是一个自定义类型，用于自动适配不同数据库的长文本类型
type LongText string

// GormDataType 根据数据库类型返回合适的字段类型
func (LongText) GormDataType() string {
	// 这里不能直接知道 db 类型，所以通常返回一个通用名，
	// 然后在 GormDBDataType 中处理
	return "long_text"
}

// GormDBDataType 在实际建表时被调用，可访问 db 对象
func (LongText) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "TEXT" // PostgreSQL 没有 longtext，TEXT 可存无限长度
	case "mysql":
		return "LONGTEXT"
	case "sqlite":
		return "TEXT"
	default:
		return "TEXT"
	}
}

// 实现 Scanner 和 Valuer 以支持数据库读写（可选，但推荐）
func (g *LongText) Scan(value interface{}) error {
	if value == nil {
		*g = ""
		return nil
	}
	if str, ok := value.([]byte); ok {
		*g = LongText(str)
		return nil
	}
	if str, ok := value.(string); ok {
		*g = LongText(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into LongText", value)
}

func (g LongText) Value() (driver.Value, error) {
	return string(g), nil
}
