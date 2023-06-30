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

package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type OidcConfig struct {
	IssuerURL             string               `gorm:"type:varchar(255)" json:"issuerUrl" validate:"required" description:"Issuer地址"`                                 //
	Scopes                ArrayString          `json:"scopes" validate:"required" description:"域"`                                                                    //
	AutoGetAuthAddress    uint                 `gorm:"type:uint;size:8;default:1" json:"autoGetAuthAddress" gorm:"type:uint;size:8;default:1" description:"自动获取填充信息"` // 自动获取填充下面的信息
	AuthorizationEndpoint string               `gorm:"type:varchar(255)" json:"authorizationEndpoint" validate:"required" description:"AuthorizationEndpoint"`        // 自动填充
	TokenEndpoint         string               `gorm:"type:varchar(255)" json:"tokenEndpoint" validate:"required" description:"TokenEndpoint"`                        // 自动填充
	UserinfoEndpoint      string               `gorm:"type:varchar(255)" json:"userinfoEndpoint" validate:"required" description:"UserinfoEndpoint"`                  // 自动填充
	JwksUri               string               `gorm:"type:varchar(255)" json:"jwksUri" validate:"required" description:"JwksUri"`                                    // 自动填充
	OpenIDConfiguration   *OpenIDConfiguration `json:"openIdConfiguration,omitempty" description:"OpenIDConfiguration"`                                               // openIdConfiguration配置
	ScopesSupported       ArrayString          `json:"scopesSupported" description:"支持的scope"`
}

func (OidcConfig) GormDataType() string {
	return "json"
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (ins *OidcConfig) Scan(value interface{}) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal OpenIDConfiguration value: ", value))
	}
	err := json.Unmarshal(byteValue, ins)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (ins OidcConfig) Value() (driver.Value, error) {
	re, err := json.Marshal(ins)
	return re, err
}
