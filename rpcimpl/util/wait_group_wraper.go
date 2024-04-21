package util

import (
	"sync"
)

type WGWrapper struct {
	sync.WaitGroup
}

func (w *WGWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
