package main

import (
	"errors"
	"io"
	"net"
	"time"
)

var errConnClosedByServer = errors.New("...Connection was closed by peer")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnet struct {
	address string
	in      io.ReadCloser
	out     io.Writer
	timeout time.Duration
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnet{address: address, timeout: timeout, in: in, out: out}
}

func (c *telnet) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	c.conn = conn
	return err
}

func (c *telnet) Send() error {
	_, err := io.Copy(c.conn, c.in)
	return err
}

func (c *telnet) Receive() error {
	w, err := io.Copy(c.out, c.conn)
	if w == 0 && err == nil {
		return errConnClosedByServer
	}
	return err
}

func (c *telnet) Close() error {
	var err1 error
	if c.conn != nil {
		err1 = c.conn.Close()
	}
	err2 := c.in.Close()
	return errors.Join(err1, err2)
}
