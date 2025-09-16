package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Variables map[string]interface{} `json:"variables,omitempty"`
	Providers map[string]ProviderCfg `json:"providers,omitempty"`
	Resources []ResourceCfg          `json:"resources"`
}

type ProviderCfg map[string]interface{}

type ResourceCfg struct {
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Provider   string                 `json:"provider"`
	Attributes map[string]interface{} `json:"attributes"`
}

func LoadConfig(path string, varFile string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	// merge var-file
	if varFile != "" {
		vb, err := os.ReadFile(varFile)
		if err != nil {
			return nil, err
		}
		var vmap map[string]interface{}
		if err := json.Unmarshal(vb, &vmap); err != nil {
			return nil, err
		}
		if cfg.Variables == nil {
			cfg.Variables = make(map[string]interface{})
		}
		for k, v := range vmap {
			cfg.Variables[k] = v
		}
	}

	for name, pcfg := range cfg.Providers {
		kv := map[string]interface{}(pcfg)
		interpolate(kv, cfg.Variables)
		cfg.Providers[name] = ProviderCfg(kv)
	}
	for i := range cfg.Resources {
		interpolate(cfg.Resources[i].Attributes, cfg.Variables)
		cfg.Resources[i].Name = interpolateString(cfg.Resources[i].Name, cfg.Variables)
	}
	return &cfg, nil
}

func interpolate(m map[string]interface{}, vars map[string]interface{}) {
	for k, v := range m {
		switch t := v.(type) {
		case string:
			m[k] = interpolateString(t, vars)
		case map[string]interface{}:
			interpolate(t, vars)
		case []interface{}:
			for i, elem := range t {
				if s, ok := elem.(string); ok {
					t[i] = interpolateString(s, vars)
				} else if mm, ok := elem.(map[string]interface{}); ok {
					interpolate(mm, vars)
				}
			}
		}
	}
}

func interpolateString(s string, vars map[string]interface{}) string {
	out := s
	if vars == nil {
		return out
	}
	for k, v := range vars {
		placeholder := fmt.Sprintf("${var.%s}", k)
		if strings.Contains(out, placeholder) {
			out = strings.ReplaceAll(out, placeholder, fmt.Sprintf("%v", v))
		}
	}
	return out
}
