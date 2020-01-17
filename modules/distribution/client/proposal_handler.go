package client

import (
	"github.com/persistenceOne/persistenceSDK/modules/distribution/client/cli"
	"github.com/persistenceOne/persistenceSDK/modules/distribution/client/rest"
	govclient "github.com/persistenceOne/persistenceSDK/modules/gov/client"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
