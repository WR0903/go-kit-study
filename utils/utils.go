package utils

import (
	"log"
	"strconv"

	"github.com/google/uuid"
	consulapi "github.com/hashicorp/consul/api"
)

var ConsulClient *consulapi.Client
var ServiceID string
var ServiceName string
var ServicePort int

func init() {
	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(config) //创建客户端
	if err != nil {
		log.Fatal(err)
	}
	ConsulClient = client
	ServiceID = "userService" + uuid.New().String()
}

func SetServiceNameAndPort(name string, port int) {
	ServiceName = name
	ServicePort = port
}

func RegService() {

	reg := consulapi.AgentServiceRegistration{}
	reg.ID = ServiceID
	reg.Name = ServiceName    //注册service的名字
	reg.Address = "127.0.0.1" //注册service的ip
	reg.Port = ServicePort    //注册service的端口
	reg.Tags = []string{"primary"}

	check := consulapi.AgentServiceCheck{}                                   //创建consul的检查器
	check.Interval = "5s"                                                    //设置consul心跳检查时间间隔
	check.HTTP = "http://127.0.0.1:" + strconv.Itoa(ServicePort) + "/health" //设置检查使用的url

	reg.Check = &check

	err := ConsulClient.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Fatal(err)
	}
}
func UnRegService() {
	ConsulClient.Agent().ServiceDeregister(ServiceName)
}
