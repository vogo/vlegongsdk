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

	"github.com/vogo/vlegongsdk/examples"
	"github.com/vogo/vlegongsdk/proofs"
	"github.com/vogo/vlegongsdk/systems"
	"github.com/vogo/vogo/vlog"
	"github.com/vogo/vogo/vos"
)

func main() {
	client := examples.LoadClient()

	// 创建系统服务（用于上传文件）
	systemService := systems.NewSystemService(client)

	// 创建完工证明服务
	proofsService := proofs.NewProofService(client)

	// 示例：上传完工附件并提交
	uploadAndSubmitCompletionProofExample(systemService, proofsService)
}

// 上传完工附件并提交示例
func uploadAndSubmitCompletionProofExample(systemService *systems.SystemService, proofsService *proofs.ProofService) {
	fmt.Println("\n===== 上传完工附件并提交示例 =====")

	// 步骤1: 上传完工附件文件
	filePath := vos.EnsureEnvString("COMPLETION_PROOF_FILE_PATH") // 替换为实际的文件路径
	fmt.Printf("上传完工附件文件: %s\n", filePath)

	// 上传文件
	fileResp, err := systemService.UploadFileFromPath(filePath, systems.FileTypeCompletionProof)
	if err != nil {
		vlog.Errorf("Failed to upload completion attachment file: %v\n", err)
		return
	}

	fmt.Printf("文件上传成功，文件ID: %s\n", fileResp.FileID)

	// 步骤2: 提交完工附件
	orderNo := vos.EnsureEnvString("LEGONG_ORDER_NO") // 替换为实际的订单号
	fmt.Printf("提交完工附件，订单号: %s\n", orderNo)

	// 创建提交完工附件请求
	req := &proofs.SubmitCompletionProofRequest{
		OrderNo: orderNo,
		FileID:  fileResp.FileID,
	}

	// 调用提交完工附件接口
	err = proofsService.SubmitCompletionProof(req)
	if err != nil {
		vlog.Errorf("Failed to submit completion attachment: %v\n", err)
		return
	}

	fmt.Println("提交完工附件成功")
}
