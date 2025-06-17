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

	"github.com/google/uuid"
	"github.com/vogo/vlegongsdk/examples"
	"github.com/vogo/vlegongsdk/settlements"
	"github.com/vogo/vogo/vlog"
	"github.com/vogo/vogo/vos"
)

func main() {
	client := examples.LoadClient()

	// 创建结算服务
	settlementService := settlements.NewSettlementService(client, nil)

	// 示例4：发起支付
	// payExample(settlementService)

	// 示例5：查询订单
	queryOrderExample(settlementService)
}

// 发起支付示例
func payExample(service *settlements.SettlementService) {
	// 生成唯一的外部订单号
	outOrderNo := fmt.Sprintf("O%s", uuid.New().String()[:12])

	// 创建支付请求
	amount := 100.00 // 支付金额100元

	req := &settlements.PayRequest{
		AccountNo:   vos.EnsureEnvString("LEGONG_BANK_CARD_NO"), // 示例银行卡号
		Amount:      amount,
		IDCardNo:    vos.EnsureEnvString("LEGONG_FREELANCER_IDCARD"), // 示例身份证号
		Name:        vos.EnsureEnvString("LEGONG_FREELANCER_NAME"),
		NotifyURL:   vos.EnsureEnvString("LEGONG_NOTIFY_URL"),
		OutOrderNo:  outOrderNo,
		PayChannel:  settlements.PayChannelBankCard,
		ProjectCode: vos.EnsureEnvString("LEGONG_PROJECT_CODE"),
		Remark:      "测试支付",
	}

	// 调用支付接口
	resp, err := service.Pay(req)
	if err != nil {
		vlog.Error(err)
		return
	}

	// 处理响应
	fmt.Printf("平台订单号: %s\n", resp.OrderNo)
	fmt.Printf("外部订单号: %s\n", resp.OutOrderNo)
	fmt.Printf("订单状态: %s (%s)\n", resp.Status, settlements.GetOrderStatusDesc(resp.Status))

	// 打印可选字段（如果存在）
	if resp.Amount != 0 {
		fmt.Printf("实际支付金额: %.2f元\n", resp.Amount)
	}
	if resp.Tax != 0 {
		fmt.Printf("税费总计: %.2f元\n", resp.Tax)
	}
	if resp.IncomeTax != 0 {
		fmt.Printf("个税: %.2f元\n", resp.IncomeTax)
	}
	if resp.ValueAddedTax != 0 {
		fmt.Printf("增值税: %.2f元\n", resp.ValueAddedTax)
	}
	if resp.AdditionalTax != 0 {
		fmt.Printf("附加税: %.2f元\n", resp.AdditionalTax)
	}
	if resp.TotalAmount != 0 {
		fmt.Printf("订单总金额: %.2f元\n", resp.TotalAmount)
	}
	if resp.ServiceChargeBearWay > 0 {
		fmt.Printf("服务费承担方: %d (%s)\n", resp.ServiceChargeBearWay, settlements.GetServiceChargeBearWayDesc(resp.ServiceChargeBearWay))
	}
	if resp.IncomeTaxBearWay > 0 {
		fmt.Printf("个税承担方: %d (%s)\n", resp.IncomeTaxBearWay, settlements.GetIncomeTaxBearWayDesc(resp.IncomeTaxBearWay))
	}
	if resp.ReceivedAmount != 0 {
		fmt.Printf("到账金额: %.2f元\n", resp.ReceivedAmount)
	}
	if resp.Desc != "" {
		fmt.Printf("状态描述: %s\n", resp.Desc)
	}
	if resp.PlatformServiceCharge != 0 {
		fmt.Printf("平台服务费: %.2f元\n", resp.PlatformServiceCharge)
	}
	if resp.PlatformServiceChargeBearWay > 0 {
		fmt.Printf("平台服务费承担方: %d (%s)\n", resp.PlatformServiceChargeBearWay, settlements.GetServiceChargeBearWayDesc(resp.PlatformServiceChargeBearWay))
	}
	if resp.SuccessTime != "" {
		fmt.Printf("支付时间: %s\n", resp.SuccessTime)
	}
}

// 查询订单示例
func queryOrderExample(service *settlements.SettlementService) {
	// 创建订单查询请求
	req := &settlements.QueryOrderRequest{
		OutOrderNo: vos.EnsureEnvString("LEGONG_OUT_ORDER_NO"), // 替换为实际的外部订单号
	}

	// 调用订单查询接口
	resp, err := service.QueryOrder(req)
	if err != nil {
		vlog.Infof("查询订单失败: %v\n", err)
		return
	}

	// 处理响应
	fmt.Printf("平台订单号: %s\n", resp.OrderNo)
	fmt.Printf("外部订单号: %s\n", resp.OutOrderNo)
	fmt.Printf("订单状态: %s (%s)\n", resp.Status, settlements.GetOrderStatusDesc(resp.Status))
	fmt.Printf("实际支付金额: %.2f元\n", resp.Amount)
	fmt.Printf("税费总计: %.2f元\n", resp.Tax)

	// 打印可选字段（如果存在）
	if resp.IncomeTax != 0 {
		fmt.Printf("个税: %.2f元\n", resp.IncomeTax)
	}
	if resp.ValueAddedTax != 0 {
		fmt.Printf("增值税: %.2f元\n", resp.ValueAddedTax)
	}
	if resp.AdditionalTax != 0 {
		fmt.Printf("附加税: %.2f元\n", resp.AdditionalTax)
	}

	fmt.Printf("订单总金额: %.2f元\n", resp.TotalAmount)
	fmt.Printf("服务费承担方: %d (%s)\n", resp.ServiceChargeBearWay, settlements.GetServiceChargeBearWayDesc(resp.ServiceChargeBearWay))
	fmt.Printf("个税承担方: %d (%s)\n", resp.IncomeTaxBearWay, settlements.GetIncomeTaxBearWayDesc(resp.IncomeTaxBearWay))
	fmt.Printf("到账金额: %.2f元\n", resp.ReceivedAmount)

	if resp.Desc != "" {
		fmt.Printf("状态描述: %s\n", resp.Desc)
	}

	if resp.PlatformServiceCharge != 0 {
		fmt.Printf("平台服务费: %.2f元\n", resp.PlatformServiceCharge)
	}

	if resp.PlatformServiceChargeBearWay > 0 {
		fmt.Printf("平台服务费承担方: %d (%s)\n", resp.PlatformServiceChargeBearWay, settlements.GetServiceChargeBearWayDesc(resp.PlatformServiceChargeBearWay))
	}

	fmt.Printf("支付时间: %s\n", resp.SuccessTime)
}
