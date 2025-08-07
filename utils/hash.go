package utils

import (
	"github.com/google/uuid"
	"github.com/spaolacci/murmur3"
)

// MurmurHash64 计算key-value的MurmurHash64值
func MurmurHash64(key, value string) uint64 {
	k := murmur3.Sum32([]byte(key))
	v := murmur3.Sum32([]byte(value))

	// 把k 和 v 拼接成一个 uint64
	return uint64(k) | (uint64(v) << 32)
}

func GenerateReleaseID() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
