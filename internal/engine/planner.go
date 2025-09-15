package engine

import (
	"fmt"

	"github.com/dev-hak/mini-terraform/internal/config"
	"github.com/dev-hak/mini-terraform/internal/providers"
	st "github.com/dev-hak/mini-terraform/internal/state"
)

type Plan struct {
	Operations []Operation
}

type Operation struct {
	Action   string // create, update, delete, noop
	Resource st.ResourceState
}

func GeneratePlan(cfg *config.Config, state *st.State) (*Plan, error) {
	if state == nil {
		state = st.NewEmptyState()
	}
	p := &Plan{}

	// configure providers
	for name, pcfg := range cfg.Providers {
		prov, err := providers.GetProvider(name)
		if err != nil {
			return nil, err
		}
		if err := prov.Configure(map[string]interface{}(pcfg)); err != nil {
			return nil, err
		}
	}

	// decide create/update/noop for resources in config
	for _, rc := range cfg.Resources {
		prov, err := providers.GetProvider(rc.Provider)
		if err != nil {
			return nil, err
		}
		_ = prov

		existing, _ := st.FindResource(state, rc.Type, rc.Name)
		if existing == nil {
			op := Operation{
				Action: "create",
				Resource: st.ResourceState{
					Type:       rc.Type,
					Name:       rc.Name,
					Provider:   rc.Provider,
					Attributes: rc.Attributes,
				},
			}
			p.Operations = append(p.Operations, op)
		} else {
			if areMapsEqual(existing.Attributes, rc.Attributes) {
				p.Operations = append(p.Operations, Operation{Action: "noop", Resource: *existing})
			} else {
				p.Operations = append(p.Operations, Operation{Action: "update", Resource: st.ResourceState{
					Type:       rc.Type,
					Name:       rc.Name,
					Provider:   rc.Provider,
					ID:         existing.ID,
					Attributes: rc.Attributes,
				}})
			}
		}
	}

	// resources in state but not in config => delete
	for _, es := range state.Resources {
		if !existsInConfig(es, cfg) {
			p.Operations = append(p.Operations, Operation{Action: "delete", Resource: es})
		}
	}
	return p, nil
}

func GeneratePlanForDestroy(cfg *config.Config, state *st.State) (*Plan, error) {
	st.EnsureState(&state)
	p := &Plan{}
	for _, rc := range cfg.Resources {
		existing, _ := st.FindResource(state, rc.Type, rc.Name)
		if existing != nil {
			p.Operations = append(p.Operations, Operation{Action: "delete", Resource: *existing})
		}
	}
	return p, nil
}

func existsInConfig(r st.ResourceState, cfg *config.Config) bool {
	for _, rc := range cfg.Resources {
		if rc.Type == r.Type && rc.Name == r.Name {
			return true
		}
	}
	return false
}

func areMapsEqual(a, b map[string]interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for k, va := range a {
		vb, ok := b[k]
		if !ok {
			return false
		}
		if fmt.Sprintf("%v", va) != fmt.Sprintf("%v", vb) {
			return false
		}
	}
	return true
}
