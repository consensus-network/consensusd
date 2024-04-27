package rpchandlers

import (
	"github.com/consensus-network/consensusd/app/appmessage"
	"github.com/consensus-network/consensusd/app/rpc/rpccontext"
	"github.com/consensus-network/consensusd/infrastructure/network/netadapter/router"
)

// HandleGetHeaders handles the respectively named RPC command
func HandleGetHeaders(context *rpccontext.Context, _ *router.Router, request appmessage.Message) (appmessage.Message, error) {
	response := &appmessage.GetHeadersResponseMessage{}
	response.Error = appmessage.RPCErrorf("not implemented")
	return response, nil
}
