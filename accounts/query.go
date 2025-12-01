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

package accounts

import (
	"fmt"
)

// BalanceQueryRequest 余额查询请求
type BalanceQueryRequest struct {
	CompanyCode       string `json:"companyCode"`       // 企业编码
	CrowdsourcingCode string `json:"crowdsourcingCode"` // 众包企业编码
}

// BalanceQueryResponse 余额查询响应
type BalanceQueryResponse struct {
	AccountList []Account `json:"accountList"` // 企业账户列表
}

// BalanceQuery 查询企业余额
func (s *AccountService) BalanceQuery(req *BalanceQueryRequest) (*BalanceQueryResponse, error) {
	var resp BalanceQueryResponse
	err := s.client.DoRequest("/settlement/accountApi/balanceQuery", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("查询企业余额失败: %w", err)
	}

	return &resp, nil
}
