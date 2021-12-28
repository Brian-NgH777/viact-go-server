package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os/exec"
	"strings"
	"time"

	t "go/server/services/token"
)

const (
	REFRESH_TOKEN            = "refreshToken"
	ACCESS_TOKEN             = "AccessToken"
	REFRESH_TOKEN_SECRET_KEY = "VSC&lz`*y=4'?E2?Zc/{1/e2vXT2QQh\"JAua)EtO3Lr2&~UfuLVEd!dBO5s`,\"}"
	ACCESS_TOKEN_SECRET_KEY  = "r&0#T)F*~T;7rer])[mHrGt\"(/a^p~]UC.;k4K1%r}A+`8E\"F#,~IIAnI#~uU3S"
	HEALTH_KEY               = "healthCheck"
)

var (
	listScanDevices []*webhookDeviceItem
)

var (
	corsAllowHeaders     = "authorization"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

// services
type services struct {
	fastHttp *router.Router
	redis    *Redis
	mongo    *MongoInstance
}

// Requests
type macReq struct {
	Name       string `json:"name"`
	MacAddress string `json:"macAddress"`
}

type snapshotReq struct {
	Name string `json:"name"`
	Rtsp string `json:"rtsp"`
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

type extendPiAccessTokenReq struct {
	RefeshToken string `json:"refeshToken"`
	AccessToken string `json:"accessToken"`
}

type piRegisterReq struct {
	DeviceID string `json:"deviceID"`
}

type piHealthReq struct {
	DeviceID string `json:"deviceID"`
}

// responses
type respModel struct {
	Data interface{} `json:"data"`
}

type authPiResp struct {
	AccessToken string `json:"accessToken"`
	RefeshToken string `json:"refeshToken"`
}

// Schema mongo
type DeviceSchema struct {
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

// Webhook Req
type webhookDeviceReq struct {
	Data []*webhookDeviceItem `json:"data"`
}

type webhookDeviceItem struct {
	Ip     string `json:"ip"`
	Mac    string `json:"mac"`
	Vendor string `json:"vendor"`
}

func New() *services {
	fastHttp := router.New()
	redis := RConnect()
	mongo := MConnect()
	return &services{fastHttp: fastHttp, redis: redis, mongo: mongo}
}

func CORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		//ctx.Response.Header.SetBytesV("Access-Control-Allow-Origin", "*")

		//ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		//ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		//ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Content-Type", "application/json")

		next(ctx)
	}
}

func (s *services) FastHttp(host string, port int) {
	service := fmt.Sprintf("%s:%d", host, port)

	s.fastHttp.GET("/ping", s.pingHandler)
	// health check pi
	s.fastHttp.GET("/api/pi/status", s.piStatusHandler)
	s.fastHttp.POST("/api/pi/health", WebhookAuth(s.piHealthHandler))

	// auth pi
	s.fastHttp.POST("/api/pi/register", s.registerPiHandler)
	s.fastHttp.POST("/api/pi/access-token/extend", s.extendPiAccessTokenHandler)

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
	s.fastHttp.POST("/webhook/devices", WebhookAuth(s.webhookDevicesHandler))

	// Webhook for snapshots Handler
	s.fastHttp.POST("/api/snapshot", s.snapshotDeviceHandler) // call pythoncli for snapshot
	s.fastHttp.POST("/webhook/snapshots", WebhookAuth(s.webhookSnapshotsHandler))

	// Serve static files
	s.fastHttp.NotFound = fasthttp.FSHandler("/home/ec2-user/viact-go-server/static", 0)

	se := &fasthttp.Server{
		Handler:            CORS(s.fastHttp.Handler),
		MaxRequestBodySize: 100 * 1024 * 1024,
	}
	se.ListenAndServe(service)
}

func (s *services) pingHandler(ctx *fasthttp.RequestCtx) {
	ctx.Write([]byte("Ping Pong Pong"))
}

func (s *services) verificationMacHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &respModel{}
	v := &macReq{}
	err := json.Unmarshal(ctx.PostBody(), v)
	if err != nil {
		//ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		//ctx.SetContentType("text/plain")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(err.Error())
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

func (s *services) createMacHandler(ctx *fasthttp.RequestCtx) {
	rep := &respModel{}
	v := &macReq{}
	err := json.Unmarshal(ctx.PostBody(), v)
	if err != nil {
		//ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		//ctx.SetContentType("text/plain")
		//ctx.SetStatusCode(fasthttp.StatusBadRequest)
		//ctx.SetBodyString(err.Error())
		errors.New(err.Error())
		return
	}
	_, err = s.redis.HSet(v.MacAddress, "macAress", fmt.Sprintf("%s-%s", v.Name, v.MacAddress))
	if err != nil {
		//ctx.Error("HSet is false", fasthttp.StatusInternalServerError)
		//ctx.SetContentType("text/plain")
		//ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		//ctx.SetBodyString(err.Error())
		errors.New(err.Error())
		return
	}

	rep.Data = true
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.Write(reply)

}

func (s *services) listScanDeviceHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &respModel{}

	rep.Data = listScanDevices
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(200)
	ctx.Write(reply)
}

func (s *services) scanDeviceHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &respModel{}

	data, err := exec.Command("/usr/local/bin/action", "find_device").Output()
	if err != nil {
		ctx.Error(fmt.Sprintf("Run Command failed! Error:%s", err.Error()), fasthttp.StatusInternalServerError)
		return
	}

	rep.Data = string(data)
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(200)
	ctx.Write(reply)
}

func (s *services) snapshotDeviceHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &respModel{}

	v := &snapshotReq{}
	err := json.Unmarshal(ctx.PostBody(), v)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	arg := fmt.Sprintf("RTSP_LINK=%s FILE_NAME=%s", v.Rtsp, v.Name)
	fmt.Println("v.Rtsp, v.Name", arg)
	data, err := exec.Command("/usr/local/bin/action", "get_first_frame", arg).Output()
	if err != nil {
		ctx.Error(fmt.Sprintf("Run Command failed! Error:%s", err.Error()), fasthttp.StatusInternalServerError)
		return
	}

	rep.Data = string(data)
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(200)
	ctx.Write(reply)
}

func (s *services) streamDeviceHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &respModel{}
	var d DeviceSchema
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

	arg := fmt.Sprintf("RTSP_LINK=%s RTMP_LINK=%s", d.High, d.RTMP)
	_, err = exec.Command("/usr/local/bin/action", "livestream", arg).Output()
	if err != nil {
		ctx.Error(fmt.Sprintf("Run Command failed! Error:%s", err.Error()), fasthttp.StatusInternalServerError)
		return
	}

	rep.Data = d.RTMP
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(200)
	ctx.Write(reply)
}

func (s *services) listDeviceHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &respModel{}
	devices := []*DeviceSchema{}

	collectionDevice := s.mongo.Db.Collection("devices")
	cur, err := collectionDevice.Find(ctx, bson.D{})
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result DeviceSchema
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

func (s *services) createDevicesHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	v := &createDeviceReq{}
	rep := &respModel{}
	device := &DeviceSchema{}
	err := json.Unmarshal(ctx.PostBody(), v)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
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
	rep.Data = true
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(201)
	ctx.Write(reply)
}

func (s *services) registerPiHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &authPiResp{}
	v := &piRegisterReq{}
	err := json.Unmarshal(ctx.PostBody(), v)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(v.DeviceID)) == 0 {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	tAccess, err := generatorAccessToken(v)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	tRefresh, err := generatorRefreshToken(v)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	rep.AccessToken = tAccess
	rep.RefeshToken = tRefresh
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(200)
	ctx.Write(reply)
}

func (s *services) extendPiAccessTokenHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &authPiResp{}
	v := &extendPiAccessTokenReq{}
	err := json.Unmarshal(ctx.PostBody(), v)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	payload, err := verifyToken(v.RefeshToken, REFRESH_TOKEN)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusUnauthorized)
		return
	}

	tAccess := v.AccessToken
	_, err = verifyToken(v.AccessToken, ACCESS_TOKEN)
	if err != nil {
		if errors.Is(err, t.ErrExpiredToken) {
			tAccess, err = generatorAccessToken(&piRegisterReq{DeviceID: payload.DeviceID})
			if err != nil {
				ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
				return
			}
		} else {
			ctx.Error(err.Error(), fasthttp.StatusUnauthorized)
			return
		}
	}

	rep.AccessToken = tAccess
	rep.RefeshToken = v.RefeshToken
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(200)
	ctx.Write(reply)
}

func (s *services) piHealthHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &respModel{}
	v := &piHealthReq{}
	err := json.Unmarshal(ctx.PostBody(), v)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	d, err := s.redis.Set(HEALTH_KEY, v.DeviceID, time.Second*5)
	if err != nil {
		ctx.Error("Set is false", fasthttp.StatusInternalServerError)
		return
	}

	rep.Data = d
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(200)
	ctx.Write(reply)
}

func (s *services) piStatusHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	rep := &respModel{}
	rep.Data = true

	_, err := s.redis.Get(HEALTH_KEY)
	if err != nil {
		rep.Data = false
	}

	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(200)
	ctx.Write(reply)
}

func (s *services) webhookDevicesHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	v := &webhookDeviceReq{}
	rep := &respModel{}
	err := json.Unmarshal(ctx.PostBody(), v)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	listScanDevices = v.Data
	rep.Data = true
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(201)
	ctx.Write(reply)
}

func (s *services) webhookSnapshotsHandler(ctx *fasthttp.RequestCtx) {
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

	rep := &respModel{}
	rep.Data = true
	reply, _ := json.Marshal(rep)
	ctx.SetStatusCode(200)
	ctx.Write(reply)
}

func generatorAccessToken(data *piRegisterReq) (string, error) {
	maker, err := t.NewJWTMaker(ACCESS_TOKEN_SECRET_KEY)

	deviceID := data.DeviceID
	duration := time.Hour * 24 * 7

	token, err := maker.CreateToken(deviceID, duration)

	return token, err
}

func generatorRefreshToken(data *piRegisterReq) (string, error) {
	maker, err := t.NewJWTMaker(REFRESH_TOKEN_SECRET_KEY)

	deviceID := data.DeviceID
	duration := time.Hour * 24 * 365

	token, err := maker.CreateToken(deviceID, duration)

	return token, err
}

func verifyToken(token string, typeAuth string) (*t.Payload, error) {
	key := ACCESS_TOKEN_SECRET_KEY
	if typeAuth == REFRESH_TOKEN {
		key = REFRESH_TOKEN_SECRET_KEY
	}
	maker, err := t.NewJWTMaker(key)
	payload, err := maker.VerifyToken(token)
	return payload, err
}

func WebhookAuth(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		auth := ctx.Request.Header.Peek("Authorization")
		if auth == nil {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
		}
		payload, err := verifyToken(string(auth), ACCESS_TOKEN)
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusUnauthorized)
		} else {
			fmt.Println("payloadpayload", payload)
			h(ctx)
			return
		}

		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
		ctx.Response.Header.Set("WWW-Authenticate", "Basic realm=Restricted")
	})
}
