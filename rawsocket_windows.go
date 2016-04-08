package rawsocket

import (
    "syscall"
)

type rawSockConn struct {
    fd syscall.Handle
}