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

package common

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/net/http2"
	"k8s.io/klog/v2"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultOrder = "id desc"
const DefaultPage = 1
const DefaultPageSize = 20
const (
	QueryTypeEqual       = "eq"
	QueryTypeLike        = "like"
	QueryTypeIn          = "in"
	ParamTypeString      = "string"
	ParamTypeNumber      = "integer"
	ParamTypeBool        = "bool"
	ParamTypeStringSlice = "stringSlice"
	ParamTypeNumberSlice = "numberSlice"
)

type ApiInfo struct {
	Tag         string
	Description string
}

var (
	ApiInfos []ApiInfo
)

func RegisterApiInfo(info ApiInfo) {
	ApiInfos = append(ApiInfos, info)
}
func (a ApiInfo) Tags() []string {
	return []string{a.Tag}
}

func GetRequestPaginationInformation(req *restful.Request) (page int, size int, order string) {
	page = String2Int(req.QueryParameter("page"), DefaultPage)
	size = String2Int(req.QueryParameter("size"), DefaultPageSize)
	order = req.QueryParameter("order")
	if len(order) == 0 {
		order = DefaultOrder
	}
	return page, size, order
}

type QueryParam struct {
	WhereQuery string
	WhereArgs  []interface{}
}

func CreateHttpClient(useHttp2 bool, timeout time.Duration) (client *http.Client) {
	if useHttp2 {
		client = &http.Client{
			Transport: &http2.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: timeout,
		}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: timeout,
		}
	}
	return client
}
func HttpRequest(client *http.Client, method, address string, headers, cookies map[string]interface{}, queries url.Values, body []byte) (response *http.Response, err error) {
	req, err := http.NewRequest(method, address, bytes.NewReader(body))
	if err != nil {
		err = fmt.Errorf("create http request failed, method: %s, address: %s, err: %s", method, address, err.Error())
		klog.Error(err)
		return response, err
	}
	for k, v := range headers {
		req.Header.Add(k, fmt.Sprintf("%s", v))
	}
	for k, v := range cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: fmt.Sprintf("%s", v)})
	}
	req.URL.RawQuery = queries.Encode()
	return client.Do(req)
}

type ResponseError struct {
	Message string `json:"message" yaml:"message"`
	Detail  string `json:"detail" yaml:"detail"`
	Alert   string `json:"alert" yaml:"alert"`
}
type AuthRedirectInfo struct {
	Message               string                 `json:"message"`
	AuthorizationEndpoint string                 `json:"authorizationEndpoint"`
	Params                map[string]interface{} `json:"params"`
	Alert                 string                 `json:"alert"`
}

type ResponseList struct {
	Data  any   `json:"data" yaml:"data"`
	Total int64 `json:"total" yaml:"total"`
}

func ResponseSuccess(resp *restful.Response, info interface{}) {
	resp.WriteAsJson(info)

}
func ResponseAuthRedirect(resp *restful.Response, bundle *i18n.Bundle, lang, message, authorizationEndpoint string,
	params map[string]interface{}, ctx context.Context) {
	resp.WriteHeader(http.StatusUnauthorized)
	var body AuthRedirectInfo
	body.Message = message
	body.Params = params
	body.AuthorizationEndpoint = authorizationEndpoint
	body.Alert, _ = GetLocaleMessage(bundle, nil, lang, "statusUnauthorized")
	_ = resp.WriteAsJson(body)
}
func ResponseErrorMessage(resp *restful.Response, bundle *i18n.Bundle, code int, lang, message, detail string, ctx context.Context) {
	resp.WriteHeader(code)
	var body ResponseError
	body.Message = message
	body.Detail = detail
	switch code {
	case http.StatusUnauthorized:
		body.Alert, _ = GetLocaleMessage(bundle, nil, lang, "statusUnauthorized")
	case http.StatusBadRequest:
		body.Alert, _ = GetLocaleMessage(bundle, nil, lang, "statusBadRequest")
	case http.StatusForbidden:
		body.Alert, _ = GetLocaleMessage(bundle, nil, lang, "statusForbidden")
	case http.StatusInternalServerError:
		body.Alert, _ = GetLocaleMessage(bundle, nil, lang, "statusInternalServerError")

	}
	_ = resp.WriteAsJson(body)
}

//RequestQuery paramType: string,number queryType: eq,like
func RequestQuery(name, paramType, queryType string, req *restful.Request, queryParam *QueryParam) {
	value := req.QueryParameter(name)
	nv := strings.TrimSpace(value)
	if nv != "" {
		//相等可以是字符或者数字
		if queryType == "" || queryType == QueryTypeEqual {
			if paramType == ParamTypeNumber {
				v := StringsToUint(nv)
				if v > 0 {
					if queryParam.WhereQuery == "" {
						queryParam.WhereQuery = fmt.Sprintf(" %s = ? ", CamelString2Snake(name))
					} else {
						queryParam.WhereQuery += fmt.Sprintf(" AND %s = ? ", CamelString2Snake(name))
					}
					queryParam.WhereArgs = append(queryParam.WhereArgs, v)
				}
			} else if paramType == ParamTypeString || paramType == "" {
				if queryParam.WhereQuery == "" {
					queryParam.WhereQuery = fmt.Sprintf(" %s = ? ", CamelString2Snake(name))
				} else {
					queryParam.WhereQuery += fmt.Sprintf(" AND %s = ? ", CamelString2Snake(name))
				}
				queryParam.WhereArgs = append(queryParam.WhereArgs, nv)
			} else if paramType == ParamTypeBool {
				if queryParam.WhereQuery == "" {
					queryParam.WhereQuery = fmt.Sprintf(" %s = ? ", CamelString2Snake(name))
				} else {
					queryParam.WhereQuery += fmt.Sprintf(" AND %s = ? ", CamelString2Snake(name))
				}
				if strings.ToUpper(nv) != "0" || strings.ToUpper(nv) != "f" || strings.ToUpper(nv) != "false" {
					queryParam.WhereArgs = append(queryParam.WhereArgs, 0)
				} else {
					queryParam.WhereArgs = append(queryParam.WhereArgs, 1)
				}
			}
			// like 只能为字符串
		} else if queryType == QueryTypeLike {
			if queryParam.WhereQuery == "" {
				queryParam.WhereQuery = fmt.Sprintf(" %s LIKE  ? ", CamelString2Snake(name))
			} else {
				queryParam.WhereQuery += fmt.Sprintf(" AND %s LIKE  ? ", CamelString2Snake(name))
			}
			queryParam.WhereArgs = append(queryParam.WhereArgs, fmt.Sprintf("%%%s%%", nv))

		} else if queryType == QueryTypeIn {
			valueSlice := req.QueryParameters(name)
			if paramType == ParamTypeStringSlice {
				if queryParam.WhereQuery == "" {
					queryParam.WhereQuery = fmt.Sprintf(" %s IN (?) ", CamelString2Snake(name))
				} else {
					queryParam.WhereQuery += fmt.Sprintf(" AND %s IN (?)", CamelString2Snake(name))
				}
				queryParam.WhereArgs = append(queryParam.WhereArgs, valueSlice)
			} else if paramType == ParamTypeNumber {
				if queryParam.WhereQuery == "" {
					queryParam.WhereQuery = fmt.Sprintf(" %s IN (?) ", CamelString2Snake(name))
				} else {
					queryParam.WhereQuery += fmt.Sprintf(" AND %s IN (?)", CamelString2Snake(name))
				}
				queryParam.WhereArgs = append(queryParam.WhereArgs, StringsToUints(valueSlice))
			}
		}
	}
}

//RequestQuerySearch queryType: eq,like use 'or' to connect
func RequestQuerySearch(value, queryType string, fields []string, queryParam *QueryParam) {
	if len(value) == 0 || len(fields) == 0 {
		return
	}
	//相等可以是字符或者数字
	if queryType == "" || queryType == QueryTypeEqual {
		for _, name := range fields {
			if queryParam.WhereQuery == "" {
				queryParam.WhereQuery = fmt.Sprintf(" %s = ? ", CamelString2Snake(name))
			} else {
				queryParam.WhereQuery += fmt.Sprintf(" OR %s = ? ", CamelString2Snake(name))
			}
			queryParam.WhereArgs = append(queryParam.WhereArgs, value)
		}
		// like 只能为字符串
	} else if queryType == QueryTypeLike {
		for _, name := range fields {
			if queryParam.WhereQuery == "" {
				queryParam.WhereQuery = fmt.Sprintf(" %s LIKE  ? ", CamelString2Snake(name))
			} else {
				queryParam.WhereQuery += fmt.Sprintf(" OR %s LIKE  ? ", CamelString2Snake(name))
			}
			queryParam.WhereArgs = append(queryParam.WhereArgs, fmt.Sprintf("%%%s%%", value))
		}
	}
}

func RequestQueryEqual(name, paramType, queryType string, value string, queryParam *QueryParam) {
	if len(strings.TrimSpace(value)) == 0 || len(strings.TrimSpace(name)) == 0 {
		return
	}
	if paramType == ParamTypeNumber {
		v := StringsToUint(value)
		if v > 0 {
			if queryParam.WhereQuery == "" {
				queryParam.WhereQuery = fmt.Sprintf(" %s = ? ", CamelString2Snake(name))
			} else {
				queryParam.WhereQuery += fmt.Sprintf(" AND %s = ? ", CamelString2Snake(name))
			}
			queryParam.WhereArgs = append(queryParam.WhereArgs, v)
		}
	} else if paramType == ParamTypeString || paramType == "" {
		if queryParam.WhereQuery == "" {
			queryParam.WhereQuery = fmt.Sprintf(" %s = ? ", CamelString2Snake(name))
		} else {
			queryParam.WhereQuery += fmt.Sprintf(" AND %s = ? ", CamelString2Snake(name))
		}
		queryParam.WhereArgs = append(queryParam.WhereArgs, value)
	}
}

func TokenErr(resp *restful.Response, typ, description string, statusCode int) error {
	data := struct {
		Error       string `json:"error"`
		Description string `json:"error_description,omitempty"`
	}{typ, description}
	resp.ResponseWriter.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(statusCode)
	_ = resp.WriteAsJson(data)
	return nil
}

func Request(method, address string, headers map[string]string, queries map[string]interface{}, body interface{}) (response *http.Response, err error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 10 * time.Second,
	}
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, address, b)
	if err != nil {
		klog.Errorf("create request failed, err: %s", err.Error())
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	queryValues := url.Values{}
	for k, v := range queries {
		queryValues.Add(k, fmt.Sprintf("%v", v))
	}
	req.URL.RawQuery = queryValues.Encode()
	return client.Do(req)
}

func NewHTTPClientWithCA(rootCA string, insecureSkipVerify bool) (client *http.Client, err error) {

	var block *pem.Block
	block, _ = pem.Decode([]byte(rootCA))
	if block == nil {
		err = errors.New("ca decode failed")
		return nil, err
	}
	// Only use PEM "CERTIFICATE" blocks without extra headers
	if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
		err = fmt.Errorf("ca decode failed, block type: %s is not CERTIFICATE", block.Type)
		return nil, err

	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		err = fmt.Errorf("ca decode failed, err: %s", err.Error())
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AddCert(cert)
	// Copied from http.DefaultTransport.
	tlsConfig := tls.Config{RootCAs: pool, InsecureSkipVerify: insecureSkipVerify}
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tlsConfig,
			Proxy:           http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	return client, err
}
