package tcpclientserver

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
)

// Server structure
type Server struct {
	address                  string // Строка типа: localhost:12345
	config                   *tls.Config
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, message string)
}

// OnNewClient called when server start listen client
func (s *Server) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// OnClientConnectionClosed method closed connection with client
func (s *Server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// OnNewMessage method processing new message from client
func (s *Server) OnNewMessage(callback func(c *Client, message string)) {
	s.onNewMessage = callback
}

// Listen starts network server
func (s *Server) Listen() {
	var listener net.Listener
	var err error
	if s.config == nil {
		listener, err = net.Listen("tcp", s.address)
	} else {
		listener, err = tls.Listen("tcp", s.address, s.config)
	}
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		client := &Client{
			conn:   conn,
			Server: s,
		}
		go client.Listen()
	}
}

// NewServer  must create install of Server
func NewServer(address string) *Server {
	log.Println("Creating server with address", address)
	server := &Server{
		address: address,
		config:  nil,
	}

	server.OnNewClient(func(c *Client) {})
	server.OnNewMessage(func(c *Client, message string) {})
	server.OnClientConnectionClosed(func(c *Client, err error) {})

	return server
}

// NewWithTLS must create instans of Server with tls encription
func NewWithTLS(address string, certFile string, keyFile string) *Server {
	log.Println("Creating server with address", address)
	cert, _ := tls.LoadX509KeyPair(certFile, keyFile)
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	server := &Server{
		address: address,
		config:  &config,
	}

	server.OnNewClient(func(c *Client) {})
	server.OnNewMessage(func(c *Client, message string) {})
	server.OnClientConnectionClosed(func(c *Client, err error) {})

	return server
}

// Client structure
type Client struct {
	conn   net.Conn
	Server *Server
}

func (c *Client) Listen() {
	c.Server.onNewClientCallback(c)
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}
		c.Server.onNewMessage(c, message)
	}
}

// Send is send string message to server
func (c *Client) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	return err
}

// SendBytes is send bytes message to server
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

// Conn get connection
func (c *Client) Conn() net.Conn {
	return c.conn
}

// Close connection with server
func (c *Client) Close() error {
	return c.conn.Close()
}
