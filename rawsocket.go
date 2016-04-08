package rawsocket

import (
    "syscall"
    "net"
)


//RawSockConn abstracts the data which is needed to create a raw socket to send packets with IP_HDRINCL
type RawSockConn struct {
    fd syscall.Handle
}

//CreateRawSocket creates a new RawSocket to send IP packets to the interface.
func CreateRawSocket(ipHdrInc bool) (*RawSockConn, error)  {
    if ipHdrInc {
        fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, 0xFF /* syscall.IPPROTO_RAW */)
        if err != nil {
            return nil, err
        }
        return &RawSockConn{ fd: fd }, nil
    }
    fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_IP)
    if err != nil {
        return nil, err
    }
    return &RawSockConn{ fd: fd }, nil
}

//WriteTo sends a packet to the destination indicated by ip. Be wary that you need to have an IP header in the buffer b!!!
func (s *RawSockConn) WriteTo(b []byte, ip net.IP) (n int, err error) {
    var sa syscall.Sockaddr
    ip4 := ip.To4()
    if ip4 == nil {
        sa6 := &syscall.SockaddrInet6 {
            Port: 0,
            ZoneId: 0,
        }
        copy(sa6.Addr[:], ip)
        sa = sa6
    } else {
        sa4 := &syscall.SockaddrInet4 {
            Port: 0,
        }
        copy(sa4.Addr[:], ip4)
        sa = sa4
    }
    err = syscall.Sendto(s.fd, b, 0, sa)
    if err != nil {
        return
    }
    n = len(b)
    return n, nil
}

//Close closes the associated resources.
func (s *RawSockConn) Close() error {
    return syscall.Close(s.fd)
}


//ReadFrom receives raw packets. Not really sure why you need this, but addr may return nil without error, when
//the returned sockaddr can't be casted to either ipv4 or ipv6. Better not use that method and use the go implementation. 
func (s * RawSockConn) ReadFrom(b []byte) (n int, addr net.IP, err error) {
    var from syscall.Sockaddr
    n, from, err = syscall.Recvfrom(s.fd, b, 0)
    
    if err != nil {
        return
    }
    
    if sa4, ok := from.(*syscall.SockaddrInet4); ok {
        return n, net.IP(sa4.Addr[:]), nil
    }
    if sa6, ok := from.(*syscall.SockaddrInet6); ok {
        return n, net.IP(sa6.Addr[:]), nil
    }
    
    return n, nil, nil
}