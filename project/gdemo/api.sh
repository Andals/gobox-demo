#! /bin/bash

cur_path=`dirname $0`

cd $cur_path/../../
prj_home=`pwd`

cd $prj_home/src

# 关闭已经启动的程序
pidFile="$prj_home/tmp/api.pid"
if [ -f $pidFile ]; then 
    pid=`cat $pidFile`
    if [ -f "/proc/$pid/stat" ]; then 
        echo -e "\ntry to kill family api exist process $pid."
        kill $pid
    else
        echo -e "\nfamily api process $pid not exist, please check it."
    fi
fi

echo -e "try to start gdemo api process\n"

# 启动新的程序
./go.sh install main/api
$prj_home/bin/api -prjHome=$prj_home &
