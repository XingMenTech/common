package utils

import (
	"regexp"
	"strings"
)

//var region *ip2region.Ip2Region

const (
	china = "中国"
	india = "印度"
)

func init() {
	//var err error
	//region, err = ip2region.New("conf/ip2region.db")
	//if err != nil {
	//	fmt.Println(fmt.Sprintf("[services.ip] load db ,err :%v", err))
	//	return
	//}
	//if region != nil {
	//	region.Close()
	//}
}

//func IPInterception(memberIP string) (pass bool) {
//	if region != nil {
//		ip, err := region.MemorySearch(memberIP)
//		//查询不到IP，放过
//		if err != nil {
//			return true
//		}
//
//		if ip.Country == china {
//			return false
//		}
//
//		return true
//	} else {
//		return true
//	}
//}

// IPValidator 提供IP地址验证功能
type IPValidator struct {
	singleIPRegex *regexp.Regexp
	ipListRegex   *regexp.Regexp
}

// NewIPValidator 创建新的IP验证器实例
func NewIPValidator() *IPValidator {
	return &IPValidator{
		//(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[1-9])
		//(25[0-5]|2[0-4][0-9]|[1-9]?[0-9][0-9]?)
		singleIPRegex: regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[1-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[1-9]))$`),
		ipListRegex:   regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[1-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[1-9]))(,\s*(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[1-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[1-9]))*$`),
	}
}

// ValidateSingleIP 验证单个IP地址
func (v *IPValidator) ValidateSingleIP(ip string) bool {
	return v.singleIPRegex.MatchString(strings.TrimSpace(ip))
}

// ValidateIPList 验证IP地址列表
func (v *IPValidator) ValidateIPList(input string) (bool, []string, []string) {
	trimmedInput := strings.TrimSpace(input)
	if trimmedInput == "" {
		return false, nil, nil
	}

	// 首先验证整体格式
	if !v.ipListRegex.MatchString(trimmedInput) {
		return false, nil, nil
	}

	// 分割并验证每个IP
	ips := strings.Split(trimmedInput, ",")
	validIPs := make([]string, 0)
	invalidIPs := make([]string, 0)

	for _, ip := range ips {
		trimmedIP := strings.TrimSpace(ip)
		if trimmedIP == "" {
			continue
		}
		if v.ValidateSingleIP(trimmedIP) {
			validIPs = append(validIPs, trimmedIP)
		} else {
			invalidIPs = append(invalidIPs, trimmedIP)
		}
	}

	return len(invalidIPs) == 0, validIPs, invalidIPs
}

// ValidateIPListStrict 严格验证IP地址列表（不允许空格）
func (v *IPValidator) ValidateIPListStrict(input string) bool {
	strictRegex := regexp.MustCompile(`^(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[1-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[1-9])(,(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[1-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[1-9]))*$`)
	return strictRegex.MatchString(input)
}

// FormatIPList 格式化IP地址列表
func FormatIPList(input string) string {
	ips := strings.Split(input, ",")
	formatted := make([]string, 0)

	for _, ip := range ips {
		trimmed := strings.TrimSpace(ip)
		if trimmed != "" {
			formatted = append(formatted, trimmed)
		}
	}

	return strings.Join(formatted, ", ")
}

// ExtractIPs 从字符串中提取所有IP地址
func ExtractIPs(input string) []string {
	ips := strings.Split(input, ",")
	result := make([]string, 0)

	for _, ip := range ips {
		trimmed := strings.TrimSpace(ip)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// CountValidIPs 统计有效IP地址数量
func CountValidIPs(input string) int {
	validator := NewIPValidator()
	ips := ExtractIPs(input)
	count := 0

	for _, ip := range ips {
		if validator.ValidateSingleIP(ip) {
			count++
		}
	}

	return count
}
