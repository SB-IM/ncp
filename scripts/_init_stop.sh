#!/usr/bin/sh
# nvidia nano GPIO

echo Stop GPIO initial.
echo 12 > /sys/class/gpio/unexport
sleep 0.5
echo 12 > /sys/class/gpio/export
sleep 1
echo out > /sys/class/gpio/gpio12/direction

