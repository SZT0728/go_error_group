package go_error_group

import (
	"fmt"
	"go_error_group/App"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)


type AServer struct {
}

func (a AServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "demo Aserver")
}

type BServer struct {
}

func (b BServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "demo BServer")
}

func main() {

	a := &AServer{}
	err := App.Run("localhost:8000", a)
	if err != nil {
		fmt.Printf("Aserver.run failed err = %v", err)
		return
	}

	b := &BServer{}
	err = App.Run("localhost:8001",b)
	if err != nil{
		fmt.Printf("Bserver.run failed err=%v",err)
		return
	}


	App.Register(syscall.SIGTERM)
	App.Register(syscall.SIGINT)
	App.Register(syscall.SIGQUIT)

	var sigChan = make(chan os.Signal, 10)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-sigChan:
		App.Shutdown()
		fmt.Printf("程序即将退出")
	}

}

