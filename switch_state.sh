#!/bin/bash

name=$1

docker exec -it $name /bin/sh  -c "touch start" 

