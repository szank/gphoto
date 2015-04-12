package main

// errors
const (
	GP_OK = 0
	//Generic Error
	GP_ERROR = -1
	//Bad parameters passed
	GP_ERROR_BAD_PARAMETERS = -2
	//Out of memory
	GP_ERROR_NO_MEMORY = -3
	//Error in the camera driver
	GP_ERROR_LIBRARY = -4
	//Unknown libgphoto2 port passed
	GP_ERROR_UNKNOWN_PORT = -5
	//Functionality not supported
	GP_ERROR_NOT_SUPPORTED = -6
	//Generic I/O error
	GP_ERROR_IO = -7
	//Buffer overflow of internal structure
	GP_ERROR_FIXED_LIMIT_EXCEEDED = -8
	//Operation timed out
	GP_ERROR_TIMEOUT = -10
	//Serial ports not supported
	GP_ERROR_IO_SUPPORTED_SERIAL = -20
	//USB ports not supported
	GP_ERROR_IO_SUPPORTED_USB = -21
	//Error initialising I/O
	GP_ERROR_IO_INIT = -31
	//I/O during read
	GP_ERROR_IO_READ = -34
	//I/O during write
	GP_ERROR_IO_WRITE = -35
	//I/O during update of settings
	GP_ERROR_IO_UPDATE = -37
	//Specified serial speed not possible.
	GP_ERROR_IO_SERIAL_SPEED = -41
	//Error during USB Clear HALT
	GP_ERROR_IO_USB_CLEAR_HALT = -51
	//Error when trying to find USB device
	GP_ERROR_IO_USB_FIND = -52
	//Error when trying to claim the USB device
	GP_ERROR_IO_USB_CLAIM = -53
	//Error when trying to lock the device
	GP_ERROR_IO_LOCK = -60
	//Unspecified error when talking to HAL
	GP_ERROR_HAL = -70
)

//widget types
const (
	// Window widget
	//This is the toplevel configuration widget. It should likely contain multiple #GP_WIDGET_SECTION entries.
	GP_WIDGET_WINDOW = iota //(0)
	//Section widget (think Tab)
	GP_WIDGET_SECTION
	//Text widget (string)
	GP_WIDGET_TEXT
	//Slider widget (float)
	GP_WIDGET_RANGE
	// Toggle widget (think check box) (int)
	GP_WIDGET_TOGGLE
	// Radio button widget (string)
	GP_WIDGET_RADIO
	// Menu widget (same as RADIO) (string)
	GP_WIDGET_MENU
	// Button press widget ( CameraWidgetCallback )
	GP_WIDGET_BUTTON
	//Date entering widget (int)
	GP_WIDGET_DATE
)

//Log level
const (
	//Log message is an error infomation
	GP_LOG_ERROR = iota
	//Log message is an verbose debug infomation
	GP_LOG_VERBOSE
	//Log message is an debug infomation
	GP_LOG_DEBUG
	//Log message is a data hex dump
	GP_LOG_DATA
)
