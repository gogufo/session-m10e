package functions

import (
	. "session/model"
	"strings"
	"time"

	. "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func SetSession(name string, isAdmin int, completed int, readonly int) (sessionToken string, exptime int, refreshToken string, refresh_exptime int, err error) {
	// Create a new random session token

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	//exptime = time.Now().Add(time.Duration(viper.GetInt("token.expiretime"))).Unix()
	exptime = int(time.Now().Unix()) + viper.GetInt("token.expiretime")
	//sn.SetErrorLog(fmt.Sprintf("exptime: %v", exptime))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":    name,
		"exipred": exptime,
	})

	// Sign and get the complete encoded token as a string using the secret
	sessionToken, err = token.SignedString([]byte(viper.GetString("token.secretKey")))

	if err != nil {
		return "", 0, "", 0, err
	}

	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	WriteTokenInRedis(sessionToken, name, isAdmin, completed, exptime, readonly)

	refresh_exptime = int(time.Now().Unix()) + viper.GetInt("token.refresh_token_expiration")

	rftoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":    name,
		"exipred": refresh_exptime,
	})

	// Sign and get the complete encoded token as a string using the secret
	refreshToken, err = rftoken.SignedString([]byte(viper.GetString("token.secretKey")))

	if err != nil {
		return "", 0, "", 0, err
	}

	refreshTokenarr := strings.Split(refreshToken, ".")

	//Save RefreshToken to DB
	rf := RefreshTokens{}
	rf.UID = name
	rf.RefreshToken = refreshTokenarr[2]
	rf.LifeTime = refresh_exptime

	db, err := ConnectDBv2()
	if err == nil {

		db.Conn.Create(&rf)

	}

	return sessionToken, exptime, refreshTokenarr[2], refresh_exptime, nil
}
