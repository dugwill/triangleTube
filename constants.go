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
	boilerStatus1           = 0
	locketStatus2           = 1
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

// Modbu Configuration
const (
	baud       = 38400
	dataLength = 8
	parity     = "n"
	stopBits   = 1
)
