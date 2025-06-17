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

	"github.com/vogo/vlegongsdk/cores"
	"github.com/vogo/vlegongsdk/settlements"
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

	// 创建结算服务
	settlementService := settlements.NewSettlementService(client, nil)

	// 示例1：发起自主签约
	startSignExample(settlementService)

	// 示例2：查询签约状态
	querySignExample(settlementService)

	// 示例3：查询签约记录
	querySignRecordExample(settlementService)
}

// 发起自主签约示例
func startSignExample(service *settlements.SettlementService) {
	// 创建自主签约请求
	req := &settlements.StartSignRequest{
		ProjectCode:  "PROJECT001",
		Name:         "张三",
		IDCardNo:     "110101199001011234", // 示例身份证号
		NoticeType:   "1",                  // 发送短信
		SignPlatform: "1",                  // 网页版短信认证签署
		RedirectURL:  "https://example.com/redirect",
		NotifyURL:    "https://example.com/notify",
	}

	// 调用自主签约接口
	resp, err := service.StartSign(req)
	if err != nil {
		vlog.Infof("发起自主签约失败: %v\n", err)
		return
	}

	// 处理响应
	fmt.Printf("签约流程ID: %s\n", resp.SignFlowID)
	fmt.Printf("签约状态: %d\n", resp.SignStatus)
	fmt.Printf("签署链接短链: %s\n", resp.SignShortURL)
	fmt.Printf("签署有效截止日期: %s\n", resp.SignValidity)
	fmt.Printf("签署描述: %s\n", resp.SignDesc)

	// 根据签约状态进行业务处理
	switch resp.SignStatus {
	case 0:
		fmt.Println("签署中，请等待用户完成签署")
	case 1:
		fmt.Println("已签约，签约成功")
	case 2:
		fmt.Println("拒签，用户拒绝签署")
	case 3:
		fmt.Println("过期，签约已过期")
	case 4:
		fmt.Println("失败，签约失败")
	}
}

// 查询签约状态示例
func querySignExample(service *settlements.SettlementService) {
	// 创建签约查询请求
	req := &settlements.QuerySignRequest{
		SignFlowID: "FLOW123456789", // 替换为实际的签约流程ID
	}

	// 调用签约查询接口
	resp, err := service.QuerySign(req)
	if err != nil {
		vlog.Infof("查询签约状态失败: %v\n", err)
		return
	}

	// 处理响应
	fmt.Printf("签约流程ID: %s\n", resp.SignFlowID)
	fmt.Printf("签约状态: %d\n", resp.SignStatus)
	fmt.Printf("签署完成时间: %s\n", resp.SignEndTime)
	fmt.Printf("签署描述: %s\n", resp.SignDesc)
}

// 查询签约记录示例
func querySignRecordExample(service *settlements.SettlementService) {
	// 创建签约记录查询请求
	req := &settlements.QuerySignRecordRequest{
		IDCardNo:    "110101199001011234", // 示例身份证号
		CompanyCode: "COMPANY001",
		ProjectCode: "PROJECT001",
		SignStatus:  []int{1}, // 查询已签署状态的协议
	}

	// 调用签约记录查询接口
	resp, err := service.QuerySignRecord(req)
	if err != nil {
		vlog.Infof("查询签约记录失败: %v\n", err)
		return
	}

	// 处理响应
	fmt.Printf("找到 %d 条签约记录\n", len(resp.SignRecordList))

	// 遍历签约记录
	for i, record := range resp.SignRecordList {
		fmt.Printf("\n记录 #%d:\n", i+1)
		fmt.Printf("协议名称: %s\n", record.AgreementName)
		fmt.Printf("签约流程ID: %s\n", record.SignFlowID)
		fmt.Printf("签约状态: %d\n", record.SignStatus)
		fmt.Printf("签署描述: %s\n", record.SignDesc)
		fmt.Printf("用工人员姓名: %s\n", record.Name)
		fmt.Printf("用工人员身份证号: %s\n", record.IDCardNo)
		fmt.Printf("税种: %d\n", record.TaxKind)
		fmt.Printf("税种描述: %s\n", record.TaxKindDesc)
		fmt.Printf("签约时间: %s\n", record.SignDate)
		fmt.Printf("协议过期时间: %s\n", record.ExpiryDate)
	}
}

// handleSignCallback 处理签约回调
func handleSignCallback(data interface{}) error {
	// 类型断言
	callbackReq, ok := data.(settlements.SignCallbackRequest)
	if !ok {
		return fmt.Errorf("无效的回调数据类型: %T", data)
	}

	// 处理签约回调
	vlog.Infof("收到签约回调: 流程ID=%s, 状态=%d(%s), 完成时间=%s, 描述=%s",
		callbackReq.SignFlowID,
		callbackReq.SignStatus,
		settlements.GetSignStatusDesc(callbackReq.SignStatus),
		callbackReq.SignEndTime,
		callbackReq.SignDesc,
	)

	// 根据签约状态进行业务处理
	switch callbackReq.SignStatus {
	case 0: // 签署中
		log.Println("签署中，等待用户完成签署")
		// 可以在这里添加业务逻辑
	case 1: // 已签约
		log.Println("签约成功，更新业务状态")
		// 可以在这里添加业务逻辑，如更新数据库中的签约状态
	case 2: // 拒签
		log.Println("用户拒签，需要进行相应处理")
		// 可以在这里添加业务逻辑
	case 3: // 过期
		log.Println("签约已过期，需要重新发起签约")
		// 可以在这里添加业务逻辑
	case 4: // 失败
		log.Println("签约失败，需要排查原因")
		// 可以在这里添加业务逻辑
	default:
		vlog.Infof("未知的签约状态: %d", callbackReq.SignStatus)
	}

	return nil
}
