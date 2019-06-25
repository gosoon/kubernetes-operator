#!/bin/bash


[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ../deploy/config.sh ] && . ../deploy/config.sh || exit

rpm -e ${define_docker_version}
if [ $? -eq 0 ];then
    action "remove ${define_docker_version} failed !!!" /bin/false
else
    action "remove ${define_docker_version} success !" /bin/true
fi
