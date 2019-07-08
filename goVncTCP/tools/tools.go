package tools

import (
	"net"
)

type NodeType string

const (
	NodesTable string   = "NODES LIST"
	Stem       NodeType = "0"
	Twig       NodeType = "1"
	NodePort   string   = "3333"
)

type Node struct {
	ADDRESS   string
	TYPE      NodeType
	PUBLICKEY string
}

type NodesList struct {
	NLIST []Node
}

type ServerConnection struct {
	NodeData      Node
	HelloReceived bool
}

func (conn *ServerConnection) CopyNode(node *Node) {
	//

	conn.NodeData = *node
}

func GetLocalIps() (addresses []string) {
	//

	ifaces, err := net.Interfaces()
	if err != nil {
		//

		return
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			//

			return
		}

		for _, addr := range addrs {
			switch v := addr.(type) {

			case *net.IPNet:
				addresses = append(addresses, v.IP.String())
			case *net.IPAddr:
				addresses = append(addresses, v.IP.String())
			}
		}
	}

	return
}

func NodeInNodes(node Node, nodes []Node) bool {
	for _, n := range nodes {
		//

		if node == n {
			//

			return true
		}
	}
	return false
}

func StringInSlice(val string, slice []string) bool {
	for _, v := range slice {
		//

		if val == v {
			//

			return true
		}
	}
	return false
}
