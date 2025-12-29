package dispatcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	valdiatex "teaching_manage/pkg/valdiate"
	"teaching_manage/pkg/wraper"
)

// ErrHandlerNotFound is returned when a route/handler name is not registered.
var ErrHandlerNotFound = errors.New("handler not found")

// ErrResultTypeMismatch indicates a typed dispatch returned an unexpected result type.
var ErrResultTypeMismatch = errors.New("result type mismatch")

// Handler is the internal interface used to execute a registered function.
// Handlers accept a raw JSON payload; typed wrappers convert JSON -> typed value.
type Handler interface {
	Serve(ctx context.Context, payload json.RawMessage) (string, error)
}

// HandlerFunc is a simple adapter to allow using functions as Handlers.
type HandlerFunc func(ctx context.Context, payload json.RawMessage) (string, error)

func (hf HandlerFunc) Serve(ctx context.Context, payload json.RawMessage) (string, error) {
	return hf(ctx, payload)
}

// Dispatcher holds registered handlers
type Dispatcher struct {
	mu       sync.RWMutex
	handlers map[string]Handler
}

// New creates an empty Dispatcher.
func New() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string]Handler),
	}
}

// Register registers a handler function under a name.
// The handler MUST have the signature: func(context.Context, *T) (R, error)
// where `*T` is the pointer type to unmarshal the request payload into.
// The dispatcher will store the raw handler; middlewares run at dispatch time.
func (d *Dispatcher) Register(name string, fn HandlerFunc) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[name] = fn
	return nil
}

// NewHandlerFunc wraps a typed handler `func(context.Context, *Req) (Res, error)`
// into a `HandlerFunc` that accepts `json.RawMessage`, unmarshals into `*Req`,
// and invokes the typed function.
func NewHandlerFunc[Req any, Res any](fn func(context.Context, *Req) (Res, error)) HandlerFunc {
	return HandlerFunc(func(ctx context.Context, payload json.RawMessage) (string, error) {
		req := new(Req)
		if len(payload) > 0 {
			if err := json.Unmarshal(payload, req); err != nil {
				return wraper.NewBadResponse("unmarshal json to object fail").ToJSON(), err
			}
		}

		// validate struct
		err := valdiatex.ValidateStruct(req)
		if err != nil {
			return wraper.NewBadResponse("request validation failed: " + err.Error()).ToJSON(), err
		}

		// call the typed function
		res, err := fn(ctx, req)
		if err != nil {
			return wraper.NewBadResponse(err.Error()).ToJSON(), err
		}
		successResp := wraper.NewSuccessResponse(res)
		successResp.Data = res
		return successResp.ToJSON(), nil
	})
}

// RegisterTyped registers a typed handler function. `fn` must be
// `func(context.Context, *Req) (Res, error)`. The dispatcher will unmarshal
// the incoming JSON into `*Req` and call `fn`.
func RegisterTyped[Req any, Res any](d *Dispatcher, name string, fn func(context.Context, *Req) (Res, error)) error {
	return d.Register(name, NewHandlerFunc(fn))
}

// NewHandlerFuncNoReq wraps a typed handler func(context.Context) (Res, error)
// into a HandlerFunc that ignores any payload and invokes fn.
func NewHandlerFuncNoReq[Res any](fn func(context.Context) (Res, error)) HandlerFunc {
	return HandlerFunc(func(ctx context.Context, payload json.RawMessage) (string, error) {
		res, err := fn(ctx)
		if err != nil {
			return wraper.NewBadResponse(err.Error()).ToJSON(), err
		}
		successResp := wraper.NewSuccessResponse(res)
		successResp.Data = res
		return successResp.ToJSON(), nil
	})
}

// RegisterNoReq registers a handler that does not take a request body.
// fn should be func(context.Context) (Res, error).
func RegisterNoReq[Res any](d *Dispatcher, name string, fn func(context.Context) (Res, error)) error {
	return d.Register(name, NewHandlerFuncNoReq(fn))
}

// Dispatch finds a handler by name and executes it with the provided JSON payload.
func (d *Dispatcher) Dispatch(ctx context.Context, name string, payload json.RawMessage) (string, error) {
	// read handler and middleware slice under lock, then release before execution
	d.mu.RLock()
	h, ok := d.handlers[name]
	d.mu.RUnlock()

	if !ok {
		return wraper.NewBadResponse(fmt.Sprintf("handler [%s] not found ", name)).ToJSON(), ErrHandlerNotFound
	}

	// call the handler (fn) now that Pre hooks have completed
	return h.Serve(ctx, payload)
}
