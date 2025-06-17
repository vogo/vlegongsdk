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

package systems

import "fmt"

// IdentityAuthRequest 实名认证请求
type IdentityAuthRequest struct {
	Name     string `json:"name"`     // 姓名
	IDCardNo string `json:"idCardNo"` // 身份证号
}

// IdentityAuthResponse 实名认证响应
type IdentityAuthResponse struct {
	AuthStatus  string `json:"authStatus"`  // 认证状态，S:成功，F:失败
	Message     string `json:"message"`     // 备注信息
	OrderNumber string `json:"orderNumber"` // 请求编号
}

// IdentityAuth 实名认证，认证姓名与身份证是否匹配
func (s *SystemService) IdentityAuth(req *IdentityAuthRequest) (*IdentityAuthResponse, error) {
	var resp IdentityAuthResponse
	err := s.client.DoRequest("/sys/authApi/iden", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("实名认证失败: %w", err)
	}

	return &resp, nil
}
