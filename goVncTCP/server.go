package tcpServer

import (
	"encoding/hex"
	"encoding/json"
	"goVncPVM/goVncTCP/client"
	"goVncPVM/goVncTCP/tools"
	"goVncPVM/helpers"
	"log"
	"net"
	"runtime"
	"time"

	"github.com/go-redis/redis"
	"github.com/tidwall/evio"
)

var redisDB *redis.Client
var clients map[tools.Node]*client.Client
var serverConnections map[string]*tools.ServerConnection
var readChannel chan string
var errorChannel chan tools.Node
var localAddresses []string
var tranChannel *chan string
var myNodeType tools.NodeType
var myPublicKey string
var myPrivateKey []byte
var myIp string

func getActiveNodes() (nodes []tools.Node, err error) {
	//

	stringCmd := redisDB.Get(tools.NodesTable)
	if helpers.IsRedisError(stringCmd) {
		//

		err = stringCmd.Err()
		return
	}

	var nodesList tools.NodesList

	err = json.Unmarshal([]byte(stringCmd.Val()), &nodesList)
	if err != nil {
		//

		log.Println("Can't marshal bytes. ", err)
		return
	}

	for _, node := range nodesList.NLIST {
		//

		if net.ParseIP(node.ADDRESS) == nil {
			//

			log.Println("Can't parse IP. ")
			continue
		}

		if len(node.PUBLICKEY) != 66 {
			//

			log.Println("Wrong publickey.")
			continue
		}

		if tools.StringInSlice(node.ADDRESS, localAddresses) {
			//

			log.Println("Our local IP.")
			if node.PUBLICKEY == myPublicKey {
				//

				myNodeType = node.TYPE

			} else {
				//

				log.Println("it's not my public key!!!")
				//panic("it's not my public key!!!")
			}

			continue
		}

		if node.PUBLICKEY == myPublicKey {
			// if we start go in docker, we can't rely on IP's

			log.Println("I've found myself! PK: ", myPublicKey)
			continue
		}

		nodes = append(nodes, node)
	}

	return
}

func handleConnection(c evio.Conn, in []byte) (out []byte, action evio.Action) {
	//

	if myNodeType == tools.Stem {
		//

		return
	}

	ip, _, _ := net.SplitHostPort(c.RemoteAddr().String())
	data := string(in)
	if node, ok := serverConnections[ip]; ok {
		//

		if node.NodeData.TYPE == tools.Stem {
			// We do not work with stem node

			return
		}

		tranType := helpers.GetRawTransactionType(data)
		if tranType == "ST" {
			//

			if node.HelloReceived {
				//

				readChannel <- data
			}

		} else if tranType == "HL" {
			//

			readChannel <- data
		}

	} else {
		//

		log.Println("Sorry, I don't know you, ", ip)
	}

	return
}

func createConnectionsWithNodes() (delay time.Duration, action evio.Action) {
	//

	log.Println("tick", time.Now().Unix())
	delay = 10 * time.Second
	nodes, _ := getActiveNodes()
	log.Println("NODES FROM SERVER: ", nodes)
	for node, client := range clients {
		//

		if !tools.NodeInNodes(node, nodes) {
			//

			log.Println("We don't need this node anymore")
			client.CloseConnection()
			delete(clients, node)
		}
	}

	for ip, connection := range serverConnections {
		//

		if !tools.NodeInNodes(connection.NodeData, nodes) {
			//

			log.Println("We don't need this connection anymore")
			delete(serverConnections, ip)
		}
	}

	for _, node := range nodes {
		//

		if _, ok := clients[node]; ok {
			//

			log.Println("Such node already exists")
			continue
		}

		if _, ok := serverConnections[node.ADDRESS]; !ok {
			//

			serverConnection := new(tools.ServerConnection)
			serverConnection.CopyNode(&node)
			serverConnections[node.ADDRESS] = serverConnection

			serverConnection.HelloReceived = true

			continue
		}

		if node.TYPE == tools.Stem {
			// We do not work with STEM

			continue
		}

		servAddr := net.JoinHostPort(node.ADDRESS, tools.NodePort)

		conn, err := net.DialTimeout("tcp", servAddr, 3*time.Second)
		if err != nil {
			//

			log.Println("Dial failed:", err.Error())
			continue
		}

		// log.Println(conn.LocalAddr(), conn.RemoteAddr(), err)
		// err = conn.SetKeepAlive(true)
		// if err != nil {
		// 	//
		//
		// 	log.Println("Set keep alive failed: ", err.Error())
		// 	continue
		// }
		//
		// err = conn.SetKeepAlivePeriod(3 * time.Second)
		// if err != nil {
		// 	//
		//
		// 	log.Println("Set keep alive timeout failed: ", err.Error())
		// 	continue
		// }

		log.Println("New connection")
		newClient := client.NewClient(node, conn, &errorChannel)
		newClient.Start()

		// Send Hello transaction with signature
		if helloTransaction, ok := helpers.CreateHelloTransaction(myPublicKey, myPrivateKey, myIp); ok {
			//

			newClient.Write(helloTransaction)
			clients[node] = newClient

		} else {
			//

			newClient.CloseConnection()
		}
	}

	return
}

func closeConnections() {
	//

	for _, client := range clients {
		//

		client.CloseConnection()
	}
}

func startServer() (err error) {
	//

	var events evio.Events
	events.NumLoops = runtime.NumCPU()
	events.Data = handleConnection
	events.Tick = createConnectionsWithNodes

	return evio.Serve(events, "tcp://"+net.JoinHostPort("0.0.0.0", tools.NodePort))
}

func readWorker() {
	// We get data, checking it validity, signature and resending to everyone
	// Also, here we chek for client errors

	for {
		//

		select {
		case data := <-readChannel:
			//

			log.Println("NEW DATA")
			tranType := helpers.GetRawTransactionType(data)
			if tranType == "HL" {
				// HELLO transaction

				if helloTransaction, err := helpers.ParseHelloTransaction(data); err == nil {
					//

					if ok := helpers.VerifyHelloTransaction(helloTransaction); !ok {
						//

						if serverConnection, ok := serverConnections[helloTransaction.ADDRESS]; ok {
							//

							serverConnection.HelloReceived = true
							serverConnections[helloTransaction.ADDRESS] = serverConnection
						}
					}
				}

			} else if tranType == "ST" {
				// Simple transaction

				if simpleTransaction, err := helpers.ParseSimpleTransaction(data); err == nil {
					//

					if transactionTime, _, ok := helpers.VerifySimpleTransaction(simpleTransaction); ok {
						//

						log.Println("RECEIVED FROM TCP SERVER: ", data)
						errRedis := redisDB.ZAdd("RAW TRANSACTIONS", redis.Z{
							Score:  float64(transactionTime),
							Member: data,
						})

						if helpers.IsRedisError(errRedis) {
							//

							log.Println(errRedis.Err())
						}
					}
				}

			} else if tranType == "AT" {
				// Applicant transaction

				if simpleTransaction, err := helpers.ParseApplicantTransaction(data); err == nil {
					//

					if transactionTime, _, ok := helpers.VerifyApplicantTransaction(simpleTransaction); ok {
						//

						log.Println("RECEIVED FROM TCP SERVER: ", data)
						errRedis := redisDB.ZAdd("RAW TRANSACTIONS", redis.Z{
							Score:  float64(transactionTime),
							Member: data,
						})

						if helpers.IsRedisError(errRedis) {
							//

							log.Println(errRedis.Err())
						}
					}
				}

			} else if tranType == "VT" {
				// Vote transaction

				if simpleTransaction, err := helpers.ParseVoteTransaction(data); err == nil {
					//

					if transactionTime, _, ok := helpers.VerifyVoteTransaction(simpleTransaction); ok {
						//

						log.Println("RECEIVED FROM TCP SERVER: ", data)
						errRedis := redisDB.ZAdd("RAW TRANSACTIONS", redis.Z{
							Score:  float64(transactionTime),
							Member: data,
						})

						if helpers.IsRedisError(errRedis) {
							//

							log.Println(errRedis.Err())
						}
					}
				}

			} else if tranType == "UAT" {
				// Unregister Apllicant transaction

				if simpleTransaction, err := helpers.ParseUATransaction(data); err == nil {
					//

					if transactionTime, _, ok := helpers.VerifyUATransaction(simpleTransaction); ok {
						//

						log.Println("RECEIVED FROM TCP SERVER: ", data)
						errRedis := redisDB.ZAdd("RAW TRANSACTIONS", redis.Z{
							Score:  float64(transactionTime),
							Member: data,
						})

						if helpers.IsRedisError(errRedis) {
							//

							log.Println(errRedis.Err())
						}
					}
				}

			} else if tranType == "UVT" {
				// Unregister Vote transaction

				if simpleTransaction, err := helpers.ParseUVTransaction(data); err == nil {
					//

					if transactionTime, _, ok := helpers.VerifyUVTransaction(simpleTransaction); ok {
						//

						log.Println("RECEIVED FROM TCP SERVER: ", data)
						errRedis := redisDB.ZAdd("RAW TRANSACTIONS", redis.Z{
							Score:  float64(transactionTime),
							Member: data,
						})

						if helpers.IsRedisError(errRedis) {
							//

							log.Println(errRedis.Err())
						}
					}
				}
			}

		case node := <-errorChannel:
			//

			log.Println("ERROR: ", node.ADDRESS)
			delete(clients, node)

		case tran := <-*tranChannel:
			//

			for _, client := range clients {
				//

				log.Println("RECEIVED FROM REST: ", tran)
				client.Write(tran)
			}
		}
	}
}

func Start(r *redis.Client, c *chan string, privateKey []byte, ip string) {
	//

	redisDB = r
	tranChannel = c
	myPrivateKey = privateKey
	myPublicKey = hex.EncodeToString(helpers.PubkeyFromSeckey(myPrivateKey))
	myIp = ip
	clients = make(map[tools.Node]*client.Client)
	serverConnections = make(map[string]*tools.ServerConnection)
	readChannel = make(chan string, 1024)
	errorChannel = make(chan tools.Node)
	localAddresses = tools.GetLocalIps()

	go readWorker()

	if err := startServer(); err != nil {
		//

		panic(err.Error())
	}

	close(readChannel)
	close(errorChannel)
}
