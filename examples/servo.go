// Servo control with hardware generated Pulse Width Modulation (PWM).
// Raspberry Pi 4: chip=pwmchip0, channel=0 on pin=GPIO18, channel=1 on pin=GPIO19
// Raspberry Pi 5: chip=pwmchip2, channel=2 on pin=GPIO18, channel=3 on pin=GPIO19
// Servo parameters:
// - period           =   20000 μs
// - min angle        =     500 μs
// - neutral position =    1500 μs
// - max angle        =    2500 μs

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Klaus-Tockloth/pwm"
)

func main() {
	chip := "pwmchip2"
	channel := "2"
	pulse := time.Duration(1500 * time.Microsecond) // neutral position
	period := time.Duration(20000 * time.Microsecond)
	waitForPermission := time.Duration(500 * time.Millisecond)

	servoPwmHandle, err := pwm.Initialize(chip, channel, pulse, period, waitForPermission)
	if err != nil {
		log.Printf("error [%v] at pwm.Initialze()", err)
		return
	}

	defer func() {
		err = pwm.Unexport(&servoPwmHandle)
		if err != nil {
			log.Printf("error [%v] at pwm.Unexport()", err)
		}
	}()

	fmt.Printf("servo neutral position ...\n")
	err = pwm.Enable(&servoPwmHandle)
	if err != nil {
		log.Printf("error [%v] at pwm.Enable()", err)
		return
	}
	time.Sleep(2 * time.Second)

	fmt.Printf("servo min position ...\n")
	err = pwm.SetPulse(&servoPwmHandle, time.Duration(500*time.Microsecond))
	if err != nil {
		log.Printf("error [%v] at pwm.SetPulse()", err)
		return
	}
	time.Sleep(2 * time.Second)

	fmt.Printf("servo max position ...\n")
	err = pwm.SetPulse(&servoPwmHandle, time.Duration(2500*time.Microsecond))
	if err != nil {
		log.Printf("error [%v] at pwm.SetPulse()", err)
		return
	}
	time.Sleep(2 * time.Second)

	err = pwm.Disable(&servoPwmHandle)
	if err != nil {
		log.Printf("error [%v] at pwm.Disable()", err)
		return
	}

	// program exit: unexport is done by defer function
	fmt.Printf("Done\n")
}
