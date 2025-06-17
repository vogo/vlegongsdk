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

package signs

import "fmt"

// QuerySignRecordRequest 签约记录查询请求
type QuerySignRecordRequest struct {
	IDCardNo    string `json:"idCard"`      // 用工人员身份证号
	CompanyCode string `json:"companyCode"` // 企业编码，用工企业在乐工平台的编码
	ProjectCode string `json:"projectCode"` // 项目编码，企业编码、项目编码最少必填其一
	SignStatus  []int  `json:"signStatus"`  // 签署状态 0 签署中 1 已签署 2 拒签 3 过期 4 失败，默认查询已签署状态的协议
}

// SignRecord 签约记录
type SignRecord struct {
	AgreementName string `json:"agreementName"` // 协议名称
	SignFlowID    string `json:"signFlowId"`    // 签约流程ID
	SignStatus    int    `json:"signStatus"`    // 签约状态 0 签署中 1 已签署 2 拒签 3 过期 4 失败
	SignDesc      string `json:"signDesc"`      // 签署描述
	Name          string `json:"name"`          // 用工人员姓名
	IDCardNo      string `json:"idCard"`        // 用工人员身份证号
	TaxKind       int    `json:"taxKind"`       // 税种 1 经营所得
	TaxKindDesc   string `json:"taxKindDesc"`   // 税种描述
	SignDate      string `json:"signDate"`      // 签约时间，格式：yyyy-MM-dd HH:mm:ss
	ExpiryDate    string `json:"expiryDate"`    // 协议过期时间，格式：yyyy-MM-dd HH:mm:ss
}

// QuerySignRecordResponse 签约记录查询响应
type QuerySignRecordResponse struct {
	SignRecordList []SignRecord `json:"signRecordList"` // 协议列表
}

// QuerySignRecord 查询签约记录
func (s *SignService) QuerySignRecord(req *QuerySignRecordRequest) (*QuerySignRecordResponse, error) {
	var resp QuerySignRecordResponse
	err := s.client.DoRequest("/settlement/signApi/querySignRecord", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("查询签约记录失败: %w", err)
	}

	return &resp, nil
}
