package state

import (
	"github.com/mohamadrezamomeni/telecraft/pkg/telecrafterror"
)

type Repo interface {
	Set(key string, state *State) error
	Get(key string) (*State, bool)
	Delete(key string) error
}

func NewRepository(repoType string) (Repo, error) {
	switch repoType {
	case "cache":
		return newCache(), nil
	}
	return nil, telecrafterror.Scope("").Input(repoType).NotFound().ErrorWrite()
}
