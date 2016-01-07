package rpctest

import (
	"github.com/eaciit/pecel/v1/client"
	"github.com/eaciit/pecel/v1/server"
	"github.com/eaciit/toolkit"
	"testing"
	"time"
)

var server *pecelserver.Server
var client *pecelclient.Client
var serverInit bool
var (
	serverSecret string = "ariefdarmawan"
)

type controller struct {
}

type Score struct {
	Subject string
	Value   int
}

func (a *controller) Hi(in toolkit.M) *toolkit.Result {
	r := toolkit.NewResult()
	name := in.GetString("name")
	r.SetBytes(struct {
		HelloMessage string
		TimeNow      time.Time
		Scores       []Score
	}{"Hello " + name, time.Now(), []Score{{"Bahasa Indonesia", 90}, {"Math", 85}}}, "gob")
	return r
}

func checkTestSkip(t *testing.T) {
	if serverInit == false {
		t.Skip()
	}
}

func TestStart(t *testing.T) {
	server = new(pecelserver.Server)
	server.RegisterRPCFunctions(new(controller))
	//server.SetSecret(serverSecret)
	server.AllowMultiLogin = true
	server.AddUser("ariefdarmawan", serverSecret)
	e := server.Start("localhost:8001")
	if e == nil {
		serverInit = true
	} else {
		t.Errorf("Fail to start server: %s", e.Error())
	}
}

func checkResult(result *toolkit.Result, t *testing.T) {
	if result.Status != toolkit.Status_OK {
		t.Error(result.Message)
	} else {
		if result.IsEncoded() == false {
			t.Logf("Result: %v", result.Data)
		} else {

			m := struct {
				HelloMessage string
				//TimeNow      time.Time
				Scores []Score
			}{}

			//m := toolkit.M{}
			e := result.GetFromBytes(&m)
			if e != nil {
				t.Errorf("Unable to decode result: %s\n", e.Error())
				return
			}
			t.Logf("Result (decoded): %s", toolkit.JsonString(m))
		}
	}
}

func TestClient(t *testing.T) {
	checkTestSkip(t)
	client = new(pecelclient.Client)
	e := client.Connect(server.Address, serverSecret, "ariefdarmawan")
	//e := client.Connect(server.Address, serverSecret+"_10", "ariefdarmawan")
	if e != nil {
		t.Error(e.Error())
		return
	}

	var result *toolkit.Result
	result = client.Call("ping", toolkit.M{})
	checkResult(result, t)
}

func TestClientDouble(t *testing.T) {
	checkTestSkip(t)
	client2 := new(pecelclient.Client)
	e := client2.Connect(server.Address, serverSecret, "ariefdarmawan")
	if e == nil {
		client2.Close()
		t.Logf("Able to connect multi")
		return
	} else {
		t.Error(e)
	}
}

func TestClientHi(t *testing.T) {
	checkTestSkip(t)
	r := client.Call("hi", toolkit.M{}.Set("name", "Arief Darmawan"))
	checkResult(r, t)
}

type Server2RPC struct {
}

func (s *Server2RPC) Hi2(in toolkit.M) *toolkit.Result {
	return toolkit.NewResult().SetData("Hi from server 2")
}

func Test2Server(t *testing.T) {
	server2 := new(pecelserver.Server)
	server2.RegisterRPCFunctions(new(Server2RPC))
	server2.AddUser("admin", "admin")
	e := server2.Start("localhost:8888")
	if e != nil {
		t.Error("Unable to start server: " + e.Error())
		return
	}

	client2 := new(pecelclient.Client)
	client2.Connect(server2.Address, "admin", "admin")
	r := client2.Call("hi2", nil)

	if r.Status == toolkit.Status_NOK {
		t.Errorf("Call fail : %s", r.Message)
		return
	}

	if r.Data.(string) != "Hi from server 2" {
		t.Errorf("Fail, got " + r.Data.(string))
	}

	client2.Close()
	server2.Stop()
}

func TestStop(t *testing.T) {
	checkTestSkip(t)
	//server.Stop()
	if client != nil {
		client.Close()
	}
	server.Stop()
}
