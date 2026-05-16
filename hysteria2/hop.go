package hysteria2

import (
	"net"
	"sync/atomic"
	"time"

	M "github.com/metacubex/sing/common/metadata"
)

type hopPacketConn struct {
	net.PacketConn
	baseAddr  M.Socksaddr
	ports     []uint16
	addrIndex atomic.Uint32
}

func newHopPacketConn(pc net.PacketConn, baseAddr M.Socksaddr, ports []uint16) *hopPacketConn {
	return &hopPacketConn{
		PacketConn: pc,
		baseAddr:   baseAddr,
		ports:      ports,
	}
}

func (c *hopPacketConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	idx := c.addrIndex.Add(1) - 1
	port := c.ports[idx%uint32(len(c.ports))]
	targetAddr := net.UDPAddr{
		IP:   c.baseAddr.UDPAddr().IP,
		Port: int(port),
	}
	return c.PacketConn.WriteTo(b, &targetAddr)
}

func (c *hopPacketConn) LocalAddr() net.Addr {
	return c.PacketConn.LocalAddr()
}

func (c *hopPacketConn) Close() error {
	return c.PacketConn.Close()
}

func (c *hopPacketConn) SetDeadline(t time.Time) error {
	return c.PacketConn.SetDeadline(t)
}

func (c *hopPacketConn) SetReadDeadline(t time.Time) error {
	return c.PacketConn.SetReadDeadline(t)
}

func (c *hopPacketConn) SetWriteDeadline(t time.Time) error {
	return c.PacketConn.SetWriteDeadline(t)
}
