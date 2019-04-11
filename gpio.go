package zl_mumble

import (
	"fmt"
	"time"

	"github.com/dchote/gpio"
	"github.com/stianeikeland/go-rpio"
)

func (b *Talkiepi) initGPIO() {
	// we need to pull in rpio to pullup our button pin
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		b.GPIOEnabled = false
		return
	} else {
		b.GPIOEnabled = true
	}

	TransmitButtonPinPullUp := rpio.Pin(TransmitButtonPin)
	TransmitButtonPinPullUp.PullUp()
	TestButtonPinPullUp := rpio.Pin(TestButtonPin)
	TestButtonPinPullUp.PullUp()

	_ = rpio.Close()

	// unfortunately the gpio watcher stuff doesnt work for me in this context, so we have to poll the button instead
	b.TransmitButton = gpio.NewInput(TransmitButtonPin)
	go func() {
		for {
			currentState, err := b.TransmitButton.Read()

			if currentState != b.TransmitButtonState && err == nil {
				b.TransmitButtonState = currentState

				if b.Stream != nil {
					if b.TransmitButtonState == 1 {
						fmt.Printf("Transmit-Button is released\n")
						b.TransmitStop()
					} else {
						fmt.Printf("Transmit-Button is pressed\n")
						b.TransmitStart()
					}
				}

			}

			time.Sleep(200 * time.Millisecond)
		}
	}()

	b.TestButton = gpio.NewInput(TestButtonPin)
	go func() {
		for {
			currentState, err := b.TestButton.Read()

			if currentState != b.TestButtonState && err == nil {
				b.TestButtonState = currentState

				if b.TestButtonState == 1 {
					fmt.Printf("Test-Button is released\n")
					b.TestStop()
				} else {
					fmt.Printf("Test-Button is pressed\n")
					b.TestStart()
				}
			}

			time.Sleep(200 * time.Millisecond)
		}
	}()

	// then we can do our gpio stuff
	b.SmallOnlineLED = gpio.NewOutput(SmallOnlineLEDPin, false)
	b.SmallParticipantsLED = gpio.NewOutput(SmallParticipantsLEDPin, false)
	b.BigParticipantsLED = gpio.NewOutput(BigParticipantsLEDPin, false)
	b.SmallTransmitLED = gpio.NewOutput(SmallTransmitLEDPin, false)
	b.BigTransmitLED = gpio.NewOutput(BigTransmitLEDPin, false)
}

func (b *Talkiepi) LEDOn(LED gpio.Pin) {
	if b.GPIOEnabled == false {
		return
	}

	_ = LED.High()
}

func (b *Talkiepi) LEDOff(LED gpio.Pin) {
	if b.GPIOEnabled == false {
		return
	}

	_ = LED.Low()
}

func (b *Talkiepi) LEDOnAll() {
	if b.GPIOEnabled == false {
		return
	}

	b.LEDOn(b.SmallOnlineLED)
	b.LEDOn(b.SmallParticipantsLED)
	b.LEDOn(b.BigParticipantsLED)
	b.LEDOn(b.SmallTransmitLED)
	b.LEDOn(b.BigTransmitLED)
}

func (b *Talkiepi) LEDOffAll() {
	if b.GPIOEnabled == false {
		return
	}

	b.LEDOff(b.SmallOnlineLED)
	b.LEDOff(b.SmallParticipantsLED)
	b.LEDOff(b.BigParticipantsLED)
	b.LEDOff(b.SmallTransmitLED)
	b.LEDOff(b.BigTransmitLED)
}
