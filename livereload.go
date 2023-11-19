package livereload

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var HandlerPath = "/__livereload__"

var (
	upgrader = websocket.Upgrader{}
	inbox    = make(chan string, 2)
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

	lastPing := time.Now()
	logr.Info("[livereload] listening")

	for !quit {
		select {
		case <-r.Context().Done():
			logr.Info("[livereload] context.Done")
			quit = true
			break
		case msg := <-inbox:
			log.Println("lr: reload received")
			sendMessage(w, "", msg)
			log.Println("lr: reload sent")

		default:
			if time.Now().Sub(lastPing) >= time.Second*10 {
				lastPing = time.Now()
				keepAlive(w)
			}
		}
	}

	return nil
}

func keepAlive(w http.ResponseWriter) {
	sendMessage(w, "", "ping")
}

func sendMessage(w http.ResponseWriter, event, data string) {
	w.Header().Set("Content-Type", "text/event-stream")

	if event != "" {
		fmt.Fprint(w, "event:", event, "\n")
	} else if event == ":" {
		fmt.Fprint(w, ":", event, "\n")
	}

	if event != ":" {
		fmt.Fprint(w, "data:", data, "\n")
		// fmt.Fprintln(w, "retry:", 500)
	}

	fmt.Fprint(w, "\n\n")

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

type dummyLogger struct{}

func (dummyLogger) Info(args ...interface{}) {
}
func (dummyLogger) Error(args ...interface{}) {
}
