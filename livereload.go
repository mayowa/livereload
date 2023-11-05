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

func HandleServerMux(mux *http.ServeMux, logger Logger) {
	reloader := &ReloadServer{}
	mux.Handle(HandlerPath, reloader)
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type ReloadServer struct {
	Log Logger
}

func (s *ReloadServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.Log == nil {
		s.Log = dummyLogger{}
	}

	err := ReloadHandler(w, r, s.Log)
	if err != nil {
		s.Log.Error()
	}
}

func ReloadHandler(w http.ResponseWriter, r *http.Request, logger Logger) error {
	quit := false
	if logger == nil {
		logger = dummyLogger{}
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
				logger.Error(fmt.Errorf("writeMessage: %w", err))
			}
		case <-r.Context().Done():
			logger.Info("done")
			quit = true
			break

		default:
			_, in, err := ws.ReadMessage()
			if err != nil {
				logger.Error(fmt.Errorf("readMessage: %w", err))
				quit = true
				break
			}
			logger.Info(string(in))
		}
	}

	return nil
}

type dummyLogger struct{}

func (dummyLogger) Info(args ...interface{}) {
}
func (dummyLogger) Error(args ...interface{}) {
}
