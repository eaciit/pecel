package appserver

import (
	"errors"
	//"fmt"
	"github.com/eaciit/errorlib"
	"github.com/eaciit/toolkit"
	"net/rpc"

	"time"
)

const (
	objClient = "Client"
)

var _timeOut time.Duration

func SetDialTimeout(t time.Duration) {
	_timeOut = t
}

func DialTimeout() time.Duration {
	if _timeOut == 0 {
		_timeOut = 90 * time.Second
	}
	return _timeOut
}

type Client struct {
	UserID    string
	LoginDate time.Time

	client    *rpc.Client
	address   string
	secret    string
	sessionID string
}

func (a *Client) Connect(address string, secret string, userid string) error {
	//("Connecting to %s@%s \n", userid, address)
	a.UserID = userid
	client, e := rpc.Dial("tcp", address)
	if e != nil {
		return errorlib.Error(packageName, objClient, "Connect", "["+address+"] Unable to connect: "+e.Error())
	}
	a.client = client
	a.LoginDate = time.Now().UTC()

	r := a.Call("addsession", toolkit.M{}.Set("auth_secret", secret).Set("auth_referenceid", a.UserID))
	if r.Status != toolkit.Status_OK {
		return errors.New("[" + address + "] Connect: " + r.Message + " User:" + a.UserID)
	}
	m := toolkit.M{}
	toolkit.FromBytes(r.Data.([]byte), "gob", &m)
	a.address = address
	a.sessionID = m.GetString("referenceid")
	a.secret = m.GetString("secret")
	return nil
}

func (a *Client) Close() {
	if a.client != nil {
		a.Call("removesession", nil)
		a.client.Close()
	}
}

func (a *Client) Call(methodName string, in toolkit.M) *toolkit.Result {
	if a.client == nil {
		return toolkit.NewResult().SetErrorTxt("Unable to call, no connection handshake")
	}
	if in == nil {
		in = toolkit.M{}
	}
	out := toolkit.NewResult()
	in["method"] = methodName
	if in.GetString("auth_referenceid") == "" {
		in["auth_referenceid"] = a.sessionID
	}
	//fmt.Println("SessionID: " + a.sessionID)
	if in.Has("auth_secret") == false {
		in.Set("auth_secret", a.secret)
	}
	e := a.client.Call("Rpc.Do", in, out)
	//_ = "breakpoint"
	if e != nil {
		return out.SetErrorTxt(a.address + "." + methodName + " Fail: " + e.Error())
	}
	return out
}
