package protowire

import (
	"github.com/consensus-network/consensusd/app/appmessage"
	"github.com/pkg/errors"
)

func (x *ConsensusdMessage_Verack) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "ConsensusdMessage_Verack is nil")
	}
	return &appmessage.MsgVerAck{}, nil
}

func (x *ConsensusdMessage_Verack) fromAppMessage(_ *appmessage.MsgVerAck) error {
	return nil
}
