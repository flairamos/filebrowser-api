package filebrowserapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func DeleteFile(config FileBrowserConfig) error {
	// 创建请求体
	filename := ""
	spl := strings.Split(config.FilePath, "/")
	for i := len(spl) - 1; i >= 0; i-- {
		if strings.Contains(spl[i], ".") {
			filename = strings.TrimSpace(spl[i])
		}
	}
	requestBody := map[string]string{
		"path": filename,
	}

	// 将请求体转换为 JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("创建请求体失败: %v", err)
	}

	// 创建请求
	req, err := http.NewRequest("DELETE", config.BaseURL+"/api/resources"+config.FilePath, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	token, err := GetToken(config)
	if err != nil {
		return err
	}
	req.Header.Set("x-auth", token)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("删除失败，状态码: %d，响应: %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}
