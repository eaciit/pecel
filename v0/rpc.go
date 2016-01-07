package appserver

import (
	"github.com/eaciit/toolkit"
	//"time"
	"errors"
	"strings"
)

//type RpcFn func(toolkit.M, *toolkit.Result) error
type RpcFn func(toolkit.M) *toolkit.Result
type RpcFns map[string]*RpcFnInfo

type RpcFnInfo struct {
	AuthRequired bool
	AuthType     string
	Fn           RpcFn
}

//type ReturnedBytes []byte

type Rpc struct {
	Fns               RpcFns
	Server            *Server
	MarshallingMethod string
}

var _marshallingMethod string

func MarshallingMethod() string {
	if _marshallingMethod == "" {
		_marshallingMethod = "gob"
	} else {
		_marshallingMethod = strings.ToLower(_marshallingMethod)
	}
	return _marshallingMethod
}

func SetMarshallingMethod(m string) {
	_marshallingMethod = m
}

func (r *Rpc) Do(in toolkit.M, out *toolkit.Result) error {
	if r.Fns == nil {
		r.Fns = map[string]*RpcFnInfo{}
	}

	//in.Set("rpc", r)
	method := in.GetString("method")
	if method == "" {
		return errors.New("Method is empty")
	}
	fninfo, fnExist := r.Fns[method]
	if !fnExist {
		return errors.New("Method " + method + " is not exist")
		/*
			" Available methodnames on " + r.Server.Address + "  are: " + strings.Join(func() []string {
				ret := []string{}
				for name, _ := range r.Fns {
					ret = append(ret, name)
				}
				return ret
			}(), ", "))
		*/
	}
	if fninfo.AuthRequired {
		referenceID := in.GetString("auth_referenceid")
		secret := in.GetString("auth_secret")
		valid := r.Server.validateSecret(fninfo.AuthType, referenceID, secret)
		if valid != string(toolkit.Status_OK) {
			return errors.New("Unauthorised to call " + method + " " + valid + "  Profile:" + referenceID)
		}
	}
	res := fninfo.Fn(in)
	if res.Status != toolkit.Status_OK {
		return errors.New("RPC Call error: " + res.Message)
	}
	//*out = toolkit.ToBytes(res.Data, MarshallingMethod())
	*out = *res
	return nil
}

func AddFntoRpc(r *Rpc, svr *Server, k string, fn RpcFn, needValidation bool, authType string) {
	//func (r *Rpc) AddFn(k string, fn RpcFn) {
	//if r.Server == nil {
	svr.Log.Info("Register " + svr.Address + "/" + k)
	r.Server = svr
	//}
	if r.Fns == nil {
		r.Fns = map[string]*RpcFnInfo{}
	}
	r.Fns[k] = &RpcFnInfo{
		AuthRequired: needValidation,
		AuthType:     authType,
		Fn:           fn,
	}
}
