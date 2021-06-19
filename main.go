package main

import (
	"flag"
	"fmt"
	"log"
	"new/endpoints"
	"new/services"
	"new/transport"
	"new/utils"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"net/http"

	kitlog "github.com/go-kit/kit/log"

	kitHttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

func main() {
	name := flag.String("name", "", "服务名称")
	port := flag.Int("port", 0, "服务端口")
	flag.Parse()
	if *name == "" {
		log.Fatal("请指定服务名称")
	}
	if *port == 0 {
		log.Fatal("请指定端口")
	}
	utils.SetServiceNameAndPort(*name, *port)
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stdout)
		logger = kitlog.WithPrefix(logger, "mykit", "1.0")
		logger = kitlog.WithPrefix(logger, "time", kitlog.DefaultTimestampUTC)
		logger = kitlog.WithPrefix(logger, "caller", kitlog.DefaultCaller)
	}
	user := services.UserService{}
	limit := rate.NewLimiter(1, 5)
	endp := endpoints.RateLimit(limit)(endpoints.UserServiceLogMiddleware(logger)(endpoints.GenUserEnPoint(user)))

	options := []kitHttp.ServerOption{
		kitHttp.ServerErrorEncoder(utils.MyErrorEncoder),
		//ServerErrorEncoder支持ErrorEncoder类型的参数 type ErrorEncoder func(ctx context.Context, err error, w http.ResponseWriter)
		//我们自定义的MyErrorEncoder只要符合ErrorEncoder类型就可以传入
	} //创建ServerOption切片

	handler := kitHttp.NewServer(endp, transport.DecodeUserRequest, transport.EncodeUserResponse, options...) //使用go kit创建server传入我们之前定义的两个解析函数
	r := mux.NewRouter()                                                                                      //使用mux来使服务支持路由
	//r.Handle(`/user/{uid:\d+}`, serverHandler) //这种写法支持多种请求方法，访问Examp: http://localhost:8080/user/121便可以访问
	r.Methods("GET").Path(`/user/{uid:\d+}`).Handler(handler) //这种写法限定了请求只支持GET方法
	r.Methods("GET").Path("/health").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-type", "application/json")
		writer.Write([]byte(`{"status":"ok"}`))
	})
	errChan := make(chan error)
	go func() {
		utils.RegService() //调用注册服务程序
		err := http.ListenAndServe(":"+strconv.Itoa(utils.ServicePort), r)
		if err != nil {
			log.Println(err)
			errChan <- err
		}
	}()
	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sigChan)
	}()
	getErr := <-errChan //只要报错 或者service关闭阻塞在这里的会进行下去
	utils.UnRegService()
	log.Println(getErr)
}
