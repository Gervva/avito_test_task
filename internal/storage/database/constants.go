package database

import "errors"

var (
	ErrTagAlreadyExists    = errors.New("tag already exists")
	ErrBannerAlreadyExists = errors.New("banner already exists")
	ErrBannerNotExist      = errors.New("banner not exist")
	ErrRepoDB              = errors.New("DB ERROR")
)
