package main

import (
	"encoding/hex"
	"flag"
	"goVncPVM/goVncRest"
	"goVncPVM/goVncTCP"
	"goVncPVM/helpers"
	"log"
	"net"
	"time"

	"github.com/go-redis/redis"
	"github.com/tkanos/gonfig"
)

var redisDB *redis.Client
var tranChannel chan string

type RedisConfiguration struct {
	RedisHost  string
	RedisPort  string
	RedisDbNum int
}

func connectToRedis() (err error) {
	//

	redisConf := RedisConfiguration{}
	err = gonfig.GetConf("config/redis.json", &redisConf)
	if err != nil {
		//

		log.Fatalln("NO CONFIG FILE. ", err)
	}

	// var redisDBNumInt int = 1
	// redisHost, ok := os.LookupEnv("REDIS_PORT_6379_TCP_ADDR")
	// if !ok {
	// 	//
	//
	// 	redisHost = "0.0.0.0"
	// }
	//
	// redisPort, ok := os.LookupEnv("REDIS_PORT_6379_TCP_PORT")
	// if !ok {
	// 	//
	//
	// 	redisPort = "6379"
	// }
	//
	// redisDBNum, ok := os.LookupEnv("REDIS_PORT_6379_DB_NUM")
	// if !ok {
	// 	//
	//
	// 	redisDBNumInt = 1
	//
	// } else {
	// 	//
	//
	// 	if redisDBNumInt64, err := strconv.ParseInt(redisDBNum, 10, 64); err != nil {
	// 		//
	//
	// 		redisDBNumInt = int(redisDBNumInt64)
	//
	// 	} else {
	// 		//
	//
	// 		redisDBNumInt = 1
	// 	}
	// }

	redisDB = redis.NewClient(&redis.Options{
		Addr:         net.JoinHostPort(redisConf.RedisHost, redisConf.RedisPort),
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		DB:           redisConf.RedisDbNum,
	})

	statusCmd := redisDB.Ping()
	if helpers.IsRedisError(statusCmd) {
		//

		log.Fatalln("No connection to REDIS. ", statusCmd.Err())
		return statusCmd.Err()
	}

	return
}

func main() {
	//

	privateKeyPtr := flag.String("privateKey",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"A private key associated with your node")

	ipPtr := flag.String("ip",
		"0.0.0.0",
		"Your external IP")

	flag.Parse()

	privateKey, err := hex.DecodeString(*privateKeyPtr)
	if err != nil {
		//

		panic("WRONG PRIVATE KEY")
	}

	tranChannel = make(chan string, 1024)

	if err := connectToRedis(); err != nil {
		//

		panic(err.Error())
	}

	defer redisDB.Close()

	go restServer.Start(redisDB, &tranChannel, *ipPtr)
	tcpServer.Start(redisDB, &tranChannel, privateKey, *ipPtr)
	//
	// consoleReader := bufio.NewReader(os.Stdin)
	// log.Println("Enter text: ") // Just for better testing. Should be refactored for production
	// consoleReader.ReadString('\n')
	//
	// close(tranChannel)
}
