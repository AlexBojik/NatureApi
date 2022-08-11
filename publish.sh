#!/bin/bash

REMOTEHOST=root@10.127.146.7
REMOTEROOT=$REMOTEHOST:/root/go/src/water-api

scp ./esia/* $REMOTEROOT/esia/
scp ./handlers/* $REMOTEROOT/handlers/
scp ./models/* $REMOTEROOT/models/
scp ./sql/* $REMOTEROOT/sql/
scp ./utils/* $REMOTEROOT/utils/
scp ./server.go $REMOTEROOT/

ssh $REMOTEHOST "cd /root/go/src/water-api; /root/go/src/water-api/rebuildstart.sh"
