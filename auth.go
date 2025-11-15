package mohuacloud

import (
    "fmt"

    "github.com/mohuatech/mohuacloud-go-sdk/types"
)

// AuthService 认证服务
type AuthService struct {
    client *Client
}

// NewAuthService 创建认证服务实例
func NewAuthService(client *Client) *AuthService {
    return &AuthService{
        client: client,
    }
}

// Login API登录
func (a *AuthService) Login(account, password string) (*types.LoginResponse, error) {
    if account == "" {
        account = a.client.config.Account
    }
    if password == "" {
        password = a.client.config.Password
    }

    if account == "" || password == "" {
        return nil, fmt.Errorf("account and password are required")
    }

    req := &types.LoginRequest{
        Account:  account,
        Password: password,
    }

    var resp types.LoginResponse
    _, err := a.client.R().
        SetBody(req).
        SetResult(&resp).
        Post("/v1/login_api")

    if err != nil {
        return nil, err
    }

    if resp.Status != 200 {
        return nil, fmt.Errorf("login failed: %s", resp.Msg)
    }

    // 保存token到client
    a.client.SetToken(resp.JWT)

    return &resp, nil
}

// SetToken 设置token
func (a *AuthService) SetToken(token string) {
    a.client.SetToken(token)
}

// GetToken 获取当前token
func (a *AuthService) GetToken() string {
    return a.client.config.JWTToken
}