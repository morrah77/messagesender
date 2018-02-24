#!/bin/sh
echo 'received '$1
if [ -z "$1" ] ; then
  echo 'Run locally...'
  ./build/commservice &
  ./build/messageservice
else
  case "$1" in
      docker)
        echo 'Run with docker...'
        docker build -t messagesenderer . ;;
      *)
        echo 'Invalid option! To bild docker image please call:\nbuild.sh docker\n'
         exit 0 ;;
    esac
fi

