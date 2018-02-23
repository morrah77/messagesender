#! /bin/sh
echo 'received '$1
if [ -z "$1" ] ; then
  echo 'Build locally...'
  go install main
else
  case "$1" in
      docker)
        echo 'Build docker image...'
        docker build -t messagesenderer . ;;
      *)
        echo 'Invalid option! To bild docker image please call:\nbuild.sh docker\n'
         exit 0 ;;
    esac
fi
echo 'Build finished.'