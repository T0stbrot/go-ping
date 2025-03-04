package goping

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type PingResult struct {
	Target  string `json:"target"`
	TTL     int    `json:"ttl"`
	LastHop string `json:"lasthop"`
	RTT     string `json:"rtt,omitempty"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func Ping4(destination string, ttl int, timeout int) PingResult {
	result := PingResult{Target: destination, TTL: ttl}

	conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		result.Message = "error"
		result.Error = fmt.Sprintf("%v", err)
		return result
	}
	defer conn.Close()

	p := ipv4.NewPacketConn(conn)
	if err := p.SetTTL(ttl); err != nil {
		result.Message = "error"
		result.Error = fmt.Sprintf("%v", err)
		return result
	}

	dst, err := net.ResolveIPAddr("ip4", destination)
	if err != nil {
		result.Message = "error"
		result.Error = fmt.Sprintf("%v", err)
		return result
	}

	icmpMessage := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("icmp"),
		},
	}

	msgBytes, err := icmpMessage.Marshal(nil)
	if err != nil {
		result.Message = "error"
		result.Error = fmt.Sprintf("%v", err)
		return result
	}

	sT := time.Now()

	if _, err := conn.WriteTo(msgBytes, dst); err != nil {
		result.Message = "error"
		result.Error = fmt.Sprintf("%v", err)
		return result
	}

	buf := make([]byte, 1492)
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))

	n, addr, err := conn.ReadFrom(buf)
	if err != nil {
		result.Message = "error"
		result.Error = fmt.Sprintf("%v", err)
		return result
	}

	eT := time.Now()
	result.RTT = fmt.Sprintf("%.3fms", float64(eT.Sub(sT).Microseconds())/1000)

	reply, err := icmp.ParseMessage(1, buf[:n])
	if err != nil {
		result.Message = "error"
		result.Error = fmt.Sprintf("%v", err)
		return result
	}

	result.LastHop = addr.String()
	switch reply.Type {
	case ipv4.ICMPTypeEchoReply:
		result.Message = "suceed"
	case ipv4.ICMPTypeTimeExceeded:
		result.Message = "timeexceed"
	default:
		result.Message = "error"
		result.Error = fmt.Sprintf("%v", reply)
	}

	return result
}
