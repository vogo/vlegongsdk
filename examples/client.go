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

package examples

import (
	"github.com/vogo/vlegongsdk/cores"
	"github.com/vogo/vogo/vlog"
	"github.com/vogo/vogo/vos"
)

func LoadClient() *cores.Client {
	// 创建配置
	config := cores.NewConfig(
		vos.EnsureEnvString("LEGONG_API_URL"),     // 替换为实际的API地址
		vos.EnsureEnvString("LEGONG_ORG_CODE"),    // 替换为实际的机构编号
		vos.EnsureEnvString("LEGONG_TENANT_CODE"), // 替换为实际的租户编码
		vos.EnsureEnvString("LEGONG_PRIVATE_KEY"), // 替换为实际的私钥
		vos.EnsureEnvString("LEGONG_PUBLIC_KEY"),  // 替换为实际的平台公钥
	)

	// 创建客户端
	client, err := cores.NewClient(config)
	if err != nil {
		vlog.Fatalf("Failed to create client: %v", err)
	}

	return client
}
