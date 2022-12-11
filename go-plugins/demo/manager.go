package demo

import (
	"context"
	"fmt"
)

type Manager struct {
	plugins map[string]*PluginCaller
}

func NewManager() *Manager {
	return &Manager{
		plugins: make(map[string]*PluginCaller),
	}
}

func (m *Manager) Register(name string, plugin *PluginCaller) {
	m.plugins[name] = plugin
}

func (m *Manager) CallPlugins(ctx context.Context, pluginNames []string, req *Request) (*Response, map[string]error) {

	pluginErrs := make(map[string]error)
	var res *Response
	for _, plgName := range pluginNames {
		caller, ok := m.plugins[plgName]
		if !ok {
			pluginErrs[plgName] = fmt.Errorf("plugin %q is not registered", plgName)
			return nil, pluginErrs
		}

		var err error
		res, err = caller.Call(ctx, req)
		if err != nil {
			pluginErrs[plgName] = err
		}

		if res != nil && !res.IsContinue {
			return res, pluginErrs
		}
	}

	return res, pluginErrs
}
