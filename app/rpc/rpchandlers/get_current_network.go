package rpchandlers

import (
	"github.com/consensus-network/consensusd/app/appmessage"
	"github.com/consensus-network/consensusd/app/rpc/rpccontext"
	"github.com/consensus-network/consensusd/infrastructure/network/netadapter/router"
)

// HandleGetCurrentNetwork handles the respectively named RPC command
func HandleGetCurrentNetwork(context *rpccontext.Context, _ *router.Router, _ appmessage.Message) (appmessage.Message, error) {
	response := appmessage.NewGetCurrentNetworkResponseMessage(context.Config.ActiveNetParams.Net.String())
	return response, nil
}
