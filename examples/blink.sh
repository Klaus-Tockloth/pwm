#!/bin/sh
# Raspberry Pi 4: chip=pwmchip0, channel=0 on pin=GPIO18, channel=1 on pin=GPIO19
# Raspberry Pi 5: chip=pwmchip2, channel=2 on pin=GPIO18, channel=3 on pin=GPIO19
# Linux kernel documentation: www.kernel.org/doc/Documentation/pwm.txt

set -o verbose

# set chip name and channel number
CHANNEL=2
CHIP=pwmchip2
PWM=pwm$CHANNEL

# export channel
echo $CHANNEL > /sys/class/pwm/$CHIP/export
sleep 1

# let something blink (e.g. LED on pin GPIO18)
DutyCycleKeptByLinux=$(cat /sys/class/pwm/$CHIP/$PWM/duty_cycle)
if [ $DutyCycleKeptByLinux != "0" ]; then
  echo 0 > /sys/class/pwm/$CHIP/$PWM/duty_cycle
fi
echo 1000000000 > /sys/class/pwm/$CHIP/$PWM/period
echo 500000000 > /sys/class/pwm/$CHIP/$PWM/duty_cycle
echo 1 > /sys/class/pwm/$CHIP/$PWM/enable
sleep 10
echo 0 > /sys/class/pwm/$CHIP/$PWM/enable

# unexport channel
echo $CHANNEL > /sys/class/pwm/$CHIP/unexport
