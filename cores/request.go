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

package cores

import (
	"encoding/json"
	"fmt"
	"time"
)

// RequestHead 请求头
type RequestHead struct {
	Charset     string `json:"charset"`
	Version     string `json:"version"`
	Sign        string `json:"sign"`
	SignType    string `json:"signType"`
	RequestID   string `json:"requestId"`
	RequestTime string `json:"requestTime"`
	Nonce       string `json:"nonce"`
	OrgCode     string `json:"orgCode"`
	EncAlgo     string `json:"encAlgo"`
	TenantCode  string `json:"tenantCode"`
}

// RequestBody 请求体
type RequestBody struct {
	Data interface{} `json:"data"`
}

// Request 请求
type Request struct {
	Head RequestHead `json:"head"`
	Body RequestBody `json:"body"`
}

// ResponseHead 响应头
type ResponseHead struct {
	Charset     string `json:"charset"`
	Version     string `json:"version"`
	SignType    string `json:"signType"`
	Sign        string `json:"sign,omitempty"`
	RequestID   string `json:"requestId"`
	RequestTime string `json:"requestTime"`
	Nonce       string `json:"nonce"`
	OrgCode     string `json:"orgCode,omitempty"`
	EncAlgo     string `json:"encAlgo"`
	TenantCode  string `json:"tenantCode"`
}

// ResponseBody 响应体
type ResponseBody struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response 响应
type Response struct {
	Head ResponseHead `json:"head"`
	Body ResponseBody `json:"body"`
}

// NewRequest 创建一个新的请求
func NewRequest(config *Config, data interface{}) *Request {
	// 生成请求ID
	requestID := fmt.Sprintf("%d%d", time.Now().UnixNano()/1e6, time.Now().Nanosecond()%1000)
	// 生成请求时间
	requestTime := time.Now().Format("20060102150405")
	// 生成随机字符串
	nonce := GenerateNonce(32)

	return &Request{
		Head: RequestHead{
			Charset:     CharsetUTF8,
			Version:     Version,
			SignType:    SignTypeRSA,
			RequestID:   requestID,
			RequestTime: requestTime,
			Nonce:       nonce,
			OrgCode:     config.OrgCode,
			EncAlgo:     EncAlgoRSA,
			TenantCode:  config.TenantCode,
		},
		Body: RequestBody{
			Data: data,
		},
	}
}

// ToMap 将请求转换为map，用于签名
func (r *Request) ToMap() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 处理head
	headBytes, err := json.Marshal(r.Head)
	if err != nil {
		return nil, err
	}

	var headMap map[string]interface{}
	if err := json.Unmarshal(headBytes, &headMap); err != nil {
		return nil, err
	}

	for k, v := range headMap {
		result[k] = v
	}

	// 处理body
	bodyBytes, err := json.Marshal(r.Body)
	if err != nil {
		return nil, err
	}

	var bodyMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &bodyMap); err != nil {
		return nil, err
	}

	for k, v := range bodyMap {
		result[k] = v
	}

	return result, nil
}

// GenerateNonce 生成指定长度的随机字符串
func GenerateNonce(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond) // 确保每次生成的字符不同
	}
	return string(result)
}
