#!/bin/bash

kernel_num_1=$(uname -r | awk -F '.' '{print $1}')
kernel_num_2=$(uname -r | awk -F '.' '{print $2}')
if [ ${kernel_num_1} -lt 3 ];then
    echo "kernel version is minor < 3.0"
    exit 1
fi

if [ ${kernel_num_1} -eq 3 -a ${kernel_num_2} -lt 10 ];then
    echo "kernel version is minor < 3.10"
    exit 1
fi
