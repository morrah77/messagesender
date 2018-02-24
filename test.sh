#!/bin/sh
echo 'received '$1
if [ -z "$1" ] ; then
  echo 'testing $1...'
  go test ./$1/
else
echo 'testing all packages...'
  for p in transport schedule
  do
    go test ./$1/
  done
fi
echo 'Test finished.'
