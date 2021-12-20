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
	mongo    *MongoInstance
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

type createDeviceReq struct {
	Ip            string `json:"ip"`
	User          string `json:"user"`
	Password      string `json:"password"`
	DeviceType    string `json:"deviceType"`
	DeviceName    string `json:"deviceName"`
	High          string `json:"high"`
	Medium        string `json:"medium"`
	Low           string `json:"low"`
	Port          string `json:"port"`
	RTSPTransport string `json:"rtspTransport"`
	HTTPPort      string `json:"httpport"`
	PTZ           string `json:"ptz"`
}

type repModel struct {
	Data interface{} `json:"data"`
}

func New() *servives {
	fastHttp := router.New()
	redis := RConnect()
	mongo := MConnect()
	return &servives{fastHttp: fastHttp, redis: redis, mongo: mongo}
}

func (s *servives) FastHttp(host string, port int) {
	service := fmt.Sprintf("%s:%d", host, port)

	s.fastHttp.GET("/ping", s.pingHandler)
	s.fastHttp.GET("/api/devices/list", s.listDeviceMacHandler)

	s.fastHttp.POST("/api/mac/verification", s.verificationMacHandler)
	s.fastHttp.POST("/api/mac/create", s.createMacHandler)
	s.fastHttp.POST("/api/devices/create", s.createDevicesHandler)

	// Webhook for find list devices
	s.fastHttp.POST("/webhook/devices", s.webhookDevicesHandler)

	// Webhook for snapshots
	s.fastHttp.POST("/webhook/snapshots", s.webhookSnapshotsHandler)

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

func (s *servives) webhookDevicesHandler(ctx *fasthttp.RequestCtx) {
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

func (s *servives) webhookSnapshotsHandler(ctx *fasthttp.RequestCtx) {
	imageByte, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	err = UploadFile(imageByte)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	rep := &repModel{}
	rep.Data = true
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(200)
	ctx.Write(reply)
}

func (s *servives) createDevicesHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	v := &createDeviceReq{}
	rep := &repModel{}
	err := json.Unmarshal(ctx.PostBody(), v)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	// call server python pi for run
	rep.Data = true
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(201)
	ctx.Write(reply)
}
