package functions

import (
	sn "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

func WriteTokenInRedis(sessionToken string, uid string, isadmin int, completed int, exptime int, readonly int) {

	n := viper.GetString("redis.host")
	conn, err := redis.DialURL(n)
	if err != nil {
		sn.SetErrorLog("session.go:128: " + err.Error())
	}

	_, err = conn.Do("HMSET", sessionToken, "expired", exptime, "uid", uid, "isadmin", isadmin, "completed", completed, "readonly", readonly) //commandName , ARG1, ARG2, ARG3
	if err != nil {
		// If there is an error in setting the cache, return an internal server error

		sn.SetErrorLog("session.go:137: " + err.Error())
	}

	_, err = conn.Do("EXPIRE", sessionToken, viper.GetInt("token.expiretime"))
	if err != nil {
		// If there is an error in setting the cache, return an internal server error

		sn.SetErrorLog("session.go:146: " + err.Error())
	}

}
