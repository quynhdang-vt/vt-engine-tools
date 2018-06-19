#!/bin/bash
set -x
if [ $# -ge 2 ]; then
  filename=$1
fi
if [ -z $filename ]; then
  filename=10_Alan_test.mp3
fi
f=${filename##*/}
base_filename=${f%.*}
c0_filename=${base_filename}_0.wav
c1_filename=${base_filename}_1.wav

r0_filename=${base_filename}_0_16K.wav
r1_filename=${base_filename}_1_16K.wav

rm -f $c0_filename
rm -f $c1_filename
rm -f $r0_filename
rm -f $r1_filename
###------------------------
## SPLIT the files first
ffmpeg -i 10_Alan_test.mp3 -map_channel 0.0.0 ${c0_filename} -map_channel 0.0.1 ${c1_filename}
ffmpeg -i $c0_filename -acodec pcm_s16le -ac 1 -ar 16000 $r0_filename
ffmpeg -i $c1_filename -acodec pcm_s16le -ac 1 -ar 16000 $r1_filename

###------------------------
### RUN THE ENGINES
ENGINE_ID=cc8f5c80-b2d9-b22b-2eb7-16a2040f5e38
#ASSET_FILE=$filename ENGINE_ID=$ENGINE_ID ./test-engine.sh
ASSET_FILE=$r0_filename ENGINE_ID=$ENGINE_ID ./test-engine.sh
ASSET_FILE=$r1_filename ENGINE_ID=$ENGINE_ID ./test-engine.sh

echo "---------------- Generate the following TDOs ---------------"
cat recording_ids.log
echo "---------------- Run ./check-job.sh later for engine outputs ---"
