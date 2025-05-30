package goping

import (
	"fmt"
	"net"
	"syscall"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type PingResult struct {
	Target  string `json:"target"`
	LastHop string `json:"lasthop"`
	RTT     string `json:"rtt,omitempty"`
	Message string `json:"message,omitempty"`
}

type PingProto struct {
	IP string
	Listen string
	Type icmp.Type
	Conn4 ipv4.PacketConn
	Conn6 ipv6.PacketConn
}

func Ping(ver int, destination string, ttl int, timeout int) PingResult {
	result := PingResult{Target: destination}

	proto := PingProto{IP: "ip4", Listen: "0.0.0.0", Type: ipv4.ICMPTypeEcho}

	if ver == 6 {
		proto.IP = "ip6"
		proto.Listen = "::"
		proto.Type = ipv6.ICMPTypeEchoRequest
	}

	conn, err := net.ListenPacket(fmt.Sprintf("%s:icmp", proto.IP), proto.Listen)
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}
	defer conn.Close()

	dst, err := net.ResolveIPAddr(proto.IP, destination)
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	icmpMessage := icmp.Message{
		Type: proto.Type,
		Code: 0,
		Body: &icmp.Echo{
			ID:   syscall.Getpid() & 0xffff,
			Seq:  1,
			Data: make([]byte, 16),
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

	n, addr, err := conn.ReadFrom(buf)
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	eT := time.Now()
	result.RTT = fmt.Sprintf("%.3f", float64(eT.Sub(sT).Microseconds())/1000)

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
	case ipv6.ICMPTypeEchoReply:
		result.Message = "succeed"
	case ipv6.ICMPTypeTimeExceeded:
		result.Message = "timeexceed"
	default:
		result.Message = fmt.Sprintf("%v", reply)
	}

	return result
}

func Ping4(destination string, ttl int, timeout int) PingResult {
	result := PingResult{Target: destination}

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
			ID:   syscall.Getpid() & 0xffff,
			Seq:  1,
			Data: make([]byte, 16),
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

	n, _, addr, err := p.ReadFrom(buf)
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	eT := time.Now()
	result.RTT = fmt.Sprintf("%.3f", float64(eT.Sub(sT).Microseconds())/1000)

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
	result := PingResult{Target: destination}

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
			ID:   syscall.Getpid() & 0xffff,
			Seq:  1,
			Data: make([]byte, 16),
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

	n, _, addr, err := p.ReadFrom(buf)
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	eT := time.Now()
	result.RTT = fmt.Sprintf("%.3f", float64(eT.Sub(sT).Microseconds())/1000)

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
