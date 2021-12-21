package services

import (
	"encoding/json"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var (
	listScanDevices []*webhookDeviceItem
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
	Thumbnail     string `json:"thumbnail"`
	CameraName    string `json:"cameraName"`
}

type repModel struct {
	Data interface{} `json:"data"`
}

type Device struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	IP            string             `json:"ip,omitempty" bson:"ip,omitempty"`
	User          string             `json:"user,omitempty" bson:"user,omitempty"`
	Password      string             `json:"password,omitempty" bson:"password,omitempty"`
	DeviceType    string             `json:"deviceType,omitempty" bson:"deviceType,omitempty"`
	DeviceName    string             `json:"deviceName,omitempty" bson:"deviceName,omitempty"`
	High          string             `json:"high,omitempty" bson:"high,omitempty"`
	Medium        string             `json:"medium,omitempty" bson:"medium,omitempty"`
	Low           string             `json:"low,omitempty" bson:"low,omitempty"`
	Port          string             `json:"port,omitempty" bson:"port,omitempty"`
	RTSPTransport string             `json:"rtspTransport,omitempty" bson:"rtspTransport,omitempty"`
	HTTPPort      string             `json:"httpPort,omitempty" bson:"httpPort,omitempty"`
	PTZ           string             `json:"ptz,omitempty" bson:"ptz,omitempty"`
	Thumbnail     string             `json:"thumbnail,omitempty" bson:"thumbnail,omitempty"`
	CameraName    string             `json:"cameraName,omitempty" bson:"cameraName,omitempty"`
	CreatedAt     time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt     time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
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
	s.fastHttp.GET("/api/scan-device/list", s.listScanDeviceMacHandler)
	s.fastHttp.GET("/api/device/list", s.listDeviceMacHandler)

	s.fastHttp.POST("/api/mac/verification", s.verificationMacHandler)
	s.fastHttp.POST("/api/mac/create", s.createMacHandler)
	s.fastHttp.POST("/api/devices/create", s.createDevicesHandler)

	// Webhook for find list devices
	s.fastHttp.POST("/webhook/devices", s.webhookDevicesHandler)

	// Webhook for snapshots
	s.fastHttp.POST("/webhook/snapshots", s.webhookSnapshotsHandler)

	s.fastHttp.NotFound = fasthttp.FSHandler("./static", 0)

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
	ctx.SetStatusCode(200)
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
	ctx.SetStatusCode(201)
	ctx.Write(reply)

}

func (s *servives) listScanDeviceMacHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &repModel{}

	rep.Data = listScanDevices
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(200)
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

	listScanDevices = v.Data
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
	path := fmt.Sprintf("%s%s", "../static/", imageByte.Filename)
	if err = fasthttp.SaveMultipartFile(imageByte, path); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	// err = UploadFile(imageByte)
	//if err != nil {
	//	ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	//	return
	//}

	rep := &repModel{}
	rep.Data = true
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(200)
	ctx.Write(reply)
}

func (s *servives) listDeviceMacHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &repModel{}
	devices := []*Device{}

	collectionBooking := s.mongo.Db.Collection("devices")
	cur, err  := collectionBooking.Find(ctx, bson.D{})
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result Device
		err = cur.Decode(&result)
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
		devices = append(devices, &result)
	}
	if err = cur.Err(); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	rep.Data = devices
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(201)

	ctx.Write(reply)
}


func (s *servives) createDevicesHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	v := &createDeviceReq{}
	rep := &repModel{}
	device := &Device{}
	err := json.Unmarshal(ctx.PostBody(), v)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	collectionBooking := s.mongo.Db.Collection("devices")
	device.IP = v.Ip
	device.User = v.User
	device.Password = v.Password
	device.DeviceType = v.DeviceType
	device.DeviceName = v.DeviceName
	device.High = v.High
	device.Medium = v.Medium
	device.Low = v.Low
	device.Port = v.Port
	device.RTSPTransport = v.RTSPTransport
	device.HTTPPort = v.HTTPPort
	device.PTZ = v.PTZ
	device.Thumbnail = v.Thumbnail
	device.CameraName = v.CameraName
	device.CreatedAt = time.Now().UTC()
	device.UpdatedAt = time.Now().UTC()

	_, err = collectionBooking.InsertOne(ctx, device)
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
