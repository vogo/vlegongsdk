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
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vogo/vogo/vencoding/vjson"
	"github.com/vogo/vogo/vlog"
)

// Client 客户端
type Client struct {
	config      *Config
	httpClient  *http.Client
	privateKey  *rsa.PrivateKey
	platformKey *rsa.PublicKey
}

// NewClient 创建一个新的客户端
func NewClient(config *Config) (*Client, error) {
	// 解析私钥
	privateKey, err := ParsePrivateKey(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %w", err)
	}

	// 解析平台公钥
	platformKey, err := ParsePublicKey(config.PlatformPublicKey)
	if err != nil {
		return nil, fmt.Errorf("解析平台公钥失败: %w", err)
	}

	// 创建HTTP客户端
	httpClient := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	return &Client{
		config:      config,
		httpClient:  httpClient,
		privateKey:  privateKey,
		platformKey: platformKey,
	}, nil
}

// DoRequest 发送请求并处理响应
func (c *Client) DoRequest(path string, reqData interface{}, respData interface{}) error {
	vlog.Infof("DoRequest path: %s, reqData: %s", path, vjson.EnsureMarshal(reqData))

	// 创建请求
	req := NewRequest(c.config, reqData)

	// 签名
	reqMap, err := req.ToMap()
	if err != nil {
		return fmt.Errorf("转换请求为map失败: %w", err)
	}

	signStr := BuildSignString(reqMap)
	sign, err := Sign(signStr, c.privateKey)
	if err != nil {
		return fmt.Errorf("签名失败: %w", err)
	}

	req.Head.Sign = sign

	// 加密请求体
	reqBodyBytes, err := json.Marshal(req.Body.Data)
	if err != nil {
		return fmt.Errorf("序列化请求体失败: %w", err)
	}

	encryptedData, err := Encrypt(string(reqBodyBytes), c.platformKey)
	if err != nil {
		return fmt.Errorf("加密请求体失败: %w", err)
	}

	// 替换加密后的数据
	req.Body.Data = encryptedData

	// 序列化请求
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	// 发送HTTP请求
	url := fmt.Sprintf("%s%s", c.config.BaseURL, path)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json;charset=utf-8")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer httpResp.Body.Close()

	// 读取响应
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	var resp Response
	if err = json.Unmarshal(respBytes, &resp); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	// 验证签名
	respMap, err := c.responseToMap(resp)
	if err != nil {
		return fmt.Errorf("转换响应为map失败: %w", err)
	}

	signStr = BuildSignString(respMap)
	if err := Verify(signStr, resp.Head.Sign, c.platformKey); err != nil {
		return fmt.Errorf("验证签名失败: %w", err)
	}

	// 检查响应码
	if resp.Body.Code != "00000" {
		return fmt.Errorf("请求失败: %s - %s", resp.Body.Code, resp.Body.Msg)
	}

	// 解密响应数据
	if encryptedResp, ok := resp.Body.Data.(string); ok {
		decryptedData, err := Decrypt(encryptedResp, c.privateKey)
		if err != nil {
			return fmt.Errorf("解密响应数据失败: %w", err)
		}

		// 解析解密后的数据
		if err := json.Unmarshal([]byte(decryptedData), respData); err != nil {
			return fmt.Errorf("解析解密后的数据失败: %w", err)
		}
	} else if resp.Body.Data != nil {
		// 直接解析响应数据
		dataBytes, err := json.Marshal(resp.Body.Data)
		if err != nil {
			return fmt.Errorf("序列化响应数据失败: %w", err)
		}

		if err := json.Unmarshal(dataBytes, respData); err != nil {
			return fmt.Errorf("解析响应数据失败: %w", err)
		}
	} else {
		return errors.New("响应数据为空")
	}

	vlog.Infof("DoRequest respData: %s", vjson.EnsureMarshal(respData))

	return nil
}

// responseToMap 将响应转换为map，用于验签
func (c *Client) responseToMap(resp Response) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 处理head
	headBytes, err := json.Marshal(resp.Head)
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
	bodyBytes, err := json.Marshal(resp.Body)
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
