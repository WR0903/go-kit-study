package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"new/endpoints"
	"strconv"

	"github.com/gorilla/mux"
)

func DecodeUserRequest(c context.Context, r *http.Request) (interface{}, error) { //这个函数决定了使用哪个request结构体来请求
	vars := mux.Vars(r)             //通过这个返回一个map，map中存放的是参数key和值，因为我们路由地址是这样的/user/{uid:\d+}，索引参数是uid,访问Examp: http://localhost:8080/user/121，所以值为121
	if uid, ok := vars["uid"]; ok { //
		uid, _ := strconv.Atoi(uid)
		return endpoints.UserRequest{Uid: uid}, nil
	}
	return nil, errors.New("参数错误")
}

func EncodeUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json") //设置响应格式为json，这样客户端接收到的值就是json，就是把我们设置的UserResponse给json化了

	return json.NewEncoder(w).Encode(response) //判断响应格式是否正确
}

func GetUserInfo_Request(_ context.Context, request *http.Request, r interface{}) error {
	user_request := r.(endpoints.UserRequest)
	request.URL.Path += "/user/" + strconv.Itoa(user_request.Uid)
	return nil
}

func GetUserInfo_Response(_ context.Context, res *http.Response) (response interface{}, err error) {
	if res.StatusCode > 400 {
		return nil, errors.New("no data")
	}
	var user_response endpoints.UserResponse
	err = json.NewDecoder(res.Body).Decode(&user_response)
	if err != nil {
		return nil, err
	}
	return user_response, err
}
