package handler

import "quiver/services"

// ReleaseHandler 命名空间控制器
type ReleaseHandler struct {
	releaseService *services.ReleaseService
}

// NewReleaseHandler 创建命名空间控制器实例
func NewReleaseHandler() *ReleaseHandler {
	return &ReleaseHandler{
		releaseService: services.NewReleaseService(),
	}
}
