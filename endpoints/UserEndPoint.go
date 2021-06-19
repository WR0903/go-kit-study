package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"

	"new/services"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
)

type UserRequest struct { //封装User请求结构体
	Uid int `json:"uid"`
}

type UserResponse struct {
	Result string `json:"result"`
}

//加入限流中间件
func RateLimit(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return nil, errors.New("too many request")
			}
			return next(ctx, request)
		}

	}
}

//日志中间件,每一个service都应该有自己的日志中间件
func UserServiceLogMiddleware(logger log.Logger) endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
	return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest) //通过类型断言获取请求结构体
			logger.Log("event", "get user", "userid", r.Uid)
			return next(ctx, request)
		}
	}
}

func GenUserEnPoint(userService services.IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		r := request.(UserRequest)           //通过类型断言获取请求结构体
		result := userService.GetName(r.Uid) //
		return UserResponse{Result: result}, nil
	}
}
