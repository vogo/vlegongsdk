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

package main

import (
	"fmt"

	"github.com/vogo/vlegongsdk/cores"
	"github.com/vogo/vlegongsdk/systems"
	"github.com/vogo/vogo/vlog"
)

func main() {
	// 创建配置
	config := cores.NewConfig(
		"https://api.example.com", // 替换为实际的API地址
		"V00001",                  // 替换为实际的机构编号
		"uptest",                  // 替换为实际的租户编码
		"-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----", // 替换为实际的私钥
		"-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----",           // 替换为实际的平台公钥
	)

	// 创建客户端
    client, err := cores.NewClient(config)
    if err != nil {
        vlog.Fatalf("failed to create client | err: %v", err)
    }

	// 创建系统服务
	systemService := systems.NewSystemService(client)

	// 创建实名认证请求
	req := &systems.IdentityAuthRequest{
		Name:     "张三",
		IDCardNo: "110101199001011234", // 示例身份证号
	}

	// 调用实名认证接口
    resp, err := systemService.IdentityAuth(req)
    if err != nil {
        vlog.Fatalf("real name authentication failed | name: %s | id_card_no: %s | err: %v", req.Name, req.IDCardNo, err)
    }

	// 处理响应
	fmt.Printf("认证状态: %s\n", resp.AuthStatus)
	fmt.Printf("备注信息: %s\n", resp.Message)
	fmt.Printf("请求编号: %s\n", resp.OrderNumber)

	// 根据认证状态进行业务处理
	if resp.AuthStatus == "S" {
		fmt.Println("实名认证成功")
	} else {
		fmt.Println("实名认证失败")
	}
}
