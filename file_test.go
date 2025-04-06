package filebrowserapi

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestUploadFileStream(t *testing.T) {
	config := NewConfig("http://test.con:8080", "admin", "admin", "/schoollove/img1.jpg")
	by, _ := os.ReadFile("./img1.jpg")
	stream, err := UploadFileStream(config, by)
	if err != nil {
		t.Error(err)
	}
	t.Log(stream)
}

func TestSimpleUpload(t *testing.T) {
	res, err := SimpleUpload("http://test.con:8080", "./img.png", "/schoollove/img.png", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoxLCJsb2NhbGUiOiJ6aC1jbiIsInZpZXdNb2RlIjoibW9zYWljIiwic2luZ2xlQ2xpY2siOmZhbHNlLCJwZXJtIjp7ImFkbWluIjp0cnVlLCJleGVjdXRlIjp0cnVlLCJjcmVhdGUiOnRydWUsInJlbmFtZSI6dHJ1ZSwibW9kaWZ5Ijp0cnVlLCJkZWxldGUiOnRydWUsInNoYXJlIjp0cnVlLCJkb3dubG9hZCI6dHJ1ZX0sImNvbW1hbmRzIjpbXSwibG9ja1Bhc3N3b3JkIjpmYWxzZSwiaGlkZURvdGZpbGVzIjpmYWxzZSwiZGF0ZUZvcm1hdCI6ZmFsc2V9LCJpc3MiOiJGaWxlIEJyb3dzZXIiLCJleHAiOjE3NDM5MzAxNzIsImlhdCI6MTc0MzkyMjk3Mn0.I4M8cKgei3j07JMOl1ERQHKE9QB", false)
	if err != nil {
		fmt.Println(res, err)
	}
}

func TestDeleteFile(t *testing.T) {
	config := NewConfig("http://test.con:8080", "admin", "admin", "/schoollove/img.png")
	err := DeleteFile(config)
	if err != nil {
		t.Error(err)
	}
	t.Log("delete success")
}

func TestCreateDir(t *testing.T) {
	// 配置参数
	baseURL := "http://test.con:8080"
	dirPath := "/test_dir"
	authToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoxLCJsb2NhbzYWljIiwic2luZ2xlQ2xpY2siOmZhbHNlLCJwZXJtIjp7ImFkbWluIjp0cnVlLCJleGVjdXRlIjp0cnVlLZSwibW9kaWZ5Ijp0cnVlLCJkZWxldGUiOnRydWUsInNoYXJlIjp0cnVlLCJkb3dubG9hZCI6dHJ1ZX0sImNvbW1hbmRzIjpbXSwibG9ja1Bhc3N3b3JkIjpmYWxzZSwiaGlkZURvdGZpbGVzIjpmYWxzZSwiZGF0ZUZvcm1hdCI6ZmFsc2V9LCJpc3MiOiJGaWxlIEJyb3dzZXIiLCJleHAiOjE3NDM5MzAxNzIsImlhdCI6MTc0MzkyMjk3Mn0.I4M8cKgei3j07JMOl1ERQHKE9QBa7Qfm7XZNoC34Tj8"

	// 创建文件夹
	fmt.Println("开始创建文件夹...")
	err := CreateDir(baseURL, dirPath, authToken)
	if err != nil {
		fmt.Printf("创建文件夹失败: %v\n", err)
		return
	}
	fmt.Println("文件夹创建成功")

	// 验证文件夹是否存在
	fmt.Println("\n验证文件夹...")
	verifyURL := fmt.Sprintf("%s/api/resources%s", baseURL, dirPath)
	req, err := http.NewRequest("GET", verifyURL, nil)
	if err != nil {
		fmt.Printf("创建验证请求失败: %v\n", err)
		return
	}
	req.Header.Set("X-Auth", authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("验证请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("验证失败: 服务器返回状态码 %d\n", resp.StatusCode)
		return
	}

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取验证响应失败: %v\n", err)
		return
	}

	fmt.Printf("服务器响应: %s\n", string(respBody))
	fmt.Println("测试完成")
}
