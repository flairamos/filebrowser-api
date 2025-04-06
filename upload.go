package filebrowserapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func checkFileType(name string) string {
	switch name {
	case "png":
		return "image/png"
	case "jpg", "jpeg":
		return "image/jpg"
	case "webp":
		return "image/webp"
	case "gif":
		return "image/gif"
	default:
		return ""
	}
}

func UploadFileStream(config FileBrowserConfig, by []byte) (string, error) {
	reader := bytes.NewReader(by)
	// 确保目标路径以/开头
	var targetPath = config.FilePath
	if !strings.HasPrefix(targetPath, "/") {
		targetPath = "/" + targetPath
	}

	// 构建请求URL
	url := fmt.Sprintf("%s/api/resources%s?override=%v", config.BaseURL, targetPath, config.Override)

	// 取文件名
	fileExt := ""
	split := strings.Split(strings.ToLower(config.FilePath), "/")
	for i := len(split) - 1; i >= 0; i-- {
		if strings.Contains(split[i], ".") {
			str := strings.Split(split[i], ".")
			fileExt = strings.TrimSpace(str[1])
			break
		}
	}
	// fmt.Println("----exe", fileExt)
	// 创建请求
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	if checkFileType(fileExt) == "" {
		return "", fmt.Errorf("不支持该文件上传: %v", err)
	}
	// 设置请求头
	req.Header.Set("Content-Type", "image/png")
	token, err := GetToken(config)
	if err != nil {
		return "", fmt.Errorf("Token请求失败: %v", err)
	}
	req.Header.Set("X-Auth", token)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("上传失败: %s (状态码: %d)", string(respBody), resp.StatusCode)
	}

	return string(respBody), nil
}

// SimpleUpload 使用简单上传方式上传文件
func SimpleUpload(baseURL, filePath, targetPath, authToken string, override bool) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 确保目标路径以/开头
	if !strings.HasPrefix(targetPath, "/") {
		targetPath = "/" + targetPath
	}

	// 构建请求URL
	url := fmt.Sprintf("%s/api/resources%s?override=%v", baseURL, targetPath, override)

	split := strings.Split(strings.ToLower(file.Name()), ".")
	fmt.Println("----exe", split)
	// 创建请求
	req, err := http.NewRequest("POST", url, file)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "image/png")
	req.Header.Set("X-Auth", authToken)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("上传失败: %s (状态码: %d)", string(respBody), resp.StatusCode)
	}

	return string(respBody), nil
}
