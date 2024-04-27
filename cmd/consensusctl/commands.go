package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/consensus-network/consensusd/infrastructure/network/netadapter/server/grpcserver/protowire"
)

var commandTypes = []reflect.Type{
	reflect.TypeOf(protowire.ConsensusdMessage_AddPeerRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_GetConnectedPeerInfoRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_GetPeerAddressesRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_GetCurrentNetworkRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_GetInfoRequest{}),

	reflect.TypeOf(protowire.ConsensusdMessage_GetBlockRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_GetBlocksRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_GetHeadersRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_GetBlockCountRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_GetBlockDagInfoRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_GetSelectedTipHashRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_GetVirtualSelectedParentBlueScoreRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_GetVirtualSelectedParentChainFromBlockRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_ResolveFinalityConflictRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_EstimateNetworkHashesPerSecondRequest{}),

	reflect.TypeOf(protowire.ConsensusdMessage_GetBlockTemplateRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_SubmitBlockRequest{}),

	reflect.TypeOf(protowire.ConsensusdMessage_GetMempoolEntryRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_GetMempoolEntriesRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_GetMempoolEntriesByAddressesRequest{}),

	reflect.TypeOf(protowire.ConsensusdMessage_SubmitTransactionRequest{}),

	reflect.TypeOf(protowire.ConsensusdMessage_GetUtxosByAddressesRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_GetBalanceByAddressRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_GetCoinSupplyRequest{}),

	reflect.TypeOf(protowire.ConsensusdMessage_BanRequest{}),
	reflect.TypeOf(protowire.ConsensusdMessage_UnbanRequest{}),
}

type commandDescription struct {
	name       string
	parameters []*parameterDescription
	typeof     reflect.Type
}

type parameterDescription struct {
	name   string
	typeof reflect.Type
}

func commandDescriptions() []*commandDescription {
	commandDescriptions := make([]*commandDescription, len(commandTypes))

	for i, commandTypeWrapped := range commandTypes {
		commandType := unwrapCommandType(commandTypeWrapped)

		name := strings.TrimSuffix(commandType.Name(), "RequestMessage")
		numFields := commandType.NumField()

		var parameters []*parameterDescription
		for i := 0; i < numFields; i++ {
			field := commandType.Field(i)

			if !isFieldExported(field) {
				continue
			}

			parameters = append(parameters, &parameterDescription{
				name:   field.Name,
				typeof: field.Type,
			})
		}
		commandDescriptions[i] = &commandDescription{
			name:       name,
			parameters: parameters,
			typeof:     commandTypeWrapped,
		}
	}

	return commandDescriptions
}

func (cd *commandDescription) help() string {
	sb := &strings.Builder{}
	sb.WriteString(cd.name)
	for _, parameter := range cd.parameters {
		_, _ = fmt.Fprintf(sb, " [%s]", parameter.name)
	}
	return sb.String()
}
