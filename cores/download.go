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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vogo/vogo/vlog"
)

// FileDownloadResponse 文件下载响应
type FileDownloadResponse struct {
	IsFile    bool   // 是否是文件，true表示是文件，false表示是错误响应
	FileName  string // 文件名
	Content   []byte // 文件内容
	ErrorResp string // 错误响应JSON
}

// DoFileDownloadRequest 发送文件下载请求并处理响应
func (c *Client) DoFileDownloadRequest(path string, reqData interface{}) (*FileDownloadResponse, error) {
	vlog.Infof("do_file_download_request | path: %s | req_data: %v", path, reqData)

	// 创建请求
	req := NewRequest(c.config)

	// 加密请求体
	reqBodyBytes, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	encryptedData, err := Encrypt(string(reqBodyBytes), c.platformKey)
	if err != nil {
		return nil, fmt.Errorf("加密请求体失败: %w", err)
	}

	// 设置加密后的数据
	req.Body.Data = encryptedData

	// 签名
	reqMap, err := req.ToMap()
	if err != nil {
		return nil, fmt.Errorf("转换请求为map失败: %w", err)
	}

	signStr := BuildSignString(reqMap)

	sign, err := Sign(signStr, c.privateKey)
	if err != nil {
		return nil, fmt.Errorf("签名失败: %w", err)
	}

	req.Head.Sign = sign

	// 序列化请求
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 发送HTTP请求
	reqUrl := fmt.Sprintf("%s%s", c.config.BaseURL, path)

	httpReq, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(reqBytes))
	if err != nil {
		vlog.Errorf("failed to create http request | body: %s | err: %v", reqBytes, err)
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json;charset=utf-8")

	httpClient := &http.Client{
		Timeout: time.Duration(c.config.Timeout) * time.Second,
	}

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		vlog.Errorf("failed to send http request | body: %s | err: %v", reqBytes, err)
		return nil, fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer func() {
		if closeErr := httpResp.Body.Close(); closeErr != nil {
			vlog.Errorf("failed to close response body | err: %v", closeErr)
		}
	}()

	// 检查状态码
	if httpResp.StatusCode != http.StatusOK {
		vlog.Errorf("response status error | status: %s", httpResp.Status)
		return nil, fmt.Errorf("请求失败: %s", httpResp.Status)
	}

	// 创建响应对象
	response := &FileDownloadResponse{}

	// 检查Content-Type
	contentType := httpResp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		// 处理JSON响应（错误情况）
		respBytes, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, fmt.Errorf("读取响应失败: %w", err)
		}

		response.IsFile = false
		response.ErrorResp = string(respBytes)

		// 解析响应以检查错误
		var resp Response
		if err = json.Unmarshal(respBytes, &resp); err != nil {
			vlog.Errorf("failed to parse response | body: %s | err: %v", respBytes, err)
			return response, fmt.Errorf("解析响应失败: %w", err)
		}

		// 验证签名
		respMap, err := c.responseToMap(resp)
		if err != nil {
			vlog.Errorf("failed to convert response to map | body: %s | err: %v", respBytes, err)
			return response, fmt.Errorf("转换响应为map失败: %w", err)
		}

		signStr = BuildSignString(respMap)
		if err := Verify(signStr, resp.Head.Sign, c.platformKey); err != nil {
			return response, fmt.Errorf("验证签名失败: %w", err)
		}

		// 检查响应码
		if resp.Body.Code != "00000" {
			return response, fmt.Errorf("请求失败: %s - %s", resp.Body.Code, resp.Body.Msg)
		}

		return response, nil
	}

	// 处理文件响应
	// 获取文件名
	contentDisposition := httpResp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		// 解析Content-Disposition头
		parts := strings.Split(contentDisposition, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "filename=") {
				fileName := strings.TrimPrefix(part, "filename=")
				// 移除可能的引号
				fileName = strings.Trim(fileName, `"'`)
				// URL解码文件名

				decodedFileName, err := url.QueryUnescape(fileName)
				if err == nil {
					fileName = decodedFileName
				}
				response.FileName = fileName
				break
			}
		}
	}

	// 读取文件内容
	content, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取文件内容失败: %w", err)
	}

	response.IsFile = true
	response.Content = content

	return response, nil
}
