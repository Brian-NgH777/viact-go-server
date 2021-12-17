package services

import (
	"encoding/json"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

var (
	listDevices []*webhookDeviceItem
)

type servives struct {
	fastHttp *router.Router
	redis    *Redis
}

type macReq struct {
	Name       string `json:"name"`
	MacAddress string `json:"macAddress"`
}

type webhookDeviceReq struct {
	Data []*webhookDeviceItem `json:"data"`
}

type webhookDeviceItem struct {
	Ip     string `json:"ip"`
	Mac    string `json:"mac"`
	Vendor string `json:"vendor"`
}

type repModel struct {
	Data interface{} `json:"data"`
}

func New() *servives {
	fastHttp := router.New()
	redis := RConnect()
	return &servives{fastHttp: fastHttp, redis: redis}
}

func (s *servives) FastHttp(host string, port int) {
	service := fmt.Sprintf("%s:%d", host,port)

	s.fastHttp.GET("/ping", s.pingHandler)
	s.fastHttp.GET("/api/device/list", s.listDeviceMacHandler)
	s.fastHttp.POST("/api/mac/verification", s.verificationMacHandler)
	s.fastHttp.POST("/api/mac/create", s.createMacHandler)
	// Webhook for find list devices
	s.fastHttp.POST("/webhook/devices", s.createDevicesHandler)

	// Webhook for live streaming
	s.fastHttp.POST("/webhook/", s.createDevicesHandler)

	fasthttp.ListenAndServe(service, s.fastHttp.Handler)
}

func (s *servives) pingHandler(ctx *fasthttp.RequestCtx) {
	ctx.Write([]byte("Ping Pong Pong"))
}

func (s *servives) verificationMacHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &repModel{}
	v := &macReq{}
	err := json.Unmarshal(ctx.PostBody(), v)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	rep.Data = true
	_, err = s.redis.HGet(v.MacAddress, "macAress")
	if err != nil {
		rep.Data = false
	}
	reply, _ := json.Marshal(rep)

	ctx.Write(reply)
}

func (s *servives) createMacHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &repModel{}
	v := &macReq{}
	err := json.Unmarshal(ctx.PostBody(), v)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	_, err = s.redis.HSet(v.MacAddress, "macAress", fmt.Sprintf("%s-%s", v.Name, v.MacAddress))
	if err != nil {
		ctx.Error("HSet is false", fasthttp.StatusInternalServerError)
	}

	rep.Data = true
	reply, _ := json.Marshal(rep)

	ctx.Write(reply)

}

func (s *servives) listDeviceMacHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &repModel{}

	rep.Data = listDevices
	reply, _ := json.Marshal(rep)

	ctx.Write(reply)
}

func (s *servives) createDevicesHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	v := &webhookDeviceReq{}
	rep := &repModel{}
	err := json.Unmarshal(ctx.PostBody(), v)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	listDevices = v.Data
	rep.Data = true
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(201)
	ctx.Write(reply)
}
