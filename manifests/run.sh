#!/usr/bin/env bash

# - the right sequence to start mongo upf nrf amf smf udr pcf udm nssf ausf webui
NF_LIST="f5gc-upf f5gc-mongodb f5gc-nrf f5gc-amf f5gc-ausf-kwaf f5gc-nssf f5gc-pcf  f5gc-smf f5gc-udm f5gc-udr f5gc-webui dbg-rockmongo"

CMD_POOL="apply|delete|show"
if [[ ! "$1" =~ $CMD_POOL ]]
then
    echo "Usage: $0 [ ${CMD_POOL//|/ | } ]"
    exit 1
fi

if [ "$1" == "show" ]; then
    for NF in ${NF_LIST}; do
        kubectl get deployment -n f5gc ${NF} --no-headers=true
    done
else
    for NF in ${NF_LIST}; do
        kubectl $1 -k ${NF}
        sleep 5
    done
fi