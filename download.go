package filebrowserapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// DownloadFile 从FileBrowser服务下载文件
func DownloadFile(opts DownloadOptions) error {
	// 构建API URL
	apiURL := fmt.Sprintf("%s/api/raw%s", opts.BaseURL, opts.FilePath)

	// 添加查询参数
	queryParams := url.Values{}
	if opts.Format != "" {
		queryParams.Set("algo", opts.Format)
	}
	if opts.Inline {
		queryParams.Set("inline", "true")
	}

	if len(queryParams) > 0 {
		apiURL += "?" + queryParams.Encode()
	}

	// 创建HTTP请求
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 添加认证头
	if opts.AuthToken != "" {
		req.Header.Set("X-Auth", opts.AuthToken)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(body))
	}

	// 确定输出文件名
	outputPath := opts.OutputPath
	if outputPath == "" {
		// 如果没有指定输出路径，使用文件原始名称
		outputPath = filepath.Base(opts.FilePath)
	}

	// 创建输出目录（如果不存在）
	outputDir := filepath.Dir(outputPath)
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("创建输出目录失败: %w", err)
		}
	}

	// 创建输出文件
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer outFile.Close()

	// 复制响应内容到输出文件
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	fmt.Printf("文件已成功下载到: %s\n", outputPath)
	return nil
}

// DownloadSharedFile 下载共享文件
func DownloadSharedFile(baseURL, hash, token, outputPath string) error {
	// 构建API URL
	apiURL := fmt.Sprintf("%s/api/public/dl/%s", baseURL, hash)

	// 添加查询参数
	queryParams := url.Values{}
	if token != "" {
		queryParams.Set("token", token)
	}

	if len(queryParams) > 0 {
		apiURL += "?" + queryParams.Encode()
	}

	// 创建HTTP请求
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(body))
	}

	// 确定输出文件名
	if outputPath == "" {
		// 如果没有指定输出路径，使用Content-Disposition中的文件名或默认名称
		contentDisposition := resp.Header.Get("Content-Disposition")
		if contentDisposition != "" {
			parts := strings.Split(contentDisposition, "filename*=utf-8''")
			if len(parts) > 1 {
				outputPath = parts[1]
				// URL解码文件名
				decodedName, err := url.QueryUnescape(outputPath)
				if err == nil {
					outputPath = decodedName
				}
			}
		}

		if outputPath == "" {
			outputPath = fmt.Sprintf("shared_file_%s", hash)
		}
	}

	// 创建输出目录（如果不存在）
	outputDir := filepath.Dir(outputPath)
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("创建输出目录失败: %w", err)
		}
	}

	// 创建输出文件
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer outFile.Close()

	// 复制响应内容到输出文件
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	fmt.Printf("共享文件已成功下载到: %s\n", outputPath)
	return nil
}

// ListFiles 列出目录中的文件
func ListFiles(baseURL, dirPath, authToken string) ([]FileInfo, error) {
	// 构建API URL
	apiURL := fmt.Sprintf("%s/api/resources%s", baseURL, dirPath)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 添加认证头
	if authToken != "" {
		req.Header.Set("X-Auth", authToken)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var result struct {
		Items []FileInfo `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return result.Items, nil
}
