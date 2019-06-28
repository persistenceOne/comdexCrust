#!/bin/bash

if [ -d "kafka_2.11-2.2.1" ]; then
echo "File exists"
else
echo "File does not exist, downloading"
wget http://mirrors.estointernet.in/apache/kafka/2.2.1/kafka_2.11-2.2.1.tgz
tar -xzf kafka_2.11-2.2.1.tgz
rm kafka_2.11-2.2.1.tgz
fi

cd kafka_2.11-2.2.1

bin/zookeeper-server-start.sh config/zookeeper.properties &

bin/kafka-server-start.sh config/server.properties &

trap 'killall $BGPID; exit' SIGINT
sleep 1024 &   
BGPID=${!}
sleep 1024  