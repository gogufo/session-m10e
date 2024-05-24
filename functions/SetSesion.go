package functions

import (
	"fmt"
	"time"

	sn "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func SetSession(name string, isAdmin int, completed int, readonly int) (sessionToken string, exptime int, err error) {
	// Create a new random session token

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	//exptime = time.Now().Add(time.Duration(viper.GetInt("token.expiretime"))).Unix()
	exptime = int(time.Now().Unix()) + viper.GetInt("token.expiretime")
	sn.SetErrorLog(fmt.Sprintf("exptime: %v", exptime))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":    name,
		"exipred": exptime,
	})

	// Sign and get the complete encoded token as a string using the secret
	sessionToken, err = token.SignedString([]byte(viper.GetString("token.secretKey")))

	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	WriteTokenInRedis(sessionToken, name, isAdmin, completed, exptime, readonly)
	return sessionToken, exptime, err
}
