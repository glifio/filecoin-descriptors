package main

import (
	"bytes"
	"context"
	"net/http"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	systemActor "github.com/filecoin-project/go-state-types/builtin/v8/system"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
)

type Lotus struct {
	api       api.FullNodeStruct
	rpcCloser jsonrpc.ClientCloser
}

func (l *Lotus) Open(url string) error {
	var err error
	l.rpcCloser, err = jsonrpc.NewMergeClient(context.Background(),
		url,
		"Filecoin",
		api.GetInternalStructs(&l.api),
		http.Header{})
	return err
}

func (l *Lotus) Close() {
	l.rpcCloser()
}

func (l *Lotus) GetActorCodeMap() (ActorCodeMap, error) {
	addr, err := address.NewFromString("f00")
	if err != nil {
		return nil, err
	}

	actor, err := l.api.StateReadState(context.Background(), addr, types.EmptyTSK)
	if err != nil {
		return nil, err
	}

	var state systemActor.State
	err = MapToInterface(actor.State, &state)
	if err != nil {
		return nil, err
	}

	object, err := l.api.ChainReadObj(context.Background(), state.BuiltinActors)
	if err != nil {
		return nil, err
	}

	var data manifest.ManifestData
	err = data.UnmarshalCBOR(bytes.NewReader(object))
	if err != nil {
		return nil, err
	}

	var actorCodeMap = ActorCodeMap{}
	for _, entry := range data.Entries {
		actorCodeMap[entry.Name] = entry.Code.String()
	}

	return actorCodeMap, nil
}
