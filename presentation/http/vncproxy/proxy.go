package vncproxy

import (
	"fmt"
	// "io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Config represents vnc proxy config
type Config struct {
	DialTimeout time.Duration
	EndPoint    string
}

type VncWebSocket struct {
	data []byte
	pos  int
	conn *websocket.Conn
}

func NewVncWebSocket(c *websocket.Conn) *VncWebSocket {
	return &VncWebSocket{
		conn: c,
	}
}

func (ws *VncWebSocket) Write(msg []byte) (int, error) {
	if ws.conn != nil {
		err := ws.conn.WriteMessage(websocket.BinaryMessage, msg)
		return len(msg), err
	}
	return 1, fmt.Errorf("websocket not ready")

}

func (ws *VncWebSocket) Read(msg []byte) (n int, err error) {
	if ws.conn != nil {
		_, ws.data, err = ws.conn.ReadMessage()
		// if ws.pos >= len(ws.data) {
		// 	return 0, io.EOF
		// }
		n = copy(msg, ws.data)
		// ws.pos += n
		return
	}
	return 1, fmt.Errorf("websocket not ready")
}

func (ws *VncWebSocket) LocalAddr() net.Addr {
	if ws.conn != nil {
		return ws.conn.LocalAddr()
	}
	return nil
}

func (ws *VncWebSocket) RemoteAddr() net.Addr {
	if ws.conn != nil {
		return ws.conn.RemoteAddr()
	}
	return nil
}

func (ws *VncWebSocket) Close() error {
	if ws.conn != nil {
		return ws.conn.Close()
	}
	return nil
}

// Proxy represents vnc proxy
type Proxy struct {
	dialTimeout time.Duration // Timeout for connecting to each target vnc server
	peers       map[*peer]struct{}
	l           sync.RWMutex
	endPoint    string
}

func New(conf *Config) *Proxy {
	return &Proxy{
		dialTimeout: conf.DialTimeout,
		peers:       make(map[*peer]struct{}),
		l:           sync.RWMutex{},
		endPoint:    conf.EndPoint,
	}
}

// ServeWS provides websocket handler
func (p *Proxy) ServeWS(conn *websocket.Conn) {
	ws := NewVncWebSocket(conn)
	peer, err := NewPeer(ws, p.endPoint, p.dialTimeout)

	if err != nil {
		return
	}

	p.addPeer(peer)
	defer func() {
		p.deletePeer(peer)
	}()

	go func() {
		if err := peer.ReadTarget(); err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			return
		}

	}()

	if err = peer.ReadSource(); err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			return
		}
		return
	}

}

func (p *Proxy) addPeer(peer *peer) {
	p.l.Lock()
	p.peers[peer] = struct{}{}
	p.l.Unlock()
}

func (p *Proxy) deletePeer(peer *peer) {
	p.l.Lock()
	delete(p.peers, peer)
	peer.Close()
	p.l.Unlock()
}

func (p *Proxy) Peers() map[*peer]struct{} {
	p.l.RLock()
	defer p.l.RUnlock()
	return p.peers
}
