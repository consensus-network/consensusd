package main

import (
	"context"
	"fmt"

	"github.com/consensus-network/consensusd/cmd/consensuswallet/daemon/client"
	"github.com/consensus-network/consensusd/cmd/consensuswallet/daemon/pb"
	"github.com/consensus-network/consensusd/cmd/consensuswallet/utils"
)

func balance(conf *balanceConfig) error {
	daemonClient, tearDown, err := client.Connect(conf.DaemonAddress)
	if err != nil {
		return err
	}
	defer tearDown()

	ctx, cancel := context.WithTimeout(context.Background(), daemonTimeout)
	defer cancel()
	response, err := daemonClient.GetBalance(ctx, &pb.GetBalanceRequest{})
	if err != nil {
		return err
	}

	pendingSuffix := ""
	if response.Pending > 0 {
		pendingSuffix = " (pending)"
	}
	if conf.Verbose {
		pendingSuffix = ""
		println("Address                                                                       Available             Pending")
		println("-----------------------------------------------------------------------------------------------------------")
		for _, addressBalance := range response.AddressBalances {
			fmt.Printf("%s %s %s\n", addressBalance.Address, utils.FormatCss(addressBalance.Available), utils.FormatCss(addressBalance.Pending))
		}
		println("-----------------------------------------------------------------------------------------------------------")
		print("                                                 ")
	}
	fmt.Printf("Total balance, CSS %s %s%s\n", utils.FormatCss(response.Available), utils.FormatCss(response.Pending), pendingSuffix)

	return nil
}
