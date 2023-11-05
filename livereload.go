package livereload

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var HandlerPath = "/__livereload__"

var (
	upgrader = websocket.Upgrader{}
	inbox    = make(chan string)
)

func Reload() {
	go func() {
		inbox <- "reload"
	}()
}

func HandleServerMux(mux *http.ServeMux, options *Options) {
	reloader := &ReloadServer{Options: options}
	mux.Handle(HandlerPath, reloader)
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type ReloadServer struct {
	Options *Options
}

func (s *ReloadServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.Options.Logger == nil {
		s.Options.Logger = dummyLogger{}
	}

	err := ReloadHandler(w, r, s.Options)
	if err != nil {
		s.Options.Logger.Error()
	}
}

type Options struct {
	Logger Logger
	Files  []*FileInfo
}

func ReloadHandler(w http.ResponseWriter, r *http.Request, options *Options) error {
	if options == nil {
		options = &Options{}
	}

	quit := false
	logr := options.Logger
	if logr == nil {
		logr = dummyLogger{}
	}

	if len(options.Files) > 0 {
		// watch for modified files
		fw := NewFileWatcher(options.Files)
		go fw.Run(r.Context(), inbox)
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	ws.SetReadLimit(512)

	for !quit {
		select {
		case msg := <-inbox:
			err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				logr.Error(fmt.Errorf("writeMessage: %w", err))
			}
		case <-r.Context().Done():
			logr.Info("done")
			quit = true
			break

		default:
			_, in, err := ws.ReadMessage()
			if err != nil {
				logr.Error(fmt.Errorf("readMessage: %w", err))
				quit = true
				break
			}
			logr.Info(string(in))
		}
	}

	return nil
}

type dummyLogger struct{}

func (dummyLogger) Info(args ...interface{}) {
}
func (dummyLogger) Error(args ...interface{}) {
}
