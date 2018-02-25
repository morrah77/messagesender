#!/bin/sh
echo 'received '$1
if [ -z "$1" ] ; then
  echo 'Run locally...'
  echo 'Run commservice...'
  ./build/commservice &
   echo $! > COMMSERVICE_PID
   echo $(<COMMSERVICE_PID)
  read CS_PID < COMMSERVICE_PID
  #TODO(h.lazar) consider to check wheather commservice is up
  sleep 2s
  echo 'Run messagesender...'
  ./build/messagesender
   echo "Kill commservice with PID $CS_PID..."
  `sudo kill -2 $CS_PID`
else
  case "$1" in
      docker)
        echo 'Run with docker...'

        echo 'Start commservice container...'
        docker run --rm -d --name=commservice --net=bridge --expose=9090 commservice
        SERVICE_IP=`docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' commservice`
        echo 'commservice started at '$SERVICE_IP':`docker ps --filter Name=commservice --format '{{.Ports}}``

        echo 'Start messagesender container...'
        docker run -d --rm --name=messagesender --net=bridge  -e "SERVICE_IP=$SERVICE_IP" messagesender
        echo 'messagesender started at '`docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' messagesender`:`docker ps --filter Name=messagesender --format '{{.Ports}}'`

        docker logs -f commservice ;;

      *)
        echo 'Invalid option! To run with docker please call:\nrun.sh docker\n'
         exit 0 ;;
    esac
fi

