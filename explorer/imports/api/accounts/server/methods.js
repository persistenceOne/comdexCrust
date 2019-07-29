import { Meteor } from 'meteor/meteor';
import { HTTP } from 'meteor/http';
import { Validators } from '/imports/api/validators/validators.js';
const fetchFromUrl = (url) => {
    try{
        let res = HTTP.get(LCD + url);
        if (res.statusCode == 200){
            return res
        };
    }
    catch (e){
        console.log(e);
    }
}

Meteor.methods({
    'accounts.getAccountDetail': function(address){
        this.unblock();
        let url = LCD + '/auth/accounts/'+ address;
        try{
            let available = HTTP.get(url);
            if (available.statusCode == 200){
                let response = JSON.parse(available.content);
                let account;
                if (response.type === 'auth/Account')
                    account = response.value;
                else if (response.type === 'auth/DelayedVestingAccount' || response.type === 'auth/ContinuousVestingAccount')
                    account = response.value.BaseVestingAccount.BaseAccount
                if (account && account.account_number != null)
                    return account
                return null
            }
        }
        catch (e){
            console.log(e)
        }
    },
    'accounts.getBalance': function(address){
        this.unblock();
        let balance = {}

        // get available atoms
        let url = LCD + '/accounts/'+ address;
        try{
            let available = HTTP.get(url);
            if (available.statusCode == 200){
                data = JSON.parse(available.content);
                coins = data.value.coins
                balance.available = coins.filter(function(value){
                    return value.denom == Meteor.settings.public.stakingDenom.toLowerCase();
                })
                balance.available = parseInt(balance.available[0].amount)
            }
        }
        catch (e){
            console.log(e)
        }

        // get delegated amnounts
        url = LCD + '/stake/delegators/'+address;
        try{
            let delegations = HTTP.get(url);
            if (delegations.statusCode == 200){
                let shares = 0 
                data = JSON.parse(delegations.content)

                data.delegations.forEach(function(value){
                    shares += parseInt(value.shares)
                    
                })
                
                balance.delegations = data.delegations;
                
            }
        }
        catch (e){
            console.log(e);
        }
        // get unbonding
        url = LCD + '/stake/delegators/'+address;
        try{
            let unbonding = HTTP.get(url);
            if (unbonding.statusCode == 200){
                let shares = 0
                data = JSON.parse(unbonding.content)
                balance.unbonding_delegations = data.unbonding_delegations;
            }
        }
        catch (e){
            console.log(e);
        }
        return balance;
    },
    'accounts.getDelegation'(address, validator){
        let url = `/stake/delegators/${address}/delegations/${validator}`;
        let delegations = fetchFromUrl(url);
        delegations = delegations && delegations.data;
        if (delegations && delegations.shares)
            delegations.shares = parseFloat(delegations.shares);

        url = `/stake/redelegations?delegator=${address}&validator_to=${validator}`;
        let relegations = fetchFromUrl(url).data;
        relegations = relegations && relegations.data;
        let completionTime;
        if (relegations) {
            relegations.forEach((relegation) => {
                let entries = relegation.entries
                let time = new Date(entries[entries.length-1].completion_time)
                if (!completionTime || time > completionTime)
                    completionTime = time
            })
            delegations.redelegationCompletionTime = completionTime;
        }

        url = `/stake/delegators/${address}/unbonding_delegations/${validator}`;
        let undelegations = fetchFromUrl(url);
        undelegations = undelegations && undelegations.data;
        if (undelegations) {
            delegations.unbonding = undelegations.entries.length;
            delegations.unbondingCompletionTime = undelegations.entries[0].completion_time;
        }
        return delegations;
    },
    'accounts.getAllDelegations'(address){
        let url = LCD + '/stake/delegators/'+address;

        try{
            let delegations = HTTP.get(url);
            if (delegations.statusCode == 200){
                delegations = JSON.parse(delegations.content);
                
                delegations = delegations.delegations
                if (delegations && delegations.length > 0){
                    delegations.forEach((delegation, i) => {
                        if (delegations[i] && delegations[i].shares)
                            delegations[i].shares = parseFloat(delegations[i].shares);
                    })
                }

                return delegations;
            };
        }
        catch (e){
            console.log(e);
        }
    },
    'accounts.getAllUnbondings'(address){
        let url = LCD + '/stake/delegators/'+address;

        try{
            let unbondings = HTTP.get(url);
            if (unbondings.statusCode == 200){
                let unbondingTokens = 0
                data = JSON.parse(unbondings.content)
                if (data.unbonding_delegations && data.unbonding_delegations.length >0){
                    data.unbonding_delegations.forEach(function(value){
                        unbondingTokens += value.shares
                        
                    })    
                }
                
                return unbondingTokens;
            };
        }
        catch (e){
            console.log(e);
        }
    },
    'accounts.getAllRedelegations'(address, validator){
        let url = `/stake/redelegations?delegator=${address}&validator_from=${validator}`;
        let result = fetchFromUrl(url);
        if (result && result.data) {
            let redelegations = {}
            result.data.forEach((redelegation) => {
                let entries = redelegation.entries;
                redelegations[redelegation.validator_dst_address] = {
                    count: entries.length,
                    completionTime: entries[0].completion_time
                }
            })
            return redelegations
        }
    }
})