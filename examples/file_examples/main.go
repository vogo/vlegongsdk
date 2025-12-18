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
	"os"

	"github.com/vogo/vlegongsdk/cores"
	"github.com/vogo/vlegongsdk/systems"
	"github.com/vogo/vogo/vlog"
	"github.com/vogo/vogo/vos"
)

func main() {
	// 创建配置
	config := cores.NewConfig(
		vos.EnsureEnvString("LEGONG_API_URL"),     // 替换为实际的API地址
		vos.EnsureEnvString("LEGONG_ORG_CODE"),    // 替换为实际的机构编号
		vos.EnsureEnvString("LEGONG_TENANT_CODE"), // 替换为实际的租户编码
		vos.EnsureEnvString("LEGONG_PRIVATE_KEY"), // 替换为实际的私钥
		vos.EnsureEnvString("LEGONG_PUBLIC_KEY"),  // 替换为实际的平台公钥                           // 替换为实际的平台公钥
	)

	// 创建客户端
    client, err := cores.NewClient(config)
    if err != nil {
        vlog.Fatalf("failed to create client | err: %v", err)
    }

	// 创建系统服务
	systemService := systems.NewSystemService(client)

	// 示例：上传文件
	uploadFileExample(systemService)

	// 示例：下载文件
	downloadFileExample(systemService)
}

// 上传文件示例
func uploadFileExample(systemService *systems.SystemService) {
	fmt.Println("\n===== 文件上传示例 =====")

	// 示例1: 从文件路径上传
	filePath := vos.EnsureEnvString("UPLOAD_FILE_PATH") // 替换为实际的文件路径
	fmt.Printf("从文件路径上传: %s\n", filePath)

	// 上传文件
    resp, err := systemService.UploadFileFromPath(filePath, systems.FileTypeIDCard)
    if err != nil {
        vlog.Errorf("failed to upload file | file_path: %s | file_type: %d | err: %v", filePath, systems.FileTypeIDCard, err)
        return
    }

	// 打印响应
	fmt.Printf("上传成功，文件ID: %s\n", resp.FileID)

	// 保存文件ID用于后续下载示例
	fileID := resp.FileID

	// 将文件ID保存到环境变量，供下载示例使用
	os.Setenv("DOWNLOAD_FILE_ID", fileID)
}

// 下载文件示例
func downloadFileExample(systemService *systems.SystemService) {
	fmt.Println("\n===== 文件下载示例 =====")

	// 获取文件ID，如果环境变量中没有，则使用默认值
	fileID := vos.EnsureEnvString("DOWNLOAD_FILE_ID")
	savePath := vos.EnsureEnvString("DOWNLOAD_FILE_PATH")

	// 下载文件
    err := systemService.DownloadFile(fileID, savePath)
    if err != nil {
        vlog.Errorf("failed to download file | file_id: %s | save_path: %s | err: %v", fileID, savePath, err)
    } else {
        println("下载成功")
    }
}
