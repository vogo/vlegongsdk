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
	"mime/multipart"
	"net/http"
	"time"

	"github.com/vogo/vogo/vencoding/vjson"
	"github.com/vogo/vogo/vlog"
)

// DoMultipartRequestWithBytes 使用字节数组发送multipart/form-data请求并处理响应
func (c *Client) DoMultipartRequestWithBytes(path string, reqData interface{}, fileBytes []byte, fileName string, respData interface{}) error {
	vlog.Infof("DoMultipartRequestWithBytes path: %s, reqData: %v, fileName: %s", path, vjson.EnsureMarshal(reqData), fileName)

	// 创建请求
	req := NewRequest(c.config)

	// 加密请求体
	reqBodyBytes, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("序列化请求体失败: %w", err)
	}

	encryptedData, err := Encrypt(string(reqBodyBytes), c.platformKey)
	if err != nil {
		return fmt.Errorf("加密请求体失败: %w", err)
	}

	// 设置加密后的数据
	req.Body.Data = encryptedData

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

	// 序列化请求头和请求体为JSON字符串
	headBytes, err := json.Marshal(req.Head)
	if err != nil {
		return fmt.Errorf("序列化请求头失败: %w", err)
	}

	bodyBytes, err := json.Marshal(req.Body)
	if err != nil {
		return fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建multipart/form-data请求
	url := fmt.Sprintf("%s%s", c.config.BaseURL, path)

	// 创建multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加head字段
	if err = writer.WriteField("head", string(headBytes)); err != nil {
		return fmt.Errorf("写入head字段失败: %w", err)
	}

	// 添加body字段
	if err = writer.WriteField("body", string(bodyBytes)); err != nil {
		return fmt.Errorf("写入body字段失败: %w", err)
	}

	// 添加文件
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return fmt.Errorf("创建文件表单字段失败: %w", err)
	}

	if _, err = part.Write(fileBytes); err != nil {
		return fmt.Errorf("写入文件内容失败: %w", err)
	}

	// 关闭writer
	if err = writer.Close(); err != nil {
		return fmt.Errorf("关闭multipart writer失败: %w", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", url, body)
	if err != nil {
		vlog.Errorf("创建HTTP请求失败: %v", err)
		return fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送HTTP请求
	httpClient := &http.Client{
		Timeout: time.Duration(c.config.Timeout) * time.Second,
	}

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		vlog.Errorf("发送HTTP请求失败: %v", err)
		return fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer httpResp.Body.Close()

	// 读取响应
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		vlog.Errorf("response status: %s, body: %s", httpResp.Status, respBytes)
		return fmt.Errorf("请求失败: %s", httpResp.Status)
	}

	// 解析响应
	var resp Response
	if err = json.Unmarshal(respBytes, &resp); err != nil {
		vlog.Errorf("解析响应失败: %v, body: %s", err, respBytes)
		return fmt.Errorf("解析响应失败: %w", err)
	}

	// 验证签名
	respMap, err := c.responseToMap(resp)
	if err != nil {
		vlog.Errorf("转换响应为map失败: %v, body: %s", err, respBytes)
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

	if respData != nil {
		// 解密响应数据
		decryptedData, err := Decrypt(resp.Body.Data, c.privateKey)
		if err != nil {
			vlog.Errorf("解密响应数据失败: %v, data: %s", err, resp.Body.Data)
			return fmt.Errorf("解密响应数据失败: %w", err)
		}

		vlog.Infof("DoMultipartRequestWithBytes respData: %s", decryptedData)

		// 解析解密后的数据
		if err := json.Unmarshal([]byte(decryptedData), respData); err != nil {
			vlog.Errorf("解析解密后的数据失败: %v, data: %s", err, decryptedData)
			return fmt.Errorf("解析解密后的数据失败: %w", err)
		}
	}

	return nil
}
