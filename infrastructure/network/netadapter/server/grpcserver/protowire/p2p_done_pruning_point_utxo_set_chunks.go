package protowire

import (
	"github.com/consensus-network/consensusd/app/appmessage"
	"github.com/pkg/errors"
)

func (x *ConsensusdMessage_DonePruningPointUtxoSetChunks) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "ConsensusdMessage_DonePruningPointUtxoSetChunks is nil")
	}
	return &appmessage.MsgDonePruningPointUTXOSetChunks{}, nil
}

func (x *ConsensusdMessage_DonePruningPointUtxoSetChunks) fromAppMessage(_ *appmessage.MsgDonePruningPointUTXOSetChunks) error {
	x.DonePruningPointUtxoSetChunks = &DonePruningPointUtxoSetChunksMessage{}
	return nil
}
