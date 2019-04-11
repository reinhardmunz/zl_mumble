package zl_mumble

import (
	"crypto/tls"

	"github.com/dchote/gpio"
	"github.com/dchote/gumble/gumble"
	"github.com/dchote/gumble/gumbleopenal"
)

// Raspberry Pi GPIO pin assignments (CPU pin definitions)
const (
	SmallOnlineLEDPin       uint = 6
	SmallParticipantsLEDPin uint = 7
	BigParticipantsLEDPin   uint = 17
	SmallTransmitLEDPin     uint = 8
	BigTransmitLEDPin       uint = 4
	TransmitButtonPin       uint = 18
	TestButtonPin           uint = 14
)

type Talkiepi struct {
	Config *gumble.Config
	Client *gumble.Client

	Address   string
	TLSConfig tls.Config

	ConnectAttempts uint

	Stream *gumbleopenal.Stream

	ChannelName    string
	IsConnected    bool
	IsTransmitting bool
	IsTesting      bool

	ParticipantCount int

	PttStartFile string
	PttStopFile	 string

	GPIOEnabled          bool
	SmallOnlineLED       gpio.Pin
	SmallParticipantsLED gpio.Pin
	BigParticipantsLED   gpio.Pin
	SmallTransmitLED     gpio.Pin
	BigTransmitLED       gpio.Pin
	TransmitButton       gpio.Pin
	TransmitButtonState  uint
	TestButton           gpio.Pin
	TestButtonState      uint
}
