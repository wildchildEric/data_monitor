//package main
//author: Lubia Yang
//create: 2013-10-15
//about: www.lubia.me

package modbusrtu

import (
	"errors"
	"os"
	"time"
)

// Read ...
//   parameters
//   int        fd:  file descripter for serial device
//   byte  addr:  slave device address
//   byte  code:  function code
//   uint16 sr:    starting register number
//   uint16 nr:    number of registers to read
//   byte data[]: memory area for read data
func Read(fd *os.File, addr, code byte, sr, nr uint16) ([]byte, error) {
	//Preparation for Sending a Packet
	var sendPacket = BuildSendPacketForRead(addr, code, sr, nr)
	// Preparation for Receiving a Packet
	var recvPacket = make([]byte, 256)
	_, err := fd.Write(sendPacket)
	if err != nil {
		return []byte{}, errors.New("MODBUS_ERROR_COMMUNICATION")
	}
	time.Sleep(300 * time.Millisecond)
	_, err = fd.Read(recvPacket)
	if err != nil {
		return []byte{}, errors.New("MODBUS_ERROR_COMMUNICATION")
	}
	return ParseRecvPacketForRead(sendPacket, recvPacket)
}

/*BuildSendPacketForRead ...*/
func BuildSendPacketForRead(addr byte, funCode byte, startRegister uint16, numRegister uint16) []byte {
	sendPacket := make([]byte, 8)
	//Packet Construction
	sendPacket[0] = addr                       // Slave Address
	sendPacket[1] = funCode                    // Function Code 0x03 = Multiple Read
	sendPacket[2] = byte(startRegister >> 8)   // Start Register (High Byte)
	sendPacket[3] = byte(startRegister & 0xff) // Start Register (Low Byte)
	sendPacket[4] = byte(numRegister >> 8)     // Number of Registers (High Byte)
	sendPacket[5] = byte(numRegister & 0xff)   // Number of Registers (Low Byte)
	//Add CRC16
	sendPacketCrc := Crc(sendPacket[:6])
	sendPacket[6] = byte(sendPacketCrc & 0xff)
	sendPacket[7] = byte(sendPacketCrc >> 8)
	return sendPacket
}

/*ParseRecvPacketForRead ...*/
func ParseRecvPacketForRead(sentPacket []byte, recvPacket []byte) ([]byte, error) {
	// Parse the Response
	if recvPacket[0] != sentPacket[0] || recvPacket[1] != sentPacket[1] {
		if recvPacket[0] == sentPacket[0] && recvPacket[1]&0x7f == sentPacket[1] {
			switch recvPacket[2] {
			case 1:
				return []byte{}, errors.New("MODBUS_ERROR_COMMUNICATION_ILLEGAL_FUNCTION")
			case 2:
				return []byte{}, errors.New("MODBUS_ERROR_COMMUNICATION_ILLEGAL_ADDRESS")
			case 3:
				return []byte{}, errors.New("MODBUS_ERROR_COMMUNICATION_ILLEGAL_VALUE")
			case 4:
				return []byte{}, errors.New("MODBUS_ERROR_COMMUNICATION_ILLEGAL_OPERATION")
			}
		}
		return []byte{}, errors.New("MODBUS_ERROR_COMMUNICATION")
	}
	//CRC check
	l := recvPacket[2]
	recvPacketCrc := Crc(recvPacket[:3+l])
	if recvPacket[3+l] != byte((recvPacketCrc&0xff)) || recvPacket[3+l+1] != byte((recvPacketCrc>>8)) {
		return []byte{}, errors.New("MODBUS_ERROR_COMMUNICATION")
	}
	return recvPacket[3 : l+3], nil
}

// Write ...
//   parameters
//   int        fd:  file descripter for serial device
//   byte  addr:  slave device address
//   byte  code:  function code
//   uint16 sr:    starting register number
//   uint16 nr:    number of registers to write
//   byte data[]: memory area for writing data
func Write(fd *os.File, addr, code byte, sr, nr uint16, data []byte) error {
	var sendPacket = BuildSendPacketForWrite(addr, code, sr, nr, data)
	// Preparation for Receiving a Packet
	var recvPacket = make([]byte, 256)
	_, err := fd.Write(sendPacket)
	if err != nil {
		return errors.New("MODBUS_ERROR_COMMUNICATION")
	}
	time.Sleep(300 * time.Millisecond)
	_, err = fd.Read(recvPacket)
	if err != nil {
		return errors.New("MODBUS_ERROR_COMMUNICATION")
	}
	return ParseRecvPacketForWrite(sendPacket, recvPacket)
}

/*BuildSendPacketForWrite ...*/
func BuildSendPacketForWrite(addr byte, code byte, sr uint16, nr uint16, data []byte) []byte {
	var sendPacket = make([]byte, 256)
	// Packet Construction
	sendPacket[0] = addr            // Slave Address
	sendPacket[1] = code            // Function Code 0x10 = Multiple Write
	sendPacket[2] = byte(sr >> 8)   // Start Register (High Byte)
	sendPacket[3] = byte(sr & 0xff) // Start Register (Low Byte)
	sendPacket[4] = byte(nr >> 8)   // Number of Registers (High Byte)
	sendPacket[5] = byte(nr & 0xff) // Number of Registers (Low Byte)
	sendPacket[6] = byte(nr * 2)
	for i := 0; i < int((nr * 2)); i++ {
		sendPacket[7+i] = data[i]
	}
	length := 7 + nr*2 + 2
	// Add CRC16
	sendPacketCrc := Crc(sendPacket[:length-2])
	sendPacket[length-2] = byte(sendPacketCrc & 0xff)
	sendPacket[length-1] = byte(sendPacketCrc >> 8)
	return sendPacket
}

/*ParseRecvPacketForWrite ...*/
func ParseRecvPacketForWrite(sentPacket []byte, recvPacket []byte) error {
	// Parse the Response
	if recvPacket[0] != sentPacket[0] || recvPacket[1] != sentPacket[1] {
		if recvPacket[0] == sentPacket[0] && recvPacket[1]&0x7f == sentPacket[1] {
			switch recvPacket[2] {
			case 1:
				return errors.New("MODBUS_ERROR_COMMUNICATION_ILLEGAL_FUNCTION")
			case 2:
				return errors.New("MODBUS_ERROR_COMMUNICATION_ILLEGAL_ADDRESS")
			case 3:
				return errors.New("MODBUS_ERROR_COMMUNICATION_ILLEGAL_VALUE")
			case 4:
				return errors.New("MODBUS_ERROR_COMMUNICATION_ILLEGAL_OPERATION")
			}
		}
		return errors.New("MODBUS_ERROR_COMMUNICATION")
	}
	//Target Data Filed Check
	if recvPacket[2] == sentPacket[2] &&
		recvPacket[3] == sentPacket[3] &&
		recvPacket[4] == sentPacket[4] &&
		recvPacket[5] == sentPacket[5] {
		//CRC check
		recvPacketCrc := Crc(recvPacket[:6])
		if recvPacket[6] == byte((recvPacketCrc&0xff)) && recvPacket[7] == byte((recvPacketCrc>>8)) {
			return nil
		}
	}
	return errors.New("MODBUS_ERROR_COMMUNICATION")
}

/*Crc ...*/
func Crc(data []byte) uint16 {
	var crc16 uint16 = 0xffff
	l := len(data)
	for i := 0; i < l; i++ {
		crc16 ^= uint16(data[i])
		for j := 0; j < 8; j++ {
			if crc16&0x0001 > 0 {
				crc16 = (crc16 >> 1) ^ 0xA001
			} else {
				crc16 >>= 1
			}
		}
	}
	return crc16
}
