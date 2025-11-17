// helper-------------------------------------
// @file      : tool.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/8 21:23
// -------------------------------------------

package helper

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
	"net"
	"net/url"
	"strings"
)

func Sha256Hex(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

func HashXX64String(input string) string {
	hash := xxhash.Sum64String(input)
	return fmt.Sprintf("%x", hash)
}

func CalculateMD5FromContent(content string) string {
	hash := md5.New()
	hash.Write([]byte(content))
	return hex.EncodeToString(hash.Sum(nil))
}

func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func RemoveArrayDuplicates(input []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, str := range input {
		if !seen[str] {
			seen[str] = true
			result = append(result, str)
		}
	}

	return result
}

func GetRootDomain(input string) (string, error) {
	u, err := SafeParseURL(input)
	if err != nil {
		return input, fmt.Errorf("URL 解析失败: %w", err)
	}

	hostname := u.Hostname()
	if hostname == "" {
		return input, fmt.Errorf("无法获取 Hostname")
	}

	// 是 IP 则直接返回
	if ip := net.ParseIP(hostname); ip != nil {
		return ip.String(), nil
	}

	// 提取有效根域名
	rootDomain, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		return input, fmt.Errorf("根域名解析错误: %w", err)
	}

	return rootDomain, nil
}

// SafeParseURL 安全解析 URL，处理非法 %、中文域名、空格等问题
func SafeParseURL(input string) (*url.URL, error) {
	if input == "" {
		return nil, fmt.Errorf("输入为空")
	}

	input = strings.TrimSpace(input)
	input = sanitizePercent(input)
	input = strings.ReplaceAll(input, " ", "%20")

	// 补充协议头
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		input = "https://" + input
	}

	u, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("URL 解析失败: %w", err)
	}

	// IDNA 处理中文域名
	asciiHost, err := idna.ToASCII(u.Hostname())
	if err != nil {
		return nil, fmt.Errorf("域名 IDNA 转换失败: %w", err)
	}

	// 补上端口（如有）
	u.Host = asciiHost + portSuffix(u.Host)

	return u, nil
}

// sanitizePercent 修复非法的 % 编码
func sanitizePercent(input string) string {
	var builder strings.Builder
	for i := 0; i < len(input); i++ {
		if input[i] == '%' {
			if i+2 >= len(input) || !isHexDigit(input[i+1]) || !isHexDigit(input[i+2]) {
				builder.WriteString("%25")
			} else {
				builder.WriteByte('%')
			}
		} else {
			builder.WriteByte(input[i])
		}
	}
	return builder.String()
}

// 判断是否为合法十六进制字符
func isHexDigit(b byte) bool {
	return ('0' <= b && b <= '9') || ('a' <= b && b <= 'f') || ('A' <= b && b <= 'F')
}

// 提取原 host 的端口号（如有）
func portSuffix(host string) string {
	if colon := strings.LastIndex(host, ":"); colon != -1 && colon < len(host)-1 {
		port := host[colon+1:]
		if _, err := fmt.Sscanf(port, "%d", new(int)); err == nil {
			return ":" + port
		}
	}
	return ""
}
