package utils

import (
	"regexp"
	"strings"
)

// ValidateEnv 验证环境字符串
func ValidateEnv(env string) bool {
	// 只有 dev, pro, uat, fat 几种环境
	return env == "dev" || env == "pro"
}

// ValidateAppID 验证应用ID格式
func ValidateAppName(appName string) bool {
	if len(appName) < 2 || len(appName) > 64 {
		return false
	}
	// 不能_ 和 - 开头
	if strings.ContainsRune("_-.", rune(appName[0])) {
		return false
	}
	// 只允许字母、数字、下划线、中划线
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_.-]+$`, appName)
	return matched
}

// ValidateClusterName 验证集群名称格式
func ValidateClusterName(clusterName string) bool {
	if len(clusterName) < 2 || len(clusterName) > 64 {
		return false
	}
	// 不能_ 和 - 开头
	if strings.ContainsRune("_-.", rune(clusterName[0])) {
		return false
	}
	// 只允许字母、数字、下划线、中划线
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_.-]+$`, clusterName)
	return matched
}

// ValidateNamespaceName 验证命名空间名称格式
func ValidateNamespaceName(namespaceName string) bool {
	if len(namespaceName) < 2 || len(namespaceName) > 64 {
		return false
	}
	// 不能_ 和 - 开头
	if strings.ContainsRune("_-.", rune(namespaceName[0])) {
		return false
	}
	// 只允许字母、数字、下划线、中划线
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_.-]+$`, namespaceName)
	return matched
}

// ValidateConfigKey 验证配置Key格式
func ValidateItemKey(key string) bool {
	if len(key) < 2 || len(key) > 255 {
		return false
	}
	// 不能_ 和 - 开头
	if strings.ContainsRune("_-.=@#+", rune(key[0])) {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_.-]+$`, key)
	return matched
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	if len(email) == 0 {
		return true // 可选字段
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	return matched
}

// ValidatePhone 验证手机号格式
func ValidatePhone(phone string) bool {
	if len(phone) == 0 {
		return true // 可选字段
	}
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}
