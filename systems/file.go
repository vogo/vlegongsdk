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

package systems

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileType 文件类型
type FileType int

const (
	// FileTypeIDCard 身份证图片
	FileTypeIDCard FileType = 1
	// FileTypeCompletionProof 完工证明文件
	FileTypeCompletionProof FileType = 2
)

// FileUploadRequest 文件上传请求
type FileUploadRequest struct {
	FileName string   `json:"fileName"` // 文件名
	FileHash string   `json:"fileHash"` // 文件hash值，根据文件内容进行MD5计算后的
	FileType FileType `json:"fileType"` // 文件类型：1-身份证图片，2-完工证明文件
}

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	FileID string `json:"fileId"` // 文件id标识
}

// UploadFile 上传文件
// fileBytes: 文件字节数组
// fileName: 文件名
// fileType: 文件类型
func (s *SystemService) UploadFile(fileBytes []byte, fileName string, fileType FileType) (*FileUploadResponse, error) {
	// 计算文件的MD5哈希值
	fileHash := calculateBytesMD5(fileBytes)

	// 构建请求
	req := &FileUploadRequest{
		FileName: fileName,
		FileHash: fileHash,
		FileType: fileType,
	}

	// 发送请求
	var resp FileUploadResponse
	err := s.client.DoMultipartRequestWithBytes("/sys/file/upload/file", req, fileBytes, fileName, &resp)
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	return &resp, nil
}

// UploadFileFromPath 从文件路径上传文件
// filePath: 本地文件路径
// fileType: 文件类型
func (s *SystemService) UploadFileFromPath(filePath string, fileType FileType) (*FileUploadResponse, error) {
	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	if fileInfo.IsDir() {
		return nil, fmt.Errorf("%s 是一个目录，不是文件", filePath)
	}

	// 读取文件内容
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件内容失败: %w", err)
	}

	// 调用UploadFile方法
	return s.UploadFile(fileBytes, filepath.Base(filePath), fileType)
}

// calculateBytesMD5 计算字节数组的MD5哈希值
func calculateBytesMD5(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// FileDownloadRequest 文件下载请求
type FileDownloadRequest struct {
	FileID string `json:"fileId"` // 文件ID
}

// DownloadFile 下载文件到指定路径
// fileID: 文件ID
// savePath: 保存路径
func (s *SystemService) DownloadFile(fileID string, savePath string) error {
	// 构建请求
	req := &FileDownloadRequest{
		FileID: fileID,
	}

	// 发送请求
	resp, err := s.client.DoFileDownloadRequest("/sys/fileApi/download", req)
	if err != nil {
		return fmt.Errorf("下载文件失败: %w", err)
	}

	// 检查是否是文件响应
	if !resp.IsFile {
		return fmt.Errorf("下载文件失败: %s", resp.ErrorResp)
	}

	// 如果文件名为空，使用文件ID作为文件名
	if resp.FileName == "" {
		resp.FileName = fileID
	}

	if err := os.WriteFile(savePath, resp.Content, 0o644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// DownloadFileToWriter 下载文件并写入到writer
// fileID: 文件ID
// writer: 输出writer
// 返回值: 文件名, error
func (s *SystemService) DownloadFileToWriter(fileID string, writer io.Writer) (string, error) {
	// 构建请求
	req := &FileDownloadRequest{
		FileID: fileID,
	}

	// 发送请求
	resp, err := s.client.DoFileDownloadRequest("/sys/fileApi/download", req)
	if err != nil {
		return "", fmt.Errorf("下载文件失败: %w", err)
	}

	// 检查是否是文件响应
	if !resp.IsFile {
		return "", fmt.Errorf("下载文件失败: %s", resp.ErrorResp)
	}

	// 如果文件名为空，使用文件ID作为文件名
	if resp.FileName == "" {
		resp.FileName = fileID
	}

	// 写入writer
	if _, err := writer.Write(resp.Content); err != nil {
		return resp.FileName, fmt.Errorf("写入数据失败: %w", err)
	}

	return resp.FileName, nil
}
