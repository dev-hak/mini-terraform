package engine

import (
	"fmt"

	"github.com/dev-hak/mini-terraform/internal/providers"
	st "github.com/dev-hak/mini-terraform/internal/state"
)

func Apply(plan *Plan, state *st.State) (*st.State, error) {
	if state == nil {
		state = st.NewEmptyState()
	}
	newState := &st.State{Version: state.Version, Resources: append([]st.ResourceState{}, state.Resources...)}
	for _, op := range plan.Operations {
		switch op.Action {
		case "noop":
			fmt.Printf("NOOP %s %s\n", op.Resource.Type, op.Resource.Name)
			continue
		case "create":
			prov, err := providers.GetProvider(op.Resource.Provider)
			if err != nil {
				return nil, err
			}
			id, attrs, err := prov.Create(op.Resource.Type, op.Resource.Name, op.Resource.Attributes)
			if err != nil {
				return nil, err
			}
			rs := st.ResourceState{
				Type:       op.Resource.Type,
				Name:       op.Resource.Name,
				Provider:   op.Resource.Provider,
				ID:         id,
				Attributes: attrs,
			}
			newState.Resources = append(newState.Resources, rs)
			fmt.Printf("CREATED %s %s -> %s\n", op.Resource.Type, op.Resource.Name, id)
		case "update":
			prov, err := providers.GetProvider(op.Resource.Provider)
			if err != nil {
				return nil, err
			}
			updatedAttrs, err := prov.Update(op.Resource.Type, op.Resource.ID, op.Resource.Attributes)
			if err != nil {
				return nil, err
			}
			for i := range newState.Resources {
				if newState.Resources[i].Type == op.Resource.Type && newState.Resources[i].Name == op.Resource.Name {
					newState.Resources[i].Attributes = updatedAttrs
					break
				}
			}
			fmt.Printf("UPDATED %s %s\n", op.Resource.Type, op.Resource.Name)
		case "delete":
			prov, err := providers.GetProvider(op.Resource.Provider)
			if err != nil {
				return nil, err
			}
			if err := prov.Delete(op.Resource.Type, op.Resource.ID); err != nil {
				return nil, err
			}

			idx := -1
			for i := range newState.Resources {
				if newState.Resources[i].Type == op.Resource.Type && newState.Resources[i].Name == op.Resource.Name {
					idx = i
					break
				}
			}
			if idx >= 0 {
				newState.Resources = append(newState.Resources[:idx], newState.Resources[idx+1:]...)
			}
			fmt.Printf("DELETED %s %s\n", op.Resource.Type, op.Resource.Name)
		default:
			return nil, fmt.Errorf("unknown action: %s", op.Action)
		}
	}
	return newState, nil
}
