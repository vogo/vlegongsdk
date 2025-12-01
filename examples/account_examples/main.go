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

	"github.com/vogo/vlegongsdk/accounts"
	"github.com/vogo/vlegongsdk/examples"
	"github.com/vogo/vogo/vlog"
	"github.com/vogo/vogo/vos"
)

func main() {
	client := examples.LoadClient()

	// 创建账户服务
	accountService := accounts.NewAccountService(client)

	// 查询余额示例
	balanceQueryExample(accountService)
}

// 查询余额示例
func balanceQueryExample(service *accounts.AccountService) {
	// 创建余额查询请求
	req := &accounts.BalanceQueryRequest{
		CompanyCode:       vos.EnsureEnvString("LEGONG_COMPANY_CODE"),       // 企业编码
		CrowdsourcingCode: vos.EnsureEnvString("LEGONG_CROWDSOURCING_CODE"), // 众包企业编码
	}

	// 调用余额查询接口
	resp, err := service.BalanceQuery(req)
	if err != nil {
		vlog.Errorf("查询企业余额失败: %v\n", err)
		return
	}

	// 处理响应
	fmt.Printf("查询到 %d 个账户:\n", len(resp.AccountList))
	for i, account := range resp.AccountList {
		fmt.Printf("\n账户 %d:\n", i+1)
		fmt.Printf("  众包编号: %s\n", account.CrowdsourcingCode)
		fmt.Printf("  众包企业名称: %s\n", account.CrowdsourcingName)
		fmt.Printf("  账户编号: %s\n", account.AccountNo)
		fmt.Printf("  账户名称: %s\n", account.AccountName)
		fmt.Printf("  账户总余额: %s 元\n", account.TotalBalance)
		fmt.Printf("  可用余额: %s 元\n", account.Balance)
		fmt.Printf("  冻结余额: %s 元\n", account.FrozenAmount)
	}
}
