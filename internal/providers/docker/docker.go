package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/dev-hak/mini-terraform/internal/providers"
)

type DockerProvider struct {
	config map[string]interface{}
}

func NewDockerProvider() providers.Provider {
	return &DockerProvider{}
}

func (d *DockerProvider) Name() string { return "docker" }

func (d *DockerProvider) Configure(cfg map[string]interface{}) error {
	d.config = cfg
	return nil
}

// supports resource type "docker_container"
func (d *DockerProvider) Create(resourceType, name string, attrs map[string]interface{}) (string, map[string]interface{}, error) {
	if resourceType != "docker_container" {
		return "", nil, errors.New("unsupported docker resource: " + resourceType)
	}
	image, _ := attrs["image"].(string)
	if image == "" {
		return "", nil, errors.New("docker image required")
	}
	args := []string{"run", "-d", "--name", name}
	if portsRaw, ok := attrs["ports"].([]interface{}); ok {
		for _, p := range portsRaw {
			args = append(args, "-p", fmt.Sprintf("%v", p))
		}
	}
	if envRaw, ok := attrs["env"].([]interface{}); ok {
		for _, e := range envRaw {
			args = append(args, "-e", fmt.Sprintf("%v", e))
		}
	}
	args = append(args, image)
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", nil, fmt.Errorf("docker run failed: %v: %s", err, string(out))
	}
	id := strings.TrimSpace(string(out))
	stateAttrs := map[string]interface{}{
		"image":        image,
		"ports":        attrs["ports"],
		"env":          attrs["env"],
		"container_id": id,
	}
	return id, stateAttrs, nil
}

func (d *DockerProvider) Read(resourceType, id string) (map[string]interface{}, error) {
	cmd := exec.Command("docker", "inspect", id)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var parsed []interface{}
	if err := json.Unmarshal(out, &parsed); err != nil {
		return nil, err
	}
	if len(parsed) == 0 {
		return nil, errors.New("not found")
	}
	m := parsed[0].(map[string]interface{})
	return map[string]interface{}{"inspect": m}, nil
}

func (d *DockerProvider) Update(resourceType, id string, attrs map[string]interface{}) (map[string]interface{}, error) {
	// naive update: remove container and recreate with new attributes (preserve name)
	// get current name from inspect
	ins, err := d.Read(resourceType, id)
	if err != nil {
		return nil, err
	}
	// try to extract name
	name := ""
	if inspect, ok := ins["inspect"].(map[string]interface{}); ok {
		if nm, ok := inspect["Name"].(string); ok {
			name = strings.TrimPrefix(nm, "/")
		}
	}
	if name == "" {
		return nil, errors.New("cannot determine container name for update")
	}
	// remove
	if err := exec.Command("docker", "rm", "-f", id).Run(); err != nil {
		return nil, err
	}
	// re-create with attributes
	_, st, err := d.Create(resourceType, name, attrs)
	return st, err
}

func (d *DockerProvider) Delete(resourceType, id string) error {
	return exec.Command("docker", "rm", "-f", id).Run()
}
