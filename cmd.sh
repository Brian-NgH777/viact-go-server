#!/bin/bash

cmdHandler () {
	if [ $1 == 'get_first_frame' ]
  then
      action get_first_frame "RTSP_LINK=$2 FILE_NAME=$3"
  elif [ $1 == 'livestream' ]
  then
      action livestream "RTSP_LINK=$2 RTMP_LINK=$3"
  else
    echo "cmdHandler not run"
  fi
}

cmdHandler $1 $2 $3
