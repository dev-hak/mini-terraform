package providers

import (
	"errors"
	"os"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/dev-hak/mini-terraform/internal/providers"
)

type VPSProvider struct {
	config map[string]interface{}
}

func NewVPSProvider() providers.Provider {
	return &VPSProvider{}
}

func (p *VPSProvider) Name() string { return "vps" }

func (p *VPSProvider) Configure(cfg map[string]interface{}) error {
	p.config = cfg
	return nil
}

// resource type: vps
// attributes expected: host, user, private_key (path), command
func (p *VPSProvider) Create(resourceType, id string, attrs map[string]interface{}) (string, map[string]interface{}, error) {
	if resourceType != "vps" {
		return "", nil, errors.New("unsupported vps resource: " + resourceType)
	}
	host, _ := attrs["host"].(string)
	user, _ := attrs["user"].(string)
	keyPath, _ := attrs["private_key"].(string)
	cmdStr, _ := attrs["command"].(string)
	if host == "" || user == "" || cmdStr == "" {
		return "", nil, errors.New("host, user, and command required")
	}
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return "", nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", nil, err
	}
	clientConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	client, err := ssh.Dial("tcp", host, clientConfig)
	if err != nil {
		return "", nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		return "", nil, err
	}
	defer session.Close()

	out, err := session.CombinedOutput(cmdStr)
	if err != nil {
		return "", nil, err
	}

	stateAttrs := map[string]interface{}{
		"host": host, "user": user, "out": string(out), "cmd": cmdStr,
	}
	return id, stateAttrs, nil
}

func (p *VPSProvider) Read(resourceType, id string) (map[string]interface{}, error) {
	return map[string]interface{}{"id": id}, nil
}

func (p *VPSProvider) Update(resourceType, id string, attrs map[string]interface{}) (map[string]interface{}, error) {
	_, updatedAttrs, err := p.Create(resourceType, id, attrs)
	return updatedAttrs, err
}

func (p *VPSProvider) Delete(resourceType, id string) error {
	// nothing by default
	return nil
}
