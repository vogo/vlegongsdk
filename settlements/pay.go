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

import (
	"fmt"
)

const (
	PayChannelBankCard = 1 // 银行卡

	// 服务费承担方
	ServiceChargeBearWayEnterprise = 1 // 企业承担
	ServiceChargeBearWayPersonal   = 2 // 个人承担

	// 个税承担方
	IncomeTaxBearWayEnterprise = 1 // 企业承担
	IncomeTaxBearWayPersonal   = 2 // 个人承担

	// 订单状态
	OrderStatusProcessing = "P" // 处理中
	OrderStatusSuccess    = "S" // 成功
	OrderStatusFailed     = "F" // 失败
)

// PayRequest 支付请求
type PayRequest struct {
	AccountNo   string  `json:"accountNo"`   // 银行卡号
	Amount      float64 `json:"amount"`      // 支付金额
	IDCardNo    string  `json:"idCard"`      // 身份证号
	Name        string  `json:"name"`        // 姓名
	NotifyURL   string  `json:"notifyUrl"`   // 通知URL
	OutOrderNo  string  `json:"outOrderNo"`  // 外部订单号
	PayChannel  int     `json:"payChannel"`  // 支付渠道 1-银行卡
	ProjectCode string  `json:"projectCode"` // 项目编码
	Remark      string  `json:"remark"`      // 备注
}

// Pay 发起支付
func (s *SettlementService) Pay(req *PayRequest) (*Order, error) {
	var resp Order
	err := s.client.DoRequest("/settlement/settleApi/pay", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("发起支付失败: %w", err)
	}

	return &resp, nil
}
