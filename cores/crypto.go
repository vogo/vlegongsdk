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
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// 加密和签名相关常量
const (
	// SignTypeRSA RSA签名类型
	SignTypeRSA = "0"
	// EncAlgoRSA RSA加密算法 0:RSA, 1: AES128
	EncAlgoRSA = "0"
	// CharsetUTF8 UTF-8字符集
	CharsetUTF8 = "01"
	// Version 接口版本
	Version = "1.0.0"
	// RSA最大加密明文大小
	maxEncryptBlock = 245
	// RSA最大解密密文大小
	maxDecryptBlock = 256
)

// ParsePrivateKey 解析私钥
func ParsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	// 移除可能的换行符
	privateKeyPEM = strings.ReplaceAll(privateKeyPEM, "\r\n", "")

	keyBytes, err := base64.StdEncoding.DecodeString(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}

	return privateKey.(*rsa.PrivateKey), nil
}

// ParsePublicKey 解析公钥
func ParsePublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	// 移除可能的换行符
	publicKeyPEM = strings.ReplaceAll(publicKeyPEM, "\r\n", "")

	keyBytes, err := base64.StdEncoding.DecodeString(publicKeyPEM)
	if err != nil {
		return nil, err
	}

	pubInterface, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, err
	}

	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return pub, nil
}

// Sign 使用私钥对数据进行签名
func Sign(data string, privateKey *rsa.PrivateKey) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(data))
	hashed := hash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// Verify 使用公钥验证签名
func Verify(data, signature string, publicKey *rsa.PublicKey) error {
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	hash := sha256.New()
	hash.Write([]byte(data))
	hashed := hash.Sum(nil)

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed, sig)
}

// Encrypt 使用公钥加密数据
func Encrypt(data string, publicKey *rsa.PublicKey) (string, error) {
	dataBytes := []byte(data)
	dataLen := len(dataBytes)

	// 分段加密
	var encryptedData []byte
	buffer := bytes.NewBuffer(nil)
	offset := 0

	for dataLen-offset > 0 {
		var chunk []byte
		if dataLen-offset > maxEncryptBlock {
			chunk = dataBytes[offset : offset+maxEncryptBlock]
		} else {
			chunk = dataBytes[offset:dataLen]
		}

		encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, chunk)
		if err != nil {
			return "", err
		}

		buffer.Write(encrypted)
		offset += maxEncryptBlock
	}

	encryptedData = buffer.Bytes()
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

// Decrypt 使用私钥解密数据
func Decrypt(data string, privateKey *rsa.PrivateKey) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	// 分段解密
	dataLen := len(ciphertext)
	buffer := bytes.NewBuffer(nil)
	offset := 0

	for dataLen-offset > 0 {
		var chunk []byte
		if dataLen-offset > maxDecryptBlock {
			chunk = ciphertext[offset : offset+maxDecryptBlock]
		} else {
			chunk = ciphertext[offset:dataLen]
		}

		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, chunk)
		if err != nil {
			return "", err
		}

		buffer.Write(decrypted)
		offset += maxDecryptBlock
	}

	return string(buffer.Bytes()), nil
}

// BuildSignString 构建待签名字符串
func BuildSignString(params map[string]interface{}) string {
	// 移除sign字段
	delete(params, "sign")

	// 创建有序map (类似Java的TreeMap)
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	// 按照ASCII码升序排序
	sort.Strings(keys)

	// 构建待签名字符串
	var builder strings.Builder
	firstParam := true

	for _, k := range keys {
		v := params[k]
		// 跳过空值 (与Java版本保持一致)
		// 先检查是否为nil
		if v == nil {
			continue
		}
		
		// 转换为字符串并检查是否为空
		valueStr := fmt.Sprintf("%v", v)
		if valueStr == "" {
			continue
		}

		// 添加分隔符
		if !firstParam {
			builder.WriteString("&")
		} else {
			firstParam = false
		}

		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(valueStr)
	}

	return builder.String()
}
