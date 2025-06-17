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

// IDCardAuthRequest 采集用工人员身份证请求
type IDCardAuthRequest struct {
	FreelancerID int    `json:"freelancerId"` // 用工人员编号
	FrontImgID   string `json:"frontImgId"`   // 人像面照片id，在文件上传接口中返回的fileId
	BackImgID    string `json:"backImgId"`    // 国徽面照片id，在文件上传接口中返回的fileId
}

// IDCardAuth 采集用工人员身份证进行ocr校验
// 注意：需要先通过【文件上传】接口上传身份证正反图片，获取图片ID
func (s *MemberService) IDCardAuth(req *IDCardAuthRequest) error {
	// 该接口没有返回具体数据，只需要处理可能的错误
	err := s.client.DoRequest("/member/freelancerApi/idCardAuth", req, nil)
	if err != nil {
		return fmt.Errorf("采集用工人员身份证失败: %w", err)
	}

	return nil
}
