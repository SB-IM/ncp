#!/bin/bash

# Allowable time error (s)
allowable_diff=300

time_diff=1000

get_time_diff(){
  time_standard=`curl http://worldtimeapi.org/api/ip.txt -c - | grep -v utc_datetime | grep datetime | cut -d " " -f 2`

  echo "time_standard:" "$time_standard"

  time_local=$(date -I"ns")
  #time_local="2019-05-25T15:13:09,900802674+08:00"

  echo "time_local:" "$time_local"

  time_diff=$(( $(date -d $time_standard +%s) - $(date -d $time_local +%s) ))
}

until [ $time_diff -lt $allowable_diff ]
do
  sleep 10
  get_time_diff
  echo "time_diff:" "$time_diff"
  echo "time_diff >" "$allowable_diff" "Try Retry"
done

echo "Wait Time Synchronized Succeed"

