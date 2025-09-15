package providers

import (
	"fmt"

	"github.com/dev-hak/mini-terraform/internal/providers"
)

type AWSProvider struct {
	config map[string]interface{}
}

func NewAWSProvider() providers.Provider {
	return &AWSProvider{}
}

func (a *AWSProvider) Name() string { return "aws" }

func (a *AWSProvider) Configure(cfg map[string]interface{}) error {
	a.config = cfg
	// Integrate AWS SDK here later
	return nil
}

func (a *AWSProvider) Create(resourceType, name string, attrs map[string]interface{}) (string, map[string]interface{}, error) {
	return "", nil, fmt.Errorf("aws provider Create not implemented - extend with AWS SDK")
}

func (a *AWSProvider) Read(resourceType, id string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("aws provider Read not implemented")
}

func (a *AWSProvider) Update(resourceType, id string, attrs map[string]interface{}) (map[string]interface{}, error) {
	return nil, fmt.Errorf("aws provider Update not implemented")
}

func (a *AWSProvider) Delete(resourceType, id string) error {
	return fmt.Errorf("aws provider Delete not implemented")
}
