package services

import (
	"encoding/json"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os/exec"
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

type snapshotReq struct {
	Name string `json:"name"`
	Rtsp string `json:"rtsp"`
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
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
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
	RTMP          string             `json:"rtmp,omitempty" bson:"rtmp,omitempty"`
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

	// Live stream Handler
	s.fastHttp.GET("/api/device/stream/{id}", s.streamDeviceHandler) // call pythoncli for live stream

	// verification Mac Handler
	s.fastHttp.POST("/api/mac/verification", s.verificationMacHandler)
	s.fastHttp.POST("/api/mac/create", s.createMacHandler)

	// Device Handler
	s.fastHttp.GET("/api/device/list", s.listDeviceHandler)
	s.fastHttp.POST("/api/devices/create", s.createDevicesHandler)

	// Webhook for find list devices Handler
	s.fastHttp.GET("/api/scan-device/list", s.listScanDeviceHandler)
	s.fastHttp.GET("/api/scan-device", s.scanDeviceHandler) // call pythoncli for scan
	s.fastHttp.POST("/webhook/devices", s.webhookDevicesHandler)

	// Webhook for snapshots Handler
	s.fastHttp.POST("/api/snapshot", s.snapshotDeviceHandler) // call pythoncli for snapshot
	s.fastHttp.POST("/webhook/snapshots", s.webhookSnapshotsHandler)

	// Serve static files
	s.fastHttp.NotFound = fasthttp.FSHandler("/home/ec2-user/viact-go-server/static", 0)
	//fasthttp.ListenAndServe(service)

	se := &fasthttp.Server{
		Handler:            s.fastHttp.Handler,
		MaxRequestBodySize: 100 * 1024 * 1024,
	}
	se.ListenAndServe(service)
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

func (s *servives) listScanDeviceHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &repModel{}

	rep.Data = listScanDevices
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(200)
	ctx.Write(reply)
}

func (s *servives) scanDeviceHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &repModel{}

	data, err := exec.Command("/usr/local/bin/action", "find_device").Output()
	if err != nil {
		ctx.Error("Run Command failed!", fasthttp.StatusInternalServerError)
		return
	}

	rep.Data = string(data)
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(200)
	ctx.Write(reply)
}

func (s *servives) snapshotDeviceHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &repModel{}

	v := &snapshotReq{}
	err := json.Unmarshal(ctx.PostBody(), v)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	data, err := exec.Command("/bin/bash", "/usr/bin/cmd.sh","get_first_frame", v.Rtsp, v.Name).Output()
	if err != nil {
		fmt.Println("errerrerrerrerr", err.Error())
		ctx.Error("Run Command failed!", fasthttp.StatusInternalServerError)
		return
	}

	rep.Data = string(data)
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(200)
	ctx.Write(reply)
}

func (s *servives) streamDeviceHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &repModel{}
	var d Device
	id := ctx.UserValue("id").(string)
	collectionDevice := s.mongo.Db.Collection("devices")
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objID}}
	err := collectionDevice.FindOne(ctx, filter).Decode(&d)
	if err != nil {
		ctx.Error("Not found device", fasthttp.StatusInternalServerError)
		//if err == mongo.ErrNoDocuments {
		//	return
		//}
		return
	}
	arg := fmt.Sprintf("livestream \"RTSP_LINK=%s RTMP_LINK=%s\"", d.High, d.RTMP)
	_, err = exec.Command("/usr/local/bin/action", arg).Output()
	if err != nil {
		ctx.Error("Run Command failed!", fasthttp.StatusInternalServerError)
		return
	}

	rep.Data = d.RTMP
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(200)
	ctx.Write(reply)
}

func (s *servives) listDeviceHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &repModel{}
	devices := []*Device{}

	collectionDevice := s.mongo.Db.Collection("devices")
	cur, err := collectionDevice.Find(ctx, bson.D{})
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

	now := time.Now()
	sec := now.Unix()

	collectionDevice := s.mongo.Db.Collection("devices")
	device.IP = v.Ip
	device.User = v.User
	device.Password = v.Password
	device.DeviceType = v.DeviceType
	device.DeviceName = v.DeviceName
	device.High = v.High
	if len(v.High) == 0 {
		device.High = "rtsp://admin:Viact123@192.168.92.111/live"
	}
	device.Medium = v.Medium
	device.Low = v.Low
	device.Port = v.Port
	device.RTSPTransport = v.RTSPTransport
	device.HTTPPort = v.HTTPPort
	device.PTZ = v.PTZ
	device.Thumbnail = v.Thumbnail
	device.CameraName = v.CameraName
	device.RTMP = fmt.Sprintf("rtmp://54.254.0.41/live/test%d", sec)
	device.CreatedAt = time.Now().UTC()
	device.UpdatedAt = time.Now().UTC()

	_, err = collectionDevice.InsertOne(ctx, device)
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

	//ex, err := os.Executable()
	//if err != nil {
	//	ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	//	return
	//}
	//filepath := path.Join(filepath.Dir(ex), fmt.Sprintf("%s%s", "/home/ec2-user/viact-go-server/static/", imageByte.Filename))
	if err = fasthttp.SaveMultipartFile(imageByte, fmt.Sprintf("%s%s", "/home/ec2-user/viact-go-server/static/", imageByte.Filename)); err != nil {
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
