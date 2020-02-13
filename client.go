package tcpclientserver

import (
	"bufio"
	"net"
)

// Client structure
type Client struct {
	conn   net.Conn
	Server *Server
}

func (c *Client) listen() {
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
