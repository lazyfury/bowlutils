package module

import (
	"context"
	"sync"
)

// Module 模块接口，所有模块都需要实现此接口
type Module interface {
	Start(ctx context.Context, wg *sync.WaitGroup) error
	Stop() error
}

// ModuleManager 模块管理器，用于管理所有模块的生命周期
type ModuleManager struct {
	modules map[string]Module
}

// NewModuleManager 创建新的模块管理器
func NewModuleManager() *ModuleManager {
	return &ModuleManager{
		modules: make(map[string]Module),
	}
}

// RegisterModule 注册模块
func (m *ModuleManager) RegisterModule(name string, module Module) {
	m.modules[name] = module
}

// StartAll 启动所有模块
func (m *ModuleManager) StartAll(ctx context.Context, wg *sync.WaitGroup) error {
	for _, module := range m.modules {
		if err := module.Start(ctx, wg); err != nil {
			return err
		}
	}
	return nil
}

// StopAll 停止所有模块
func (m *ModuleManager) StopAll() error {
	for _, module := range m.modules {
		if err := module.Stop(); err != nil {
			return err
		}
	}
	return nil
}
