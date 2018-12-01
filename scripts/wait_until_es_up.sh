#!/bin/bash
(while true; do curl -s 127.0.0.1:9200/_cat/health | awk '{ print $4 }' | grep -qc 'green' && exit 0; done) &
jpid="$!"
echo "Checking if elasticsearch is up for a minute"
cnt=1
while true; do
    cnt="$((cnt + 1))"
    if kill -0 "${jpid}"; then
        echo "${jpid} still running, ES cluster status is not green"
    else
        echo "cluster status green"
        exit 0
    fi

    if [ $cnt -gt 60 ]; then
        echo "failed to get a cluster in a minute, killing"
        kill -9 "${jpid}"
        exit 1
    fi
    sleep 1
done