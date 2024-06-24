package functions

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func GenAccessToken(name string, isAdmin int, completed int, readonly int) (accessToken string, exptime int, err error) {
	// Create a new random session token

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	//exptime = time.Now().Add(time.Duration(viper.GetInt("token.expiretime"))).Unix()
	exptime = int(time.Now().Unix()) + viper.GetInt("token.expiretime")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":    name,
		"exipred": exptime,
	})

	// Sign and get the complete encoded token as a string using the secret
	accessToken, err = token.SignedString([]byte(viper.GetString("token.secretKey")))

	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	//	WriteTokenInRedis(sessionToken, name, isAdmin, completed, exptime, readonly)
	return accessToken, exptime, err
}

func GenRefreshToken(name string, isAdmin int, completed int, readonly int) (refreshToken string, exptime int, err error) {
	// Create a new random session token

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	//exptime = time.Now().Add(time.Duration(viper.GetInt("token.expiretime"))).Unix()
	exptime = int(time.Now().Unix()) + viper.GetInt("token.refresh_token_expiration")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":    name,
		"exipred": exptime,
	})

	// Sign and get the complete encoded token as a string using the secret
	refreshToken, err = token.SignedString([]byte(viper.GetString("token.secretKey")))

	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	//	WriteTokenInRedis(sessionToken, name, isAdmin, completed, exptime, readonly)
	return refreshToken, exptime, err
}
