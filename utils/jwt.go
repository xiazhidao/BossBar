package utils

import (
	"BossBar/conf"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

/// 指定加密密钥
var jwtSecret = []byte(conf.SecretKey)

//Claim是一些实体（通常指的用户）的状态和额外的元数据
type Claims struct {
	Id       int `json:"data"`
	AccessId int `json:"data1"`
	jwt.StandardClaims
}

// 根据产生token
func GenerateToken(merchantId int, expiresAt int64) (string, error) {
	//设置token有效时间

	claims := Claims{
		Id: merchantId,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expiresAt,
			// 指定token发行人
			Issuer:  "acceptance-frontend",
			Subject: time.Now().Format(conf.DateTimeLayout),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// 产生访问token
func GenerateAccessToken(curMerchantId, accessId int, expiresAt int64) (string, error) {
	//设置token有效时间
	claims := Claims{
		Id:       curMerchantId,
		AccessId: accessId,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expiresAt,
			// 指定token发行人
			Issuer: "acceptance-frontend",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// 根据传入的token值获取到Claims对象信息，（进而获取其中的用户名和密码）
func ParseToken(token string) (*Claims, error) {

	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err

}
