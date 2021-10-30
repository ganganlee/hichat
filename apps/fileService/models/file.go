package models

import (
	"github.com/xormplus/xorm"
	"time"
)

type (
	Files struct {
		Id         int64
		Uuid       string    `json:"uuid" xorm:"varchar(125) notnull"`
		Name       string    `json:"name" xorm:"varchar(125) notnull"`
		Type       string    `json:"type" xorm:"varchar(50) notnull"`
		Path       string    `json:"path" xorm:"varchar(255) unique"`
		Size       uint16    `json:"size" xorm:"int default(0)"`
		CreateTime time.Time `json:"-" xorm:"created"`
		UpdateTime time.Time `json:"-" xorm:"updated"`
	}
	FilesModel struct {
		engine *xorm.Engine
	}
)

func NewFilesModel(e *xorm.Engine) *FilesModel {
	return &FilesModel{
		e,
	}
}

//插入数据
func (f *FilesModel) Insert(file *Files) (err error) {
	_, err = f.engine.Insert(file)
	return err
}

//根据文件名查找数据
func (f *FilesModel) FindByName(path string) (file *Files, err error) {
	var exist bool
	file = new(Files)

	if exist, err = f.engine.Where("path=?", path).Get(file); err != nil {
		return nil, err
	}

	if !exist {
		return nil, nil
	}

	return file, nil
}
