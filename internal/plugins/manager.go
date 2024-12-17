package plugins

import (
	"fmt"
	"plugin"
	"sync"
)

type PluginManager struct {
	plugins map[string]Plugin
	mu      sync.RWMutex
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
	}
}

func (pm *PluginManager) LoadPlugin(path string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Load the plugin
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %w", path, err)
	}

	// Look up the Plugin symbol
	symPlugin, err := p.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin %s does not export 'Plugin' symbol: %w", path, err)
	}

	// Assert that the symbol is a Plugin
	plugin, ok := symPlugin.(Plugin)
	if !ok {
		return fmt.Errorf("plugin %s does not implement Plugin interface", path)
	}

	// Store the plugin
	pm.plugins[path] = plugin
	return nil
}

func (pm *PluginManager) ExecuteHook(hookType HookType, data map[string]interface{}) (map[string]interface{}, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := data
	var err error

	for _, p := range pm.plugins {
		switch hookType {
		case ENTER_RULE:
			result, err = p.OnEnterRule(result)
		case ENTER_LINE:
			result, err = p.OnEnterLine(result)
		case EXIT_LINE:
			result, err = p.OnExitLine(result)
		case EXIT_RULE:
			result, err = p.OnExitRule(result)
		}

		if err != nil {
			return nil, fmt.Errorf("plugin execution failed: %w", err)
		}
	}

	return result, nil
}
