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

import "time"

type K8sTokenPayload struct {
	Iss                string `json:"iss"`
	Namespace          string `json:"kubernetes.io/serviceaccount/namespace"`
	SecretName         string `json:"kubernetes.io/serviceaccount/secret.name"`
	ServiceAccountName string `json:"kubernetes.io/serviceaccount/service-account.name"`
	ServiceAccountUid  string `json:"kubernetes.io/serviceaccount/service-account.uid"`
	Sub                string `json:"sub"`
}

func (receiver K8sTokenPayload) Valid() error {
	return nil
}

type ApplicationInfo struct {
	Application    string            `json:"application"`
	GoVersion      string            `json:"goVersion"`
	Commit         string            `json:"commit"`
	BuildDate      string            `json:"buildDate"`
	KubernetesInfo *K8sVersion       `json:"kubernetesInfo,omitempty"`
	OS             string            `json:"os"`
	Arch           string            `json:"arch"`
	CpuCores       int               `json:"cpuCores"`
	PhysicalInfo   *PhysicalInfo     `json:"physicalInfo,omitempty"`
	Alert          string            `json:"alert,omitempty"`
	Error          string            `json:"error,omitempty"`
	Time           time.Time         `json:"time"`
	Data           string            `json:"data"`
	Extend         map[string]string `json:"extend,omitempty"`
	Developer      string            `json:"developer,omitempty"` //
	MachineID      string            `json:"machineId,omitempty"` //
}
type MachineInformation struct {
	OS         string          `json:"os"`
	Arch       string          `json:"arch"`
	CpuCores   int             `json:"cpuCores"`
	Kubernetes *KubernetesInfo `json:"kubernetes"`
	Physical   *PhysicalInfo   `json:"physical"`
}

type PhysicalInfo struct {
	MachineID  string `json:"machineId"`
	ServerPort string `json:"serverPort"`
}
type KubernetesInfo struct {
	CA        string      `json:"ca"`
	Namespace string      `json:"namespace"`
	Server    string      `json:"server"`
	Port      string      `json:"port"`
	Version   *K8sVersion `json:"version"`
}
type K8sVersion struct {
	Namespace    string    `json:"namespace"`
	Major        string    `json:"major"`
	Minor        string    `json:"minor"`
	GitVersion   string    `json:"gitVersion"`
	GitCommit    string    `json:"gitCommit"`
	GitTreeState string    `json:"gitTreeState"`
	BuildDate    time.Time `json:"buildDate"`
	GoVersion    string    `json:"goVersion"`
	Compiler     string    `json:"compiler"`
	Platform     string    `json:"platform"`
}
