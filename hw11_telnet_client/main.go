package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

var (
	isConnClosedByServer atomic.Bool
	isClientExited       atomic.Bool
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout used to connect to server")
	flag.Parse()
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "invalid arguments")
		return
	}

	addrArgs := os.Args[len(os.Args)-2:]
	address := net.JoinHostPort(addrArgs[0], addrArgs[1])
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	in := &bytes.Buffer{}
	tc := NewTelnetClient(address, timeout, io.NopCloser(in), os.Stdout)
	defer func() {
		err := tc.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err := tc.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	stdin := stdinScan(in)
	go receiving(ctx, tc)
	sending(ctx, tc, stdin)
}

func receiving(ctx context.Context, tc TelnetClient) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := tc.Receive()
			if err != nil {
				if errors.Is(err, errConnClosedByServer) {
					isConnClosedByServer.Store(true)
					return
				} else if isClientExited.Load() {
					return
				}
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func sending(ctx context.Context, tc TelnetClient, stdin <-chan struct{}) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-stdin:
			if isClientExited.Load() {
				return
			}
			err := tc.Send()
			if err != nil {
				if isConnClosedByServer.Load() {
					fmt.Fprintln(os.Stderr, errConnClosedByServer)
					return
				}
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func stdinScan(in *bytes.Buffer) <-chan struct{} {
	out := make(chan struct{})
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			in.WriteString(scanner.Text() + "\n")
			out <- struct{}{}
		}
		fmt.Fprintln(os.Stderr, "...EOF")
		isClientExited.Store(true)
		close(out)
	}()
	return out
}
