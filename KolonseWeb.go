// KolonseWeb project KolonseWeb.go
package KolonseWeb

import (
	"fmt"
	. "github.com/kolonse/KolonseWeb/Type"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	Handlers *Handlers
	Server   *http.Server
}

func (app *App) Use(do DoStep) {
	app.Handlers.Use(do)
}

func (app *App) Get(patter string, do DoStep) {
	app.Handlers.Get(patter, do)
}

func (app *App) Post(patter string, do DoStep) {
	app.Handlers.Post(patter, do)
}

func (app *App) Listen(host string, port int) {
	BeeLogger.Info("listen on %s:%d", host, port)
	endRunning := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		app.Server.Addr = addr
		app.Server.Handler = app.Handlers
		err := app.Server.ListenAndServe()
		if err != nil {
			BeeLogger.Critical("ListenAndServe: ", err)
			endRunning <- true
			return
		}
		endRunning <- true
	}()
	go func() {
		<-sigs
		endRunning <- true
	}()
	<-endRunning
}

// NewApp returns a new beego application.
func NewApp() *App {
	h := NewHandler()
	app := &App{
		Handlers: h,
		Server:   &http.Server{},
	}
	return app
}

var DefaultApp *App

func init() {
	DefaultApp = NewApp()
}
