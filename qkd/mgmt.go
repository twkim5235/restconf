package qkd

import (
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
)

// Manage is a bridge from model to the car application.  This is the only place where you
// couple your application code to FreeCONF.
func Manage(c *QkdnList) node.Node {

	// Powerful combination, we're letting reflect do a lot of the CRUD
	// when the yang file matches the field names.  But we extend reflection
	// to add as much custom behavior as we want
	return &nodeutil.Extend{

		// Reflection
		Base: nodeutil.ReflectChild(c),

		// CRUD drilling into child objects defined by yang file
		OnChild: func(p node.Node, r node.ChildRequest) (node.Node, error) {
			return nil, nil
		},

		// Functions
		OnAction: func(p node.Node, r node.ActionRequest) (node.Node, error) {
			return nil, nil
		},

		// Events
		OnNotify: func(p node.Node, r node.NotifyRequest) (node.NotifyCloser, error) {
			return nil, nil
		},

		// override OnEndEdit just to just to know when car has been created and
		// fully initialized so we can start the car running
		OnEndEdit: func(p node.Node, r node.NodeRequest) error {
			// allow reflection node handler to finish, this is where defaults
			// get set.
			return nil
		},
	}
}
