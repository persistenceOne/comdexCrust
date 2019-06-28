package rest

import (
	"fmt"
	"net/http"
	
	"github.com/gorilla/mux"
	
	"github.com/comdex-blockchain/wire"
	dbm "github.com/tendermint/tendermint/libs/db"
)

// SetTicketIDtoDB : initiates ticketid in Database
func SetTicketIDtoDB(ticketID Ticket, kafkaDB *dbm.GoLevelDB, cdc *wire.Codec, msg []byte) {
	
	ticketid, err := cdc.MarshalJSON(ticketID)
	if err != nil {
		panic(err)
	}
	
	kafkaDB.Set(ticketid, msg)
	return
}

// AddResponseToDB : Updates response to DB
func AddResponseToDB(ticketID Ticket, response []byte, kafkaDB *dbm.GoLevelDB, cdc *wire.Codec) {
	ticketid, err := cdc.MarshalJSON(ticketID)
	if err != nil {
		panic(err)
	}
	
	kafkaDB.SetSync(ticketid, response)
	return
}

// GetResponseFromDB : gives the response from DB
func GetResponseFromDB(ticketID Ticket, kafkaDB *dbm.GoLevelDB, cdc *wire.Codec) []byte {
	ticketid, err := cdc.MarshalJSON(ticketID)
	if err != nil {
		panic(err)
	}
	
	return kafkaDB.Get(ticketid)
}

// QueryDB : REST outputs info from DB
func QueryDB(cdc *wire.Codec, r *mux.Router, kafkaDB *dbm.GoLevelDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		
		iDByte, err := cdc.MarshalJSON(vars["ticketid"])
		if err != nil {
			panic(err)
		}
		var response []byte
		if kafkaDB.Has(iDByte) == true {
			response = GetResponseFromDB(Ticket(vars["ticketid"]), kafkaDB, cdc)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			output, err := cdc.MarshalJSON("The ticket ID does not exist, it must have been deleted, Query the chain to know")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("ticket ID does not exist. Error: %s", err.Error())))
				return
			}
			w.Write(output)
			return
		}
		
		w.WriteHeader(http.StatusAccepted)
		w.Write(response)
	}
}
