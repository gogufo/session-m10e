package post

import (
	"fmt"
	. "session/functions"
	"strconv"

	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/microcosm-cc/bluemonday"
)

func SetSessionapi(t *pb.Request) (response *pb.Response) {

	args := ToMapStringInterface(t.Args)
	p := bluemonday.UGCPolicy()

	if args["uid"] == nil || args["isadm"] == nil || args["iscomp"] == nil || args["readon"] == nil {
		return ErrorReturn(t, 406, "000012", "Missing important data")
	}

	uid := p.Sanitize(fmt.Sprintf("%v", args["uid"]))
	isadm := 0
	iscomp := 0
	readon := 0

	isadm, _ = strconv.Atoi(fmt.Sprintf("%v", args["isadm"]))
	iscomp, _ = strconv.Atoi(fmt.Sprintf("%v", args["iscomp"]))
	readon, _ = strconv.Atoi(fmt.Sprintf("%v", args["readon"]))

	//Generate Access Token
	sessionToken, exptime, refreshToken, refresh_exptime, err := SetSession(uid, isadm, iscomp, readon)

	if err != nil {
		return ErrorReturn(t, 400, "000013", err.Error())
	}

	ans := make(map[string]interface{})
	ans["access_token"] = sessionToken
	ans["refresh_token"] = refreshToken
	ans["at_lifetime"] = exptime
	ans["rt_lifetime"] = refresh_exptime

	response = Interfacetoresponse(t, ans)
	return response

}
