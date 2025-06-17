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

// DismissRequest 用工人员注销请求
type DismissRequest struct {
	CompanyCode  string `json:"companyCode"`  // 企业编码，用工企业在乐工平台的编码
	FreelancerID string `json:"freelancerId"` // 用工人员编号
}

// DismissResponse 用工人员注销响应
type DismissResponse struct {
	FreelancerID string `json:"freelancerId"` // 用工人员编号
	IsDismissed  string `json:"isDismissed"`  // 注销结果，0：注销失败，1：注销成功
	Remark       string `json:"remark"`       // 备注信息，注销失败原因
}

// Dismiss 用工人员注销
func (s *MemberService) Dismiss(req *DismissRequest) (*DismissResponse, error) {
	var resp DismissResponse
	err := s.client.DoRequest("/member/freelancerSpecialApi/dismiss", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("注销用工人员失败: %w", err)
	}

	return &resp, nil
}
