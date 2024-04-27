package protowire

import (
	"github.com/consensus-network/consensusd/app/appmessage"
	"github.com/pkg/errors"
)

func (x *ConsensusdMessage_Ready) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "ConsensusdMessage_Ready is nil")
	}
	return &appmessage.MsgReady{}, nil
}

func (x *ConsensusdMessage_Ready) fromAppMessage(_ *appmessage.MsgReady) error {
	return nil
}
