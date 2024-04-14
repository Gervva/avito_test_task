package repository

import "errors"

var ErrBannerAlreadyExists = errors.New("banenr with such id already exists")
var ErrBannerNotExist = errors.New("banenr with such id not exists")