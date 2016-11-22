package synapse

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
)

type WatcherStatic struct {
	WatcherCommon
	Name string
}

func NewWatcherStatic() *WatcherStatic {
	w := &WatcherStatic{}
	return w
}

func (w *WatcherStatic) GetServiceName() string {
	return w.Name
}

func (w *WatcherStatic) Init(service *Service) error {
	if err := w.CommonInit(service); err != nil {
		return errs.WithEF(err, w.fields, "Failed to init")
	}
	return nil
}

func (w *WatcherStatic) Watch(context *ContextImpl, events chan<- ServiceReport, s *Service) {
	context.doneWaiter.Add(1)
	defer context.doneWaiter.Done()
	w.service.synapse.watcherFailures.WithLabelValues(w.service.Name, PrometheusLabelWatch).Set(0)

	reportsStop := make(chan struct{})
	go w.changedToReport(reportsStop, events, s)

	<-context.stop
	logs.WithF(w.fields).Debug("Stopping watcher")
	close(reportsStop)
	logs.WithF(w.fields).Debug("Watcher stopped")
}
