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
	"github.com/vogo/vlegongsdk/cores"
)

// SettlementService 结算服务
type SettlementService struct {
	client *cores.Client
	// 支付回调处理器
	paymentCallbackHandler *PaymentCallbackHandler
}

// NewSettlementService 创建一个新的结算服务
func NewSettlementService(client *cores.Client, paymentCallbackHandler cores.CallbackHandler[PaymentCallbackRequest]) *SettlementService {
	service := &SettlementService{
		client: client,
	}

	// 如果提供了回调处理器，则创建支付回调处理器
	if paymentCallbackHandler != nil {
		service.paymentCallbackHandler = NewPaymentCallbackHandler(client, paymentCallbackHandler)
	}

	return service
}

// GetPaymentCallbackHandler 获取支付回调处理器
func (s *SettlementService) GetPaymentCallbackHandler() *PaymentCallbackHandler {
	return s.paymentCallbackHandler
}
