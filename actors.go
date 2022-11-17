package main

import (
	accountState "github.com/filecoin-project/go-state-types/builtin/v8/account"
	cronState "github.com/filecoin-project/go-state-types/builtin/v8/cron"
	eamState "github.com/filecoin-project/go-state-types/builtin/v8/eam"
	initState "github.com/filecoin-project/go-state-types/builtin/v8/init"
	marketState "github.com/filecoin-project/go-state-types/builtin/v8/market"
	minerState "github.com/filecoin-project/go-state-types/builtin/v8/miner"
	multisigState "github.com/filecoin-project/go-state-types/builtin/v8/multisig"
	paychState "github.com/filecoin-project/go-state-types/builtin/v8/paych"
	powerState "github.com/filecoin-project/go-state-types/builtin/v8/power"
	rewardState "github.com/filecoin-project/go-state-types/builtin/v8/reward"
	systemState "github.com/filecoin-project/go-state-types/builtin/v8/system"
	verifregState "github.com/filecoin-project/go-state-types/builtin/v8/verifreg"
	accountActor "github.com/filecoin-project/specs-actors/v8/actors/builtin/account"
	cronActor "github.com/filecoin-project/specs-actors/v8/actors/builtin/cron"
	initActor "github.com/filecoin-project/specs-actors/v8/actors/builtin/init"
	marketActor "github.com/filecoin-project/specs-actors/v8/actors/builtin/market"
	minerActor "github.com/filecoin-project/specs-actors/v8/actors/builtin/miner"
	multisigActor "github.com/filecoin-project/specs-actors/v8/actors/builtin/multisig"
	paychActor "github.com/filecoin-project/specs-actors/v8/actors/builtin/paych"
	powerActor "github.com/filecoin-project/specs-actors/v8/actors/builtin/power"
	rewardActor "github.com/filecoin-project/specs-actors/v8/actors/builtin/reward"
	verifregActor "github.com/filecoin-project/specs-actors/v8/actors/builtin/verifreg"
)

type ReflectableActor struct {
	State   interface{}
	Methods map[uint64]interface{}
}

type CustomMethod struct {
	Name   string
	Param  interface{}
	Return interface{}
}

var reflectableActors = map[ActorName]ReflectableActor{
	"account": {
		State: (*accountState.State)(nil),
		Methods: map[uint64]interface{}{
			1: accountActor.Actor.Constructor,
			2: accountActor.Actor.PubkeyAddress,
		},
	},
	"cron": {
		State: (*cronState.State)(nil),
		Methods: map[uint64]interface{}{
			1: cronActor.Actor.Constructor,
			2: cronActor.Actor.EpochTick,
		},
	},
	"eam": {
		State: nil,
		Methods: map[uint64]interface{}{
			2: CustomMethod{
				Name:   "Create",
				Param:  (*eamState.CreateParams)(nil),
				Return: (*eamState.CreateReturn)(nil),
			},
			3: CustomMethod{
				Name:   "Create2",
				Param:  (*eamState.Create2Params)(nil),
				Return: (*eamState.Create2Return)(nil),
			},
		},
	},
	"init": {
		State: (*initState.State)(nil),
		Methods: map[uint64]interface{}{
			1: initActor.Actor.Constructor,
			2: initActor.Actor.Exec,
		},
	},
	"multisig": {
		State: (*multisigState.State)(nil),
		Methods: map[uint64]interface{}{
			1: multisigActor.Actor.Constructor,
			2: multisigActor.Actor.Propose,
			3: multisigActor.Actor.Approve,
			4: multisigActor.Actor.Cancel,
			5: multisigActor.Actor.AddSigner,
			6: multisigActor.Actor.RemoveSigner,
			7: multisigActor.Actor.SwapSigner,
			8: multisigActor.Actor.ChangeNumApprovalsThreshold,
			9: multisigActor.Actor.LockBalance,
		},
	},
	"paymentchannel": {
		State: (*paychState.State)(nil),
		Methods: map[uint64]interface{}{
			1: (*paychActor.Actor).Constructor,
			2: paychActor.Actor.UpdateChannelState,
			3: paychActor.Actor.Settle,
			4: paychActor.Actor.Collect,
		},
	},
	"reward": {
		State: (*rewardState.State)(nil),
		Methods: map[uint64]interface{}{
			1: rewardActor.Actor.Constructor,
			2: rewardActor.Actor.AwardBlockReward,
			3: rewardActor.Actor.ThisEpochReward,
			4: rewardActor.Actor.UpdateNetworkKPI,
		},
	},
	"storagemarket": {
		State: (*marketState.State)(nil),
		Methods: map[uint64]interface{}{
			1: marketActor.Actor.Constructor,
			2: marketActor.Actor.AddBalance,
			3: marketActor.Actor.WithdrawBalance,
			4: marketActor.Actor.PublishStorageDeals,
			5: marketActor.Actor.VerifyDealsForActivation,
			6: marketActor.Actor.ActivateDeals,
			7: marketActor.Actor.OnMinerSectorsTerminate,
			8: marketActor.Actor.ComputeDataCommitment,
			9: marketActor.Actor.CronTick,
		},
	},
	"storageminer": {
		State: (*minerState.State)(nil),
		Methods: map[uint64]interface{}{
			1:  minerActor.Actor.Constructor,
			2:  minerActor.Actor.ControlAddresses,
			3:  minerActor.Actor.ChangeWorkerAddress,
			4:  minerActor.Actor.ChangePeerID,
			5:  minerActor.Actor.SubmitWindowedPoSt,
			6:  minerActor.Actor.PreCommitSector,
			7:  minerActor.Actor.ProveCommitSector,
			8:  minerActor.Actor.ExtendSectorExpiration,
			9:  minerActor.Actor.TerminateSectors,
			10: minerActor.Actor.DeclareFaults,
			11: minerActor.Actor.DeclareFaultsRecovered,
			12: minerActor.Actor.OnDeferredCronEvent,
			13: minerActor.Actor.CheckSectorProven,
			14: minerActor.Actor.ApplyRewards,
			15: minerActor.Actor.ReportConsensusFault,
			16: minerActor.Actor.WithdrawBalance,
			17: minerActor.Actor.ConfirmSectorProofsValid,
			18: minerActor.Actor.ChangeMultiaddrs,
			19: minerActor.Actor.CompactPartitions,
			20: minerActor.Actor.CompactSectorNumbers,
			21: minerActor.Actor.ConfirmUpdateWorkerKey,
			22: minerActor.Actor.RepayDebt,
			23: minerActor.Actor.ChangeOwnerAddress,
			24: minerActor.Actor.DisputeWindowedPoSt,
			25: minerActor.Actor.PreCommitSectorBatch,
			26: minerActor.Actor.ProveCommitAggregate,
			27: minerActor.Actor.ProveReplicaUpdates,
		},
	},
	"storagepower": {
		State: (*powerState.State)(nil),
		Methods: map[uint64]interface{}{
			1: powerActor.Actor.Constructor,
			2: powerActor.Actor.CreateMiner,
			3: powerActor.Actor.UpdateClaimedPower,
			4: powerActor.Actor.EnrollCronEvent,
			5: powerActor.Actor.CronTick,
			6: powerActor.Actor.UpdatePledgeTotal,
			8: powerActor.Actor.SubmitPoRepForBulkVerify,
			9: powerActor.Actor.CurrentTotalPower,
		},
	},
	"system": {
		State:   (*systemState.State)(nil),
		Methods: map[uint64]interface{}{},
	},
	"verifiedregistry": {
		State: (*verifregState.State)(nil),
		Methods: map[uint64]interface{}{
			1: verifregActor.Actor.Constructor,
			2: verifregActor.Actor.AddVerifier,
			3: verifregActor.Actor.RemoveVerifier,
			4: verifregActor.Actor.AddVerifiedClient,
			5: verifregActor.Actor.UseBytes,
			6: verifregActor.Actor.RestoreBytes,
			7: verifregActor.Actor.RemoveVerifiedClientDataCap,
		},
	},
}
