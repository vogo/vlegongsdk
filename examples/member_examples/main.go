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
	"time"

	"github.com/vogo/vlegongsdk/cores"
	"github.com/vogo/vlegongsdk/members"
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
		vlog.Fatalf("创建客户端失败: %v", err)
	}

	// 创建成员服务
	memberService := members.NewMemberService(client)

	// 示例：注册用工人员
	registerExample(memberService)

	// 示例：查询用工人员信息
	getInfoExample(memberService)

	// 示例：绑定银行卡
	addBankCardExample(memberService)

	// 示例：解绑银行卡
	unbindBankCardExample(memberService)

	// 示例：采集用工人员身份证
	idCardAuthExample(memberService)

	// 示例：注销用工人员
	dismissExample(memberService)
}

// 注册用工人员示例
func registerExample(memberService *members.MemberService) {
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
		vlog.Infof("注册用工人员失败: %v\n", err)
		return
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

// 查询用工人员信息示例
func getInfoExample(memberService *members.MemberService) {
	// 创建查询用工人员信息请求
	req := &members.InfoRequest{
		CompanyCode:  "COMPANY001", // 替换为实际的企业编码
		FreelancerID: "12345",      // 替换为实际的用工人员编号
	}

	// 发送请求
	resp, err := memberService.GetInfo(req)
	if err != nil {
		vlog.Infof("查询用工人员信息失败: %v\n", err)
		return
	}

	// 打印响应
	fmt.Printf("用工人员编号: %s\n", resp.FreelancerID)
	fmt.Printf("用工人员状态: %s\n", resp.Status)
	fmt.Printf("导入时间: %s\n", resp.CreateTime)
	fmt.Printf("用工人员姓名: %s\n", resp.Name)
	fmt.Printf("实名校验状态: %s\n", resp.ValidateStatus)
	if resp.Remark != "" {
		fmt.Printf("备注信息: %s\n", resp.Remark)
	}
}

// 绑定银行卡示例
func addBankCardExample(memberService *members.MemberService) {
	// 创建绑定银行卡请求
	req := &members.AddBankCardRequest{
		CompanyCode: "COMPANY001",          // 替换为实际的企业编码
		IDCardNo:    "110101199001011234",  // 替换为实际的身份证号
		BankCardNo:  "6222021234567890123", // 替换为实际的银行卡号
		Bank:        "中国工商银行",              // 替换为实际的开户行名称
	}

	// 发送请求
	err := memberService.AddBankCard(req)
	if err != nil {
		vlog.Infof("绑定银行卡失败: %v\n", err)
		return
	}

	fmt.Println("绑定银行卡成功")
}

// 解绑银行卡示例
func unbindBankCardExample(memberService *members.MemberService) {
	// 创建解绑银行卡请求
	req := &members.UnbindBankCardRequest{
		CompanyCode: "COMPANY001",          // 替换为实际的企业编码
		IDCardNo:    "110101199001011234",  // 替换为实际的身份证号
		BankCardNo:  "6222021234567890123", // 替换为实际的银行卡号
	}

	// 发送请求
	err := memberService.UnbindBankCard(req)
	if err != nil {
		vlog.Infof("解绑银行卡失败: %v\n", err)
		return
	}

	fmt.Println("解绑银行卡成功")
}

// 采集用工人员身份证示例
func idCardAuthExample(memberService *members.MemberService) {
	// 创建采集用工人员身份证请求
	// 注意：需要先通过文件上传接口上传身份证正反面照片，获取图片ID
	req := &members.IDCardAuthRequest{
		FreelancerID: "12345",        // 替换为实际的用工人员编号
		FrontImgID:   "front_img_id", // 替换为实际的身份证人像面照片ID
		BackImgID:    "back_img_id",  // 替换为实际的身份证国徽面照片ID
	}

	// 发送请求
	err := memberService.IDCardAuth(req)
	if err != nil {
		vlog.Infof("采集用工人员身份证失败: %v\n", err)
		return
	}

	fmt.Println("采集用工人员身份证成功")
}

// 注销用工人员示例
func dismissExample(memberService *members.MemberService) {
	// 创建注销用工人员请求
	req := &members.DismissRequest{
		CompanyCode:  "COMPANY001", // 替换为实际的企业编码
		FreelancerID: "12345",      // 替换为实际的用工人员编号
	}

	// 发送请求
	resp, err := memberService.Dismiss(req)
	if err != nil {
		vlog.Infof("注销用工人员失败: %v\n", err)
		return
	}

	// 打印响应
	fmt.Printf("用工人员编号: %s\n", resp.FreelancerID)
	if resp.IsDismissed == "1" {
		fmt.Println("注销成功")
	} else {
		fmt.Printf("注销失败，原因: %s\n", resp.Remark)
	}
}
