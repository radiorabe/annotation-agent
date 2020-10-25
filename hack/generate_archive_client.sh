#!/bin/bash

pushd `dirname $0`/..

swagger generate client \
  --spec=annotator/archive/raar/swagger.json \
  --target=annotator/archive/raar/ \
  --default-scheme=https \
  --additional-initialism=RAAR \
  --model=AudioFile \
  --model=Broadcast \
  --model=Show \
  --model=User \
  --operation=GetBroadcastsID \
  --operation=GetBroadcastsBroadcastIDAudioFiles \
  --operation=GetBroadcastsYearMonthDay \
  --operation=PostLogin

popd
