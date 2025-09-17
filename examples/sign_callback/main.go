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

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vogo/vlegongsdk/cores"
	"github.com/vogo/vlegongsdk/examples"
	"github.com/vogo/vlegongsdk/signs"
	"github.com/vogo/vogo/vlog"
)

func main() {
	client := examples.LoadClient()

	// 创建自定义回调处理函数
	callbackHandler := cores.CallbackHandlerFunc[signs.SignCallbackRequest](handleSignCallback)

	// 创建结算服务，并传入回调处理函数
	service := signs.NewSignService(client, callbackHandler)

	// 获取签约回调处理器
	signCallbackHandler := service.GetSignCallbackHandler()
	if signCallbackHandler == nil {
		vlog.Fatal("Sign callback handler is empty")
	}

	// 注册HTTP处理器
	http.Handle("/api/callback/sign", signCallbackHandler)

	// 启动HTTP服务器
	port := 8080
	vlog.Infof("Starting HTTP server, listening on port: %d", port)
	vlog.Infof("Sign callback URL: http://localhost:%d/api/callback/sign", port)
	vlog.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// handleSignCallback 处理签约回调
func handleSignCallback(data signs.SignCallbackRequest) error {
	// 处理签约回调
	vlog.Infof("Received sign callback: ProcessID=%s, Status=%d(%s), CompletionTime=%s, Description=%s",
		data.SignFlowID,
		data.SignStatus,
		signs.GetSignStatusDesc(data.SignStatus),
		data.SignEndTime,
		data.SignDesc,
	)

	// 根据签约状态进行业务处理
	switch data.SignStatus {
	case 0: // 签署中
		log.Println("签署中，等待用户完成签署")
		// 可以在这里添加业务逻辑
	case 1: // 已签约
		log.Println("签约成功，更新业务状态")
		// 可以在这里添加业务逻辑，如更新数据库中的签约状态
	case 2: // 拒签
		log.Println("用户拒签，需要进行相应处理")
		// 可以在这里添加业务逻辑
	case 3: // 过期
		log.Println("签约已过期，需要重新发起签约")
		// 可以在这里添加业务逻辑
	case 4: // 失败
		log.Println("签约失败，需要排查原因")
		// 可以在这里添加业务逻辑
	default:
		vlog.Infof("Unknown sign status: %d", data.SignStatus)
	}

	return nil
}
