package glow

import (
	"github.com/spiral/roadrunner/service/rpc"
)

// ID contains default service name.
const ID = "glow"

// Service provides ability to forward stderr of the workers to stdout of roadrunner.
type Service struct {
}

// Init initializes the service.
func (s *Service) Init(r *rpc.Service) (ok bool, err error) {
	if r != nil {
		if err := r.Register(ID, &Glow{s}); err != nil {
			return false, err
		}
	}
	return true, nil
}
