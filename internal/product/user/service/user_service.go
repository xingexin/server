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

// jwtSecret JWT令牌的签名密钥
// 注意：此密钥硬编码在代码中，生产环境应从配置文件或环境变量中读取
var jwtSecret = []byte("gee")

// Claims JWT令牌的自定义声明结构
// 包含用户身份信息和标准JWT声明
type Claims struct {
	UserID  int    `json:"user_id"`  // 用户ID，用于标识用户身份
	Account string `json:"account"`   // 用户账号，用于显示用户信息
	jwt.RegisteredClaims               // JWT标准声明（签发时间、过期时间等）
}

// Login 用户登录，验证账号密码并生成JWT令牌
// 业务流程：
// 1. 根据账号从数据库查询用户信息
// 2. 使用bcrypt验证密码是否正确
// 3. 生成包含用户信息的JWT令牌（有效期999小时）
// 4. 返回JWT令牌字符串
//
// 安全说明：
// - 密码使用bcrypt加密存储，验证时使用bcrypt.CompareHashAndPassword
// - JWT令牌使用HS256算法签名
// - 令牌有效期为999小时（约41天）
//
// 返回值：
// - string: JWT令牌字符串，客户端需在后续请求的Authorization头中携带
// - error: 账号不存在或密码错误时返回错误
func (s *UserService) Login(account, password string) (string, error) {
	// 根据账号查询用户
	user, err := s.uRepo.FindUserByAccount(account)
	if err != nil {
		return "", errors.New("invalid account or password")
	}

	// 使用bcrypt验证密码（将数据库中的加密密码与用户输入的明文密码对比）
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("invalid account or password")
	}

	// 构建JWT声明（包含用户ID、账号、签发时间、过期时间等）
	claims := Claims{
		UserID:  user.Uid,
		Account: user.Account,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),                           // 签发时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 999)),      // 过期时间（999小时后）
			Issuer:    "gee",                                                     // 签发者
		},
	}

	// 使用HS256算法生成JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用密钥对令牌进行签名
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Register 用户注册，对密码进行加密并创建新用户
// 业务流程：
// 1. 使用bcrypt对密码进行加密（cost=10）
// 2. 构建用户模型对象（设置账号、加密后的密码、姓名、创建时间）
// 3. 调用Repository层创建用户记录
//
// 安全说明：
// - 密码使用bcrypt加密，cost=10（2^10次迭代）
// - bcrypt会自动生成随机盐值，相同密码的哈希值也不同
// - 加密后的密码约60个字符，数据库字段需足够长
//
// 错误情况：
// - 账号已存在：Repository层会返回唯一约束错误
// - 密码加密失败：返回bcrypt错误
//
// 参数：
// - account: 用户账号，需唯一
// - password: 明文密码，将被bcrypt加密后存储
// - name: 用户姓名
func (s *UserService) Register(account, password, name string) error {
	// 使用bcrypt对密码进行加密（cost=10表示2^10次迭代）
	// 返回的哈希值包含：算法版本、cost、盐值、密码哈希
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	// 构建用户模型对象
	u := model.User{
		Account:   account,
		Password:  string(passwordHash), // 存储加密后的密码
		Name:      name,
		CreatedAt: time.Now(),
	}

	// 调用Repository层创建用户记录
	err = s.uRepo.CreateUser(&u)
	if err != nil {
		return err
	}

	return nil
}
