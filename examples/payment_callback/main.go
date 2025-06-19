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
	"github.com/vogo/vlegongsdk/settlements"
	"github.com/vogo/vogo/vlog"
)

func main() {
	client := examples.LoadClient()

	// 创建自定义回调处理函数
	callbackHandler := cores.CallbackHandlerFunc[*settlements.PaymentCallbackRequest](handlePaymentCallback)

	// 创建结算服务，并传入回调处理函数
	service := settlements.NewSettlementService(client, callbackHandler)

	// 获取支付回调处理器
	paymentCallbackHandler := service.GetPaymentCallbackHandler()
	if paymentCallbackHandler == nil {
		vlog.Fatal("支付回调处理器为空")
	}

	// 注册HTTP处理器
	http.Handle("/api/callback/payment", paymentCallbackHandler)

	// 启动HTTP服务器
	port := 8080
	vlog.Infof("启动HTTP服务器，监听端口: %d", port)
	vlog.Infof("支付回调URL: http://localhost:%d/api/callback/payment", port)
	vlog.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// handlePaymentCallback 处理支付回调
func handlePaymentCallback(data *settlements.PaymentCallbackRequest) error {
	// 处理支付回调
	vlog.Infof("收到支付回调: 订单号=%s, 外部订单号=%s, 状态=%s, 金额=%.2f, 到账金额=%.2f, 完成时间=%s",
		data.OrderNo,
		data.OutOrderNo,
		data.Status,
		data.Amount,
		data.ReceivedAmount,
		data.SuccessTime,
	)

	// 根据支付状态进行业务处理
	switch data.Status {
	case "S": // 支付成功
		log.Println("支付成功，更新业务状态")
		// 可以在这里添加业务逻辑，如更新数据库中的支付状态
		vlog.Infof("支付渠道: %d, 税费: %.2f, 服务费承担方: %s",
			data.PayChannel,
			data.Tax,
			settlements.GetBearWayDesc(data.ServiceChargeBearWay),
		)
	case "F": // 支付失败
		log.Println("支付失败，需要排查原因")
		// 可以在这里添加业务逻辑
		vlog.Infof("失败原因: %s", data.Desc)
	case "P": // 处理中
		log.Println("支付处理中，等待最终结果")
		// 可以在这里添加业务逻辑
	default:
		vlog.Infof("未知的支付状态: %s", data.Status)
	}

	return nil
}
