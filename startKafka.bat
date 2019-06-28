@echo off
cd kafka_2.12-2.1.0
title startKafka
echo "this file is for windows and is WIP. it wont work if you dont have kafka and zookeeper ready and configured"

echo "preparing to clean kafka servers"
echo "starting zookeeper"
start powershell.exe -noexit -file ".\bin\windows\zookeeper-server-start .\config\zookeeper.properties"

echo "starting kafka"
timeout /t 2 /nobreak
start powershell.exe -noexit -file ".\bin\windows\kafka-server-start .\config\server.properties"


timeout /t 2 /nobreak
echo "deleting topics"
.\bin\windows\kafka-topics --delete --zookeeper localhost:2181 --topic "Topic" 
echo "marked for deletion"

echo "stopping kafka"
.\bin\windows\kafka-server-stop
timeout /t 2 /nobreak

echo "stopping zookeeper"
zkserver stop
timeout /t 2 /nobreak

timeout /t 2 /nobreak
echo "starting kafka fresh"
start powershell.exe -noexit -file ".\bin\windows\zookeeper-server-start .\config\zookeeper.properties"
start powershell.exe -noexit -file ".\bin\windows\kafka-server-start .\config\server.properties"
cd ..