package tracker

import (
	"fmt"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/scionproto/scion/go/lib/snet"
	"net"
)

type Peer struct {
	IsScionPeer bool
	udpAddr     *snet.UDPAddr
	IP          net.IP
	Port        int
	ID          []byte
}

// Set from the non-compact form in BEP 3.
func (p *Peer) fromDictInterface(d map[string]interface{}) {
	if _, ok := d["ia"]; ok {
		p.parseScionPeer(d)
		return
	}
	p.parseDefaultPeer(d)
}

func (p *Peer) parseDefaultPeer(d map[string]interface{}) {
	p.IP = net.ParseIP(d["ip"].(string))
	if _, ok := d["peer id"]; ok {
		p.ID = []byte(d["peer id"].(string))
	}
	p.Port = int(d["port"].(int64))
	p.IsScionPeer = false
}

func (p *Peer) parseScionPeer(d map[string]interface{}) {
	udpStr := fmt.Sprintf("%s,[%s]:%d", d["ia"], d["ip"], d["port"])
	udpAddr, err := snet.ParseUDPAddr(udpStr)
	if err != nil {
		p.IsScionPeer = true
		p.udpAddr = udpAddr
	}
	if _, ok := d["peer id"]; ok {
		p.ID = []byte(d["peer id"].(string))
	}
}

func (p Peer) FromNodeAddr(na krpc.NodeAddr) Peer {
	p.IP = na.IP
	p.Port = na.Port
	return p
}
