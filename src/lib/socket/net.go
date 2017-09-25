package socket

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"lib/utils"
	"strings"
)

type Packet struct {
	proto   uint32
	count   uint32
	content []byte
}

func (p *Packet) SetProto(proto uint32) {
	p.proto = proto
}

func (p *Packet) SetContent(content []byte) {
	p.content = content
}

func (p *Packet) GetProto() uint32 {
	return p.proto
}
func (p *Packet) GetContent() []byte {
	return p.content
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 9 * time.Second
	maxMessageSize = 10 * 1024
	//连接建立后5秒内没有收到登陆请求，断开socket
	waitForLogin = time.Second * 5
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512 * 30,
	WriteBufferSize: 512 * 30,
}

func wSHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorln(err)
		}
	}()
	if r.Method != "GET" {
		glog.Errorln("Method is get, dropped")
		return
	}

	sock, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Errorln(err)
		return
	}

	ip := sock.RemoteAddr().String()
	iparr := strings.Split(ip, ":")
	var iip uint32
	if len(iparr) > 0 {
		ip = iparr[0]
		iip = utils.InetToaton(ip)
	}
	c := newConnection(sock, iip)
	go c.Reader(c.ReadChan)
	go c.LoginTimeout()
	go c.WritePump()
	c.ReadPump()
}
