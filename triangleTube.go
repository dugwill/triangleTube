package triangleTube

import (
	"encoding/binary"
	"fmt"
	"log"
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
	BoilerSupplyTemp       uint16
	BoilerReturnTemp       uint16
	BoilerFlueTemp         uint16
	OutdoorTemp            uint16
	FlameIonizationCurrent uint16
	BoilerFiringRate       uint16
	BoilerSetpoint         uint16

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

	results, err := b.client.ReadDiscreteInputs(boilerSupplyTemp, 2)

	if err != nil {
		fmt.Println("Error reading Boiler temp")
		return
	}

	b.BoilerSupplyTemp = makeUint(results)

	return nil
}

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

func (b *TriangleTube) getBoilerSupplyTemp() uint16 {
	return b.BoilerSupplyTemp
}

func (b *TriangleTube) getBoilerReturnTemp() uint16 {
	return b.BoilerReturnTemp
}

func (b *TriangleTube) getOutdoorTemp() uint16 {
	return b.BoilerSupplyTemp
}

func makeBytes(u uint16) (b []byte) {
	binary.LittleEndian.PutUint16(b, uint16(u))
	return b
}

func makeUint(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)

}
