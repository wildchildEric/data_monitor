package dtuserver

import (
	"log"
	"net"
	"time"
)

/*Request ... */
func Request(address string, message string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	conn.Write([]byte(message + "\n"))
	readAndLog(conn)

	for i := 0; i < 60; i++ {
		time.Sleep(10 * time.Second)
		conn.Write([]byte(message + "\n"))
		readAndLog(conn)
	}
}

func readAndLog(conn net.Conn) {
	buffer := make([]byte, 256)
	conn.Read(buffer)
	log.Print(string(buffer))
}
