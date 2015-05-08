package gphoto

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
	gpWidgetWindow = iota //(0)
	gpWidgetSection
	gpWidgetText
	gpWidgetRange
	gpWidgetToggle
	gpWidgetRadio
	gpWidgetMenu
	gpWidgetButton
	gpWidgetDate
)

//widget types
const (
	//WidgetWindow is the toplevel configuration widget. It should likely contain multiple #WidgetSection entries.
	WidgetWindow WidgetType = "window"
	//WidgetSection : Section widget (think Tab)
	WidgetSection WidgetType = "section"
	//WidgetText : Text widget (string)
	WidgetText WidgetType = "text"
	//WidgetRange : Slider widget (float)
	WidgetRange WidgetType = "range"
	//WidgetToggle : Toggle widget (think check box) (int)
	WidgetToggle WidgetType = "toggle"
	//WidgetRadio : Radio button widget (string)
	WidgetRadio WidgetType = "radio"
	//WidgetMenu : Menu widget (same as RADIO) (string)
	WidgetMenu WidgetType = "menu"
	//WidgetButton : Button press widget ( CameraWidgetCallback )
	WidgetButton WidgetType = "button"
	//WidgetDate : Date entering widget (int)
	WidgetDate WidgetType = "date"
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

//File types
const (
	//FileTypePreview is a preview of an image
	FileTypePreview = iota
	//FileTypeNormal is regular normal data of a file
	FileTypeNormal
	//FileTypeRaw usually the same as FileTypeNormal for modern cameras ( left for compatibility purposes)
	FileTypeRaw
	//FileTypeAudio is a audio view of a file. Perhaps an embedded comment or similar
	FileTypeAudio
	//FileTypeExif is the  embedded EXIF data of an image
	FileTypeExif
	//FileTypeMetadata is the metadata of a file, like Metadata of files on MTP devices
	FileTypeMetadata
)
