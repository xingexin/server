package repository

import (
	"server/internal/product/model"

	"gorm.io/gorm"
)

type UserWriter interface {
	CreateUser(user *model.User) error
	DeleteUser(uid int) error
	UpdatePassword(uid int, password string) error
	UpdateName(uid int, password string) error
}

type UserReader interface {
	FindUserByUid(uid int) (*model.User, error)
	FindUserByAccount(account string) (*model.User, error)
}

type UserRepository interface {
	UserReader
	UserWriter
} //user操作

type gormUserRepository struct {
	gormDB *gorm.DB
}

func NewUserRepository(gDB *gorm.DB) UserRepository {
	return &gormUserRepository{gormDB: gDB}
}

func (uRepo *gormUserRepository) CreateUser(user *model.User) error {
	err := uRepo.gormDB.Create(user).Error
	return err
}

func (uRepo *gormUserRepository) DeleteUser(uid int) error {
	err := uRepo.gormDB.Delete(&model.User{}, uid).Error
	return err
}

func (uRepo *gormUserRepository) UpdatePassword(uid int, password string) error {
	err := uRepo.gormDB.Model(&model.User{}).Where(uid).Update("password", password).Error
	return err
}

func (uRepo *gormUserRepository) UpdateName(uid int, name string) error {
	err := uRepo.gormDB.Model(&model.User{}).Where(uid).Update("password", name).Error
	return err
}

func (uRepo *gormUserRepository) FindUserByUid(uid int) (*model.User, error) {
	var user model.User
	err := uRepo.gormDB.Find(&user, uid).Error
	return &user, err
}

func (uRepo *gormUserRepository) FindUserByAccount(account string) (*model.User, error) {
	var user model.User
	err := uRepo.gormDB.Where("account=?", account).Find(&user).Error
	return &user, err
}
