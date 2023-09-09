package network

import (
	"log"
	"net"
	"sync"
)

type ConnectionHandlerFunc func(conn net.Conn, s *Server)

type Server struct {
	listener                net.Listener
	quit                    chan interface{}
	wg                      sync.WaitGroup
	handleConnectionFunc    ConnectionHandlerFunc
}

func NewServer(addr string, connHandler ConnectionHandlerFunc) *Server {
	s := &Server{
		quit: make(chan interface{}),
        handleConnectionFunc: connHandler,
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Server failed on listen(): %v", err)
	}
	s.listener = l
	s.wg.Add(1)

    go s.serve()

    return s
}

func (s *Server) serve() {
    defer s.wg.Done()

    for {
        conn, err := s.listener.Accept()

        if err != nil {
            // Select chooses b/w multiple channels
            // Checks in a non-blocking way
            select {
            case <- s.quit:
                return
            default:
                log.Println("Accept() error", err)
            }
        } else {
            s.wg.Add(1)
            go func() {
                s.handleConnectionFunc(conn, s)
                s.wg.Done()
            }()
        }
        
    }
}

func (s *Server) Stop() {
    close(s.quit)       // close the channel
    s.listener.Close()
    s.wg.Wait()
}