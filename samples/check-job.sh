#!/bin/bash

NC='\033[0m'
red()      { TSTAMP=$(date +"%Y-%m-%d %H:%M:%S"); echo -e "\033[0;31m$TSTAMP $@${NC}";}
blue()     { TSTAMP=$(date +"%Y-%m-%d %H:%M:%S"); echo -e "\033[0;96m$TSTAMP $@${NC}";}
green()    { TSTAMP=$(date +"%Y-%m-%d %H:%M:%S"); echo -e "\033[1;32m$TSTAMP $@${NC}";}
yellow()   { TSTAMP=$(date +"%Y-%m-%d %H:%M:%S"); echo -e "\033[0;33m$TSTAMP $@${NC}";}
darkGreen(){ TSTAMP=$(date +"%Y-%m-%d %H:%M:%S"); echo -e "\033[38;5;002m$TSTAMP $@${NC}";}

#AIWARE_GRAPHIQL_TOKEN=
GRAPHQL_HOST=api.veritone.com
BATCH_RETRIEVE=1
RECORDING_IDS_FILE=./recording_ids.log

CURL_OPTS="-s -k"

#Retrieve v-vlf and ttml..
function verify  {
    local recording_id=$1

    GRAPHQL_API_URL=https://${GRAPHQL_HOST}/v3/graphql

    mkdir -p $OUTDIR
    local recording_output_file=$OUTDIR/${recording_id}-status.json
    JOB_DONE=true
set -x
    curl ${CURL_OPTS} -H "Authorization: Bearer $TOKEN"  -H "Content-type: application/json"  -o ${recording_output_file}  \
-d '{"query":"{vlf:temporalDataObject(id: \"'${recording_id}'\") {id assets(assetType:\"v-vlf\"){ count records{id assetType contentType signedUri jsondata}} } transcript:temporalDataObject(id: \"'${recording_id}'\") {id assets(assetType:\"transcript\"){ count records{id assetType contentType signedUri jsondata}} }}"}' \
  ${GRAPHQL_API_URL}
set +x
    cat ${recording_output_file} | jq ''

    # v-vlf
    nAssets=$(cat ${recording_output_file} | jq '.data.vlf.assets.count'  | sed 's/\"//g')
    (( nAssets -- ))
    for i in $(seq 0 $nAssets); do
      assetSignedUri=`cat ${recording_output_file} | jq ".data.vlf.assets.records[${i}].signedUri" | sed 's/\"//g'`
      set -x
      curl -o $OUTDIR/${recording_id}_${i}_v-vlf.json $assetSignedUri
      set +x
    done
    
    # transcript (ttml)
    nAssets=$(cat ${recording_output_file} | jq '.data.transcript.assets.count'  | sed 's/\"//g')
    (( nAssets -- ))
    for i in $(seq 0 $nAssets); do
      assetSignedUri=`cat ${recording_output_file} | jq ".data.transcript.assets.records[${i}].signedUri" | sed 's/\"//g'`
      set -x
      curl -o $OUTDIR/${recording_id}_${i}_ttml.xml $assetSignedUri
      set +x
    done
}

help()
{
  echo "Usage: $0 RECORDING_ID"
  echo ""
  echo "For example: "
  echo "     $0 165"
  exit 1
}
help2(){
  echo "Usage: BATCH_RETRIEVE=1 $0 RECORDING_ID_FILE"
  echo ""
  echo "For example: "
  echo "    BATCH_RETRIEVE=1 $0 /opt/veritone/bundle/tests/conductor/20180523_074913/conductor-verify.csv"
  exit 1
}
# ------ START HERE.......

if [ -z $AIWARE_GRAPHIQL_TOKEN ]; then
    echo "Please define $AIWARE_GRAPHIQL_TOKEN"
    exit 1
fi
TOKEN="${AIWARE_GRAPHIQL_TOKEN}"


OUTDIR="results"


if [ -z $BATCH_RETRIEVE ]; then
    if [ $# -lt 1 ]; then
      help
    fi

    verify $1 $GRAPHQL_HOST
else
    if [ -z $RECORDING_IDS_FILE ]; then
       if [ $# -lt 1 ]; then
          help2
       fi
       RECORDING_IDS_FILE=$1
    fi

    blue "--------------------------------------------------------------------"
    blue "Checking jobs for recordings as read from file $RECORDING_IDS_FILE...."
    blue "--------------------------------------------------------------------"
    while read -r astring; do
        recording_id=$(echo $astring| awk '{print $1}' | sed 's/,*$//g')
        verify $recording_id $GRAPHQL_HOST
    done <<< "$(cat $RECORDING_IDS_FILE)"
fi

echo "THE END.."
