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

package settlements

import "fmt"

// QuerySignRequest 签约查询请求
type QuerySignRequest struct {
	SignFlowID string `json:"signFlowId"` // 签约流程ID
}

// QuerySignResponse 签约查询响应
type QuerySignResponse struct {
	SignFlowID  string `json:"signFlowId"`  // 签约流程ID
	SignStatus  int    `json:"signStatus"`  // 签约状态 0 签署中 1 已签约 2 拒签 3 过期 4 失败
	SignEndTime string `json:"signEndTime"` // 签署完成时间，格式：yyyy-MM-dd HH:mm:ss
	SignDesc    string `json:"signDesc"`    // 签署描述
}

// QuerySign 查询签约状态
func (s *SettlementService) QuerySign(req *QuerySignRequest) (*QuerySignResponse, error) {
	var resp QuerySignResponse
	err := s.client.DoRequest("/settlement/signApi/query", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("查询签约状态失败: %w", err)
	}

	return &resp, nil
}
