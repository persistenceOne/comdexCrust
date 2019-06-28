cd %GOPATH%\src\github.com\comdex-blockchain
echo "to build and to update deps run makeAll file"

echo "cleaning"

rmdir /Q /S  C:\.maincli
rmdir /Q /S  C:\.maind

rmdir /Q /S  C:\.fiatcli
rmdir /Q /S  C:\.fiatd

rmdir /Q /S  C:\.assetcli
rmdir /Q /S  C:\.assetd

@echo off 
echo "------------------------------------------------------------------------------------------"
echo "initialising"
maind init --name main --chain-id comdex-main
timeout /t 2 /nobreak

echo "------------------------------------------------------------------------------------------"
start maind start --address tcp://0.0.0.0:36656 --rpc.laddr tcp://0.0.0.0:36657 --p2p.laddr tcp://0.0.0.0:36658 
echo "Main chain is started with rpc 36657.."
timeout /t 7 /nobreak

echo "Staring rest-servers"
echo "------------------------------------------------------------------------------------------"
start maincli rest-server --node tcp://0.0.0.0:36657 --laddr tcp://0.0.0.0:31118 --chain-id comdex-main
timeout /t 2 /nobreak

echo "Staring blockExplorer"
echo "------------------------------------------------------------------------------------------"
start blockExplorer
timeout /t 2 /nobreak