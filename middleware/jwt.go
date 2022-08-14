package middleware

import (
	jwt "github.com/dgrijalva/jwt-go"
	"hm/setting"
	"time"
)

var jwtSecret = []byte(setting.AppSettings.JwtSecret)

type  Claims struct {
	UserID 	 string `json:"userid"`
	Authority     string     `json:"authority"`
	jwt.StandardClaims
}

// GenerateToken 产生token的函数
func GenerateToken(userid, authority string)(string,error){
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * 24 * time.Hour)

	claims:=Claims{
		userid,
		authority,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer: "gin-blog",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken 验证token的函数
func ParseToken(token string)(*Claims,error){
	tokenClaims, err := jwt.ParseWithClaims(token,&Claims{},func(token *jwt.Token)(interface{},error){
		return jwtSecret,nil
	})

	if tokenClaims != nil{
		if claims, ok := tokenClaims.Claims.(*Claims);ok && tokenClaims.Valid{
			return claims, nil
		}
	}

	return nil, err
}
