package goping

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type PingResult struct {
	Target  string `json:"target"`
	TTL     int    `json:"ttl,omitempty"`
	LastHop string `json:"lasthop"`
	RTT     string `json:"rtt,omitempty"`
	Message string `json:"message,omitempty"`
}

func Ping4(destination string, ttl int, timeout int) PingResult {
	result := PingResult{Target: destination, TTL: ttl}

	conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}
	defer conn.Close()

	p := ipv4.NewPacketConn(conn)
	if err := p.SetTTL(ttl); err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	dst, err := net.ResolveIPAddr("ip4", destination)
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
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
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	sT := time.Now()

	if _, err := conn.WriteTo(msgBytes, dst); err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	buf := make([]byte, 1280)
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))

	n, cm, addr, err := p.ReadFrom(buf)
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	eT := time.Now()
	result.RTT = fmt.Sprintf("%.3f", float64(eT.Sub(sT).Microseconds())/1000)

	if cm != nil {
		result.TTL = cm.TTL
	}

	reply, err := icmp.ParseMessage(1, buf[:n])
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	result.LastHop = addr.String()
	switch reply.Type {
	case ipv4.ICMPTypeEchoReply:
		result.Message = "suceed"
	case ipv4.ICMPTypeTimeExceeded:
		result.Message = "timeexceed"
	default:
		result.Message = fmt.Sprintf("%v", reply)
	}

	return result
}

func Ping6(destination string, ttl int, timeout int) PingResult {
	result := PingResult{Target: destination, TTL: ttl}

	conn, err := net.ListenPacket("ip6:ipv6-icmp", "::")
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}
	defer conn.Close()

	p := ipv6.NewPacketConn(conn)
	if err := p.SetHopLimit(ttl); err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	dst, err := net.ResolveIPAddr("ip6", destination)
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	icmpMessage := icmp.Message{
		Type: ipv6.ICMPTypeEchoRequest,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("icmp6"),
		},
	}

	msgBytes, err := icmpMessage.Marshal(nil)
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	sT := time.Now()

	if _, err := conn.WriteTo(msgBytes, dst); err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	buf := make([]byte, 1280)
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))

	n, cm, addr, err := p.ReadFrom(buf)
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	eT := time.Now()
	result.RTT = fmt.Sprintf("%.3f", float64(eT.Sub(sT).Microseconds())/1000)

	if cm != nil {
		result.TTL = cm.HopLimit
	}

	reply, err := icmp.ParseMessage(58, buf[:n])
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	result.LastHop = addr.String()
	switch reply.Type {
	case ipv6.ICMPTypeEchoReply:
		result.Message = "succeed"
	case ipv6.ICMPTypeTimeExceeded:
		result.Message = "timeexceed"
	default:
		result.Message = fmt.Sprintf("%v", reply)
	}

	return result
}
