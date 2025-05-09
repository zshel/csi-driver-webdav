/*
Copyright 2023.

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

package webdav

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"
)

func ParseEndpoint(ep string) (string, string, error) {
	if strings.HasPrefix(strings.ToLower(ep), "unix://") || strings.HasPrefix(strings.ToLower(ep), "tcp://") {
		s := strings.SplitN(ep, "://", 2)
		if s[1] != "" {
			return s[0], s[1], nil
		}
	}
	return "", "", fmt.Errorf("invalid endpoint: %v", ep)
}

func getLogLevel(method string) int32 {
	if method == "/csi.v1.Identity/Probe" ||
		method == "/csi.v1.Node/NodeGetCapabilities" ||
		method == "/csi.v1.Node/NodeGetVolumeStats" {
		return 8
	}
	return 2
}

func logGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	level := klog.Level(getLogLevel(info.FullMethod))
	klog.V(level).Infof("GRPC call: %s", info.FullMethod)
	klog.V(level).Infof("GRPC request: %s", req)

	resp, err := handler(ctx, req)
	if err != nil {
		klog.Errorf("GRPC error: %v", err)
	} else {
		klog.V(level).Infof("GRPC response: %s", resp)
	}
	return resp, err
}

func NewControllerServiceCapability(cap csi.ControllerServiceCapability_RPC_Type) *csi.ControllerServiceCapability {
	return &csi.ControllerServiceCapability{
		Type: &csi.ControllerServiceCapability_Rpc{
			Rpc: &csi.ControllerServiceCapability_RPC{
				Type: cap,
			},
		},
	}
}

func NewNodeServiceCapability(cap csi.NodeServiceCapability_RPC_Type) *csi.NodeServiceCapability {
	return &csi.NodeServiceCapability{
		Type: &csi.NodeServiceCapability_Rpc{
			Rpc: &csi.NodeServiceCapability_RPC{
				Type: cap,
			},
		},
	}
}

func MakeVolumeId(webdavSharePath, volumeName string, parameters map[string]string) string {
	var paramStrings []string
	if parameters != nil && len(parameters) > 0 {
		for k, v := range parameters {
			encoded := base64.RawURLEncoding.EncodeToString([]byte(v))
			paramStrings = append(paramStrings, fmt.Sprintf("%s=%s", k, encoded))
		}
	}
	if len(paramStrings) > 0 {
		return fmt.Sprintf("%s#%s#%s", webdavSharePath, volumeName, strings.Join(paramStrings, ":"))
	}
	return fmt.Sprintf("%s#%s", webdavSharePath, volumeName)
}

func ParseVolumeId(volumeId string) (webdavSharePath, subDir string, parameters map[string]string, err error) {
	arr := strings.Split(volumeId, "#")
	parsedParameters := make(map[string]string)
	if len(arr) < 2 {
		return "", "", nil, errors.New("invalid volumeId")
	}
	if len(arr) > 2 {
		paramPairs := strings.Split(arr[2], ":")
		for _, pair := range paramPairs {
			kv := strings.Split(pair, "=")
			if len(kv) == 2 {
				decoded, err := base64.RawURLEncoding.DecodeString(kv[1])
				if err != nil {
					return "", "", nil, fmt.Errorf("failed to decode parameter value: %v", err)
				}
				parsedParameters[kv[0]] = string(decoded)
			}
		}
	}
	return arr[0], arr[1], parsedParameters, nil
}
