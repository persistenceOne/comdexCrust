#!/bin/bash
sh ./cleanup.sh

if [ ! -e  ${GOBIN}/maind ] || [ ! -e ${GOBIN}/maincli ] || 
	[ ! -e  ${GOBIN}/assetd ] || [ ! -e ${GOBIN}/assetcli ] || 
		[ ! -e  ${GOBIN}/fiatd ] || [ ! -e ${GOBIN}/fiatcli ]; then
			make all
fi 
echo "Intiating the chains ..."
maind init --name main --chain-id comdex-main --seed "pudding inform torch tourist cherry rebuild quarter latin flip grid fit clip label gallery nominee drift exit blast company student gather shoe velvet post"
sleep 2s

assetd init --name asset --chain-id comdex-asset
sleep 2s

fiatd init --name fiat --chain-id comdex-fiat
sleep 2s
echo "chains are initated."

echo "starting the chains.."

maind start --address tcp://0.0.0.0:36656 --rpc.laddr tcp://0.0.0.0:36657 --p2p.laddr tcp://0.0.0.0:36658  &
echo "Main chain is started with rpc 36657.."
sleep 3s

assetd start --address tcp://0.0.0.0:46656 --rpc.laddr tcp://0.0.0.0:46657 --p2p.laddr tcp://0.0.0.0:46658 &
echo "Asset chain is started with rpc 46657.."
sleep 3s

fiatd start --address tcp://0.0.0.0:56656 --rpc.laddr tcp://0.0.0.0:56657 --p2p.laddr tcp://0.0.0.0:56658 &
echo "Fiat chain is started with rpc 56657.."
sleep 3s


echo "Staring rest-servers"

maincli rest-server --node tcp://0.0.0.0:36657 --laddr tcp://0.0.0.0:31118 --chain-id comdex-main &
sleep 2s

assetcli rest-server --node tcp://0.0.0.0:46657 --laddr tcp://0.0.0.0:41118 --chain-id comdex-asset &
sleep 2s

fiatcli rest-server --node tcp://0.0.0.0:56657 --laddr tcp://0.0.0.0:51118  --chain-id comdex-fiat &
sleep 2	s

${GOBIN}/blockExplorer
echo "Block Explorer is started at port 2259..."
sleep 2s

trap 'killall $BGPID; exit' SIGINT
sleep 1024 &   
BGPID=${!}
sleep 1024  
