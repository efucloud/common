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

package license

import (
	"encoding/json"
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"github.com/efucloud/common"
	"github.com/golang-jwt/jwt/v4"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"
)

const (
	k8sPath              = "/var/run/secrets/kubernetes.io/serviceaccount"
	kubernetesServerAddr = "KUBERNETES_PORT_443_TCP_ADDR"
	kubernetesServerPort = "KUBERNETES_SERVICE_PORT"
	dockerEnv            = "/.dockerenv"
)

// GetMachineInformation 根据部署来生成机器信息
func GetMachineInformation(appName string) (applicationInfo common.ApplicationInfo) {
	var (
		info common.MachineInformation
		ca   []byte
		err  error
	)
	info.OS = runtime.GOOS
	info.Arch = runtime.GOARCH
	info.CpuCores = runtime.GOMAXPROCS(0)
	applicationInfo.Application = appName
	applicationInfo.OS = runtime.GOOS
	applicationInfo.Arch = runtime.GOARCH
	applicationInfo.CpuCores = runtime.GOMAXPROCS(0)
	applicationInfo.Time = time.Now().Local()

	// 判断是否在k8s集群中运行
	ca, err = os.ReadFile(path.Join(k8sPath, "ca.crt"))
	if err == nil {
		applicationInfo.KubernetesInfo = new(common.K8sVersion)
		info.Kubernetes = new(common.KubernetesInfo)
		info.Kubernetes.CA = string(ca)
		if ns, err := os.ReadFile(path.Join(k8sPath, "namespace")); err == nil {
			info.Kubernetes.Namespace = string(ns)
			applicationInfo.KubernetesInfo.Namespace = info.Kubernetes.Namespace
		} else {
			applicationInfo.Error = err.Error()
			return
		}
		var k8sTokenPayload common.K8sTokenPayload
		if token, err := os.ReadFile(path.Join(k8sPath, "token")); err == nil {
			if t, err := jwt.ParseWithClaims(string(token), &k8sTokenPayload, nil); err == nil {
				if payload, ok := t.Claims.(*common.K8sTokenPayload); ok {
					if payload.Namespace != applicationInfo.KubernetesInfo.Namespace {
						applicationInfo.Error = "namespace is not right"
						return
					} else if payload.Iss != "kubernetes/serviceaccount" {
						applicationInfo.Error = "application is nor run in kubernetes"
						return
					}
				} else {
					applicationInfo.Error = "application is nor run in kubernetes"
					return
				}
			} else {
				applicationInfo.Error = err.Error()
				return
			}
		} else {
			applicationInfo.Error = err.Error()
			return
		}
		info.Kubernetes.Server = os.Getenv(kubernetesServerAddr)
		info.Kubernetes.Port = os.Getenv(kubernetesServerPort)
		//获取k8s版本信息
		verAddr := fmt.Sprintf("https://%s:%s/version", info.Kubernetes.Server, info.Kubernetes.Port)
		if response, err := common.Request(http.MethodGet, verAddr, nil, nil, nil); err == nil {
			if response.StatusCode == http.StatusOK {
				err = json.NewDecoder(response.Body).Decode(info.Kubernetes.Version)
				if err != nil {
					applicationInfo.Error = err.Error()
					return
				}
			}
		} else {
			applicationInfo.Error = err.Error()
			return
		}
	} else {
		klog.Infof("current run system is: %s", runtime.GOOS)
		// 只判断为linux时判断是否docker运行
		if runtime.GOOS == "linux" {
			//只要是linux就认为是容器内部
			if common.PathExists(dockerEnv) {
				applicationInfo.Error = fmt.Sprintf("application not support running in docker")
			}
		} else {
			info.Physical = new(common.PhysicalInfo)
			applicationInfo.PhysicalInfo = new(common.PhysicalInfo)
			info.Physical.MachineID, err = machineid.ProtectedID(appName)
			applicationInfo.PhysicalInfo = info.Physical
			if err == nil {
				mid, _ := machineid.ID()
				if na, ok := common.WhiteList[mid]; ok {
					applicationInfo.Extend = make(map[string]string)
					applicationInfo.Extend[na] = mid
				}
			}
		}
	}

	return
}
