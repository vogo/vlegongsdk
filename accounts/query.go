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
	"github.com/vogo/vogo/vstrconv"
)

// BalanceQueryRequest 余额查询请求
type BalanceQueryRequest struct {
	CompanyCode       string `json:"companyCode"`       // 企业编码
	CrowdsourcingCode string `json:"crowdsourcingCode"` // 众包企业编码
}

// BalancAccountData 企业账户信息
type BalancAccountData struct {
	CrowdsourcingCode string `json:"crowdsourcingCode"` // 众包编号
	CrowdsourcingName string `json:"crowdsourcingName"` // 众包企业名称
	AccountNo         string `json:"accountNo"`         // 企业账户编号
	AccountName       string `json:"accountName"`       // 企业账户名称
	TotalBalance      string `json:"totalBalance"`      // 企业账户余额（元）
	Balance           string `json:"balance"`           // 企业可用余额（元）
	FrozenAmount      string `json:"frozenAmount"`      // 企业冻结余额（元）
}

// BalanceQueryResponse 余额查询响应
type BalanceQueryResponse struct {
	AccountList []BalancAccountData `json:"accountList"` // 企业账户列表
}

type BalancAccount struct {
	CrowdsourcingCode string  `json:"crowdsourcingCode"` // 众包编号
	CrowdsourcingName string  `json:"crowdsourcingName"` // 众包企业名称
	AccountNo         string  `json:"accountNo"`         // 企业账户编号
	AccountName       string  `json:"accountName"`       // 企业账户名称
	TotalBalance      float64 `json:"totalBalance"`      // 企业账户余额（元）
	Balance           float64 `json:"balance"`           // 企业可用余额（元）
	FrozenAmount      float64 `json:"frozenAmount"`      // 企业冻结余额（元）
}

func DataToAccount(data BalancAccountData) *BalancAccount {
	return &BalancAccount{
		CrowdsourcingCode: data.CrowdsourcingCode,
		CrowdsourcingName: data.CrowdsourcingName,
		AccountNo:         data.AccountNo,
		AccountName:       data.AccountName,
		TotalBalance:      vstrconv.Float64(data.TotalBalance),
		Balance:           vstrconv.Float64(data.Balance),
		FrozenAmount:      vstrconv.Float64(data.FrozenAmount),
	}
}

// BalanceQuery 查询企业余额
func (s *AccountService) BalanceQuery(req *BalanceQueryRequest) ([]*BalancAccount, error) {
	var resp BalanceQueryResponse
	err := s.client.DoRequest("/settlement/accountApi/balanceQuery", req, &resp)
	if err != nil {
		return nil, err
	}

	accounts := make([]*BalancAccount, 0, len(resp.AccountList))
	for _, account := range resp.AccountList {
		accounts = append(accounts, DataToAccount(account))
	}

	return accounts, nil
}
