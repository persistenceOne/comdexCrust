import { Meteor } from 'meteor/meteor';
import { HTTP } from 'meteor/http';
import { getAddress } from 'tendermint/lib/pubkey.js';
import { Chain, ChainStates } from '../chain.js';
import { Validators } from '../../validators/validators.js';
import { VotingPowerHistory } from '../../voting-power/history.js';

findVotingPower = (validator, genValidators) => {
    for (let v in genValidators){
        if (validator.pub_key.value == genValidators[v].pub_key.value){
            return parseInt(genValidators[v].power);
        }
    }
}

Meteor.methods({
    'chain.getConsensusState': function(){
        this.unblock();
        let url = RPC+'/dump_consensus_state';
        try{
            let response = HTTP.get(url);
            let consensus = JSON.parse(response.content);
            consensus = consensus.result;
            let height = consensus.round_state.height;
            let round = consensus.round_state.round;
            let step = consensus.round_state.step;
            let votedPower = Math.round(parseFloat(consensus.round_state.votes[round].prevotes_bit_array.split(" ")[3])*100);

            Chain.update({chainId:Meteor.settings.public.chainId}, {$set:{
                votingHeight: height, 
                votingRound: round, 
                votingStep: step, 
                votedPower: votedPower,
                proposerAddress: consensus.round_state.validators.proposer.address,
                prevotes: consensus.round_state.votes[round].prevotes,
                precommits: consensus.round_state.votes[round].precommits
            }});
        }
        catch(e){
            console.log(e);
        }
    },
    'chain.updateStatus': function(){
        this.unblock();
        let url = RPC+'/status';
        try{
            let response = HTTP.get(url);
            let status = JSON.parse(response.content);
            status = status.result;
            let chain = {};
            chain.chainId = status.node_info.network;
            chain.latestBlockHeight = status.sync_info.latest_block_height;
            chain.latestBlockTime = status.sync_info.latest_block_time;

            url = RPC+'/validators';
            response = HTTP.get(url);
            let validators = JSON.parse(response.content);
            validators = validators.result.validators;
            chain.validators = validators.length;
            let activeVP = 0;
            for (v in validators){
                activeVP += parseInt(validators[v].voting_power);
            }
            chain.activeVotingPower = activeVP;

            // Get chain states
            if (parseInt(chain.latestBlockHeight) > 0){
                let chainStates = {};
                chainStates.height = parseInt(status.sync_info.latest_block_height);
                chainStates.time = new Date(status.sync_info.latest_block_time);

                url = LCD + '/stake/pool';
                try{
                    response = HTTP.get(url);
                    let bonding = JSON.parse(response.content);
                    // chain.bondedTokens = bonding.bonded_tokens;
                    // chain.notBondedTokens = bonding.not_bonded_tokens;
                    chainStates.bondedTokens = parseInt(bonding.bonded_tokens);
                    chainStates.notBondedTokens = parseInt(bonding.loose_tokens);
                }
                catch(e){
                    console.log(e);
                }

                // url = LCD + '/distribution/community_pool';
                // try {
                //     response = HTTP.get(url);
                //     let pool = JSON.parse(response.content);
                //     if (pool && pool.length > 0){
                //         chainStates.communityPool = [];
                //         pool.forEach((amount, i) => {
                //             chainStates.communityPool.push({
                //                 denom: amount.denom,
                //                 amount: parseFloat(amount.amount)
                //             })
                //         })
                //     }
                // }
                // catch (e){
                //     console.log(e)
                // }

                // url = LCD + '/minting/inflation';
                // try{
                //     response = HTTP.get(url);
                //     let inflation = JSON.parse(response.content);
                //     if (inflation){
                //         chainStates.inflation = parseFloat(inflation)
                //     }
                // }
                // catch(e){
                //     console.log(e);
                // }

                // url = LCD + '/minting/annual-provisions';
                // try{
                //     response = HTTP.get(url);
                //     let provisions = JSON.parse(response.content);
                //     if (provisions){
                //         chainStates.annualProvisions = parseFloat(provisions)
                //     }
                // }
                // catch(e){
                //     console.log(e);
                // }

                ChainStates.insert(chainStates);
            }

            // chain.totalVotingPower = totalVP;

            Chain.update({chainId:chain.chainId}, {$set:chain}, {upsert: true});

            // validators = Validators.find({}).fetch();
            // console.log(validators);
            return chain.latestBlockHeight;
        }
        catch (e){
            console.log(e);
            return "Error getting chain status.";
        }
    },
    'chain.getLatestStatus': function(){
        Chain.find().sort({created:-1}).limit(1);
    },
    'chain.genesis': function(){
        let chain = Chain.findOne({chainId: Meteor.settings.public.chainId});
        
        if (chain && chain.readGenesis){
            console.log('Genesis file has been processed');
        }
        else{
            console.log('=== Start processing genesis file ===');
            let response = HTTP.get(Meteor.settings.genesisFile);
            let genesis = JSON.parse(response.content);
            genesis = genesis.result.genesis
            let chainParams = {
                chainId: genesis.chain_id,
                genesisTime: genesis.genesis_time,
                consensusParams: genesis.consensus_params,
                auth: genesis.app_state.auth,
                bank: genesis.app_state.bank,
                staking: {
                    pool: genesis.app_state.stake.pool,
                    params: genesis.app_state.stake.params
                },
              
                gov: {
                    startingProposalId: genesis.app_state.gov.starting_proposalID,
                    depositParams: genesis.app_state.gov.deposit_period,
                    votingParams: genesis.app_state.gov.voting_period,
                    tallyParams: genesis.app_state.gov.tallying_procedure
                }
                
            }

            let totalVotingPower = 0;
            // read validator in genesis
            if(genesis.app_state.stake.validators && (genesis.app_state.stake.validators.length >0)){
                for (i in genesis.app_state.stake.validators){
                    let validator = {
                        description: genesis.app_state.stake.validators[i].description,
                        commission: genesis.app_state.stake.validators[i].commission,
                        commission_max : genesis.app_state.stake.validators[i].commission_max,
                        commission_change_rate : genesis.app_state.stake.validators[i].commission_change_rate,
                        commission_change_today : genesis.app_state.stake.validators[i].commission_change_today,
                        operator_address: genesis.app_state.stake.validators[i].operator,
                        delegator_address: genesis.app_state.stake.bonds[i].delegator_addr,
                        delegator_shares :Math.floor(parseInt(genesis.app_state.stake.validators[i].delegator_shares) / Meteor.settings.public.stakingFraction),
                        voting_power: Math.floor(parseInt(genesis.app_state.stake.validators[i].delegator_shares) / Meteor.settings.public.stakingFraction),
                        jailed: false,
                        status: 2,

                    }
                    totalVotingPower += validator.voting_power;
                    validator.pub_key = {
                        "type":"tendermint/PubKeyEd25519",
                        "value":genesis.app_state.stake.validators[i].pub_key.value
                    }

                    validator.address = getAddress(validator.pub_key);
                    validator.accpub = Meteor.call('pubkeyToBech32', validator.pub_key, Meteor.settings.public.bech32PrefixAccPub);
                    validator.operator_pubkey = Meteor.call('pubkeyToBech32', validator.pub_key, Meteor.settings.public.bech32PrefixValPub);
                    validator.consensus_pubkey = Meteor.call('pubkeyToBech32', validator.pub_key, Meteor.settings.public.bech32PrefixConsPub);

                    VotingPowerHistory.insert({
                        address: validator.address,
                        prev_voting_power: 0,
                        voting_power: validator.voting_power,
                        type: 'add',
                        height: 0,
                        block_time: genesis.genesis_time
                    });
                    Validators.insert(validator);
                }
            }
                
            chainParams.readGenesis = true;
            chainParams.activeVotingPower = totalVotingPower;
            let result = Chain.upsert({chainId:chainParams.chainId}, {$set:chainParams});


            console.log('=== Finished processing genesis file ===');

        }
        
        return true;
    }
})