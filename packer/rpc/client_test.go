package rpc

import (
	"net"
	"testing"
)

func testConn(t *testing.T) (net.Conn, net.Conn) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var serverConn net.Conn
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		defer l.Close()
		var err error
		serverConn, err = l.Accept()
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	clientConn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	<-doneCh

	return clientConn, serverConn
}

func testClientServer(t *testing.T) (*Client, *Server) {
	clientConn, serverConn := testConn(t)

	server := NewServer(serverConn)
	go server.Serve()

	client, err := NewClient(clientConn)
	if err != nil {
		server.Close()
		t.Fatalf("err: %s", err)
	}

	return client, server
}
