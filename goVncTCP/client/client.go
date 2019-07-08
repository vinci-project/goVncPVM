package client

import (
	"bufio"
	"goVncPVM/goVncTCP/tools"
	"log"
	"net"
	"time"
)

type Client struct {
	node          tools.Node
	outcomingData chan string
	reader        *bufio.Reader
	writer        *bufio.Writer
	connection    net.Conn
	errorChannel  *chan tools.Node
}

func NewClient(node tools.Node, connection net.Conn, errorChannel *chan tools.Node) *Client {
	//

	client := new(Client)
	client.node = node
	client.outcomingData = make(chan string)
	client.reader = bufio.NewReader(connection)
	client.writer = bufio.NewWriter(connection)
	client.connection = connection
	client.errorChannel = errorChannel

	return client
}

func (client *Client) write() {
	//

	for data := range client.outcomingData {
		//

		log.Println("SEND TO CLIENT: ", data)
		err := client.connection.SetWriteDeadline(time.Now().Add(3 * time.Second))
		if err != nil {
			//

			client.connection.Close()
			*client.errorChannel <- client.node
			return
		}

		_, err = client.writer.WriteString(data)
		if err != nil {
			//

			client.connection.Close()
			*client.errorChannel <- client.node
			return
		}

		err = client.writer.Flush()
		if err != nil {
			//

			client.connection.Close()
			*client.errorChannel <- client.node
			return
		}

		log.Println("SEND TO CLIENT SUCCESFULL")
	}
}

func (client *Client) Write(data string) {
	//

	client.outcomingData <- data
}

func (client *Client) Start() {
	//

	go client.write()
}

func (client *Client) CloseConnection() {
	//

	client.connection.Close()
}
