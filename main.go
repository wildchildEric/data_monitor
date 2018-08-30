package main

import (
	"data_monitor/dtuserver"
	"data_monitor/web"
	"log"
	"os"
)

func main() {
	startDtuServer()
	web.StartServe(8080)
}

func startDtuServer() {
	go func() {
		err := dtuserver.ServeAndListenTCP(6666)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}()
}

// func startClients() {
// 	for i := 0; i < 500; i++ {
// 		time.Sleep(1 * time.Millisecond)
// 		go func() {
// 			// dtuserver.Request("114.215.42.42:6666", fmt.Sprintf("Hello from client %d", i))
// 			dtuserver.Request("localhost:6666", fmt.Sprintf("Hello from client %d", i))
// 		}()
// 	}
// }
