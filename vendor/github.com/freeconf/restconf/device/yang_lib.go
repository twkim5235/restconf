package device

import (
	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/val"
)

type ResolveModule interface {
	ResolveModuleHnd(hnd ModuleHnd) (*meta.Module, error)
}

func LoadModules(ietfYangLib *node.Browser, resolver ResolveModule) (map[string]*meta.Module, error) {
	mods := make(map[string]*meta.Module)
	n := loadModulesListNode(mods, resolver)
	if err := ietfYangLib.Root().Find("modules-state/module").InsertInto(n).LastErr; err != nil {
		return nil, err
	}
	return mods, nil
}

func loadModulesListNode(mods map[string]*meta.Module, resolver ResolveModule) node.Node {
	return &nodeutil.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			key := r.Key
			if r.New {
				hnd := ModuleHnd{Name: r.Key[0].String()}
				return loadModuleNode(mods, resolver, &hnd), key, nil
			}
			return nil, nil, nil
		},
	}
}

func loadModuleNode(mods map[string]*meta.Module, resolver ResolveModule, hnd *ModuleHnd) node.Node {
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(hnd),
		OnEndEdit: func(p node.Node, r node.NodeRequest) error {
			if err := p.EndEdit(r); err != nil {
				return err
			}
			mod, err := resolver.ResolveModuleHnd(*hnd)
			if err != nil {
				return err
			}
			mods[mod.Ident()] = mod
			return nil
		},
	}
}

type ModuleHnd struct {
	Name      string
	Schema    string
	Revision  string
	Namespace string
}
