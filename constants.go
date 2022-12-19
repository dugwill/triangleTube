package triangleTube

// Holding Registers
const (
	chDemand           = 512
	maxfriingRate      = 513
	chSetpoint         = 514
	ch1MaxSetpoint     = 1280
	dhwStorageSetpoint = 1281
)

// Input Registers
const (
	boilerStatus            = 0
	locketStatus            = 1
	boilerSupplyTemp        = 768
	boilerReturnTemp        = 769
	boilerFlueTemp          = 771
	outdoorTemp             = 772
	flameIonizationCurrent  = 774
	boilerCascadeFiringRate = 775
	boilerSetpoint          = 776
)

// Supported Commands
const (
	readHoldingRegisters   = 3
	readInputRegisters     = 4
	writeSingeRegister     = 6
	writeMultipleRegisters = 16
	reportSlaveID          = 17
)

// Modbus Configuration
const (
	baud       = 38400
	dataLength = 8
	parity     = "N"
	stopBits   = 1
)

// Status Bits
var manual = uint16(1)
var dhwMode = uint16(2)
var chMode = uint16(4)
var freezeMode = uint16(8)
var flame = uint16(16)
var chPump = uint16(32)
var dhwPump = uint16(64)
var system = uint16(128)
