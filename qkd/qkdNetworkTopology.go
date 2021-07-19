package qkd

import (
	"container/list"
)

type QkdnList struct {
	QkdnId int64

	listeners *list.List
}

type QkdnListListner func(q QkdnList)

func Mew() *QkdnList {
	q := &QkdnList{
		listeners: list.New(),
	}
	q.newQkdn()
	return q
}

func (q *QkdnList) newQkdn() {
	q.QkdnId = 10000001
}
