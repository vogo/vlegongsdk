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
	"io"
	"net/http"

	"github.com/vogo/vogo/vlog"
)

// CallbackHandler 回调处理器接口
type CallbackHandler[T any] interface {
	// HandleCallback 处理回调请求
	HandleCallback(data T) error
}

// CallbackHandlerFunc 回调处理函数类型
type CallbackHandlerFunc[T any] func(data T) error

// HandleCallback 实现CallbackHandler接口
func (f CallbackHandlerFunc[T]) HandleCallback(data T) error {
	return f(data)
}

// DecodeCallbackRequest 解码回调请求
func DecodeCallbackRequest[T any](client *Client, r *http.Request, data *T) error {
	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("读取请求体失败: %w", err)
	}
	defer r.Body.Close()

	vlog.Infof("收到回调请求: %s", string(body))

	// 解析请求
	var req Request
	if err = json.Unmarshal(body, &req); err != nil {
		return fmt.Errorf("解析请求失败: %w", err)
	}

	// 验证签名
	reqMap, err := req.ToMap()
	if err != nil {
		return fmt.Errorf("转换请求为map失败: %w", err)
	}

	signStr := BuildSignString(reqMap)
	if err = Verify(signStr, req.Head.Sign, client.platformKey); err != nil {
		return fmt.Errorf("验证签名失败: %w", err)
	}

	// 解密请求体
	decryptedData, err := Decrypt(req.Body.Data, client.privateKey)
	if err != nil {
		return fmt.Errorf("解密请求数据失败: %w", err)
	}

	// 解析解密后的数据
	if err := json.Unmarshal([]byte(decryptedData), data); err != nil {
		return fmt.Errorf("解析解密后的数据失败: %w", err)
	}

	vlog.Infof("解析回调数据: %+v", data)

	return nil
}

// WriteCallbackResponse 写入回调响应
func WriteCallbackResponse(w http.ResponseWriter) {
	// 根据API文档，回调响应必须是字符串"success"
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("success"))
}
