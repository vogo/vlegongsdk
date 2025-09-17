# 乐工-共享经济服务系统 golang sdk

官网: https://ilegong.cn

这是一个用于对接乐工共享经济服务系统API的Go语言SDK，提供了基础的加解密、签名验签、请求处理等功能，以及用工人员注册、签约、支付、文件上传和完工证明等业务接口的封装。

## 目录结构

```
├── cores/       # 核心功能，包括客户端、配置、加解密、签名验签、请求处理等
├── members/     # 用工人员相关功能，包括注册、身份认证、银行卡绑定等
├── settlements/ # 结算相关功能，包括支付、订单查询和支付回调处理等
├── signs/       # 签约相关功能，包括发起签约、查询签约状态和签约回调处理等
├── systems/     # 系统相关功能，包括文件上传下载、实名认证等
├── proofs/      # 完工证明相关功能，包括提交完工附件等
├── examples/    # 使用示例
└── doc/         # API文档
```

## 安装

```bash
go get github.com/vogo/vlegongsdk
```

## 使用示例

### 初始化客户端

```go
// 创建配置
config := cores.NewConfig(
    "https://labor-cloud-gate-stab.tenserpay.xyz/api", // API基础URL
    "V00001",                  // 机构编号
    "uptest",                  // 租户编码
    "PRIVATE KEY BASE64",      // 私钥
    "PUBLIC KEY BASE64",       // 平台公钥
)

// 创建客户端
client, err := cores.NewClient(config)
if err != nil {
    vlog.Fatalf("Failed to create client: %v", err)
}
```

### 用工人员注册

```go
// 创建成员服务
memberService := members.NewMemberService(client)

// 创建用工人员注册请求
req := &members.RegisterRequest{
    CompanyCode:    "COMPANY001",                      // 企业编码
    FreelancerName: "张三",                           // 用工人员姓名
    MobilePhone:    "13800138000",                    // 手机号码
    IDCardNo:       "110101199001011234",             // 身份证号
    FreelancerType: "1",                             // 1:自由职业者, 2:雇员
    CreateTime:     time.Now().Format("20060102150405"), // 当前时间
}

// 发送请求
resp, err := memberService.Register(req)
if err != nil {
    vlog.Fatalf("Failed to register member: %v", err)
}

// 打印响应
fmt.Printf("注册成功，用工人员编号: %d\n", resp.FreelancerID)
```

### 上传文件

```go
// 创建系统服务
systemService := systems.NewSystemService(client)

// 上传身份证文件
resp, err := systemService.UploadFileFromPath("./idcard.jpg", systems.FileTypeIDCard)
if err != nil {
    vlog.Fatalf("Failed to upload file: %v", err)
}

// 打印文件ID
fmt.Printf("上传成功，文件ID: %s\n", resp.FileID)
```

### 发起支付

```go
// 创建结算服务
settlementService := settlements.NewSettlementService(client, nil)

// 创建支付请求
req := &settlements.PayRequest{
    AccountNo:   "6222021234567890123",  // 银行卡号
    Amount:      100.00,                // 支付金额（元）
    IDCardNo:    "110101199001011234", // 身份证号
    Name:        "张三",               // 姓名
    NotifyURL:   "https://example.com/payment/callback", // 支付结果通知URL
    OutOrderNo:  "ORDER123456",        // 外部订单号
    PayChannel:  settlements.PayChannelBankCard, // 支付渠道（银行卡）
    ProjectCode: "PROJECT001",         // 项目编码
    Remark:      "测试支付",           // 备注
}

// 发起支付
order, err := settlementService.Pay(req)
if err != nil {
    vlog.Fatalf("Failed to initiate payment: %v", err)
}

// 打印订单信息
fmt.Printf("支付成功，订单号: %s, 状态: %s\n", order.OrderNo, order.Status)
```

### 发起签约

```go
// 创建签约服务
signService := signs.NewSignService(client, nil)

// 创建签约请求
req := &signs.StartSignRequest{
    ProjectCode:  "PROJECT001",         // 项目编号
    Name:         "张三",               // 用工人员姓名
    IDCardNo:     "110101199001011234", // 身份证号
    NoticeType:   "1",                 // 发送短信
    SignPlatform: "1",                 // 网页版短信认证签署
    RedirectURL:  "https://example.com/redirect", // 签署完成后跳转地址
    NotifyURL:    "https://example.com/sign/callback", // 签署结果通知URL
}

// 发起签约
resp, err := signService.StartSign(req)
if err != nil {
    vlog.Fatalf("Failed to initiate signing: %v", err)
}

// 打印签约信息
fmt.Printf("签约流程ID: %s, 签约链接: %s\n", resp.SignFlowID, resp.SignShortURL)
```

### 提交完工附件

```go
// 创建系统服务和完工证明服务
systemService := systems.NewSystemService(client)
proofsService := proofs.NewProofsService(client)

// 上传完工证明文件
fileResp, err := systemService.UploadFileFromPath("./completion_proof.pdf", systems.FileTypeCompletionProof)
if err != nil {
    vlog.Fatalf("Failed to upload completion proof file: %v", err)
}

// 提交完工附件
req := &proofs.SubmitCompletionProofRequest{
    OrderNo: "ORDER123456", // 订单号
    FileID:  fileResp.FileID, // 文件ID
}

err = proofsService.SubmitCompletionProof(req)
if err != nil {
    vlog.Fatalf("Failed to submit completion attachment: %v", err)
}

fmt.Println("完工附件提交成功")
```

## 功能特性

- RSA加密和SHA256withRSA签名验签机制
- 自动处理请求头和签名
- 自动加密请求体和解密响应体
- 提供友好的接口封装
- 支持用工人员全生命周期管理（注册、实名认证、银行卡绑定、注销等）
- 支持自主签约和签约状态查询
- 支持结算支付和订单查询
- 支持文件上传下载（身份证、完工证明等）
- 支持完工附件提交
- 提供回调处理机制

## 注意事项

- 请确保提供正确的私钥和平台公钥
- 请确保在生产环境中使用HTTPS进行通信，保证数据传输安全
- 私钥请妥善保管，不要泄露给第三方
- 回调处理接口需要部署在可公网访问的服务器上
- 支付接口调用后，需要通过查询接口或等待回调通知来确认最终支付结果
- 请确保提供正确的机构编号和租户编码
- 请确保API基础URL正确

## 贡献代码

我们非常欢迎您为乐工SDK贡献代码！如果您有任何改进或新功能的想法，请按照以下步骤参与：

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

### 问题反馈

如果您在使用过程中遇到任何问题或有任何建议，欢迎通过以下方式反馈：

- 提交 [GitHub Issues](https://github.com/vogo/vlegongsdk/issues)
- 在 Pull Request 中详细描述您的问题或建议

我们会尽快回复并解决您提出的问题。

## 许可证

[Apache License Version 2.0](LICENSE)