#! /bin/sh
exitWithStatus() {
  if [ $1 != 0 ]; then
    echo 'Build failed!'
  else
    echo 'Build finished successfully.'
  fi
  exit $1
}
echo 'received '$1
if [ -z "$1" ] ; then
  echo 'Build messagesender locally...'
  go build -o build/messagesender main.go
  BUILD_STATUS=$?
  exitWithStatus $BUILD_STATUS
else
  case "$1" in
      docker)
        echo 'Build messagesender and commservice docker images...'
        docker build -t commservice -f Dockerfile_commservice .
        BUILD_STATUS=$?
        if [ $BUILD_STATUS != 0 ]; then
          exitWithStatus $BUILD_STATUS
        fi
        docker build -t messagesender .
        BUILD_STATUS=$?
        exitWithStatus $BUILD_STATUS ;;
      *)
        echo 'Invalid option! To bild docker image please call:\nbuild.sh docker\n'
         exit 0 ;;
    esac
fi
