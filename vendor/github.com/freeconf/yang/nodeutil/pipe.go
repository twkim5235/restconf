package nodeutil

import (
	"errors"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/val"
)

type Pipe struct {
	messages chan *pipeMessage
	position *pipeMessage
}

type tok int

const (
	PipeSelect tok = iota + 1
	PipeListItem
	PipeLeaf
	PipeEnd
)

func NewPipe() *Pipe {
	return &Pipe{}
}

func (self *Pipe) peek() *pipeMessage {
	if self.position == nil {
		self.position = <-self.messages
		if self.position.tok == PipeEnd {
			close(self.messages)
		}
	}
	return self.position
}

type pipeMessage struct {
	tok   tok
	ident string
	val   val.Value
	key   []val.Value
	err   error
}

func (self *Pipe) consume() {
	self.position = nil
}

func (self *Pipe) Close(err error) {
	defer func() {
		if r := recover(); r != nil {
			// channel was probably already closed so log err if there was one
			if err != nil {
				fc.Err.Printf(err.Error())
			}
		}
	}()
	self.messages <- &pipeMessage{
		tok: PipeEnd,
		err: err,
	}
}

func (self *Pipe) PullPush() (node.Node, node.Node) {
	self.messages = make(chan *pipeMessage)
	pull := &Basic{}
	push := &Basic{}
	pull.OnChild = func(r node.ChildRequest) (node.Node, error) {
		if r.New {
			return nil, errors.New("Not a writer")
		}
		msg := self.peek()
		if msg.tok != PipeSelect || msg.ident != r.Meta.Ident() {
			return nil, msg.err
		}
		defer self.consume()
		return pull, msg.err
	}
	pull.OnNext = func(r node.ListRequest) (node.Node, []val.Value, error) {
		if r.New {
			return nil, nil, errors.New("Not a writer")
		}
		msg := self.peek()
		if msg.tok != PipeListItem {
			return nil, nil, msg.err
		}
		defer self.consume()
		return pull, msg.key, msg.err
	}
	pull.OnField = func(r node.FieldRequest, hnd *node.ValueHandle) (err error) {
		msg := self.peek()
		if msg.tok != PipeLeaf || msg.ident != r.Meta.Ident() {
			return msg.err
		}
		defer self.consume()
		hnd.Val = msg.val
		return msg.err
	}
	push.OnChild = func(r node.ChildRequest) (node.Node, error) {
		if !r.New {
			return nil, nil
		}
		self.messages <- &pipeMessage{
			tok:   PipeSelect,
			ident: r.Meta.Ident(),
		}
		return push, nil
	}
	push.OnField = func(r node.FieldRequest, hnd *node.ValueHandle) (err error) {
		self.messages <- &pipeMessage{
			tok:   PipeLeaf,
			val:   hnd.Val,
			ident: r.Meta.Ident(),
		}
		return nil
	}
	push.OnNext = func(r node.ListRequest) (node.Node, []val.Value, error) {
		if !r.New {
			return nil, nil, nil
		}
		self.messages <- &pipeMessage{
			tok:   PipeListItem,
			key:   r.Key,
			ident: r.Meta.Ident(),
		}
		return push, r.Key, nil
	}
	return pull, push
}
