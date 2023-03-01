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

var WhiteList = map[string]string{}

func init() {
	WhiteList = make(map[string]string)
	WhiteList["3171BCDA-B314-5D58-9B7A-5A791CA9EFD1"] = "cloudy"
	WhiteList["63a9b468c7434f0d9035285aa0d43f2b"] = "aliyun"
	WhiteList["290B530A-97EF-56CE-A0F8-991B5EF4CFBD"] = "wenxiang"
}
