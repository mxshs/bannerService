#!/bin/bash

host="localhost"
port="3000"
resource="banner"
arbitrary=""

# perf test with get all banners by default, if "rps" is passed as argument, switch to user_banner endpoint tests
if [ "$1" == "rps" ]
then
    resource="user_banner?tag_id=1&feature_id=1"
fi

# wrk flags
if [ "$#" > 1 ]
then
    arbitrary=${@:2:$#}
fi

wrk -t12 -c99 -d1m http://$host:$port/$resource $arbitrary
