package pairmodels

import (
	"github.com/strongo/db"
	"github.com/strongo/app"
)

const UserKind = "User"

type User struct {
	db.StringID
	*UserEntity
}

type UserEntity = strongo.AppUserBase

var _ db.EntityHolder = (*PairsPlayer)(nil)

func (User) Kind() string {
	return PairsPlayerKind
}

func (user User) Entity() interface{} {
	return user.UserEntity
}

func (User) NewEntity() interface{} {
	return new(UserEntity)
}

func (user User) SetEntity(entity interface{}) {
	user.UserEntity = entity.(*UserEntity)
}
