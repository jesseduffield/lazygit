package gocui

import "context"

type PendingTask struct {
	id         int
	Ctx        context.Context
	CancelFunc context.CancelFunc
}
