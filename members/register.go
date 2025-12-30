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

// RegisterRequest 用工人员注册请求
type RegisterRequest struct {
	CompanyCode    string `json:"companyCode"`              // 企业编码，用工企业在乐工平台的编码
	FreelancerName string `json:"freelancerName"`           // 用工人员姓名
	MobilePhone    string `json:"mobilePhone"`              // 手机号码
	IDCardNo       string `json:"idCardNo"`                 // 用工人员身份证号
	FreelancerType string `json:"freelancerType,omitempty"` // 用工类型，1:自由职业者, 2:雇员
	CreateTime     string `json:"createTime,omitempty"`     // 创建时间，yyyyMMddHHmmss
}

// RegisterResponse 用工人员注册响应
type RegisterResponse struct {
	FreelancerID   int    `json:"freelancerId"`   // 用工人员编号
	Status         int    `json:"status"`         // 用工人员状态，0：离职，1：待上岗，2：已上岗
	CreateTime     string `json:"createTime"`     // 导入时间，yyyyMMddHHmmss
	Name           string `json:"name"`           // 用工人员姓名
	ValidateStatus int    `json:"validateStatus"` // 实名校验状态，1：成功 2：失败
	Remark         string `json:"remark"`         // 备注信息，实名校验失败的原因
}

// Register 用工人员注册
func (s *MemberService) Register(req *RegisterRequest) (*RegisterResponse, error) {
	var resp RegisterResponse
	err := s.client.DoRequest("/member/freelancerSpecialApi/register", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
