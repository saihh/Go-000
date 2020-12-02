package dao

import (
    "github.com/jinzhu/gorm"
    "github.com/pkg/errors"
)

type User {
    ID string
    Name string
    //xxx
}

var db *gorm.DB

func GetUserByID(id string) (*User, error) {
    user := &User{}
    err := db.Table("xxx").Where("id=?", id).Find(user).Error
    if err != nil {
        return errors.Wrap(err, "user not found")
    }
    return user, nil
}




