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
	"github.com/vogo/vlegongsdk/cores"
)

// SignService 签约服务
type SignService struct {
	client *cores.Client
	// 签约回调处理器
	signCallbackHandler *SignCallbackHandler
}

// NewSignService 创建一个新的签约服务
func NewSignService(client *cores.Client, signCallbackHandler cores.CallbackHandler[SignCallbackRequest]) *SignService {
	service := &SignService{
		client: client,
	}

	// 如果提供了回调处理器，则创建签约回调处理器
	if signCallbackHandler != nil {
		service.signCallbackHandler = NewSignCallbackHandler(client, signCallbackHandler)
	}

	return service
}

// GetSignCallbackHandler 获取签约回调处理器
func (s *SignService) GetSignCallbackHandler() *SignCallbackHandler {
	return s.signCallbackHandler
}
