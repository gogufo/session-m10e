package functions

import (
	sn "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

func DelSession(sessionToken string) {

	n := viper.GetString("redis.host")
	conn, err := redis.DialURL(n)
	if err != nil {
		sn.SetErrorLog("session.go:93: " + err.Error())
	}

	response, err := redis.Values(conn.Do("HMGET", sessionToken, "expired", "uid", "isadmin")) //commandName , ARG1, ARG2, ARG3
	if err != nil {
		// If there is an error in setting the cache, return an internal server error

		sn.SetErrorLog("session.go:100: " + err.Error())
	}
	var exptime int64
	var uid string
	var isadmin int

	if _, err := redis.Scan(response, &exptime, &uid, &isadmin); err != nil {
		// handle error
		sn.SetErrorLog("session.go:108: " + err.Error())
	}

	if uid == "" {

		return
	}

	_, err = conn.Do("DEL", sessionToken)
	if err != nil {

		return
	}
}
