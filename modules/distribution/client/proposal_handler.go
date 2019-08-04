package client

import (
	"github.com/commitHub/commitBlockchain/modules/distribution/client/cli"
	"github.com/commitHub/commitBlockchain/modules/distribution/client/rest"
	govclient "github.com/commitHub/commitBlockchain/modules/gov/client"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
