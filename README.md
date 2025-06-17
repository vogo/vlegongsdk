# 乐工-共享经济服务系统 golang sdk

官网: https://ilegong.cn

这是一个用于对接乐工共享经济服务系统API的Go语言SDK，提供了基础的加解密、签名验签、请求处理等功能，以及用工人员注册等业务接口的封装。

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
    vlog.Fatalf("创建客户端失败: %v", err)
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
    vlog.Fatalf("注册用工人员失败: %v", err)
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
    vlog.Fatalf("上传文件失败: %v", err)
}

// 打印文件ID
fmt.Printf("上传成功，文件ID: %s\n", resp.FileID)
```

### 发起支付

```go
// 创建结算服务
settlementService := settlements.NewSettlementService(client)

// 创建支付请求
req := &settlements.PayRequest{
    AccountNo:   "6222021234567890123",           // 银行卡号
    Amount:      100.50,                         // 支付金额
    IDCardNo:    "110101199001011234",           // 身份证号
    Name:        "张三",                         // 姓名
    NotifyURL:   "https://example.com/notify",   // 通知URL
    OutOrderNo:  "ORDER_20230101_001",          // 外部订单号
    PayChannel:  settlements.PayChannelBankCard, // 支付渠道
    ProjectCode: "PROJECT001",                   // 项目编码
    Remark:      "测试支付",                     // 备注
}

// 发送请求
order, err := settlementService.Pay(req)
if err != nil {
    vlog.Fatalf("发起支付失败: %v", err)
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
    ProjectCode:  "PROJECT001",                   // 项目编号
    Name:         "张三",                         // 用工人员姓名
    IDCardNo:     "110101199001011234",           // 用工人员身份证号
    NoticeType:   "1",                           // 通知方式 1：发送短信
    SignPlatform: "1",                           // 签署平台 1-网页版短信认证签署
    RedirectURL:  "https://example.com/redirect", // 签署完成后重定向地址
    NotifyURL:    "https://example.com/notify",   // 签署结果通知地址
}

// 发送请求
resp, err := signService.StartSign(req)
if err != nil {
    vlog.Fatalf("发起签约失败: %v", err)
}

// 打印签约信息
fmt.Printf("签约流程ID: %s, 签约状态: %d, 签约链接: %s\n", 
    resp.SignFlowID, resp.SignStatus, resp.SignShortURL)
```

### 提交完工附件

```go
// 创建系统服务和完工证明服务
systemService := systems.NewSystemService(client)
proofsService := proofs.NewProofsService(client)

// 上传完工证明文件
fileResp, err := systemService.UploadFileFromPath("./completion_proof.pdf", systems.FileTypeCompletionProof)
if err != nil {
    vlog.Fatalf("上传完工证明文件失败: %v", err)
}

// 提交完工附件
req := &proofs.SubmitCompletionProofRequest{
    OrderNo: "ORDER123456789", // 支付返回的订单号
    FileID:  fileResp.FileID,   // 完工附件文件ID
}

err = proofsService.SubmitCompletionProof(req)
if err != nil {
    vlog.Fatalf("提交完工附件失败: %v", err)
}

fmt.Println("提交完工附件成功")
```

## 目录结构

- `cores/`: 核心包，包含配置、加解密、签名验签、请求处理、回调处理等基础功能
- `members/`: 成员服务包，包含用工人员注册、注销、信息查询、实名认证、银行卡绑定等接口
- `settlements/`: 结算服务包，包含支付、订单查询等结算相关接口
- `signs/`: 签约服务包，包含自主签约、签约状态查询、签约回调处理等接口
- `systems/`: 系统服务包，包含文件上传下载、系统认证等接口
- `proofs/`: 完工证明服务包，包含提交完工附件等接口
- `examples/`: 各功能模块的示例代码

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
- 请确保提供正确的机构编号和租户编码
- 请确保API基础URL正确

## 许可证

[Apache License Version 2.0](LICENSE)