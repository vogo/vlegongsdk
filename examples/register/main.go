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
	"log"
	"time"

	"github.com/vogo/vlegongsdk/cores"
	"github.com/vogo/vlegongsdk/members"
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
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 创建成员服务
	memberService := members.NewMemberService(client)

	// 创建用工人员注册请求
	req := &members.RegisterRequest{
		CompanyCode:    "COMPANY001",                        // 替换为实际的企业编码
		FreelancerName: "张三",                                // 替换为实际的用工人员姓名
		MobilePhone:    "13800138000",                       // 替换为实际的手机号码
		IDCardNo:       "110101199001011234",                // 替换为实际的身份证号
		FreelancerType: "1",                                 // 1:自由职业者, 2:雇员
		CreateTime:     time.Now().Format("20060102150405"), // 当前时间，格式：yyyyMMddHHmmss
	}

	// 发送请求
	resp, err := memberService.Register(req)
	if err != nil {
		log.Fatalf("注册用工人员失败: %v", err)
	}

	// 打印响应
	fmt.Printf("注册成功，用工人员编号: %s\n", resp.FreelancerID)
	fmt.Printf("用工人员状态: %s\n", resp.Status)
	fmt.Printf("导入时间: %s\n", resp.CreateTime)
	fmt.Printf("用工人员姓名: %s\n", resp.Name)
	fmt.Printf("实名校验状态: %s\n", resp.ValidateStatus)
	if resp.Remark != "" {
		fmt.Printf("备注信息: %s\n", resp.Remark)
	}
}
