package services

type IUserService interface {
	GetName(userid int) string
}

type UserService struct{}

func (this UserService) GetName(userid int) string {
	if userid == 101 {
		return "jerry"
	}
	return "guest"
}
