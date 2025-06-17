# 共享经济服务系统 SDK

这是一个用于对接共享经济服务系统API的Go语言SDK，提供了基础的加解密、签名验签、请求处理等功能，以及用工人员注册等业务接口的封装。

## 安装

```bash
go get github.com/vogo/vlegongsdk
```

## 使用示例

### 初始化客户端

```go
// 创建配置
config := cores.NewConfig(
    "https://api.example.com", // API基础URL
    "V00001",                  // 机构编号
    "uptest",                  // 租户编码
    "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----", // 私钥
    "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----",       // 平台公钥
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
fmt.Printf("注册成功，用工人员编号: %s\n", resp.FreelancerID)
```

## 目录结构

- `cores/`: 核心包，包含配置、加解密、签名验签、请求处理等基础功能
- `members/`: 成员服务包，包含用工人员相关接口
- `examples/`: 示例代码

## 功能特性

- RSA加密和SHA256withRSA签名
- 自动处理请求头和签名
- 自动加密请求体和解密响应体
- 提供友好的接口封装

## 注意事项

- 请确保提供正确的私钥和平台公钥
- 请确保提供正确的机构编号和租户编码
- 请确保API基础URL正确

## 许可证

[Apache License Version 2.0](LICENSE)