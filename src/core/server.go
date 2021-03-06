package core

import (
	"log"
)

const (
	MAXCLIENTS = 50
)

type IServer interface {
	Listen(port string)
	Accept() (tc *TCPClient)
	Close()
}

type ClientTable map[IClient]*Client

type Server struct {
	s IServer
	// listener net.Listener
	clients ClientTable
	pending chan IClient
	// quiting  chan net.Conn
	incoming chan string
	outgoing chan string
}

func CreateServer() (server *Server) {
	server = &Server{
		s:       &TCPServer{},
		clients: make(ClientTable),
		pending: make(chan IClient),
		// quiting:  make(chan net.Conn),
		incoming: make(chan string),
		outgoing: make(chan string),
	}
	return
}

func (self *Server) Listen(port string) {
	go func() {
		for {
			select {
			case msg := <-self.incoming:
				log.Println(msg)
			case client := <-self.pending:
				self.Join(client)
				//
				// case conn := <-server.quiting:
			}
		}
	}()
	self.s.Listen(port)
}

func (self *Server) Join(ic IClient) {
	client := CreateClient(ic)
	self.clients[ic] = client

	go func() {
		for msg := range client.incoming {
			// package msg whish conn
			// msg = fmt.Sprintf("format string", a ...interface{})
			if !MsgRoute(client, msg) {
				client.PutOutgoing("command error, Usage:'chatroom join 1','chatroom send hello'")
				// self.incoming <- msg
			}
		}
	}()
}

func (self *Server) Start(port string) {
	self.Listen(port)
	logger.Println("server start")
	// l, _ := net.Listen("tcp", fmt.Sprintf(":%s", port))
	// self.listener = l
	// defer self.listener.Close()
	defer self.s.Close()
	// chan listen

	for {
		self.pending <- self.s.Accept()
		// if conn, err := self.listener.Accept(); err == nil {
		// 	self.pending <- conn
		// go func(c net.Conn) {
		// 	buf := make([]byte, 1024)
		// 	for {
		// 		cn, err := c.Read(buf)
		// 		if err != nil {
		// 			c.Close()
		// 			break
		// 		}
		// 		log.Println(cn, string(buf[:cn]))
		// 	}
		// }(conn)
		// }
	}
}
