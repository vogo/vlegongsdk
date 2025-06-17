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

package proofs

import (
	"fmt"
)

// SubmitCompletionProofRequest 提交完工附件请求
type SubmitCompletionProofRequest struct {
	OrderNo string `json:"orderNo"` // 支付返回的订单号
	FileID  string `json:"fileId"`  // 完工附件文件id，在文件上传接口中返回的fileId
}

// SubmitCompletionProof 提交完工附件
// 提交支付对应的完工附件。需先上传完工附件获取文件id，再调用此接口
func (s *ProofService) SubmitCompletionProof(req *SubmitCompletionProofRequest) error {
	err := s.client.DoRequest("/settlement/completionProofApi/submit", req, nil)
	if err != nil {
		return fmt.Errorf("提交完工附件失败: %w", err)
	}

	return nil
}
