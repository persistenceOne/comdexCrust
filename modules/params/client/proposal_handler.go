package client

import (
	govclient "github.com/persistenceOne/comdexCrust/modules/gov/client"
	"github.com/persistenceOne/comdexCrust/modules/params/client/cli"
	"github.com/persistenceOne/comdexCrust/modules/params/client/rest"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
