package service

import (
	"cron_master/model"
	"cron_master/model/request"
)

func Login(u *model.User) (*model.User, error) {
	return model.AuthLogin(u)
}

func GetUserList(search *request.ComPageInfo) ([]model.User, int, error) {
	return model.GetUserList(search)
}

func SaveUser(u *model.User) error {
	return model.SaveUser(u)
}

func UpdatePwd(u *request.ModifyPwd) (err error) {
	return model.UpdatePwd(u.UserId, u.NewPwd, u.OldPwd)
}
