package main

import (
	_ "embed"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js wasm_exec.js
//go:embed wasm_exec.js
var wasmExec []byte

func main() {
	e := echo.New()

	e.File("/", "./index.html")

	e.GET("/wasm_exec.js", func(c echo.Context) error {
		minifier := minify.New()
		minifier.AddFunc("text/javascript", js.Minify)

		minifiedContent, err := minifier.Bytes("text/javascript", wasmExec)
		if err != nil {
			fmt.Println(err)
			return err
		}

		c.Response().Header().Set(echo.HeaderContentType, "text/javascript")
		c.Response().Header().Set(echo.HeaderContentLength, fmt.Sprint(len(minifiedContent)))

		return c.Blob(http.StatusOK, "text/javascript", []byte(minifiedContent))
	})

	e.File("/websshclient.wasm", "./wasm/wasm")

	e.Logger.Fatal(e.Start(":8080"))
}
