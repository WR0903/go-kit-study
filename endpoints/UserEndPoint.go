package endpoints

import (
	"context"

	"new/services"

	"github.com/go-kit/kit/endpoint"
)

type UserRequest struct { //封装User请求结构体
	Uid int `json:"uid"`
}

type UserResponse struct {
	Result string `json:"result"`
}

func GenUserEnPoint(userService services.IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)           //通过类型断言获取请求结构体
		result := userService.GetName(r.Uid) //
		return UserResponse{Result: result}, nil
	}
}
