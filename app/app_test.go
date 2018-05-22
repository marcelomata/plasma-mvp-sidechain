package app

import (
	"os"
	"testing"
	//"fmt" //for debugging
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
	types "plasma-mvp-sidechain/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	//rlp "github.com/ethereum/go-ethereum/rlp"
)

func newChildChain() *ChildChain {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	return NewChildChain(logger, db)
}

func TestDepositMsg(t *testing.T) {
	cc := newChildChain()

	confirmSigs := [2]crypto.Signature{crypto.SignatureSecp256k1{}, crypto.SignatureSecp256k1{}}
	privKeyA := ethcrypto.GenerateKey()
	privKeyB := ethcrypto.GenerateKey()

	// Construct a SpendMsg
	var msg = types.SpendMsg{
		Blknum1:      0,
		Txindex1:     0,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       types.EthPrivKeyToSDKAddress(privKeyA),
		ConfirmSigs1: confirmSigs,
		Blknum2:      0,
		Txindex2:     0,
		Oindex2:      0,
		DepositNum2:  0,
		Owner2:       crypto.Address([]byte("")),
		ConfirmSigs2: confirmSigs,
		Newowner1:    types.EthPrivKeyToSDKAddress(privKeyB),
		Denom1:       1000,
		Newowner2:    crypto.Address([]byte("")),
		Denom2:       0,
		Fee:          1,
	}

	sig := ethcrypto.Sign(msg.GetSignBytes(), privKeyA).(crypto.Signature)
	tx := types.NewBaseTx(msg, []sdk.StdSignature{{
		PubKey:    privKeyA.Public().(crypto.PubKey()),
		Signature: sig,
	}})

	cdc := MakeCodec()
	txBytes, err := cdc.MarshalBinary(tx)

	require.NoError(t, err)

	// Run a check
	cres := cc.CheckTx(txBytes)
	assert.Equal(t, sdk.CodeType(6),
		sdk.CodeType(cres.Code), cres.Log)

	// Simulate a Block
	cc.BeginBlock(abci.RequestBeginBlock{})
	dres := cc.DeliverTx(txBytes)
	assert.Equal(t, sdk.CodeType(6), sdk.CodeType(dres.Code), dres.Log)

}
