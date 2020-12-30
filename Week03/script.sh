#!/bin/bash
int=1
while [ $int -lt 5 ]
do
    echo $int
    int=`expr $int + 1`
done


# while true
# do
#     sleep 1s
#     echo $(date) >> date.txt
#     echo $(date)
# done
