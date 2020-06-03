package client

import (
	govclient "github.com/commitHub/commitBlockchain/modules/gov/client"
	"github.com/commitHub/commitBlockchain/modules/params/client/cli"
	"github.com/commitHub/commitBlockchain/modules/params/client/rest"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
