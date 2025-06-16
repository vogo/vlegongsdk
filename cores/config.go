/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cores

// Config 配置对象，用于初始化客户端
type Config struct {
	// BaseURL 共享经济服务平台的基础URL
	BaseURL string
	// OrgCode 机构编号，平台分配给接入方的唯一标识号
	OrgCode string
	// TenantCode 租户编码
	TenantCode string
	// PrivateKey 接入方私钥，用于签名
	PrivateKey string
	// PlatformPublicKey 平台公钥，用于验签和加密
	PlatformPublicKey string
	// Timeout 请求超时时间，单位：秒
	Timeout int
}

// NewConfig 创建一个新的配置对象
func NewConfig(baseURL, orgCode, tenantCode, privateKey, platformPublicKey string) *Config {
	return &Config{
		BaseURL:           baseURL,
		OrgCode:           orgCode,
		TenantCode:        tenantCode,
		PrivateKey:        privateKey,
		PlatformPublicKey: platformPublicKey,
		Timeout:           30, // 默认30秒超时
	}
}
