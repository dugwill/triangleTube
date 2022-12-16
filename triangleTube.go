package triangleTube

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/dugwill/modbus-1"
)

type TriangleTube struct {
	// Modbus
	modBusID int
	client   modbus.Client

	// Input Registers
	BoilerStatus           int8
	LockoutStatus1         uint16
	LockoutStatus2         uint16
	BoilerSupplyTemp       float32
	BoilerReturnTemp       float32
	BoilerFlueTemp         float32
	OutdoorTemp            float32
	FlameIonizationCurrent uint16
	BoilerFiringRate       uint16
	BoilerSetpoint         float32

	// Holding registers
	ChDemand           uint16
	MaxFiringRate      uint16
	ChSetpoint         uint16
	Ch1MaxSetpoint     uint16
	DhwStorageSetpoint uint16
}

func NewBoiler(ID int) (b *TriangleTube, err error) {
	b = new(TriangleTube)
	// Set initial status
	b.modBusID = ID

	// Modbus RTU/ASCII
	//handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
	handler := modbus.NewRTUClientHandler("com5")
	handler.BaudRate = baud
	handler.DataBits = dataLength
	handler.Parity = parity
	handler.StopBits = stopBits
	handler.SlaveId = byte(b.modBusID)
	handler.Timeout = 5 * time.Second

	if err = handler.Connect(); err != nil {
		err = fmt.Errorf("error opening boiler port: %v", err)
		return
	}

	//defer handler.Close()

	b.client = modbus.NewClient(handler)

	return
}

func (b *TriangleTube) Update() (err error) {

	results, err := b.client.ReadInputRegisters(boilerSupplyTemp, 1)

	if err != nil {
		fmt.Println("Error reading Boiler temp")
		return
	}

	b.BoilerSupplyTemp = makeTemp(results)

	return nil
}

/*
func (b *TriangleTube) ProcessCommand(c []byte) (r []byte, err error) {

	// validate command input
	if len(c) < 3 {
		return nil, fmt.Errorf("malformed command")
	}

	command := c[0]
	registers := c[1:]
	// if there are an odd number of byte, throw error
	if len(registers)%2 != 0 {
		return nil, fmt.Errorf("incorrect register request")
	}
	var registerValue uint16
	for i := 0; i < len(registers)/2; i += 2 {
		reg := registers[i : i+2]
		// Process command
		switch command {
		case readInputRegisters:
			registerValue, err = b.readInputRegister(makeUint(reg))
			log.Println("Reg Value", registerValue)

		case readHoldingRegisters:

		default:
			return nil, fmt.Errorf("command value not supported")
		}
	}
	return

}


func (b *TriangleTube) readInputRegister(r uint16) (v uint16, err error) {
	switch r {
	case boilerSupplyTemp:
		return b.getBoilerSupplyTemp(), nil
	case boilerReturnTemp:
		return b.getBoilerReturnTemp(), nil
	case outdoorTemp:
		return b.getOutdoorTemp(), nil
	default:
		err = fmt.Errorf("register value not available")
		return 0, err
	}
}
*/

func (b *TriangleTube) getBoilerSupplyTemp() float32 {
	return b.BoilerSupplyTemp
}

func (b *TriangleTube) getBoilerReturnTemp() float32 {
	return b.BoilerReturnTemp
}

func (b *TriangleTube) getOutdoorTemp() float32 {
	return b.BoilerSupplyTemp
}

func makeBytes(u uint16) (b []byte) {
	binary.LittleEndian.PutUint16(b, uint16(u))
	return b
}

func makeUint(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func CtoF(c uint16) float32 {
	return float32((c)*9/5) + 32
}

func makeTemp(b []byte) float32 {
	celcius := makeUint(b) / 10
	return CtoF(celcius)
}
