package service

import (
	"errors"
	"server/internal/product/user/model"
	"server/internal/product/user/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserService 提供用户相关的业务逻辑服务
type UserService struct {
	uRepo repository.UserRepository
}

// NewUserService 创建一个新的用户服务实例
func NewUserService(repository repository.UserRepository) *UserService {
	return &UserService{uRepo: repository}
}

var jwtSecret = []byte("gee")

// Claims JWT令牌的自定义声明结构
type Claims struct {
	UserID  int    `json:"user_id"`
	Account string `json:"account"`
	jwt.RegisteredClaims
}

// Login 用户登录，验证账号密码并生成JWT令牌
func (s *UserService) Login(account, password string) (string, error) {

	user, err := s.uRepo.FindUserByAccount(account)
	if err != nil {
		return "", errors.New("invalid account or password")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("invalid account or password")
	}
	//对比密码
	claims := Claims{
		UserID:  user.Uid,
		Account: user.Account,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 999)),
			Issuer:    "gee",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //生成JWT_TOKEN
	tokenString, err := token.SignedString(jwtSecret)          //签名
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Register 用户注册，对密码进行加密并创建新用户
func (s *UserService) Register(account, password, name string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	u := model.User{Account: account, Password: string(passwordHash), Name: name, CreatedAt: time.Now()}
	err = s.uRepo.CreateUser(&u)
	if err != nil {
		return err
	}
	return nil
}
