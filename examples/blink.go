// Blinking LED with hardware generated Pulse Width Modulation (PWM).
// Raspberry Pi 4: chip=pwmchip0, channel=0 on pin=GPIO18, channel=1 on pin=GPIO19
// Raspberry Pi 5: chip=pwmchip2, channel=2 on pin=GPIO18, channel=3 on pin=GPIO19

package main

import (
	"RaspberryPi/pwm"
	"fmt"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	chip := "pwmchip2"
	channel := "2"
	pulse := time.Duration(500 * time.Millisecond)
	period := time.Duration(1000 * time.Millisecond)
	waitForPermission := time.Duration(500 * time.Millisecond)

	channelHandle, err := pwm.Initialize(chip, channel, pulse, period, waitForPermission)
	if err != nil {
		log.Printf("error [%v] at pwm.Initialze()", err)
		return
	}

	defer func() {
		err = pwm.Unexport(&channelHandle)
		if err != nil {
			log.Printf("error [%v] at pwm.Unexport(), channelHandle = %#v", err, channelHandle)
		}
	}()

	err = pwm.Enable(&channelHandle)
	if err != nil {
		log.Printf("error [%v] at pwm.Enable(), channelHandle = %#v", err, channelHandle)
		return
	}

	fmt.Printf("LED should blink ...\n")
	time.Sleep(10 * time.Second)

	// adjust period width for given frequency in hz
	err = pwm.SetPeriod(&channelHandle, pwm.FrequencyToPeriod(0.5))
	if err != nil {
		log.Printf("error [%v] at pwm.SetPeriod(), channelHandle = %#v", err, channelHandle)
		return
	}

	fmt.Printf("LED should blink slower ...\n")
	time.Sleep(10 * time.Second)

	// adjust pulse width for given duty cyle in %
	err = pwm.SetPulse(&channelHandle, pwm.DutyCycleToPulse(&channelHandle, 75.0))
	if err != nil {
		log.Printf("error [%v] at pwm.SetPulse(), channelHandle = %#v", err, channelHandle)
		return
	}

	fmt.Printf("LED should blink longer ...\n")
	time.Sleep(10 * time.Second)

	err = pwm.Disable(&channelHandle)
	if err != nil {
		log.Printf("error [%v] at pwm.Disable(), channelHandle = %#v", err, channelHandle)
		return
	}

	// program exit: unexport is done by defer function
}
