#!/bin/sh
echo 'received '$1
if [ -z "$1" ] ; then
  echo 'testing all packages...'
  for p in transport schedule
  do
    echo $p
    go test ./$p/ -race
  done
else
  echo testing $1...
  go test ./$1/ -race
fi
echo 'Test finished.'
