package log

import (
	"context"

	"github.com/clcng/bitcoin-wallet/pkg/util"
	"github.com/rs/zerolog"
)

type logEvent struct {
	*zerolog.Event
}

func LogWrite(logger *zerolog.Event) *logEvent {
	return &logEvent{
		Event: logger,
	}
}

func (ev *logEvent) Build(fcList ...func(*logEvent) *logEvent) *logEvent {
	if fcList == nil {
		return ev
	}
	for _, f := range fcList {
		if f == nil {
			continue
		}
		ev = f(ev)
	}
	return ev
}

func (ev *logEvent) Pkg(pkg string) *logEvent {
	ev.Event = ev.Event.Str("pkg", pkg)
	return ev
}

func (ev *logEvent) Action(action string) *logEvent {
	ev.Event = ev.Event.Str("action", action)
	return ev
}

func (ev *logEvent) CtxId(ctx context.Context) *logEvent {
	ev.Event = ev.Event.Interface("ctxId", ctx.Value("id"))
	return ev
}

func (ev *logEvent) JobId(jobId string) *logEvent {
	ev.Event = ev.Event.Str("jobId", jobId)
	return ev
}

func (ev *logEvent) FilePath(path string) *logEvent {
	ev.Event = ev.Event.Str("filePath", path)
	return ev
}

func (ev *logEvent) Param(param interface{}) *logEvent {
	field := "param"
	if param == nil {
		ev.Event = ev.Event.Str(field, "<empty>")
		return ev
	}

	str, err := util.ParseInterfaceToString(param)
	if err != nil {
		ev.Event = ev.Event.Str(field, "<unable to marshal result> err : "+err.Error()).Interface("rawParam", param)
	} else {
		ev.Event = ev.Event.Str(field, str)
	}
	return ev
}

func LogInputString(param interface{}) interface{} {
	str, err := util.ParseInterfaceToString(param)
	if err != nil {
		return param
	}
	return str
}

func (ev *logEvent) Err(err error) *logEvent {
	ev.Event = ev.Event.Err(err)
	return ev
}

func (ev *logEvent) Str(f, v string) *logEvent {
	ev.Event = ev.Event.Str(f, v)
	return ev
}
