package providers

import "errors"

type Provider interface {
	Name() string
	Configure(map[string]interface{}) error
	Create(resourceType, name string, attrs map[string]interface{}) (string, map[string]interface{}, error)
	Read(resourceType, id string) (map[string]interface{}, error)
	Update(resourceType, id string, attrs map[string]interface{}) (map[string]interface{}, error)
	Delete(resourceType, id string) error
}

var registry = make(map[string]Provider)

func RegisterProvider(name string, p Provider) {
	registry[name] = p
}

func GetProvider(name string) (Provider, error) {
	if p, ok := registry[name]; ok {
		return p, nil
	}
	return nil, errors.New("provider not found: " + name)
}
