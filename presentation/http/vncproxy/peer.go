package vncproxy

import (
	"fmt"
	"net"
	"time"

	"github.com/evangwt/go-bufcopy"
)

const (
	defaultDialTimeout = 5 * time.Second
)

var (
	bcopy = bufcopy.New()
)

// peer represents a vnc proxy peer
// with a websocket connection and a vnc backend connection
type peer struct {
	source *VncWebSocket
	target net.Conn
}

func NewPeer(ws *VncWebSocket, addr string, dialTimeout time.Duration) (*peer, error) {
	if ws == nil {
		return nil, fmt.Errorf("websocket connection is nil")
	}

	if len(addr) == 0 {
		return nil, fmt.Errorf("addr is empty")
	}

	if dialTimeout <= 0 {
		dialTimeout = defaultDialTimeout
	}
	c, err := net.DialTimeout("tcp", addr, dialTimeout)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to vnc backend")
	}

	err = c.(*net.TCPConn).SetKeepAlive(true)
	if err != nil {
		return nil, fmt.Errorf("enable vnc backend connection keepalive failed")
	}

	err = c.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("set vnc backend connection keepalive period failed")
	}

	return &peer{
		source: ws,
		target: c,
	}, nil
}

// ReadSource copy source stream to target connection
func (p *peer) ReadSource() error {

	if _, err := bcopy.Copy(p.target, p.source); err != nil {
		return fmt.Errorf("copy source(%v) => target(%v) failed", p.source.RemoteAddr(), p.target.RemoteAddr())
	}
	return nil
}

// ReadTarget copys target stream to source connection
func (p *peer) ReadTarget() error {

	if _, err := bcopy.Copy(p.source, p.target); err != nil {
		return fmt.Errorf("copy target(%v) => source(%v) failed", p.target.RemoteAddr(), p.source.RemoteAddr())
	}
	return nil
}

// Close close the websocket connection and the vnc backend connection
func (p *peer) Close() {
	p.source.Close()
	p.target.Close()
}
