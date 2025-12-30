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

package members

// AddBankCardRequest 绑定银行卡请求
type AddBankCardRequest struct {
	CompanyCode string `json:"companyCode"` // 用工企业编码
	IDCardNo    string `json:"idCardNo"`    // 用工人员身份证号
	BankCardNo  string `json:"bankCardNo"`  // 银行卡号
	Bank        string `json:"bank"`        // 开户行名称，可选
}

// AddBankCard 为用工人员绑定银行卡
func (s *MemberService) AddBankCard(req *AddBankCardRequest) error {
	// 该接口没有返回具体数据，只需要处理可能的错误
	err := s.client.DoRequest("/member/freelancerApi/addBankCard", req, nil)
	if err != nil {
		return err
	}

	return nil
}

// UnbindBankCardRequest 解绑银行卡请求
type UnbindBankCardRequest struct {
	CompanyCode string `json:"companyCode"` // 用工企业编码
	IDCardNo    string `json:"idCardNo"`    // 用工人员身份证号
	BankCardNo  string `json:"bankCardNo"`  // 银行卡号
}

// UnbindBankCard 为用工人员解绑银行卡
func (s *MemberService) UnbindBankCard(req *UnbindBankCardRequest) error {
	// 该接口没有返回具体数据，只需要处理可能的错误
	err := s.client.DoRequest("/member/freelancerApi/unbindBankCard", req, nil)
	if err != nil {
		return err
	}

	return nil
}
