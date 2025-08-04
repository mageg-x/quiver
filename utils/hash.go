package utils

import (
	"github.com/spaolacci/murmur3"
)

// MurmurHash64 计算key-value的MurmurHash64值
func MurmurHash64(key, value string) int64 {
	// 将key和value组合
	data := key + ":" + value

	// 计算hash值
	hash := murmur3.Sum64([]byte(data))

	// 转换为int64（有符号）
	return int64(hash)
}

// MurmurHash64Bytes 直接计算字节数组的hash值
func MurmurHash64Bytes(data []byte) int64 {
	hash := murmur3.Sum64(data)
	return int64(hash)
}
