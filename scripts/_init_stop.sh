#!/usr/bin/sh

echo Stop GPIO initial.
# RPI的GPIO5控制软急停，为防止出错，先unexport
echo 5 > /sys/class/gpio/unexport
sleep 0.5
echo 5 > /sys/class/gpio/export
# 需等待1秒，参考https://elinux.org/RPi_GPIO_Code_Samples#Shell
sleep 1
echo out > /sys/class/gpio/gpio5/direction

