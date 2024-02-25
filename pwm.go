/*
Purpose:
- Hardware generated Pulse Width Modulation (PWM) for Raspberry Pi.

Releases:
- v1.0.0 - 2024/02/04: initial release

Author:
- Klaus Tockloth

Copyright:
- © 2024 | Klaus Tockloth

Contact:
- klaus.tockloth@googlemail.com

Remarks:
- Lint: golangci-lint run --no-config --enable gocritic
- Vulnerability detection: govulncheck ./...
- Raspberry Pi 4: chip=pwmchip0, channels=0, 1, default pins=GPIO18, GPIO19
- Raspberry Pi 5: chip=pwmchip2, channels=2, 3, default pins=GPIO18, GPIO19
- Linux kernel documentation: www.kernel.org/doc/Documentation/pwm.txt
- This library tries to stay as close as possible to the sysfs PWM interface.
*/

package pwm

import (
	"fmt"
	"os"
	"time"
)

// ChannelHandle represents current channel state.
type ChannelHandle struct {
	Chip      string
	Channel   string
	Pulse     time.Duration
	Period    time.Duration
	IsEnabled bool
}

/*
Export exports PWM channel for given chip and channel. Also sets 'Chip' and 'Channel' value in handle.
*/
func Export(chip, channel string) (ChannelHandle, error) {
	var ch ChannelHandle
	sysfsFile := fmt.Sprintf("/sys/class/pwm/%s/export", chip)

	err := os.WriteFile(sysfsFile, []byte(channel), 0644)
	if err != nil {
		return ch, fmt.Errorf("error exporting PWM channel, chip=[%v], channel=[%v], sysfs=[%v], error=[%v]", chip, channel, sysfsFile, err)
	}

	ch.Chip = chip
	ch.Channel = channel

	return ch, nil
}

/*
Unexport unexports PWM channel for given handle.
*/
func Unexport(ch *ChannelHandle) error {
	sysfsFile := fmt.Sprintf("/sys/class/pwm/%s/unexport", ch.Chip)

	err := os.WriteFile(sysfsFile, []byte(ch.Channel), 0644)
	if err != nil {
		return fmt.Errorf("error unexporting PWM channel, ch=[%v], sysfs=[%v], error=[%v]", ch, sysfsFile, err)
	}

	return nil
}

/*
SetPulse sets pulse width (in nanoseconds) for PWM channel for given handle. Pulse must be less than period. Also sets 'Pulse' value in handle.
*/
func SetPulse(ch *ChannelHandle, pulse time.Duration) error {
	sysfsFile := fmt.Sprintf("/sys/class/pwm/%s/pwm%s/duty_cycle", ch.Chip, ch.Channel)
	value := fmt.Sprintf("%d", pulse)

	err := os.WriteFile(sysfsFile, []byte(value), 0644)
	if err != nil {
		return fmt.Errorf("error setting pulse width for PWM channel, ch=[%v], sysfs=[%v], value=[%v], error=[%v]", ch, sysfsFile, value, err)
	}

	ch.Pulse = pulse

	return nil
}

/*
GetPulseFromOS gets pulse width (in nanoseconds) kept by Operating System (Linux) for given PWM channel. Also sets 'Pulse' value in handle.
For special purposes, use handle after PWM channel initialization.
*/
func GetPulseFromOS(ch *ChannelHandle) (time.Duration, error) {
	var pulse time.Duration
	sysfsFile := fmt.Sprintf("/sys/class/pwm/%s/pwm%s/duty_cycle", ch.Chip, ch.Channel)

	data, err := os.ReadFile(sysfsFile)
	if err != nil {
		return pulse, fmt.Errorf("error getting pulse width for PWM channel from Operating System, ch=[%v], sysfs=[%v], error=[%v]", ch, sysfsFile, err)
	}

	_, err = fmt.Sscan(string(data), &pulse)
	if err != nil {
		return pulse, fmt.Errorf("error converting pulse width for PWM channel, ch=[%v], data=[%v], error=[%v]", ch, data, err)
	}

	ch.Pulse = pulse

	return pulse, nil
}

/*
SetPeriod sets period width (in nanoseconds) for PWM channel for given handle. Pulse must be less than period. Also sets 'Period' value in handle.
*/
func SetPeriod(ch *ChannelHandle, period time.Duration) error {
	sysfsFile := fmt.Sprintf("/sys/class/pwm/%s/pwm%s/period", ch.Chip, ch.Channel)
	value := fmt.Sprintf("%d", period)

	err := os.WriteFile(sysfsFile, []byte(value), 0644)
	if err != nil {
		return fmt.Errorf("error setting period width for PWM channel, ch=[%v], sysfs=[%v], value=[%v], error=[%v]", ch, sysfsFile, value, err)
	}

	ch.Period = period

	return nil
}

/*
GetPeriodFromOS gets period width (in nanoseconds) kept by Operating System (Linux) for given PWM channel. Also sets 'Period' value in handle.
For special purposes, use handle after PWM channel initialization.
*/
func GetPeriodFromOS(ch *ChannelHandle) (time.Duration, error) {
	var period time.Duration
	sysfsFile := fmt.Sprintf("/sys/class/pwm/%s/pwm%s/period", ch.Chip, ch.Channel)

	data, err := os.ReadFile(sysfsFile)
	if err != nil {
		return period, fmt.Errorf("error getting period width for PWM channel from Operating System, ch=[%v], sysfs=[%v], error=[%v]", ch, sysfsFile, err)
	}

	_, err = fmt.Sscan(string(data), &period)
	if err != nil {
		return period, fmt.Errorf("error converting period width for PWM channel, ch=[%v], data=[%v], error=[%v]", ch, data, err)
	}

	ch.Period = period

	return period, nil
}

/*
Enable enables modulation on PWM channel for given handle. Also sets 'IsEnabled' value in handle.
*/
func Enable(ch *ChannelHandle) error {
	sysfsFile := fmt.Sprintf("/sys/class/pwm/%s/pwm%s/enable", ch.Chip, ch.Channel)
	value := "1"

	err := os.WriteFile(sysfsFile, []byte(value), 0644)
	if err != nil {
		return fmt.Errorf("error enabling PWM channel, ch=[%v], sysfs=[%v], value=[%v], error=[%v]", ch, sysfsFile, value, err)
	}

	ch.IsEnabled = true

	return nil
}

/*
Disable disables modulation on PWM channel for given handle. Also sets 'IsEnabled' value in handle.
*/
func Disable(ch *ChannelHandle) error {
	sysfsFile := fmt.Sprintf("/sys/class/pwm/%s/pwm%s/enable", ch.Chip, ch.Channel)
	value := "0"

	err := os.WriteFile(sysfsFile, []byte(value), 0644)
	if err != nil {
		return fmt.Errorf("error disabling PWM channel, ch=[%v], sysfs=[%v], value=[%v], error=[%v]", ch, sysfsFile, value, err)
	}

	ch.IsEnabled = false

	return nil
}

/*
GetIsEnabledFromOS gets enable/disable status kept by Operating System (Linux) for given PWM channel. Also sets 'IsEnabled' value in handle.
For special purposes, use handle after PWM channel initialization.
*/
func GetIsEnabledFromOS(ch *ChannelHandle) (bool, error) {
	isEnabled := false
	sysfsFile := fmt.Sprintf("/sys/class/pwm/%s/pwm%s/enable", ch.Chip, ch.Channel)

	data, err := os.ReadFile(sysfsFile)
	if err != nil {
		return isEnabled, fmt.Errorf("error getting enable/disable status for PWM channel from Operating System, ch=[%v], sysfs=[%v], error=[%v]", ch, sysfsFile, err)
	}

	var status int
	_, err = fmt.Sscan(string(data), &status)
	if err != nil {
		return isEnabled, fmt.Errorf("error converting enable/disable status for PWM channel, ch=[%v], data=[%v], error=[%v]", ch, data, err)
	}

	if status > 0 {
		isEnabled = true
	}
	ch.IsEnabled = isEnabled

	return isEnabled, nil
}

/*
FrequencyToPeriod calculates period width for given frequency (hertz).
*/
func FrequencyToPeriod(frequency float64) time.Duration {
	if frequency <= 0.0 {
		return time.Duration(0)
	}

	period := time.Duration(1_000_000_000.0 / frequency)

	return period
}

/*
DutyCycleToPulse calculates pulse width for given channel handle and duty cycle (percent).
Ensures that pulse is less than period.
*/
func DutyCycleToPulse(ch *ChannelHandle, dutyCycle float64) time.Duration {
	pulse := time.Duration(float64(ch.Period) * dutyCycle / 100.0)

	if pulse >= ch.Period {
		pulse = ch.Period - 1
	}

	return pulse
}

/*
Initialize initializes PWM channel. The smallest value for pulse is 0. The smallest value
for period is 1 on Pi 5 (kernel 6.1), and 15 on Pi 4 (kernel 6.1).
*/
func Initialize(chip, channel string, pulse, period, waitForPermission time.Duration) (ChannelHandle, error) {
	channelHandle, err := Export(chip, channel)
	if err != nil {
		return channelHandle, err
	}

	// give Operating System (Linux) some time to give us the permission to use the exported PWM channel
	time.Sleep(waitForPermission)

	pulseKeptByOS, err := GetPulseFromOS(&channelHandle)
	if err != nil {
		return channelHandle, err
	}

	if pulseKeptByOS > 0 {
		// set pulse to 0; this ensures to set any allowed value for period
		err = SetPulse(&channelHandle, 0)
		if err != nil {
			return channelHandle, err
		}
	}

	err = SetPeriod(&channelHandle, period)
	if err != nil {
		return channelHandle, err
	}

	err = SetPulse(&channelHandle, pulse)
	if err != nil {
		return channelHandle, err
	}

	return channelHandle, nil
}
