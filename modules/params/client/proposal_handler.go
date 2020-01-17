package client

import (
	govclient "github.com/persistenceOne/persistenceSDK/modules/gov/client"
	"github.com/persistenceOne/persistenceSDK/modules/params/client/cli"
	"github.com/persistenceOne/persistenceSDK/modules/params/client/rest"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
