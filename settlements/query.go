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

// Package settlements 提供结算相关功能
package settlements

// QueryOrderRequest 订单查询请求
type QueryOrderRequest struct {
	OutOrderNo string `json:"outOrderNo"` // 外部订单号
}

// QueryOrder 查询订单
func (s *SettlementService) QueryOrder(req *QueryOrderRequest) (*Order, error) {
	var resp Order
	err := s.client.DoRequest("/settlement/settleApi/query", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
