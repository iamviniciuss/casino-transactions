#!/bin/bash

start=`date +%s`

air -c .air.consumer.toml

end=`date +%s`

runtime=$(echo "$end - $start" | bc -l)

echo "#### Runtime: $runtime seconds ####"

echo "Hit CTRL+C"
tail -f /dev/null