package triangleTube

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dugwill/modbus-1"
)

type TriangleTube struct {
	// Modbus
	ModBusID int           `json:"modbusId"`
	Client   modbus.Client `json:"-"`
	SlaveID  uint16        `json:"-"`
	ReadTime string        `json:"time"`

	// Boiler Status
	PcManualMode         bool `json:"pcManualMode"`
	DhwMode              bool `json:"dhwMode"`
	ChMode               bool `json:"chMode"`
	FreezeProtectionMode bool `json:"freezeProtMode"`
	FlamePresent         bool `json:"flamePresent"`
	Ch1Pump              bool `json:"ch1Pump"`
	DhwPump              bool `json:"dhwPump"`
	System               bool `json:"system"`

	// Input Registers
	BoilerStatus           int8    `json:"-"`
	LockoutStatus1         uint16  `json:"-"`
	LockoutStatus2         uint16  `json:"-"`
	BoilerSupplyTemp       float32 `json:"supplyTemp"`
	BoilerReturnTemp       float32 `json:"returnTemp"`
	DHWTemp                float32 `json:"dhwTemp"`
	BoilerFlueTemp         float32 `json:"flueTemp"`
	OutdoorTemp            float32 `json:"outdoorTemp"`
	FlameIonizationCurrent uint16  `json:"ionizationCurrent"`
	BoilerFiringRate       uint16  `json:"firingRate"`
	BoilerSetpoint         float32 `json:"setPoint"`

	// Holding registers
	ChDemand           uint16 `json:"chDemand"`
	MaxFiringRate      uint16 `json:"maxFiringRate"`
	ChSetpoint         uint16 `json:"chSetPoint"`
	Ch1MaxSetpoint     uint16 `json:"ch1MaxSetPoint"`
	DhwStorageSetpoint uint16 `json:"dhwStorageSetPoint"`

	// Zones
	ZoneMap map[string]*zone `json:"zoneMap"`
}

func NewBoiler(ID int, comPort string) (b *TriangleTube, err error) {
	b = new(TriangleTube)
	// Set initial status
	b.ModBusID = ID

	// Modbus RTU/ASCII
	//handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
	handler := modbus.NewRTUClientHandler(comPort)
	handler.BaudRate = baud
	handler.DataBits = dataLength
	handler.Parity = parity
	handler.StopBits = stopBits
	handler.SlaveId = byte(b.ModBusID)
	handler.Timeout = 5 * time.Second

	if err = handler.Connect(); err != nil {
		err = fmt.Errorf("error opening boiler port: %v", err)
		return
	}

	b.Client = modbus.NewClient(handler)

	initGpio()

	return
}

func (b *TriangleTube) Update() (err error) {

	var results []byte

	//Read Status Byte
	_ = b.ReadBoilerStatus()

	//Read lockout status

	// Read Temps, FI Current, Fire rate and Setpoint
	results, err = b.Client.ReadInputRegisters(boilerSupplyTemp, 9)
	if err != nil {
		fmt.Println("Error reading Boiler temps")
		return
	}
	fmt.Printf("MultiRead Results: %x\n", results)

	//fmt.Printf("Boiler Supply Temp Raw: %x\n", results[0:2])
	b.BoilerSupplyTemp = makeTemp(results[0:2])
	//fmt.Printf("Boiler  ReturnTemp Raw: %x\n", results[2:4])
	b.BoilerReturnTemp = makeTemp2(results[2:4])
	//fmt.Printf("Hot Water Temp Raw: %x\n", results[4:6])
	b.DHWTemp = makeTemp2(results[4:6])
	//fmt.Printf("Boiler Flue Temp Raw: %x\n", results[6:8])
	b.BoilerFlueTemp = makeTemp2(results[6:8])
	//fmt.Printf("Outdoor Temp Raw: %x\n", results[8:10])
	b.OutdoorTemp = makeTemp2(results[8:10])
	//Skip bytes 10 and 11 byte pair for future use
	b.FlameIonizationCurrent = makeUint(results[12:14])
	b.BoilerFiringRate = makeUint(results[14:16])
	b.BoilerSetpoint = makeTemp2(results[16:18])

	b.ReadTime = time.Now().String()

	return nil
}

func (b *TriangleTube) ReportSlaveID() (err error) {
	results, err := b.Client.ReadInputRegisters(reportSlaveID, 1)
	if err != nil {
		fmt.Println("Error reading Boiler temp")
		return
	}

	b.SlaveID = makeUint(results)

	return nil
}

func (b *TriangleTube) GetBoilerSupplyTemp() float32 {
	return b.BoilerSupplyTemp
}

func (b *TriangleTube) GetBoilerReturnTemp() float32 {
	return b.BoilerReturnTemp
}

func (b *TriangleTube) GetDHWTemp() float32 {
	return b.DHWTemp
}

func (b *TriangleTube) GetBoilerFlueTemp() float32 {
	return b.BoilerFlueTemp
}

func (b *TriangleTube) GetOutdoorTemp() float32 {
	return b.OutdoorTemp
}

func (b *TriangleTube) GetFlameIonizationCurrent() uint16 {
	return b.FlameIonizationCurrent
}

func (b *TriangleTube) GetBoilerFiringRate() uint16 {
	return b.BoilerFiringRate
}

func (b *TriangleTube) GetBoilerSetPoint() float32 {
	return b.BoilerSetpoint
}

/**
func makeBytes(u uint16) (b []byte) {
	binary.LittleEndian.PutUint16(b, uint16(u))
	return b
}
**/

func makeUint(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func CtoF(c uint16) float32 {
	// Convert from C to F (cnvt c to signed)
	return float32((int16(c))*9/5) + 32
}

func makeTemp(b []byte) float32 {
	celcius := makeUint(b) / 10
	return CtoF(celcius)
}

func makeTemp2(b []byte) float32 {
	celcius := makeUint(b)
	return CtoF(celcius)
}

func (b *TriangleTube) ReadBoilerStatus() error {
	results, err := b.Client.ReadInputRegisters(boilerStatus, 1)
	if err != nil {
		err = fmt.Errorf("error reading boiler status: %v", err)
		return err
	}
	//PC Manual Mode status
	fmt.Printf("Boiler Status: %x\n", results)
	if makeUint(results)&manual == manual {
		b.PcManualMode = true
	} else {
		b.DhwMode = false
	}
	// DHW Mode Status
	if makeUint(results)&dhwMode == dhwMode {
		b.DhwMode = true
	} else {
		b.DhwMode = false
	}
	// CH Mode Status
	if makeUint(results)&chMode == chMode {
		b.ChMode = true
	} else {
		b.ChMode = false
	}
	// FreezeProtection Mode Status
	if makeUint(results)&freezeMode == freezeMode {
		b.FreezeProtectionMode = true
	} else {
		b.FreezeProtectionMode = false
	}
	// Flame Present Status
	if makeUint(results)&flame == flame {
		b.FlamePresent = true
	} else {
		b.FlamePresent = false
	}
	// CH Pump Status
	if makeUint(results)&chPump == chPump {
		b.Ch1Pump = true
	} else {
		b.Ch1Pump = false
	}
	// Systme Status
	if makeUint(results)&system == system {
		b.System = true
	} else {
		b.System = false
	}

	return nil
}

func (b *TriangleTube) PrintStatus() {

	fmt.Printf("Manual Mode: %v\n", b.PcManualMode)
	fmt.Printf("DHW Mode: %v\n", b.DhwMode)
	fmt.Printf("CH Mode: %v\n", b.ChMode)
	fmt.Printf("Freeze Protection Mode: %v\n", b.FreezeProtectionMode)
	fmt.Printf("Flame Present: %v\n", b.FlamePresent)
	fmt.Printf("CH Pump On: %v\n", b.Ch1Pump)
	fmt.Printf("DHW Pump On: %v\n", b.DhwPump)
	if b.System {
		fmt.Printf("System: Running\n")
	} else {
		fmt.Printf("System: Standby\n")
	}
	fmt.Printf("Boiler Supply Temp: %f\n", b.GetBoilerSupplyTemp())
	fmt.Printf("Boiler Return Temp: %v\n", b.GetBoilerReturnTemp())
	fmt.Printf("DHW Temp: %v\n", b.GetDHWTemp())
	fmt.Printf("Boiler Flue Temp: %v\n", b.GetBoilerFlueTemp())
	fmt.Printf("Outdoor Temp: %v\n", b.GetOutdoorTemp())
	fmt.Printf("Flame Ionization Current: %v uA\n", b.GetFlameIonizationCurrent())
	fmt.Printf("Boiler Firing Rate: %v\n", b.GetBoilerFiringRate())
	fmt.Printf("Boiler Set Point: %v\n", b.GetBoilerSetPoint())
}

func (b *TriangleTube) PrintJsonIndent() {
	var boilerJson []byte
	var err error

	if boilerJson, err = json.MarshalIndent(b, "  ", "  "); err != nil {
		fmt.Println("Error marshaling boiler data: ", err)
	}

	fmt.Printf("Boiler Data: \n%s\n", boilerJson)
}

func (b *TriangleTube) PrintJson() {
	var boilerJson []byte
	var err error

	if boilerJson, err = json.Marshal(b); err != nil {
		fmt.Println("Error marshaling boiler data: ", err)
	}

	fmt.Printf("Boiler Data: \n%s\n", boilerJson)
}
