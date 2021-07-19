package restconf

import (
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/val"
)

func CallHomeNode(ch *CallHome) node.Node {
	options := ch.Options()
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(&options),
		OnField: func(p node.Node, r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "registered":
				hnd.Val = val.Bool(ch.Registered)
			default:
				return p.Field(r, hnd)
			}
			return nil
		},
		OnEndEdit: func(p node.Node, r node.NodeRequest) error {
			if err := p.EndEdit(r); err != nil {
				return err
			}
			return ch.ApplyOptions(options)
		},
	}
}
