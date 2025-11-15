package types

// LoginRequest 登录请求参数
type LoginRequest struct {
    Account  string `json:"account"`
    Password string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
    JWT    string `json:"jwt"`
    Status int    `json:"status"`
    Msg    string `json:"msg"`
}

// DomainInfo 域名信息
type DomainInfo struct {
    ID        int    `json:"id"`
    HostID    int    `json:"host_id"`
    UID       int    `json:"uid"`
    Domain    string `json:"domain"`
    SSLCertID int    `json:"ssl_cert_id"`
    SSLForce  int    `json:"ssl_force"`
}

// ListDomainResponse 域名列表响应
type ListDomainResponse struct {
    Status int          `json:"status"`
    Msg    string       `json:"msg"`
    Data   []DomainInfo `json:"data"`
}

// SetSSLRequest 设置SSL请求
type SetSSLRequest struct {
    Func     string `json:"func"`
    ID       int    `json:"id"`
    SSLForce string `json:"ssl_force"`
    SSLCert  string `json:"sslCert"`
    SSLKey   string `json:"sslKey"`
}

// BaseResponse 基础响应
type BaseResponse struct {
    Status int    `json:"status"`
    Msg    string `json:"msg"`
}

// ClientConfig 客户端配置
type ClientConfig struct {
    BaseURL    string
    Account    string
    Password   string
    JWTToken   string
}