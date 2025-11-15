package mohuacloud

import (
    "os"
	"fmt"
    "strings"
	"net/url"

    "github.com/go-resty/resty/v2"
    "github.com/mohuatech/mohuacloud-go-sdk/types"
)

// Client 主客户端
type Client struct {
    client *resty.Client
    config *types.ClientConfig
    
    // 服务
    Auth        *AuthService
    VirtualHost *VirtualHostService
}

// VirtualHostService 虚拟主机服务
type VirtualHostService struct {
    client *Client
}

// Option 客户端配置选项
type Option func(*types.ClientConfig)

// NewClient 创建新的客户端实例
func NewClient(options ...Option) *Client {
    config := &types.ClientConfig{
        BaseURL: "https://cloud.mhjz1.cn",
    }

    // 应用选项
    for _, option := range options {
        option(config)
    }

    // 从环境变量读取配置（如果未设置）
    if config.BaseURL == "" {
        if envURL := os.Getenv("MOHUACLOUD_BASE_URL"); envURL != "" {
            config.BaseURL = envURL
        }
    }

    if config.Account == "" {
        if envAccount := os.Getenv("MOHUACLOUD_ACCOUNT"); envAccount != "" {
            config.Account = envAccount
        }
    }

    if config.Password == "" {
        if envPassword := os.Getenv("MOHUACLOUD_PASSWORD"); envPassword != "" {
            config.Password = envPassword
        }
    }

    // 创建resty客户端
    restyClient := resty.New().
        SetBaseURL(config.BaseURL).
        SetHeader("Content-Type", "application/json")
		restyClient.SetHeader("User-Agent", "mohuacloud-go-sdk")

    // 如果有token，设置认证头
    if config.JWTToken != "" {
        restyClient.SetHeader("Authorization", "JWT "+config.JWTToken)
    }

    client := &Client{
        client: restyClient,
        config: config,
    }

    // 初始化服务
    client.Auth = NewAuthService(client)
    client.VirtualHost = NewVirtualHostService(client)

    return client
}

// WithBaseURL 设置基础URL选项
func WithBaseURL(baseURL string) Option {
    return func(c *types.ClientConfig) {
        c.BaseURL = strings.TrimSuffix(baseURL, "/")
    }
}

// WithCredentials 设置认证信息选项
func WithCredentials(account, password string) Option {
    return func(c *types.ClientConfig) {
        c.Account = account
        c.Password = password
    }
}

// WithToken 设置token选项
func WithToken(token string) Option {
    return func(c *types.ClientConfig) {
        c.JWTToken = token
    }
}

// R 返回resty请求实例
func (c *Client) R() *resty.Request {
    return c.client.R()
}

// SetToken 设置JWT token
func (c *Client) SetToken(token string) {
    c.config.JWTToken = token
    c.client.SetHeader("Authorization", "JWT "+token)
}

// GetConfig 获取客户端配置
func (c *Client) GetConfig() *types.ClientConfig {
    return c.config
}

// NewVirtualHostService 创建虚拟主机服务实例
func NewVirtualHostService(client *Client) *VirtualHostService {
    return &VirtualHostService{
        client: client,
    }
}

// ListDomains 获取域名列表
func (v *VirtualHostService) ListDomains(hostID string) (*types.ListDomainResponse, error) {
    if hostID == "" {
        return nil, fmt.Errorf("hostID is required")
    }

    req := map[string]string{
        "func": "ListDomain",
    }

    var resp types.ListDomainResponse
    _, err := v.client.R().
        SetBody(req).
        SetResult(&resp).
        Post(fmt.Sprintf("/provision/custom/%s", hostID))

    if err != nil {
        return nil, err
    }

    if resp.Status != 200 {
        return nil, fmt.Errorf("failed to list domains: %s", resp.Msg)
    }

    return &resp, nil
}

// SetSSL 设置SSL证书
func (v *VirtualHostService) SetSSL(hostID string, req *types.SetSSLRequest) (*types.BaseResponse, error) {
    if hostID == "" {
        return nil, fmt.Errorf("hostID is required")
    }

    req.Func = "SetSSL"

    encodedReq := map[string]interface{}{
        "func":     req.Func,
        "id":       req.ID,
        "ssl_force": req.SSLForce,
        "sslCert":  url.QueryEscape(req.SSLCert),  // 添加URL编码
        "sslKey":   url.QueryEscape(req.SSLKey),   // 添加URL编码
    }
    var resp types.BaseResponse
    _, err := v.client.R().
        SetBody(encodedReq).
        SetResult(&resp).
        Post(fmt.Sprintf("/provision/custom/%s", hostID))

    if err != nil {
        return nil, err
    }

    if resp.Status != 200 {
        return nil, fmt.Errorf("failed to set SSL: %s", resp.Msg)
    }

    return &resp, nil
}