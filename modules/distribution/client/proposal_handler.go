package client

import (
	"github.com/persistenceOne/comdexCrust/modules/distribution/client/cli"
	"github.com/persistenceOne/comdexCrust/modules/distribution/client/rest"
	govclient "github.com/persistenceOne/comdexCrust/modules/gov/client"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
