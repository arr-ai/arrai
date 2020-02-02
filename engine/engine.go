package engine

import (
	"sync/atomic"

	"github.com/go-errors/errors"
	"github.com/sirupsen/logrus"

	"github.com/arr-ai/arrai/rel"
)

var (
	// Root holds the variable name used for the database variable.
	Root = "$"

	lastID = uint64(0)
)

var globals = map[string]func(rel.Value) rel.Value{
	"str": func(value rel.Value) rel.Value {
		return rel.NewString([]rune(value.String()))
	},
}

// Engine holds a database variable, allowing updates and publishing changes.
type Engine struct {
	updateDb      chan updateRequest
	addWatcher    chan *watcher
	removeWatcher chan uint64
	stop          chan struct{}
	hangup        chan struct{}
}

// Start starts a new Engine.
func Start() *Engine {
	e := &Engine{
		make(chan updateRequest),
		make(chan *watcher),
		make(chan uint64),
		make(chan struct{}),
		make(chan struct{}),
	}

	// The engine
	go func() {
		logrus.Info("Engine starting")
		defer logrus.Info("Engine dying")
		watchers := map[uint64]*watcher{}

		closeAllWatchers := func() {
			for _, w := range watchers {
				w.close()
			}
			watchers = map[uint64]*watcher{}
		}

		defer closeAllWatchers()

		global := rel.EmptyScope.With(Root, rel.None)
		for name, fn := range globals {
			global = global.With(name, rel.NewNativeFunction(name, fn))
		}

		for {
			logrus.Info("Engine select")
			select {
			case w := <-e.addWatcher:
				logrus.Info("Add watcher")
				watchers[w.id] = w
				w.update(global)
			case id := <-e.removeWatcher:
				logrus.Info("Remove watcher")
				watchers[id].close()
				delete(watchers, id)
			case req := <-e.updateDb:
				logrus.Info("Update DB")
				logrus.Infof("-> %#v", req.expr)
				logrus.Infof("-> %s", req.expr)
				value, err := req.expr.Eval(global)
				if err != nil {
					req.failed <- err
					continue
				}
				req.failed <- nil
				global = global.With(Root, value)
				for i, w := range watchers {
					logrus.Infof("Update watcher %d", i)
					w.update(global)
				}
			case <-e.stop:
				logrus.Infof("Stop")
				return
			case <-e.hangup:
				logrus.Infof("Hangup")
				closeAllWatchers()
			}
		}
	}()

	return e
}

// Stop stops the engine.
func (e *Engine) Stop() {
	e.stop <- struct{}{}
}

// Hangup hangs up on all observers.
func (e *Engine) Hangup() {
	e.hangup <- struct{}{}
}

// Update updates the database variable to equal the given expression.
func (e *Engine) Update(expr rel.Expr) error {
	failed := make(chan error)
	e.updateDb <- updateRequest{expr, failed}
	return <-failed
}

// Observe registers and returns an Observation on the given expression.
func (e *Engine) Observe(
	expr rel.Expr,
	onupdate func(rel.Value) error,
	onclose func(error),
) func() {
	id := atomic.AddUint64(&lastID, 1)
	cancel := func() {
		e.removeWatcher <- id
	}
	e.addWatcher <- &watcher{id, cancel, expr, onupdate, onclose}
	return cancel
}

type updateRequest struct {
	expr   rel.Expr
	failed chan<- error
}

type watcher struct {
	id       uint64
	cancel   func()
	expr     rel.Expr
	onupdate func(rel.Value) error
	onclose  func(error)
}

func (w *watcher) update(global rel.Scope) {
	defer func() {
		if err := recover(); err != nil {
			w.onclose(errors.WrapPrefix(err, "update panic", 0))
		}
	}()

	value, err := w.expr.Eval(global)
	if err != nil {
		w.cancel()
		w.onclose(err)
		return
	}

	if err = w.onupdate(value.(rel.Value)); err != nil {
		w.cancel()
	}
}

func (w *watcher) close() {
	w.onclose(nil)
}
