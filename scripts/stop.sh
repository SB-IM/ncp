#!/usr/bin/sh

echo Run STOP

echo 1 > /sys/class/gpio/gpio5/value
# 急停后等待2秒再恢复
sleep 2
echo 0 > /sys/class/gpio/gpio5/value

