package node

import (
	"fmt"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/meta"
)

type ContentConstraint int

const (
	ContentAll ContentConstraint = iota
	ContentOperational
	ContentConfig
)

func NewContentConstraint(initialPath *Path, expression string) (c ContentConstraint, err error) {
	switch expression {
	case "config":
		return ContentConfig, nil
	case "nonconfig":
		return ContentOperational, nil
	case "all":
		return ContentAll, nil
	}
	return ContentAll, fmt.Errorf("%w. Invalid content constraint: '%s'", fc.BadRequestError, expression)
}

func (self ContentConstraint) CheckContainerPreConstraints(r *ChildRequest) (bool, error) {
	// config containers may have operational fields so always pass on operational
	if r.IsNavigation() || self == ContentAll || self == ContentOperational {
		return true, nil
	}

	var isConfig bool
	// meta.Module does not implement HasDetails, but spec implies yes
	if d, hasDets := r.Meta.(meta.HasDetails); !hasDets {
		isConfig = true
	} else {
		isConfig = d.Config()
	}
	return isConfig, nil
}

func (self ContentConstraint) CheckFieldPreConstraints(r *FieldRequest, hnd *ValueHandle) (bool, error) {
	if r.IsNavigation() || self == ContentAll {
		return true, nil
	}
	isConfig := r.Meta.(meta.HasDetails).Config()
	return (isConfig && self == ContentConfig) || (!isConfig && self == ContentOperational), nil
}
