package mypprof

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func RunPprof(port int) {
	//引用 _ "net/http/pprof"
	if port == 0 {
		port = 10000
	}
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil) //开启调试端口
	fmt.Printf("pprof ==> http://localhost:%d/debug/pprof/heap\n", port)
	//go tool pprof http://localhost:10000/debug/pprof/heap
	//go tool pprof http://127.0.0.1:10000/debug/pprof/goroutine
}
