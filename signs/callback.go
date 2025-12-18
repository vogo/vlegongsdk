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

import (
	"fmt"
	"net/http"

	"github.com/vogo/vlegongsdk/cores"
	"github.com/vogo/vogo/vlog"
)

// SignCallbackRequest 签约结果通知回调请求
type SignCallbackRequest struct {
	SignFlowID  string `json:"signFlowId"`  // 签约流程ID
	SignStatus  int    `json:"signStatus"`  // 签约状态 0 签署中 1 已签约 2 拒签 3 过期 4 失败
	SignEndTime string `json:"signEndTime"` // 签署完成时间，格式：yyyy-MM-dd HH:mm:ss
	SignDesc    string `json:"signDesc"`    // 签署描述
}

// SignCallbackHandler 签约回调处理器
type SignCallbackHandler struct {
	client  *cores.Client
	handler cores.CallbackHandler[SignCallbackRequest]
}

// NewSignCallbackHandler 创建签约回调处理器
func NewSignCallbackHandler(client *cores.Client, handler cores.CallbackHandler[SignCallbackRequest]) *SignCallbackHandler {
	return &SignCallbackHandler{
		client:  client,
		handler: handler,
	}
}

// ServeHTTP 实现http.Handler接口，处理签约结果通知回调
func (h *SignCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 只接受POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析回调请求
	var callbackReq SignCallbackRequest
    if err := cores.DecodeCallbackRequest(h.client, r, &callbackReq); err != nil {
        vlog.Errorf("failed to parse sign callback request | err: %v", err)
        http.Error(w, fmt.Sprintf("解析签约回调请求失败: %v", err), http.StatusBadRequest)
        return
    }

	// 处理回调请求
    if err := h.handler.HandleCallback(callbackReq); err != nil {
        vlog.Errorf("failed to handle sign callback request | err: %v", err)
        http.Error(w, fmt.Sprintf("处理签约回调请求失败: %v", err), http.StatusInternalServerError)
        return
    }

	// 返回成功响应
	cores.WriteCallbackResponse(w)
}

// GetSignStatusDesc 获取签约状态描述
func GetSignStatusDesc(status int) string {
	switch status {
	case 0:
		return "签署中"
	case 1:
		return "已签约"
	case 2:
		return "拒签"
	case 3:
		return "过期"
	case 4:
		return "失败"
	default:
		return "未知状态"
	}
}
