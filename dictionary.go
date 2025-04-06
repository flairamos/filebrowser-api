package filebrowserapi

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// CreateDir 创建文件夹
func CreateDir(baseURL, dirPath, authToken string) error {
	// 确保路径以/结尾
	if !strings.HasSuffix(dirPath, "/") {
		dirPath += "/"
	}

	// 构建请求URL
	url := fmt.Sprintf("%s/api/resources%s", baseURL, dirPath)

	// 创建请求
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("X-Auth", authToken)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("创建文件夹失败: %s (状态码: %d)", string(body), resp.StatusCode)
	}
	return nil
}
