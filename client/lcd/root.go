package lcd

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
	
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	rpcserver "github.com/tendermint/tendermint/rpc/lib/server"
	
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	
	keys2 "github.com/commitHub/commitBlockchain/client/keys"
	
	"github.com/commitHub/commitBlockchain/client/rest"
	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/kafka"
	"github.com/commitHub/commitBlockchain/main/app"
	
	keybase "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"
	// unnamed import of statik for swagger UI support
	// _ "github.com/cosmos/cosmos-sdk/client/lcd/statik"
)

// RestServer represents the Light Client Rest server
type RestServer struct {
	Mux     *mux.Router
	CliCtx  context.CLIContext
	KeyBase keybase.Keybase
	
	log      log.Logger
	listener net.Listener
}

// NewRestServer creates a new rest server instance
func NewRestServer(cdc *codec.Codec) *RestServer {
	r := mux.NewRouter()
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "rest-server")
	
	return &RestServer{
		Mux:    r,
		CliCtx: cliCtx,
		log:    logger,
	}
}

// Start starts the rest server
func (rs *RestServer) Start(listenAddr string, maxOpen int, readTimeout, writeTimeout uint, kafkaBool bool, kafkaState kafka.KafkaState) (err error) {
	server.TrapSignal(func() {
		if kafkaBool == true {
			err = kafkaState.Producer.Close()
			
			err = kafkaState.Admin.Close()
			for _, consumer := range kafkaState.Consumers {
				err = consumer.Close()
			}
			err = kafkaState.Consumer.Close()
		}
		err := rs.listener.Close()
		rs.log.Error("error closing listener", "err", err)
		
	})
	
	cfg := &rpcserver.Config{
		MaxOpenConnections: maxOpen,
		ReadTimeout:        time.Duration(readTimeout) * time.Second,
		WriteTimeout:       time.Duration(writeTimeout) * time.Second,
	}
	
	rs.listener, err = rpcserver.Listen(listenAddr, cfg)
	if err != nil {
		return
	}
	rs.log.Info(
		fmt.Sprintf(
			"Starting application REST service (chain-id: %q)...",
			viper.GetString(flags.FlagChainID),
		),
	)
	
	return rpcserver.StartHTTPServer(rs.listener, rs.Mux, rs.log, cfg)
}

// ServeCommand will start the application REST service as a blocking process. It
// takes a codec to create a RestServer object and a function to register all
// necessary routes.
func ServeCommand(cdc *codec.Codec) *cobra.Command {
	flagKafka := "kafka"
	kafkaPorts := "kafkaPort"
	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rs := NewRestServer(cdc)
			
			kafkaBool := viper.GetBool(flagKafka)
			
			var kafkaState kafka.KafkaState
			
			if kafkaBool == true {
				kafkaPort := viper.GetString(kafkaPorts)
				kafkaPort = strings.Trim(kafkaPort, "\" ")
				kafkaPorts := strings.Split(kafkaPort, " ")
				kafkaState = kafka.NewKafkaState(kafkaPorts)
				rs.Mux.HandleFunc("/response/{ticketid}", kafka.QueryDB(cdc, rs.Mux, kafkaState.KafkaDB)).Methods("GET")
			}
			registerRoutes(rs, kafkaBool, kafkaState)
			
			if kafkaBool == true {
				go func() {
					for {
						rest.KafkaConsumerMsgs(rs.CliCtx, kafkaState)
						time.Sleep(kafka.SleepRoutine)
					}
				}()
				
			}
			
			// Start the rest server and return error if one exists
			err = rs.Start(
				viper.GetString(flags.FlagListenAddr),
				viper.GetInt(flags.FlagMaxOpenConnections),
				uint(viper.GetInt(flags.FlagRPCReadTimeout)),
				uint(viper.GetInt(flags.FlagRPCWriteTimeout)),
				kafkaBool,
				kafkaState,
			)
			
			return err
		},
	}
	cmd.Flags().Bool(flagKafka, false, "Whether have kafka running")
	cmd.Flags().String(kafkaPorts, "localhost:9092", "Space seperated addresses in quotes of the kafka listening node: example: --kafkaPort \"addr1 addr2\" ")
	
	return flags.RegisterRestServerFlags(cmd)
}

func registerRoutes(rs *RestServer, kafkaBool bool, kafkaState kafka.KafkaState) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux, kafkaBool, kafkaState)
	keys2.RegisterRoutes(rs.Mux)
	
}

//
// func (rs *RestServer) registerSwaggerUI() {
// 	statikFS, err := fs.New()
// 	if err != nil {
// 		panic(err)
// 	}
// 	staticServer := http.FileServer(statikFS)
// 	rs.Mux.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", staticServer))
// }
