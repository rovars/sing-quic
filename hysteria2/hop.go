package hysteria2

import (
	"net"
	"sync/atomic"
	"time"

	M "github.com/metacubex/sing/common/metadata"
)

type rrPacketConn struct {
	net.PacketConn
	baseAddr  M.Socksaddr
	ports     []uint16
	addrIndex atomic.Uint32
}

func newRrPacketConn(pc net.PacketConn, baseAddr M.Socksaddr, ports []uint16) *rrPacketConn {
	return &rrPacketConn{
		PacketConn: pc,
		baseAddr:   baseAddr,
		ports:      ports,
	}
}

func (c *rrPacketConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	idx := c.addrIndex.Add(1) - 1
	port := c.ports[idx%uint32(len(c.ports))]
	targetAddr := net.UDPAddr{
		IP:   c.baseAddr.UDPAddr().IP,
		Port: int(port),
	}
	return c.PacketConn.WriteTo(b, &targetAddr)
}

func (c *rrPacketConn) LocalAddr() net.Addr {
	return c.PacketConn.LocalAddr()
}

func (c *rrPacketConn) Close() error {
	return c.PacketConn.Close()
}

func (c *rrPacketConn) SetDeadline(t time.Time) error {
	return c.PacketConn.SetDeadline(t)
}

func (c *rrPacketConn) SetReadDeadline(t time.Time) error {
	return c.PacketConn.SetReadDeadline(t)
}

func (c *rrPacketConn) SetWriteDeadline(t time.Time) error {
	return c.PacketConn.SetWriteDeadline(t)
}
