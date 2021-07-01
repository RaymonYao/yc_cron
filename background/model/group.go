package model

import (
	"errors"
	"strings"
	"time"
	"yc_cron/model/request"
)

type Group struct {
	GroupId      uint      `json:"group_id" gorm:"PRIMARY_KEY"`
	GroupName    string    `json:"group_name"`
	Description  string    `json:"description"`
	CreateUserid int       `json:"create_userid"`
	UpdateUserid int       `json:"update_userid"`
	CreateAt     time.Time `json:"create_at"`
	UpdateAt     time.Time `json:"update_at"`
	Refer        string    // 关联外键
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

func SaveGroup(p *Group) (err error) {
	if p.GroupId == 0 {
		p.CreateAt = time.Now()
		p.UpdateAt = time.Now()
		err = mdb.Save(p).Error
	} else {
		err = mdb.Model(p).Updates(p).Error
	}
	if err != nil {
		if strings.Index(err.Error(), "uni_group_name") != -1 {
			err = errors.New("该分组名称已存在")
		}
	}
	return
}
