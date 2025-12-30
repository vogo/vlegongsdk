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

// StartSignRequest 自主签约请求
type StartSignRequest struct {
	ProjectCode  string `json:"projectCode"`  // 项目编号
	Name         string `json:"name"`         // 用工人员姓名
	IDCardNo     string `json:"idCardNo"`     // 用工人员身份证号
	NoticeType   string `json:"noticeType"`   // 通知方式 0：不发短信 1：发送短信
	SignPlatform string `json:"signPlatform"` // 签署平台，默认值1。1-网页版短信认证签署；2-跳转支付宝(移动端)或支付宝扫码进行签署；3-微信小程序人脸认证签署；4-支付宝小程序人脸认证签署；5-网页版人脸认证签署。
	RedirectURL  string `json:"redirectUrl"`  // 签署完成后，重定向跳转地址
	NotifyURL    string `json:"notifyUrl"`    // 签署结果通知
}

// StartSignResponse 自主签约响应
type StartSignResponse struct {
	SignFlowID   string `json:"signFlowId"`   // 签约流程ID
	SignStatus   int    `json:"signStatus"`   // 签约状态 0 签署中 1 已签约 2 拒签 3 过期 4 失败
	SignShortURL string `json:"signShortUrl"` // 签署链接短链
	SignValidity string `json:"signValidity"` // 签署有效截止日期，格式：yyyy-MM-dd HH:mm:ss
	SignDesc     string `json:"signDesc"`     // 签署描述
}

// StartSign 发起自主签约
func (s *SignService) StartSign(req *StartSignRequest) (*StartSignResponse, error) {
	var resp StartSignResponse
	err := s.client.DoRequest("/settlement/signApi/startSign", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
