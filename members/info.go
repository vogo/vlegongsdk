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

import "fmt"

// InfoRequest 根据id查询用工人员请求
type InfoRequest struct {
	CompanyCode  string `json:"companyCode"`  // 企业编码，用工企业在乐工平台的编码
	FreelancerID string `json:"freelancerId"` // 用工人员编号
}

// InfoResponse 用工人员信息响应
type InfoResponse struct {
	FreelancerID   string `json:"freelancerId"`   // 用工人员编号
	Status         string `json:"status"`         // 用工人员状态，0：离职，1：待上岗，2：已上岗
	CreateTime     string `json:"createTime"`     // 导入时间，yyyyMMddHHmmss
	Name           string `json:"name"`           // 用工人员姓名
	ValidateStatus string `json:"validateStatus"` // 实名校验状态，1：成功 2：失败
	Remark         string `json:"remark"`         // 备注信息，实名校验失败的原因
}

// GetInfo 根据id查询用工人员信息
func (s *MemberService) GetInfo(req *InfoRequest) (*InfoResponse, error) {
	var resp InfoResponse
	err := s.client.DoRequest("/member/freelancerSpecialApi/info", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("查询用工人员信息失败: %w", err)
	}

	return &resp, nil
}
