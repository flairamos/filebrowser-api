package filebrowserapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// 获取认证令牌
func GetToken(config FileBrowserConfig) (string, error) {
	// 创建登录请求
	loginData := fmt.Sprintf(`{"username":"%s","password":"%s"}`, config.Username, config.Password)

	// 创建请求
	req, err := http.NewRequest("POST", config.BaseURL+"/api/login", bytes.NewBufferString(loginData))
	if err != nil {
		return "", fmt.Errorf("创建登录请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送登录请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应内容失败: %v", err)
	}

	// 打印原始响应内容
	// fmt.Printf("原始响应内容: %s\n", string(bodyBytes))
	return string(bodyBytes), nil
}
