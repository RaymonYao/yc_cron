package model

import (
	"errors"
	"strings"
	"time"
	"yc_cron/model/request"
)

type Group struct {
	GroupId      uint   `json:"group_id" gorm:"PRIMARY_KEY"`
	GroupName    string `json:"group_name"`
	Description  string `json:"description"`
	CreateUserid int    `json:"create_userid"`
	UpdateUserid int    `json:"update_userid"`
	CreateTime   int64  `json:"create_time"`
	UpdateTime   int64  `json:"update_time"`
}

func GetGroupList(search *request.BasePageInfo) (groupList []Group, total int, err error) {
	db := mdb
	if search.Condition != "" && search.SearchValue != "" {
		db = db.Where(search.Condition+" like ?", "%"+search.SearchValue+"%")
	}
	if err = db.Model(&groupList).Count(&total).Error; err != nil {
		return
	}
	if search.PageSize == 0 {
		err = db.Find(&groupList).Error
	} else {
		err = db.Limit(search.PageSize).Offset(search.PageSize * (search.CurrentPage - 1)).Find(&groupList).Error
	}
	return
}

func DelGroup(groupId int) (err error) {
	var group Group
	if err = mdb.Delete(&group, groupId).Error; err != nil {
		return
	}
	return
}

func SaveGroup(group *Group) (err error) {
	nowTime := time.Now().Unix()
	if group.GroupId == 0 {
		group.CreateTime = nowTime
		group.UpdateTime = nowTime
		err = mdb.Save(group).Error
	} else {
		err = mdb.Model(group).Updates(group).Error
	}
	if err != nil {
		if strings.Index(err.Error(), "uni_group_name") != -1 {
			err = errors.New("该分组名称已存在")
		}
	}
	return
}
