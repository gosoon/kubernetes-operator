#! /bin/bash
export LC_ALL=en_US.UTF-8
export LANG=en_US.UTF-8

ROOT=`cd $(dirname $0); pwd`
BIN=${ROOT}/kubernetes-operator
LOGDIR=${ROOT}/logs/

mkdir -p ${LOGDIR}
process="kubernetes-operator"

#使用pgrep判断程序是否执行
function check_pid(){
    run_pid=$(pgrep $process)
        echo $run_pid
}
pid=$(check_pid)

function status(){
    if [ "x_$pid" != "x_" ]; then
        echo "$process running with pid: $pid"
    else
        echo "ERROR: kubernetes-operator may not running!"
    fi
}

start() {
    echo -n "starting kubernetes-operator... "
    exec ${ROOT}/kubernetes-operator 
}

stop() {
    echo -n "stopping kubernetes-operator... "
    kill `cat ${pid}`
    return 0
}

restart() {
    echo -n "restarting kubernetes-operator... "
    stop
    start
    echo "finished, plz check by urself"
}

case "$1" in
start)
    start
    ;;

stop)
    stop
    ;;

restart)
    restart
    ;;
*)
    echo "Usage: $0 {start|stop|restart}"
    exit 1
    esac
