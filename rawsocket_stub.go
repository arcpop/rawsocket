
// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package rawsocket

type rawSockConn struct {
    fd int
}