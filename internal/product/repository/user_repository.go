package repository

import (
	"server/internal/product/model"

	"gorm.io/gorm"
)

// UserWriter 定义用户写操作接口
type UserWriter interface {
	CreateUser(user *model.User) error
	DeleteUser(uid int) error
	UpdatePassword(uid int, password string) error
	UpdateName(uid int, password string) error
}

// UserReader 定义用户读操作接口
type UserReader interface {
	FindUserByUid(uid int) (*model.User, error)
	FindUserByAccount(account string) (*model.User, error)
}

// UserRepository 用户操作的数据访问接口，组合了读写操作
type UserRepository interface {
	UserReader
	UserWriter
} //user操作

type gormUserRepository struct {
	gormDB *gorm.DB
}

// NewUserRepository 创建一个新的用户仓储实例
func NewUserRepository(gDB *gorm.DB) UserRepository {
	return &gormUserRepository{gormDB: gDB}
}

// CreateUser 在数据库中创建新用户记录
func (uRepo *gormUserRepository) CreateUser(user *model.User) error {
	err := uRepo.gormDB.Create(user).Error
	return err
}

// DeleteUser 根据用户ID从数据库中删除用户记录
func (uRepo *gormUserRepository) DeleteUser(uid int) error {
	err := uRepo.gormDB.Delete(&model.User{}, uid).Error
	return err
}

// UpdatePassword 更新用户密码
func (uRepo *gormUserRepository) UpdatePassword(uid int, password string) error {
	err := uRepo.gormDB.Model(&model.User{}).Where(uid).Update("password", password).Error
	return err
}

// UpdateName 更新用户名称
func (uRepo *gormUserRepository) UpdateName(uid int, name string) error {
	err := uRepo.gormDB.Model(&model.User{}).Where(uid).Update("password", name).Error
	return err
}

// FindUserByUid 根据用户ID从数据库中查找用户
func (uRepo *gormUserRepository) FindUserByUid(uid int) (*model.User, error) {
	var user model.User
	err := uRepo.gormDB.Find(&user, uid).Error
	return &user, err
}

// FindUserByAccount 根据账号从数据库中查找用户
func (uRepo *gormUserRepository) FindUserByAccount(account string) (*model.User, error) {
	var user model.User
	err := uRepo.gormDB.Where("account=?", account).Find(&user).Error
	return &user, err
}
