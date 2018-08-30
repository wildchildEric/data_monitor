package dtuserver

import (
	"bufio"
	"data_monitor/dtus"
	"data_monitor/modbusrtu"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	/*MDtu ... */
	MDtu = dtus.NewManagedDtu()
)

/*ServeAndListenTCP start the tcp server
智能仪表 <--modbus rtu--> PLC <--modbus rtu---> DTU <---TCP IP---> DTU TCP Server
*/
func ServeAndListenTCP(port int) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	log.Printf("Dtu Tcp Server Started at %d", port)
	scheduleClean(90 * time.Second)
	scheduleWriteModbusCommand(60 * time.Second)
	ch := handleNewClient(listener)
	readClient(ch)
	return nil
}

func handleNewClient(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		defer func() {
			if p := recover(); p != nil {
				log.Printf("internal error: %v", p)
			}
		}()
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			i++
			log.Printf("%d: (%v) <-> (%v)\n", i, conn.LocalAddr(), conn.RemoteAddr())
			ch <- conn
		}
	}()
	return ch
}

func readClient(ch chan net.Conn) {
	for {
		conn := <-ch
		go func() {
			defer func() {
				if p := recover(); p != nil {
					log.Printf("internal error: %v", p)
				}
			}()
			r := bufio.NewReader(conn)
			for {
				buffer := make([]byte, 512)
				n, err := r.Read(buffer)
				if n > 0 {
					onNewMessage(conn, buffer[:n])
				}
				if err != nil { // EOF, or worse
					log.Println(err)
					MDtu.RemoveDtu(conn)
					break
				}
			}
		}()
	}
}

func onNewMessage(conn net.Conn, message []byte) {
	str := string(message)
	log.Printf("Server received []byte: %v \n", message)
	log.Printf("            raw string: %q \n", str)
	log.Printf("                string: %s \n", str)
	log.Printf("                length: %d \n", len(message))
	if len(message) == 24 && strings.Index(str, "reg") == 0 {
		//Dtu注册信息格式： "reg_xxxxxxxx_xxxxxxxxxxx"
		log.Println("Dtu registered!")
		tmp := strings.Split(str, "_")
		if len(tmp) != 3 || !isValidDtuID(tmp[1]) {
			log.Println("Invalid Connection")
			conn.Close()
			return
		}
		dtu := dtus.NewDtu()
		dtu.SetConn(conn)
		dtu.SetID(tmp[1])
		dtu.SetSimNum(tmp[2])
		MDtu.RegisterDtu(dtu)
	} else {
		dtu, ok := MDtu.RetriveDtu(conn)
		if !ok { //Make Sure refuse invalid connection
			conn.Close()
			log.Println("Closed unregistered connection")
			return
		}
		if len(message) == 1 && str == "!" {
			//Dtu 心跳字符： ！
			dtu.SetLastHeartbeatedTime(time.Now())
			log.Printf("%s heartbeated!\n", dtu)
		} else if len(message) <= 256 && len(dtu.LastSent()) > 0 {
			//Modbus response
			dtu.SetLastReceived(message)
			values, err := modbusrtu.ParseRecvPacketForRead(dtu.LastSent(), message)
			if err != nil {
				log.Println(err)
				return
			}
			log.Print("Modbus response data: ", values)
			//TODO: handle the modbus response data
			// var mySlice = []byte{8, 152}
			// data := binary.BigEndian.Uint16(mySlice)
		} else {
			log.Println("Unknow data")
		}
	}
}

func isValidDtuID(id string) bool {
	//TODO: add valid id check logic
	return true
}

func scheduleClean(d time.Duration) {
	ticker := time.NewTicker(d)
	go func() {
		for t := range ticker.C {
			log.Printf("scheduleCleanConn Tick at %s managedDtu Size: %d\n", t, MDtu.Size())
			MDtu.Clean(2 * time.Minute)
			log.Printf("After clean managedDtu Size: %d\n", MDtu.Size())
		}
	}()
}

func scheduleWriteModbusCommand(d time.Duration) {
	ticker := time.NewTicker(d)
	go func() {
		for t := range ticker.C {
			log.Println("scheduleWriteModbusCommand Tick at", t)
			packet := modbusrtu.BuildSendPacketForRead(0x01, 0x03, 0x01, 0x06)
			err := MDtu.Send(packet)
			if err != nil {
				log.Println(err)
			}
		}
	}()
}
