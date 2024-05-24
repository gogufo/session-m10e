package functions

import (
	. "session/model"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	sn "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

func UpdateSession(sessionToken string) map[string]interface{} {

	//Get sesssion token
	ans := make(map[string]interface{})

	tokenarray := strings.Split(sessionToken, " ")
	tokenlen := len(tokenarray)

	if tokenlen != 2 {
		ans["error"] = "Wrong Token initialisation"
		return ans
	}

	tokentype := tokenarray[0]
	token := tokenarray[1]

	// Check Session in Redis
	n := sn.ConfigString("redis.host")
	conn, err := redis.DialURL(n)
	if err != nil {
		sn.SetErrorLog("session.go:59 " + err.Error())
	}

	response, err := redis.Values(conn.Do("HMGET", token, "expired", "uid", "isadmin", "completed", "readonly")) //commandName , ARG1, ARG2, ARG3
	if err != nil {
		// If there is an error in setting the cache, return an internal server error
		sn.SetErrorLog("session.go:62: " + err.Error())
	}
	var exptime int
	var uid string
	var isadmin int
	var completed int
	var readonly int

	if _, err := redis.Scan(response, &exptime, &uid, &isadmin, &completed, &readonly); err != nil {
		// handle error
		sn.SetErrorLog("session.go:70: " + err.Error())
	}

	if uid == "" {
		//Check Session in DB

		//Check DB and table config
		db, err := sn.ConnectDBv2()
		if err != nil {
			if viper.GetBool("server.sentry") {
				sentry.CaptureException(err)
			} else {
				sn.SetErrorLog(err.Error())
			}

		}

		if tokentype == "APP" {
			tokentable := APITokens{}
			impersonatetable := ImpersonateTokens{}

			rows := db.Conn.Debug().Where(`token = ?`, token).First(&tokentable)

			if rows.RowsAffected == 0 {
				//Check impersonate token
				rowss := db.Conn.Debug().Where(`token = ?`, token).First(&impersonatetable)
				if rowss.RowsAffected == 0 {
					sn.SetErrorLog("No uid")
					ans["error"] = "000011" // you are not authorised
					return ans
				}

				uid = impersonatetable.UID
				completed = 1
				readonly = 0
				redisexptime := int(time.Now().Unix()) + viper.GetInt("token.expiretime")

				//Write session into Redis
				WriteTokenInRedis(token, uid, isadmin, completed, redisexptime, readonly)

			} else {

				if !tokentable.Status {
					sn.SetErrorLog("No uid")
					ans["error"] = "000011" // you are not authorised
					return ans
				}

				if tokentable.Expiration != 0 && int64(tokentable.Expiration) < time.Now().Unix() {
					sn.SetErrorLog("No uid")
					ans["error"] = "000011" // you are not authorised
					return ans
				}

				// Check Doues User is Admin in case of Token Admin Satatus
				/*
					if tokentable.IsAdmin {
						userExist := Users{}

						db.Conn.Debug().Where(`uid = ?`, tokentable.UID).First(&userExist)
						isadmin = 0
						if userExist.IsAdmin {
							isadmin = 1
						}

					}
				*/
				exptime = tokentable.Expiration
				uid = tokentable.UID
				completed = 1
				readonly = 0
				if tokentable.Readonly {
					readonly = 1
				}

				redisexptime := int(time.Now().Unix()) + viper.GetInt("token.expiretime")

				//Write session into Redis
				WriteTokenInRedis(token, uid, isadmin, completed, redisexptime, readonly)
			}
		} else {
			sn.SetErrorLog("No uid")
			ans["error"] = "000011" // you are not authorised
			return ans
		}

	}

	//updates session
	newexptime := int(time.Now().Unix()) + viper.GetInt("token.expiretime")
	WriteTokenInRedis(token, uid, isadmin, completed, newexptime, readonly)

	ans["uid"] = uid
	ans["isadmin"] = isadmin
	ans["session_expired"] = newexptime
	ans["completed"] = completed
	ans["readonly"] = readonly
	ans["token"] = token
	ans["token_type"] = tokentype
	sn.SetErrorLog(uid)
	sn.SetErrorLog(sessionToken)
	return ans

}
