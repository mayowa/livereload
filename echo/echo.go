package echo

import (
	"github.com/labstack/echo/v4"
	"github.com/mayowa/livereload"
)

func HandleEcho(e *echo.Echo, logger livereload.Logger) {
	e.GET(livereload.HandlerPath, func(c echo.Context) error {
		return livereload.ReloadHandler(c.Response(), c.Request(), logger)
	})
}
