# livereload

Livereload is a  Go package that helps you reload your page in the browser 
when you code is rerun during development

## Usage
Copy js/livereload.js to the folder where your assets are stored then
add a script tag for it in <head>

```html
<head>
	<meta charset="UTF-8">
	<meta name="viewport"
		content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">

	<script src="/public/js/livereload.js"></script>
</head>
```

Next use livereload in your Go project

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/mayowa/livereload"
)

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("assets/"))
	mux.Handle("/public/", http.StripPrefix("/public/", fs))
	
	livereload.HandleServerMux(mux, nil)
	
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		html := `
		<head>
            <meta charset="UTF-8">
            <meta name="viewport"
                content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
        
            <script src="/public/js/livereload.js"></script>
        </head>
		<body>
		  <h1>Hello World</h1>
		</body>
		`
		fmt.Fprintln(w, html)
	})
}

```

Livereload comes with an adapter for Echo

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	lre "github.com/mayowa/livereload/echo"
)

func main() {
	e := echo.New()
	e.Static("/public", "public")

	lre.HandleEcho(e, nil)

	e.GET("/", func(c echo.Context) error {
		html := `
		<head>
            <meta charset="UTF-8">
            <meta name="viewport"
                content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
        
            <script src="/public/js/livereload.js"></script>
        </head>
		<body>
		  <h1>Hello World</h1>
		</body>
		`
		return c.HTML(http.StatusOK, html)
	})
}

```