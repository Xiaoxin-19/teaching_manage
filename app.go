package main

import (
	"context"
	"fmt"
	"teaching_manage/pkg/dispatcher"
)

// App struct
type App struct {
	ctx        context.Context
	dispatcher *dispatcher.Dispatcher
}

// NewApp creates a new App application struct
func NewApp(dis *dispatcher.Dispatcher) *App {
	return &App{
		dispatcher: dis,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Dispatch dispatches a method with payload to the registered handlers
func (a *App) Dispatch(router string, payload string) string {
	resp, err := a.dispatcher.Dispatch(a.ctx, router, []byte(payload))
	if err != nil {
		// Log the error
	}
	return resp
}
