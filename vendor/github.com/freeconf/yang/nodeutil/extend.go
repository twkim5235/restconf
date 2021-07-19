package nodeutil

import (
	"context"

	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/val"
)

// Extend let's you alter any Node behavior including the nodeutil it creates.
type Extend struct {
	Base        node.Node
	OnNext      ExtendNextFunc
	OnChild     ExtendChildFunc
	OnField     ExtendFieldFunc
	OnChoose    ExtendChooseFunc
	OnAction    ExtendActionFunc
	OnNotify    ExtendNotifyFunc
	OnExtend    ExtendFunc
	OnPeek      ExtendPeekFunc
	OnBeginEdit ExtendBeginEditFunc
	OnEndEdit   ExtendEndEditFunc
	OnDelete    ExtendDeleteFunc
	OnContext   ExtendContextFunc
}

func (e *Extend) Child(r node.ChildRequest) (node.Node, error) {
	var err error
	var child node.Node
	if e.OnChild == nil {
		child, err = e.Base.Child(r)
	} else {
		child, err = e.OnChild(e.Base, r)
	}
	if child == nil || err != nil {
		return child, err
	}
	if e.OnExtend != nil {
		child, err = e.OnExtend(e, r.Selection, r.Meta, child)
	}
	return child, err
}

func (e *Extend) Next(r node.ListRequest) (child node.Node, key []val.Value, err error) {
	if e.OnNext == nil {
		child, key, err = e.Base.Next(r)
	} else {
		child, key, err = e.OnNext(e.Base, r)
	}
	if child == nil || err != nil {
		return
	}
	if e.OnExtend != nil {
		child, err = e.OnExtend(e, r.Selection, r.Meta, child)
	}
	return
}

func (e *Extend) Extend(n node.Node) node.Node {
	extendedChild := *e
	extendedChild.Base = n
	return &extendedChild
}

func (e *Extend) Field(r node.FieldRequest, hnd *node.ValueHandle) error {
	if e.OnField == nil {
		return e.Base.Field(r, hnd)
	} else {
		return e.OnField(e.Base, r, hnd)
	}
}

func (e *Extend) Choose(sel node.Selection, choice *meta.Choice) (*meta.ChoiceCase, error) {
	if e.OnChoose == nil {
		return e.Base.Choose(sel, choice)
	} else {
		return e.OnChoose(e.Base, sel, choice)
	}
}

func (e *Extend) Action(r node.ActionRequest) (output node.Node, err error) {
	if e.OnAction == nil {
		return e.Base.Action(r)
	} else {
		return e.OnAction(e.Base, r)
	}
}

func (e *Extend) Notify(r node.NotifyRequest) (closer node.NotifyCloser, err error) {
	if e.OnNotify == nil {
		return e.Base.Notify(r)
	} else {
		return e.OnNotify(e.Base, r)
	}
}

func (e *Extend) Delete(r node.NodeRequest) error {
	if e.OnDelete == nil {
		return e.Base.Delete(r)
	}
	return e.OnDelete(e.Base, r)
}

func (e *Extend) BeginEdit(r node.NodeRequest) error {
	if e.OnBeginEdit == nil {
		return e.Base.BeginEdit(r)
	}
	return e.OnBeginEdit(e.Base, r)
}

func (e *Extend) EndEdit(r node.NodeRequest) error {
	if e.OnEndEdit == nil {
		return e.Base.EndEdit(r)
	}
	return e.OnEndEdit(e.Base, r)
}

func (e *Extend) Context(sel node.Selection) context.Context {
	if e.OnContext == nil {
		return e.Base.Context(sel)
	}
	return e.OnContext(e.Base, sel)
}

func (e *Extend) Peek(sel node.Selection, consumer interface{}) interface{} {
	if e.OnPeek == nil {
		return e.Base.Peek(sel, consumer)
	}
	return e.OnPeek(e.Base, sel, consumer)
}

type ExtendNextFunc func(parent node.Node, r node.ListRequest) (next node.Node, key []val.Value, err error)
type ExtendChildFunc func(parent node.Node, r node.ChildRequest) (child node.Node, err error)
type ExtendFieldFunc func(parent node.Node, r node.FieldRequest, hnd *node.ValueHandle) error
type ExtendChooseFunc func(parent node.Node, sel node.Selection, choice *meta.Choice) (m *meta.ChoiceCase, err error)
type ExtendActionFunc func(parent node.Node, r node.ActionRequest) (output node.Node, err error)
type ExtendNotifyFunc func(parent node.Node, r node.NotifyRequest) (closer node.NotifyCloser, err error)
type ExtendFunc func(e *Extend, sel node.Selection, m meta.HasDefinitions, child node.Node) (node.Node, error)
type ExtendPeekFunc func(parent node.Node, sel node.Selection, consumer interface{}) interface{}
type ExtendBeginEditFunc func(parent node.Node, r node.NodeRequest) error
type ExtendEndEditFunc func(parent node.Node, r node.NodeRequest) error
type ExtendDeleteFunc func(parent node.Node, r node.NodeRequest) error
type ExtendContextFunc func(parent node.Node, s node.Selection) context.Context
