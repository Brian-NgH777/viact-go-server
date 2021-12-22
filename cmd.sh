#!/bin/bash

cmdHandler () {
  echo $2
  echo $3

	if [ $1 == 'get_first_frame' ]
  then
     echo $2
     echo $3
#    action get_first_frame "RTSP_LINK=rtsp://admin:Viact123@192.168.92.110:554/live FILE_NAME=luna.jpg"
  elif [ $1 == 'livestream' ]
  then
#    action livestream "RTSP_LINK=rtsp://admin:Viact123@192.168.92.111/live RTMP_LINK=rtmp://54.254.0.41/live/test7"
  else
    echo "cmdHandler not run"
  fi
}

cmdHandler $1 $2 $3
