package lcd

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/commitHub/commitBlockchain/client"
	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/client/keys"
	"github.com/commitHub/commitBlockchain/client/rpc"
	"github.com/commitHub/commitBlockchain/client/tx"
	"github.com/commitHub/commitBlockchain/client/utils"
	"github.com/commitHub/commitBlockchain/rest"
	"github.com/commitHub/commitBlockchain/wire"
	acl "github.com/commitHub/commitBlockchain/x/acl/client/rest"
	assetFactory "github.com/commitHub/commitBlockchain/x/assetFactory/client/rest"
	authcmd "github.com/commitHub/commitBlockchain/x/auth/client/cli"
	auth "github.com/commitHub/commitBlockchain/x/auth/client/rest"
	bank "github.com/commitHub/commitBlockchain/x/bank/client/rest"
	fiatFactory "github.com/commitHub/commitBlockchain/x/fiatFactory/client/rest"
	gov "github.com/commitHub/commitBlockchain/x/gov/client/rest"
	ibc "github.com/commitHub/commitBlockchain/x/ibc/client/rest"
	negotiation "github.com/commitHub/commitBlockchain/x/negotiation/client/rest"
	order "github.com/commitHub/commitBlockchain/x/order/client/rest"
	reputation "github.com/commitHub/commitBlockchain/x/reputation/client/rest"
	slashing "github.com/commitHub/commitBlockchain/x/slashing/client/rest"
	stake "github.com/commitHub/commitBlockchain/x/stake/client/rest"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tmserver "github.com/tendermint/tendermint/rpc/lib/server"
)

// ServeCommand will generate a long-running rest server
// (aka Light Client Daemon) that exposes functionality similar
// to the cli, but over rest
func ServeCommand(cdc *wire.Codec) *cobra.Command {
	flagListenAddr := "laddr"
	flagCORS := "cors"
	flagMaxOpenConnections := "max-open"
	flagKafka := "kafka"
	kafkaPorts := "portKafka"
	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) error {
			listenAddr := viper.GetString(flagListenAddr)
			kafkaBool := viper.GetBool(flagKafka)
			kafkaPort := viper.GetString(kafkaPorts)
			var handler http.Handler
			var kafkaState rest.KafkaState

			handler, kafkaState = createHandlerRest(cdc, kafkaPort, kafkaBool)

			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "rest-server")
			maxOpen := viper.GetInt(flagMaxOpenConnections)
			listener, err := tmserver.StartHTTPServer(
				listenAddr, handler, logger,
				tmserver.Config{MaxOpenConnections: maxOpen},
			)
			if err != nil {
				return err
			}
			logger.Info("REST server started")

			// wait forever and cleanup
			cmn.TrapSignal(func() {

				err := listener.Close()
				logger.Error("error closing listener", "err", err)
				if kafkaBool == true {
					kafkaState.Producer.Close()
					kafkaState.Admin.Close()
					for _, consumer := range kafkaState.Consumers {
						consumer.Close()
					}
					kafkaState.Consumer.Close()
				}
			})

			return nil
		},
	}

	cmd.Flags().String(flagListenAddr, "tcp://localhost:1317", "The address for the server to listen on")
	cmd.Flags().String(flagCORS, "", "Set the domains that can make CORS requests (* for all)")
	cmd.Flags().String(client.FlagChainID, "", "The chain ID to connect to")
	cmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "Address of the node to connect to")
	cmd.Flags().Int(flagMaxOpenConnections, 1000, "The number of maximum open connections")
	cmd.Flags().Bool(client.FlagTrustNode, false, "Whether trust connected full node")
	cmd.Flags().Bool(flagKafka, false, "Whether have kafka running")
	cmd.Flags().String(kafkaPorts, "localhost:9092", "Space seperated addresses in quotes of the kafka listening node: example: --kafkaPort \"addr1 addr2\" ")

	return cmd
}

func createHandlerRest(cdc *wire.Codec, kafkaPort string, kafka bool) (http.Handler, rest.KafkaState) {
	r := mux.NewRouter()

	kb, err := keys.GetKeyBase()
	if err != nil {
		panic(err)
	}

	cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

	var kafkaState rest.KafkaState
	if kafka == true {
		kafkaPort = strings.Trim(kafkaPort, "\" ")
		kafkaPorts := strings.Split(kafkaPort, " ")
		kafkaState = rest.NewKafkaState(kafkaPorts)
	}

	// TODO: make more functional? aka r = keys.RegisterRouteTransferRequestHandlers(r)
	r.HandleFunc("/version", CLIVersionRequestHandler).Methods("GET")
	r.HandleFunc("/node_version", NodeVersionRequestHandler(cliCtx)).Methods("GET")
	if kafka == true {
		r.HandleFunc("/response/{ticketid}", rest.QueryDB(cdc, r, kafkaState.KafkaDB)).Methods("GET")
	}

	keys.RegisterRoutes(r)
	rpc.RegisterRoutes(cliCtx, r)
	tx.RegisterRoutes(cliCtx, r, cdc)
	auth.RegisterRoutes(cliCtx, r, cdc, "acc")
	bank.RegisterRoutes(cliCtx, r, cdc, kb, kafka, kafkaState, "asset", "fiat")
	ibc.RegisterRoutes(cliCtx, r, cdc, kb, kafka, kafkaState)
	stake.RegisterRoutes(cliCtx, r, cdc, kb)
	slashing.RegisterRoutes(cliCtx, r, cdc, kb)
	gov.RegisterRoutes(cliCtx, r, cdc)
	acl.RegisterRoutes(cliCtx, r, cdc, "acl")
	assetFactory.RegisterRoutes(cliCtx, r, cdc, "asset", kb, kafka, kafkaState)
	fiatFactory.RegisterRoutes(cliCtx, r, cdc, "fiat", kb, kafka, kafkaState)
	negotiation.RegisterRoutes(cliCtx, r, cdc, "negotiation", kb, kafka, kafkaState)
	order.RegisterRoutes(cliCtx, r, cdc, "order", kb)
	reputation.RegisterRoutes(cliCtx, r, cdc, "reputation", kb, kafka, kafkaState)

	if kafka == true {
		go func() {
			for {
				utils.KafkaConsumerMsgs(cliCtx, cdc, kafkaState)
				time.Sleep(rest.SleepRoutine)
			}
		}()

	}

	return r, kafkaState
}
