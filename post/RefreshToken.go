package post

import (
	"fmt"
	. "session/model"
	"time"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

func RefreshSesion(t *pb.Request) (response *pb.Response) {

	args := ToMapStringInterface(t.Args)
	p := bluemonday.UGCPolicy()

	var rft string

	if args["refresh_token"] == nil {
		return ErrorReturn(t, 406, "000012", "Missing important data")
	}

	rft = p.Sanitize(fmt.Sprintf("%v", args["refresh_token"]))

	rf := RefreshTokens{}

	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}

		return ErrorReturn(t, 500, "000027", err.Error())
	}

	rows := db.Conn.Debug().Where(`refresh_token = ?`, rft).First(&rf)

	if rows.RowsAffected == 0 {
		// return error. user name is exist in db users
		return ErrorReturn(t, 400, "000003", "Token Not Found")
	}

	if int(time.Now().Unix()) > rf.LifeTime {

		return ErrorReturn(t, 400, "000009", "Hash expired")

	}

	ans := make(map[string]interface{})
	ans["uid"] = rf.UID

	response = Interfacetoresponse(t, ans)
	return response

}
