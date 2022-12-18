// Copyright 2019 Inspur Technologies Co.,Ltd.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use
// this file except in compliance with the License.  You may obtain a copy of the
// License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations under the License.

package oss

import (
	"strings"
)

// Refresh refreshes ak, sk and securityToken for OSSClient.
func (OSSClient OSSClient) Refresh(ak, sk, securityToken string) {
	for _, sp := range OSSClient.conf.securityProviders {
		if bsp, ok := sp.(*BasicSecurityProvider); ok {
			bsp.refresh(strings.TrimSpace(ak), strings.TrimSpace(sk), strings.TrimSpace(securityToken))
			break
		}
	}
}

func (OSSClient OSSClient) getSecurity() securityHolder {
	if OSSClient.conf.securityProviders != nil {
		for _, sp := range OSSClient.conf.securityProviders {
			if sp == nil {
				continue
			}
			sh := sp.getSecurity()
			if sh.ak != "" && sh.sk != "" {
				return sh
			}
		}
	}
	return emptySecurityHolder
}

// Close closes OSSClient.
func (OSSClient *OSSClient) Close() {
	OSSClient.httpClient = nil
	OSSClient.conf.transport.CloseIdleConnections()
	OSSClient.conf = nil
}
