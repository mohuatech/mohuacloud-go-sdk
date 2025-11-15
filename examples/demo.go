package main

import (
    "fmt"
    "log"

    sdk "github.com/mohuatech/mohuacloud-go-sdk"
    "github.com/mohuatech/mohuacloud-go-sdk/types"
)

func main() {
    // 方式1: 使用环境变量（推荐用于生产环境）
    // 设置环境变量:
    // export MOHUACLOUD_ACCOUNT=your_account
    // export MOHUACLOUD_PASSWORD=your_password
    // export MOHUACLOUD_BASE_URL=https://cloud.mhjz1.cn (可选)
    client1 := sdk.NewClient()

    // 方式2: 直接在代码中配置
    client2 := sdk.NewClient(
        sdk.WithCredentials("your_account", "your_password"),
        sdk.WithBaseURL("https://cloud.mhjz1.cn"), // 可选
    )

    // 使用client1进行演示
    client := client1

    // 1. 登录
    fmt.Println("=== 登录 ===")
    loginResp, err := client.Auth.Login("", "") // 使用配置的账号密码
    if err != nil {
        log.Fatalf("登录失败: %v", err)
    }
    fmt.Printf("登录成功: %s\n", loginResp.Msg)
    fmt.Printf("Token: %s\n", loginResp.JWT)

    // 2. 获取域名列表
    fmt.Println("\n=== 获取域名列表 ===")
    hostID := "your_host_id" // 替换为实际的主机ID
    domainsResp, err := client.VirtualHost.ListDomains(hostID)
    if err != nil {
        log.Fatalf("获取域名列表失败: %v", err)
    }
    fmt.Printf("获取到 %d 个域名:\n", len(domainsResp.Data))
    for _, domain := range domainsResp.Data {
        fmt.Printf(" - ID: %d, 域名: %s, SSL强制: %d\n", domain.ID, domain.Domain, domain.SSLForce)
    }

    // 3. 设置SSL证书（示例）
    if len(domainsResp.Data) > 0 {
        fmt.Println("\n=== 设置SSL证书 ===")
        domain := domainsResp.Data[0]
        
        sslRequest := &types.SetSSLRequest{
            ID:       domain.ID,
            SSLForce: fmt.Sprintf("%d", domain.SSLForce), // 使用原来的值，或者自定义 "0" 或 "1"
            SSLCert:  "-----BEGIN CERTIFICATE-----\n你的证书内容\n-----END CERTIFICATE-----",
            SSLKey:   "-----BEGIN PRIVATE KEY-----\n你的私钥内容\n-----END PRIVATE KEY-----",
        }

        sslResp, err := client.VirtualHost.SetSSL(hostID, sslRequest)
        if err != nil {
            log.Printf("设置SSL失败: %v", err)
        } else {
            fmt.Printf("设置SSL成功: %s\n", sslResp.Msg)
        }
    }

    // 4. 手动设置token（如果已有token）
    fmt.Println("\n=== 手动设置Token ===")
    client.Auth.SetToken("your_existing_jwt_token")
    fmt.Println("Token已设置")
}