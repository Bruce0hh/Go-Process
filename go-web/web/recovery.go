package web

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%v", err)
				log.Printf("%v\n", trace(message))
				ctx.String(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		ctx.Next()
	}
}

func trace(message string) string {
	var pcs [32]uintptr
	callers := runtime.Callers(3, pcs[:]) // 跳过头3个caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback: ")
	for _, pc := range pcs[:callers] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%v:%v", file, line))
	}
	return str.String()
}
