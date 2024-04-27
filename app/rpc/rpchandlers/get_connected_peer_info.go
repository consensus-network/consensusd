package rpchandlers

import (
	"github.com/consensus-network/consensusd/app/appmessage"
	"github.com/consensus-network/consensusd/app/rpc/rpccontext"
	"github.com/consensus-network/consensusd/infrastructure/network/netadapter/router"
)

// HandleGetConnectedPeerInfo handles the respectively named RPC command
func HandleGetConnectedPeerInfo(context *rpccontext.Context, _ *router.Router, _ appmessage.Message) (appmessage.Message, error) {
	peers := context.ProtocolManager.Peers()
	ibdPeer := context.ProtocolManager.IBDPeer()
	infos := make([]*appmessage.GetConnectedPeerInfoMessage, 0, len(peers))
	for _, peer := range peers {
		info := &appmessage.GetConnectedPeerInfoMessage{
			ID:                        peer.ID().String(),
			Address:                   peer.Address(),
			LastPingDuration:          peer.LastPingDuration().Milliseconds(),
			IsOutbound:                peer.IsOutbound(),
			TimeOffset:                peer.TimeOffset().Milliseconds(),
			UserAgent:                 peer.UserAgent(),
			AdvertisedProtocolVersion: peer.AdvertisedProtocolVersion(),
			TimeConnected:             peer.TimeConnected().Milliseconds(),
			IsIBDPeer:                 peer == ibdPeer,
		}
		infos = append(infos, info)
	}
	response := appmessage.NewGetConnectedPeerInfoResponseMessage(infos)
	return response, nil
}
