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

package eauth

import "github.com/golang-jwt/jwt/v4"

type OIDCConfig struct {
	// 提供商的地址，如https://gitlab.com,在不配置Certificate的情况下，程序会根据https://gitlab.com/.well-known/openid-configuration获取token的公钥
	Issuer string `json:"issuer" yaml:"issuer"`
	// 应用的ClientID
	ClientID string `json:"clientId" yaml:"clientId"`
	// 应用的ClientSecret
	ClientSecret string `json:"clientSecret" yaml:"clientSecret"`
	// 跳转到认证的页面，如https://gitlab.com/oauth/authorize，该信息会返回给前端用于前端组成认证重定向地址
	AuthorizationEndpoint string `json:"authorizationEndpoint" yaml:"authorizationEndpoint"`
	// 认证完成后的重定向地址，用于接收返回的code，如gitlab认证成功后返回的code,state或者err信息,
	// 前后端分离模式下，该地址为前端地址，可由前端自行拼接
	RedirectURI string `json:"redirectUri" yaml:"redirectUri"`
	// 获取eauth Token的地址
	TokenEndpoint string `json:"tokenEndpoint" yaml:"tokenEndpoint"`
	// 获取用户信息的地址
	UserinfoEndPoint string `json:"userinfoEndPoint" yaml:"userinfoEndPoint"`
	// 提供商的ca信息，可以不提供，
	IssuerCA       string `json:"issuerCa" yaml:"issuerCa"`
	UsernameClaim  string `json:"usernameClaim" yaml:"usernameClaim"`
	UsernamePrefix string `json:"usernamePrefix" yaml:"usernamePrefix"`
	// token校验的公钥信息，若不配置，应用需要根据Issuer+/.well-known/openid-configuration去获取
	// 若以gitlab为例https://gitlab.com/.well-known/openid-configuration
	Certificate string `json:"certificate" yaml:"certificate"`
}
type UserInfo struct {
	Subject          string                 `json:"sub"`
	Profile          string                 `json:"profile"`
	Email            string                 `json:"email"`
	EmailVerified    bool                   `json:"email_verified"`
	Org              string                 `json:"org"`
	OrgCustoms       map[string]interface{} `json:"orgCustoms"` // 组织自定义属性
	Providers        []string               `json:"providers"`
	Groups           []string               `json:"groups"`
	RegistrationFrom string                 `json:"registrationFrom"` // 注册渠道
	AuthProvider     string                 `json:"authProvider"`     // 认证提供商
	Username         string                 `json:"username"`         // 用户名 组织内唯一必须由DNS-1123标签格式的单元组成
	Nickname         string                 `json:"nickname"`         // 昵称，如中文名
	OrgRole          string                 `json:"orgRole"`          //组织角色
	Phone            string                 `json:"phone"`
	ID               uint                   `json:"id"`
}

type ApplicationSyncAccountInfo struct {
	AuthId           uint                   `json:"authId" yaml:"authId"`
	Organization     string                 `json:"organization" validate:"required"` // 组织编码
	Username         string                 `json:"username" validate:"dns1123"`      // 用户名 组织内唯一必须由DNS-1123标签格式的单元组成
	Nickname         string                 `json:"nickname"`                         // 昵称，如中文名
	AdminApps        []string               `json:"adminApps"`                        // 应用管理员
	Enable           uint                   `json:"enable" validate:"oneof=0 1"`      // 是否有效，组织管理员不能设置为无效
	OrgCustoms       map[string]interface{} `json:"orgCustoms"`                       // 组织自定义属性
	Hash             string                 `json:"hash"`                             // 组织:用户名的Hash
	RegistrationFrom string                 `json:"registrationFrom"`                 // 注册渠道
	Language         string                 `json:"language" validate:"oneof=en zh"`  // 语言
	Email            string                 `json:"email" yaml:"email"`
	Phone            string                 `json:"phone" yaml:"phone"`
	Groups           []string               `json:"groups" yaml:"groups"`
}
type AccountClaims struct {
	AccountID    uint     `json:"accountId" `
	Org          string   `json:"org"`
	AuthProvider string   `json:"authProvider"`
	Username     string   `json:"username"` // 用户名 组织内唯一必须由DNS-1123标签格式的单元组成
	Nickname     string   `json:"nickname"` // 昵称，如中文名
	OrgRole      string   `json:"orgRole"`  // 组织角色
	Nonce        string   `json:"nonce"`
	Email        string   `json:"email"`
	Phone        string   `json:"phone"`
	Groups       []string `json:"groups"`
	AppCode      string   `json:"appCode"`
	AppClientID  string   `json:"appClientId"`
	AppOwner     bool     `json:"appOwner"`
	jwt.RegisteredClaims
}

// LocalLoginParam 本地登录请求
type LocalLoginParam struct {
	Method      string `json:"method" validate:"oneof=password phoneCode emailCode"` // 登录类型，用户名密码/手机验证码/邮箱验证码/
	Username    string `json:"username"`
	Password    string `json:"password"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	ValidCode   string `json:"validCode"`
	Code        string `json:"code"`
	State       string `json:"state"`
	RedirectUri string `json:"redirectUri" validate:"required"`
	Bind        string `json:"bind"`
}
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token"`
}

type RefreshToken struct {
	Org       string `json:"org"`
	App       string `json:"app"`
	ExpiresIn int64  `json:"expiresIn"`
	AccountID uint   `json:"accountId"`
	Issuer    string `json:"issuer"`
	Provider  string `json:"provider"`
}

type AccountSync struct {
	CronJob string `json:"cronJob"`
	Address string `json:"address"`
}
