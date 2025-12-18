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
        vlog.Fatal("payment callback handler is empty")
    }

	// 注册HTTP处理器
	http.Handle("/api/callback/payment", paymentCallbackHandler)

	// 启动HTTP服务器
	port := 8080
    vlog.Infof("starting http server | port: %d", port)
    vlog.Infof("payment callback url | url: http://localhost:%d/api/callback/payment", port)
    vlog.Fatalf("http server error | err: %v", http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// handlePaymentCallback 处理支付回调
func handlePaymentCallback(data *settlements.PaymentCallbackRequest) error {
	// 处理支付回调
    vlog.Infof("received payment callback | order_no: %s | out_order_no: %s | status: %s | amount: %.2f | received_amount: %.2f | success_time: %s",
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
        vlog.Infof("payment success detail | pay_channel: %d | tax: %.2f | service_fee_bearer: %s",
            data.PayChannel,
            data.Tax,
            settlements.GetBearWayDesc(data.ServiceChargeBearWay),
        )
	case "F": // 支付失败
		log.Println("支付失败，需要排查原因")
		// 可以在这里添加业务逻辑
        vlog.Infof("payment failed | reason: %s", data.Desc)
	case "P": // 处理中
		log.Println("支付处理中，等待最终结果")
		// 可以在这里添加业务逻辑
	default:
        vlog.Infof("unknown payment status | status: %s", data.Status)
	}

	return nil
}
