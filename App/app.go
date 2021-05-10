package App

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"time"
)

/*

基于errgroup实现一个http_server的启动和关闭
以及Linux信号的注册和处理，要保证能够一个退出，全部注销退出:

1、对于linux信号的各种不同处理
2、捕捉到一个信号之后能够退出所有的协程和关闭server


*/

type AppServer struct {
	G         *errgroup.Group
	Ctx       context.Context
	AllServer []*http.Server
}

var appServer *AppServer

func init() {
	appServer = &AppServer{}
	appServer.G, appServer.Ctx = errgroup.WithContext(context.Background())

}

func Run(addr string, handler http.Handler) error {
	appServer.G.Go(func() error {
		svr := http.Server{
			Addr:    addr,
			Handler: handler,
		}
		appServer.AllServer = append(appServer.AllServer, &svr)
		err := svr.ListenAndServe()
		return err
	})
	return nil
}

func Register(sig os.Signal) {
	appServer.G.Go(func() error {
		var sigChan = make(chan os.Signal)
		signal.Notify(sigChan, sig)
		select {
		case <-sigChan:
			fmt.Printf("信号捕捉到signal\n")
			Shutdown()
			return nil
		case <-appServer.Ctx.Done():
			fmt.Printf("Sig = %s routinue exit\n", sig.String())

		}
		return nil
	})
}

func Shutdown() {
	appServer.G.Go(func() error {
		ctx, cancel := context.WithTimeout(appServer.Ctx, time.Second)
		defer cancel()
		for _, svr := range appServer.AllServer {
			err := svr.Shutdown(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	err := appServer.G.Wait()
	if err != nil {
		fmt.Printf("wait err = %v\n", err)
	}
	fmt.Printf("wait exit\n")
}

func Wait() {
	<-appServer.Ctx.Done()
	fmt.Printf("finish")
}
