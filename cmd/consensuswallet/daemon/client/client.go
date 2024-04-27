package client

import (
	"context"
	"time"

	"github.com/consensus-network/consensusd/cmd/consensuswallet/daemon/server"

	"github.com/pkg/errors"

	"github.com/consensus-network/consensusd/cmd/consensuswallet/daemon/pb"
	"google.golang.org/grpc"
)

// Connect connects to the consensuswalletd server, and returns the client instance
func Connect(address string) (pb.ConsensuswalletdClient, func(), error) {
	// Connection is local, so 1 second timeout is sufficient
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(server.MaxDaemonSendMsgSize)))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, nil, errors.New("consensuswallet daemon is not running, start it with `consensuswallet start-daemon`")
		}
		return nil, nil, err
	}

	return pb.NewConsensuswalletdClient(conn), func() {
		conn.Close()
	}, nil
}
