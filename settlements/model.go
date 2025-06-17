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

package settlements

const (
	PayChannelBankCard = 1 // 银行卡

	// 服务费承担方
	BearWayEnterprise = 1 // 企业承担
	BearWayPersonal   = 2 // 个人承担

	// 订单状态
	OrderStatusProcessing = "P" // 处理中
	OrderStatusSuccess    = "S" // 成功
	OrderStatusFailed     = "F" // 失败
)

// Order 支付订单
type Order struct {
	OrderNo                      string  `json:"orderNo"`                                // 平台订单号
	OutOrderNo                   string  `json:"outOrderNo"`                             // 外部订单号
	Amount                       float64 `json:"amount,omitempty"`                       // 实际支付金额(元)
	Status                       string  `json:"status"`                                 // 订单状态 S:支付成功 F:支付失败 P:提交成功，订单处理中
	Tax                          float64 `json:"tax,omitempty"`                          // 税费总计(元)
	IncomeTax                    float64 `json:"inComeTax,omitempty"`                    // 个税(元)
	ValueAddedTax                float64 `json:"valueAddedTax,omitempty"`                // 增值税(元)
	AdditionalTax                float64 `json:"additionalTax,omitempty"`                // 附加税(元)
	TotalAmount                  float64 `json:"totalAmount,omitempty"`                  // 订单总金额(元)
	ServiceChargeBearWay         int     `json:"serviceChargeBearWay,omitempty"`         // 服务费承担方, 1:企业承担,2:个人承担
	IncomeTaxBearWay             int     `json:"incomeTaxBearWay,omitempty"`             // 个税承担方, 1:企业承担,2:个人承担
	ReceivedAmount               float64 `json:"receivedAmount,omitempty"`               // 到账金额(元)
	Desc                         string  `json:"desc,omitempty"`                         // 状态描述
	PlatformServiceCharge        float64 `json:"platformServiceCharge,omitempty"`        // 平台服务费(元)
	PlatformServiceChargeBearWay int     `json:"platformServiceChargeBearWay,omitempty"` // 平台服务费承担方, 1:企业承担,2:个人承担
	SuccessTime                  string  `json:"successTime,omitempty"`                  // 支付成功时间，格式：yyyy-MM-dd HH:mm:ss
	PayChannel                   int     `json:"payChannel,omitempty"`                   // 支付渠道 0-未知 1-银行卡 2:微信 3:支付宝
}

// GetOrderStatusDesc 获取订单状态描述
func GetOrderStatusDesc(status string) string {
	switch status {
	case OrderStatusProcessing:
		return "处理中"
	case OrderStatusSuccess:
		return "成功"
	case OrderStatusFailed:
		return "失败"
	default:
		return "未知状态"
	}
}

// GetBearWayDesc 获取承担方描述
func GetBearWayDesc(bearWay int) string {
	switch bearWay {
	case BearWayEnterprise:
		return "企业承担"
	case BearWayPersonal:
		return "个人承担"
	default:
		return "未知"
	}
}
