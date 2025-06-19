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
	"net/http"

	"github.com/vogo/vlegongsdk/cores"
	"github.com/vogo/vogo/vlog"
)

// PaymentCallbackRequest 支付结果通知回调请求
type PaymentCallbackRequest struct {
	OrderNo                      string  `json:"orderNo"`                                // 平台内部订单号
	OutOrderNo                   string  `json:"outOrderNo"`                             // 外部订单号
	Status                       string  `json:"status"`                                 // 支付状态 S:支付成功 F:支付失败 P:提交成功，订单处理中
	Amount                       float64 `json:"amount"`                                 // 实际支付金额(元)
	Tax                          float64 `json:"tax"`                                    // 税费总计(元)
	IncomeTax                    float64 `json:"inComeTax,omitempty"`                    // 个税(元)
	ValueAddedTax                float64 `json:"valueAddedTax,omitempty"`                // 增值税(元)
	AdditionalTax                float64 `json:"additionalTax,omitempty"`                // 附加税(元)
	TotalAmount                  float64 `json:"totalAmount"`                            // 订单总金额(元)
	ServiceChargeBearWay         int     `json:"serviceChargeBearWay"`                   // 服务费承担方, 1:企业承担,2:个人承担
	IncomeTaxBearWay             int     `json:"incomeTaxBearWay"`                       // 个税承担方, 1:企业承担,2:个人承担
	ReceivedAmount               float64 `json:"receivedAmount"`                         // 到账金额(元)
	Desc                         string  `json:"desc,omitempty"`                         // 状态描述
	SuccessTime                  string  `json:"successTime,omitempty"`                  // 支付完成时间，格式：yyyyMMddHHmmss，支付成功返回
	PayChannel                   int     `json:"payChannel"`                             // 支付渠道 1:银行卡 2:微信支付 3:支付宝
	PlatformServiceCharge        float64 `json:"platformServiceCharge,omitempty"`        // 平台服务费(元)
	PlatformServiceChargeBearWay int     `json:"platformServiceChargeBearWay,omitempty"` // 平台服务费承担方, 1:企业承担,2:个人承担
}

// PaymentCallbackHandler 支付回调处理器
type PaymentCallbackHandler struct {
	client  *cores.Client
	handler cores.CallbackHandler[*PaymentCallbackRequest]
}

// NewPaymentCallbackHandler 创建支付回调处理器
func NewPaymentCallbackHandler(client *cores.Client, handler cores.CallbackHandler[*PaymentCallbackRequest]) *PaymentCallbackHandler {
	return &PaymentCallbackHandler{
		client:  client,
		handler: handler,
	}
}

// ServeHTTP 实现http.Handler接口，处理支付结果通知回调
func (h *PaymentCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 只接受POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析回调请求
	var callbackReq PaymentCallbackRequest
	if err := cores.DecodeCallbackRequest(h.client, r, &callbackReq); err != nil {
		vlog.Errorf("解析支付回调请求失败: %v", err)
		http.Error(w, fmt.Sprintf("解析支付回调请求失败: %v", err), http.StatusBadRequest)
		return
	}

	// 处理回调请求
	if err := h.handler.HandleCallback(&callbackReq); err != nil {
		vlog.Errorf("处理支付回调请求失败: %v", err)
		http.Error(w, fmt.Sprintf("处理支付回调请求失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	cores.WriteCallbackResponse(w)
}
