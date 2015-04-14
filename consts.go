package main

// gphoto2 errors
const (
	//gpOk means no error
	gpOk = 0
	//Error is a Generic Error
	Error = -1
	//ErrorBadParameters : Bad parameters passed
	ErrorBadParameters = -2
	//ErrorNoMemory : Out of memory
	ErrorNoMemory = -3
	//ErrorLibrary : Error in the camera driver
	ErrorLibrary = -4
	//ErrorUnknownPort : Unknown libgphoto2 port passed
	ErrorUnknownPort = -5
	//ErrorNotSupported : Functionality not supported
	ErrorNotSupported = -6
	//ErrorIO : Generic I/O error
	ErrorIO = -7
	//ErrorFixedLimitExceeded : Buffer overflow of internal structure
	ErrorFixedLimitExceeded = -8
	//ErrorTimeout : Operation timed out
	ErrorTimeout = -10
	//ErrorIOSupportedSerial : Serial ports not supported
	ErrorIOSupportedSerial = -20
	//ErrorIOSupportedUsb : USB ports not supported
	ErrorIOSupportedUsb = -21
	//ErrorIOInit : Error initialising I/O
	ErrorIOInit = -31
	//ErrorIORead : I/O during read
	ErrorIORead = -34
	//ErrorIOWrite : I/O during write
	ErrorIOWrite = -35
	//ErrorIOUpdate : I/O during update of settings
	ErrorIOUpdate = -37
	//ErrorIOSerialSpeed : Specified serial speed not possible.
	ErrorIOSerialSpeed = -41
	//ErrorIOUSBClearHalt : Error during USB Clear HALT
	ErrorIOUSBClearHalt = -51
	//ErrorIOUSBFind : Error when trying to find USB device
	ErrorIOUSBFind = -52
	//ErrorIOUSBClaim : Error when trying to claim the USB device
	ErrorIOUSBClaim = -53
	//ErrorIOLock : Error when trying to lock the device
	ErrorIOLock = -60
	//ErrorHal : Unspecified error when talking to HAL
	ErrorHal = -70
)

//widget types
const (
	//WidgetWindow is the toplevel configuration widget. It should likely contain multiple #WidgetSection entries.
	WidgetWindow = iota //(0)
	//WidgetSection : Section widget (think Tab)
	WidgetSection
	//WidgetText : Text widget (string)
	WidgetText
	//WidgetRange : Slider widget (float)
	WidgetRange
	//WidgetToggle : Toggle widget (think check box) (int)
	WidgetToggle
	//WidgetRadio : Radio button widget (string)
	WidgetRadio
	//WidgetMenu : Menu widget (same as RADIO) (string)
	WidgetMenu
	//WidgetButton : Button press widget ( CameraWidgetCallback )
	WidgetButton
	//WidgetDate : Date entering widget (int)
	WidgetDate
)

//Log level
const (
	//LogError : Log message is an error infomation
	LogError = iota
	//LogVerbose : Log message is an verbose debug infomation
	LogVerbose
	//LogDebug : Log message is an debug infomation
	LogDebug
	//LogData : Log message is a data hex dump
	LogData
)
