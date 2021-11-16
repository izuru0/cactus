package main_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger-labs/weaver-dlt-interoperability/common/protos-go/common"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	sa "github.com/hyperledger-labs/weaver-dlt-interoperability/samples/fabric/simpleassettransfer"
	"github.com/stretchr/testify/require"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	wtest "github.com/hyperledger-labs/weaver-dlt-interoperability/core/network/fabric-interop-cc/libs/testutils"
	wtestmocks "github.com/hyperledger-labs/weaver-dlt-interoperability/core/network/fabric-interop-cc/libs/testutils/mocks"
)

const (
	defaultAssetType    = "BearerBonds"
	defaultAssetId      = "asset1"
	defaultAssetOwner   = "Alice"
	defaultAssetIssuer  = "Treasury"
	defaultFaceValue    = 10000
	sourceNetworkID     = "sourcenetwork"
	destNetworkID       = "destinationnetwork"
	localNetworkIdKey   = "localNetworkID"
)

func TestInitBondAssetLedger(t *testing.T) {
	transactionContext, chaincodeStub := wtest.PrepMockStub()
	simpleAsset := sa.SmartContract{}
	simpleAsset.ConfigureInterop("interopcc")

	err := simpleAsset.InitBondAssetLedger(transactionContext, sourceNetworkID)
	require.NoError(t, err)

	chaincodeStub.PutStateReturns(fmt.Errorf("failed inserting key"))
	err = simpleAsset.InitBondAssetLedger(transactionContext, sourceNetworkID)
	require.EqualError(t, err, "failed to put to world state. failed inserting key")
}

func TestCreateAsset(t *testing.T) {
	transactionContext, chaincodeStub := wtest.PrepMockStub()
	simpleAsset := sa.SmartContract{}
	simpleAsset.ConfigureInterop("interopcc")

	err := simpleAsset.CreateAsset(transactionContext, "", "", "", "", 0, "02 Jan 26 15:04 MST")
	require.Error(t, err)

	err = simpleAsset.CreateAsset(transactionContext, defaultAssetType, "", "", "", 0, "02 Jan 26 15:04 MST")
	require.Error(t, err)

	err = simpleAsset.CreateAsset(transactionContext, defaultAssetType, defaultAssetId, "", "", 0, "02 Jan 26 15:04 MST")
	require.Error(t, err)

	err = simpleAsset.CreateAsset(transactionContext, defaultAssetType, defaultAssetId, defaultAssetOwner, "", 0, "02 Jan 26 15:04 MST")
	require.NoError(t, err)

	err = simpleAsset.CreateAsset(transactionContext, defaultAssetType, defaultAssetId, "", defaultAssetIssuer, 0, "02 Jan 26 15:04 MST")
	require.NoError(t, err)

	err = simpleAsset.CreateAsset(transactionContext, defaultAssetType, defaultAssetId, defaultAssetOwner, "", 0, "02 Jan 06 15:04 MST")
	require.EqualError(t, err, "maturity date can not be in past.")

	err = simpleAsset.CreateAsset(transactionContext, defaultAssetType, defaultAssetId, defaultAssetOwner, "", 0, "")
	require.EqualError(t, err, "maturity date provided is not in correct format, please use this format: 02 Jan 06 15:04 MST")

	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = simpleAsset.CreateAsset(transactionContext, defaultAssetType, defaultAssetId, defaultAssetOwner, "", 0, "")
	require.EqualError(t, err, "the asset asset1 already exists")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = simpleAsset.CreateAsset(transactionContext, defaultAssetType, defaultAssetId, defaultAssetOwner, "", 0, "")
	require.EqualError(t, err, "failed to read asset record from world state: unable to retrieve asset")
}

func TestReadAsset(t *testing.T) {
	transactionContext, chaincodeStub := wtest.PrepMockStub()
	simpleAsset := sa.SmartContract{}
	simpleAsset.ConfigureInterop("interopcc")

	expectedAsset := &sa.BondAsset{ID: "asset1"}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	asset, err := simpleAsset.ReadAsset(transactionContext, "", "")
	require.NoError(t, err)
	require.Equal(t, expectedAsset, asset)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	_, err = simpleAsset.ReadAsset(transactionContext, "", "")
	require.EqualError(t, err, "failed to read asset record from world state: unable to retrieve asset")

	chaincodeStub.GetStateReturns(nil, nil)
	asset, err = simpleAsset.ReadAsset(transactionContext, "", "asset1")
	require.EqualError(t, err, "the asset asset1 does not exist")
	require.Nil(t, asset)
}

func TestUpdateFaceValue(t *testing.T) {
	transactionContext, chaincodeStub := wtest.PrepMockStub()
	simpleAsset := sa.SmartContract{}
	simpleAsset.ConfigureInterop("interopcc")

	expectedAsset := &sa.BondAsset{ID: "asset1"}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	err = simpleAsset.UpdateFaceValue(transactionContext, "", "", 0)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = simpleAsset.UpdateFaceValue(transactionContext, "", "asset1", 0)
	require.EqualError(t, err, "the asset asset1 does not exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = simpleAsset.UpdateFaceValue(transactionContext, "", "asset1", 0)
	require.EqualError(t, err, "failed to read asset record from world state: unable to retrieve asset")
}

func TestUpdateMaturityDate(t *testing.T) {
	transactionContext, chaincodeStub := wtest.PrepMockStub()
	simpleAsset := sa.SmartContract{}
	simpleAsset.ConfigureInterop("interopcc")

	expectedAsset := &sa.BondAsset{ID: "asset1"}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	err = simpleAsset.UpdateMaturityDate(transactionContext, "", "", time.Now())
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = simpleAsset.UpdateMaturityDate(transactionContext, "", "asset1", time.Now())
	require.EqualError(t, err, "the asset asset1 does not exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = simpleAsset.UpdateMaturityDate(transactionContext, "", "asset1", time.Now())
	require.EqualError(t, err, "failed to read asset record from world state: unable to retrieve asset")
}

func TestDeleteAsset(t *testing.T) {
	transactionContext, chaincodeStub := wtest.PrepMockStub()
	simpleAsset := sa.SmartContract{}
	simpleAsset.ConfigureInterop("interopcc")

	asset := &sa.BondAsset{ID: "asset1"}
	bytes, err := json.Marshal(asset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	chaincodeStub.DelStateReturns(nil)
	err = simpleAsset.DeleteAsset(transactionContext, "", "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = simpleAsset.DeleteAsset(transactionContext, "", "asset1")
	require.EqualError(t, err, "the asset asset1 does not exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = simpleAsset.DeleteAsset(transactionContext, "", "")
	require.EqualError(t, err, "failed to read asset record from world state: unable to retrieve asset")
}

func TestUpdateOwner(t *testing.T) {
	transactionContext, chaincodeStub := wtest.PrepMockStub()
	simpleAsset := sa.SmartContract{}
	simpleAsset.ConfigureInterop("interopcc")

	asset := &sa.BondAsset{ID: "asset1"}
	bytes, err := json.Marshal(asset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	err = simpleAsset.UpdateOwner(transactionContext, "", "", "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = simpleAsset.UpdateOwner(transactionContext, "", "", "")
	require.EqualError(t, err, "failed to read asset record from world state: unable to retrieve asset")
}

func TestGetMyAssets(t *testing.T) {
	transactionContext, chaincodeStub := wtest.PrepMockStub()
	simpleAsset := sa.SmartContract{}
	simpleAsset.ConfigureInterop("interopcc")
	iterator := &wtestmocks.StateQueryIterator{}

	asset := &sa.BondAsset{ID: "asset1", Owner: getTestTxCreatorECertBase64()}
	bytes, err := json.Marshal(asset)
	require.NoError(t, err)

	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, false)
	iterator.NextReturns(&queryresult.KV{Value: bytes}, nil)

	chaincodeStub.GetCreatorReturns([]byte(getCreator()), nil)

	chaincodeStub.GetStateByRangeReturns(iterator, nil)
	assets, err := simpleAsset.GetAllAssets(transactionContext)
	require.NoError(t, err)
	require.Equal(t, []*sa.BondAsset{asset}, assets)

	iterator.HasNextReturns(true)
	iterator.NextReturns(nil, fmt.Errorf("failed retrieving next item"))
	assets, err = simpleAsset.GetAllAssets(transactionContext)
	require.EqualError(t, err, "failed retrieving next item")
	require.Nil(t, assets)

	chaincodeStub.GetStateByRangeReturns(nil, fmt.Errorf("failed retrieving all assets"))
	assets, err = simpleAsset.GetAllAssets(transactionContext)
	require.EqualError(t, err, "failed retrieving all assets")
	require.Nil(t, assets)
}

func TestGetAllAssets(t *testing.T) {
	transactionContext, chaincodeStub := wtest.PrepMockStub()
	simpleAsset := sa.SmartContract{}
	simpleAsset.ConfigureInterop("interopcc")
	iterator := &wtestmocks.StateQueryIterator{}

	asset := &sa.BondAsset{ID: "asset1"}
	bytes, err := json.Marshal(asset)
	require.NoError(t, err)

	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, false)
	iterator.NextReturns(&queryresult.KV{Value: bytes}, nil)

	chaincodeStub.GetStateByRangeReturns(iterator, nil)
	assets, err := simpleAsset.GetAllAssets(transactionContext)
	require.NoError(t, err)
	require.Equal(t, []*sa.BondAsset{asset}, assets)

	iterator.HasNextReturns(true)
	iterator.NextReturns(nil, fmt.Errorf("failed retrieving next item"))
	assets, err = simpleAsset.GetAllAssets(transactionContext)
	require.EqualError(t, err, "failed retrieving next item")
	require.Nil(t, assets)

	chaincodeStub.GetStateByRangeReturns(nil, fmt.Errorf("failed retrieving all assets"))
	assets, err = simpleAsset.GetAllAssets(transactionContext)
	require.EqualError(t, err, "failed retrieving all assets")
	require.Nil(t, assets)
}

func TestPledgeAsset(t *testing.T) {
	transactionContext, chaincodeStub := wtest.PrepMockStub()
	simpleAsset := sa.SmartContract{}
	simpleAsset.ConfigureInterop("interopcc")

	// Pledge non-existent asset
	expiry := uint64(time.Now().Unix()) + (5 * 60)      // Expires 5 minutes from now
	err := simpleAsset.PledgeAsset(transactionContext, defaultAssetType, defaultAssetId, destNetworkID, getRecipientECertBase64(), expiry)
	require.Error(t, err)

	maturityDate := "02 Jan 26 15:04 MST"
	err = simpleAsset.CreateAsset(transactionContext, defaultAssetType, defaultAssetId, defaultAssetOwner, "", 0, maturityDate)
	require.NoError(t, err)

	bondAssetKey := defaultAssetType + defaultAssetId
	md_time, err := time.Parse(time.RFC822, maturityDate)
	bondAsset := sa.BondAsset{
		Type: defaultAssetType,
		ID: defaultAssetId,
		Owner: defaultAssetOwner,
		Issuer: defaultAssetIssuer,
		FaceValue: defaultFaceValue,
		MaturityDate: md_time,
	}
	bondAssetJSON, _ := json.Marshal(bondAsset)
	chaincodeStub.GetStateReturnsForKey(bondAssetKey, bondAssetJSON, nil)
	chaincodeStub.GetCreatorReturns([]byte(getCreatorInContext("locker")), nil)
	err = simpleAsset.PledgeAsset(transactionContext, defaultAssetType, defaultAssetId, destNetworkID, getRecipientECertBase64(), expiry)
	require.Error(t, err)       // Asset owner is not the pledger

	bondAsset.Owner = getLockerECertBase64()
	bondAssetJSON, _ = json.Marshal(bondAsset)
	chaincodeStub.GetStateReturnsForKey(bondAssetKey, bondAssetJSON, nil)
	chaincodeStub.InvokeChaincodeReturns(shim.Success([]byte("true")))
	err = simpleAsset.PledgeAsset(transactionContext, defaultAssetType, defaultAssetId, destNetworkID, getRecipientECertBase64(), expiry)
	require.Error(t, err)       // Already locked asset cannot be pledged

	bondAssetPledgeKey := "Pledged_" + defaultAssetType + defaultAssetId
	bondAssetPledge := &common.AssetPledge{
		AssetDetails: bondAssetJSON,
		LocalNetworkID: sourceNetworkID,
		RemoteNetworkID: destNetworkID,
		Recipient: getRecipientECertBase64(),
		ExpiryTimeSecs: expiry,
	}
	bondAssetPledgeBytes, _ := proto.Marshal(bondAssetPledge)
	chaincodeStub.GetStateReturnsForKey(bondAssetPledgeKey, bondAssetPledgeBytes, nil)
	chaincodeStub.InvokeChaincodeReturns(shim.Success([]byte("false")))
	err = simpleAsset.PledgeAsset(transactionContext, defaultAssetType, defaultAssetId, "someremoteNetwork", getRecipientECertBase64(), expiry)
	require.Error(t, err)       // Already pledged asset cannot be pledged if the pledge attributes don't match the recorded value

	err = simpleAsset.PledgeAsset(transactionContext, defaultAssetType, defaultAssetId, destNetworkID, getRecipientECertBase64(), expiry)
	require.NoError(t, err)     // Asset is already pledged, so there is nothing more to be done

	chaincodeStub.GetStateClearForKey(bondAssetPledgeKey)
	err = simpleAsset.PledgeAsset(transactionContext, defaultAssetType, defaultAssetId, destNetworkID, getRecipientECertBase64(), expiry - (10 * 60))
	require.Error(t, err)       // Invalid pledge as its expiry time in the past

	chaincodeStub.GetStateReturnsForKey(localNetworkIdKey, []byte(sourceNetworkID), nil)
	chaincodeStub.PutStateReturns(nil)
	chaincodeStub.DelStateReturns(nil)
	err = simpleAsset.PledgeAsset(transactionContext, defaultAssetType, defaultAssetId, destNetworkID, getRecipientECertBase64(), expiry)
	require.NoError(t, err)     // Asset pledge is recorded
}

func TestClaimRemoteAsset(t *testing.T) {
	transactionContext, chaincodeStub := wtest.PrepMockStub()
	simpleAsset := sa.SmartContract{}
	simpleAsset.ConfigureInterop("interopcc")

	maturityDate := "02 Jan 26 15:04 MST"
	md_time, err := time.Parse(time.RFC822, maturityDate)
	bondAsset := sa.BondAsset{
		Type: defaultAssetType,
		ID: defaultAssetId,
		Owner: getLockerECertBase64(),
		Issuer: defaultAssetIssuer,
		FaceValue: defaultFaceValue,
		MaturityDate: md_time,
	}
	bondAssetJSON, _ := json.Marshal(bondAsset)

	expiry := uint64(time.Now().Unix()) - (5 * 60)
	bondAssetPledge := &common.AssetPledge{
		AssetDetails: bondAssetJSON,
		LocalNetworkID: sourceNetworkID,
		RemoteNetworkID: destNetworkID,
		Recipient: getRecipientECertBase64(),
		ExpiryTimeSecs: expiry,
	}
	bondAssetPledgeBytes, _ := proto.Marshal(bondAssetPledge)

	chaincodeStub.GetCreatorReturns([]byte(getCreatorInContext("recipient")), nil)
	err = simpleAsset.ClaimRemoteAsset(transactionContext, defaultAssetType, defaultAssetId, getLockerECertBase64(), sourceNetworkID, bondAssetPledgeBytes)
	require.Error(t, err)       // Expired pledge

	bondAssetPledge.ExpiryTimeSecs = bondAssetPledge.ExpiryTimeSecs + (10 * 60)
	bondAssetPledgeBytes, _ = proto.Marshal(bondAssetPledge)

	err = simpleAsset.ClaimRemoteAsset(transactionContext, defaultAssetType, defaultAssetId, getRecipientECertBase64(), sourceNetworkID, bondAssetPledgeBytes)
	require.Error(t, err)       // Unexpected pledged asset owner

	bondAssetPledge.Recipient = getLockerECertBase64()
	bondAssetPledgeBytes, _ = proto.Marshal(bondAssetPledge)
	err = simpleAsset.ClaimRemoteAsset(transactionContext, defaultAssetType, defaultAssetId, getLockerECertBase64(), sourceNetworkID, bondAssetPledgeBytes)
	require.Error(t, err)       // Claimer doesn't match pledge recipient

	bondAssetPledge.Recipient = getRecipientECertBase64()
	bondAssetPledgeBytes, _ = proto.Marshal(bondAssetPledge)
	err = simpleAsset.ClaimRemoteAsset(transactionContext, defaultAssetType, defaultAssetId, getLockerECertBase64(), destNetworkID, bondAssetPledgeBytes)
	require.Error(t, err)       // Pledge not made for the claiming network

	chaincodeStub.GetStateReturnsForKey(localNetworkIdKey, []byte(sourceNetworkID), nil)
	err = simpleAsset.ClaimRemoteAsset(transactionContext, defaultAssetType, defaultAssetId, getLockerECertBase64(), sourceNetworkID, bondAssetPledgeBytes)
	require.Error(t, err)       // Pledge claimed from the wrong network

	chaincodeStub.GetStateReturnsForKey(localNetworkIdKey, []byte(destNetworkID), nil)
	chaincodeStub.PutStateReturns(nil)
	err = simpleAsset.ClaimRemoteAsset(transactionContext, defaultAssetType, defaultAssetId, getLockerECertBase64(), sourceNetworkID, bondAssetPledgeBytes)
	require.NoError(t, err)     // Asset claim is recorded
}

func TestReclaimAsset(t *testing.T) {
	transactionContext, chaincodeStub := wtest.PrepMockStub()
	simpleAsset := sa.SmartContract{}
	simpleAsset.ConfigureInterop("interopcc")

	maturityDate := "02 Jan 26 15:04 MST"
	md_time, err := time.Parse(time.RFC822, maturityDate)
	bondAsset := sa.BondAsset{
		Type: defaultAssetType,
		ID: defaultAssetId,
		Owner: getLockerECertBase64(),
		Issuer: defaultAssetIssuer,
		FaceValue: defaultFaceValue,
		MaturityDate: md_time,
	}
	bondAssetJSON, _ := json.Marshal(bondAsset)

	expiry := uint64(time.Now().Unix()) + (5 * 60)
	bondAssetPledge := &common.AssetPledge{
		AssetDetails: bondAssetJSON,
		LocalNetworkID: sourceNetworkID,
		RemoteNetworkID: destNetworkID,
		Recipient: getRecipientECertBase64(),
		ExpiryTimeSecs: expiry,
	}
	bondAssetPledgeBytes, _ := proto.Marshal(bondAssetPledge)

	claimStatus := &common.AssetClaimStatus{
		AssetDetails: bondAssetJSON,
		LocalNetworkID: destNetworkID,
		RemoteNetworkID: sourceNetworkID,
		Recipient: getRecipientECertBase64(),
		ClaimStatus: false,
		ExpiryTimeSecs: expiry,
		ExpirationStatus: false,
	}
	claimStatusBytes, _ := proto.Marshal(claimStatus)

	chaincodeStub.GetCreatorReturns([]byte(getCreatorInContext("locker")), nil)
	err = simpleAsset.ReclaimAsset(transactionContext, defaultAssetType, defaultAssetId, getRecipientECertBase64(), destNetworkID, claimStatusBytes)
	require.Error(t, err)       // no pledge recorded

	bondAssetPledgeKey := "Pledged_" + defaultAssetType + defaultAssetId
	chaincodeStub.GetStateReturnsForKey(bondAssetPledgeKey, bondAssetPledgeBytes, nil)
	err = simpleAsset.ReclaimAsset(transactionContext, defaultAssetType, defaultAssetId, getRecipientECertBase64(), destNetworkID, claimStatusBytes)
	require.Error(t, err)       // pledge has not expired yet

	bondAssetPledge.ExpiryTimeSecs = expiry - (10 * 60)
	bondAssetPledgeBytes, _ = proto.Marshal(bondAssetPledge)
	chaincodeStub.GetStateReturnsForKey(bondAssetPledgeKey, bondAssetPledgeBytes, nil)
	claimStatus.ExpiryTimeSecs = expiry - (10 * 60)
	claimStatus.ExpirationStatus = false
	claimStatusBytes, _ = proto.Marshal(claimStatus)
	err = simpleAsset.ReclaimAsset(transactionContext, defaultAssetType, defaultAssetId, getRecipientECertBase64(), destNetworkID, claimStatusBytes)
	require.Error(t, err)       // claim probe time was before expiration time

	claimStatus.ClaimStatus = true
	claimStatusBytes, _ = proto.Marshal(claimStatus)
	err = simpleAsset.ReclaimAsset(transactionContext, defaultAssetType, defaultAssetId, getRecipientECertBase64(), destNetworkID, claimStatusBytes)
	require.Error(t, err)       // claim was successfully made

	claimStatus.ExpirationStatus = true
	bondAsset.ID = "someid"
	bondAssetJSON, _ = json.Marshal(bondAsset)
	claimStatus.AssetDetails = bondAssetJSON
	claimStatus.ClaimStatus = false
	claimStatusBytes, _ = proto.Marshal(claimStatus)
	err = simpleAsset.ReclaimAsset(transactionContext, defaultAssetType, defaultAssetId, getRecipientECertBase64(), destNetworkID, claimStatusBytes)
	require.Error(t, err)       // claim was for a different asset

	bondAsset.ID = defaultAssetId
	bondAssetJSON, _ = json.Marshal(bondAsset)
	claimStatus.AssetDetails = bondAssetJSON
	claimStatusBytes, _ = proto.Marshal(claimStatus)
	err = simpleAsset.ReclaimAsset(transactionContext, defaultAssetType, defaultAssetId, getRecipientECertBase64(), "somenetworkid", claimStatusBytes)
	require.Error(t, err)       // claim was probed in a different network than expected

	err = simpleAsset.ReclaimAsset(transactionContext, defaultAssetType, defaultAssetId, getLockerECertBase64(), destNetworkID, claimStatusBytes)
	require.Error(t, err)       // claim recipient was different than expected

	chaincodeStub.GetStateReturnsForKey(localNetworkIdKey, []byte(destNetworkID), nil)
	err = simpleAsset.ReclaimAsset(transactionContext, defaultAssetType, defaultAssetId, getRecipientECertBase64(), destNetworkID, claimStatusBytes)
	require.Error(t, err)       // claim was not made for an asset in my network

	chaincodeStub.GetStateReturnsForKey(localNetworkIdKey, []byte(sourceNetworkID), nil)
	chaincodeStub.PutStateReturns(nil)
	chaincodeStub.DelStateReturns(nil)
	err = simpleAsset.ReclaimAsset(transactionContext, defaultAssetType, defaultAssetId, getRecipientECertBase64(), destNetworkID, claimStatusBytes)
	require.NoError(t, err)     // Asset is reclaimed
}

func TestAssetTransferQueries(t *testing.T) {
	transactionContext, chaincodeStub := wtest.PrepMockStub()
	simpleAsset := sa.SmartContract{}
	simpleAsset.ConfigureInterop("interopcc")

	maturityDate := "02 Jan 26 15:04 MST"
	md_time, err := time.Parse(time.RFC822, maturityDate)
	bondAsset := sa.BondAsset{
		Type: defaultAssetType,
		ID: defaultAssetId,
		Owner: getLockerECertBase64(),
		Issuer: defaultAssetIssuer,
		FaceValue: defaultFaceValue,
		MaturityDate: md_time,
	}
	bondAssetJSON, _ := json.Marshal(bondAsset)

	expiry := uint64(time.Now().Unix()) + (5 * 60)
	bondAssetPledge := &common.AssetPledge{
		AssetDetails: bondAssetJSON,
		LocalNetworkID: sourceNetworkID,
		RemoteNetworkID: destNetworkID,
		Recipient: getRecipientECertBase64(),
		ExpiryTimeSecs: expiry,
	}
	bondAssetPledgeBytes, _ := proto.Marshal(bondAssetPledge)

	claimStatus := &common.AssetClaimStatus{
		AssetDetails: bondAssetJSON,
		LocalNetworkID: destNetworkID,
		RemoteNetworkID: sourceNetworkID,
		Recipient: getRecipientECertBase64(),
		ClaimStatus: true,
		ExpiryTimeSecs: expiry,
		ExpirationStatus: false,
	}
	claimStatusBytes, _ := proto.Marshal(claimStatus)

	// Query for pledge when none exists
	chaincodeStub.GetCreatorReturns([]byte(getCreatorInContext("locker")), nil)
	pledgeStatus, err := simpleAsset.GetAssetPledgeStatus(transactionContext, defaultAssetType, defaultAssetId, getLockerECertBase64(), destNetworkID,
		getRecipientECertBase64())
	require.NoError(t, err)
	lookupPledge := &common.AssetPledge{}
	proto.Unmarshal(pledgeStatus, lookupPledge)
	var lookupPledgeAsset sa.BondAsset
	json.Unmarshal([]byte(lookupPledge.AssetDetails), &lookupPledgeAsset)
	require.Equal(t, "", lookupPledgeAsset.Type)
	require.Equal(t, "", lookupPledgeAsset.ID)
	require.Equal(t, "", lookupPledgeAsset.Owner)
	require.Equal(t, "", lookupPledgeAsset.Issuer)
	require.Equal(t, "", lookupPledge.LocalNetworkID)
	require.Equal(t, "", lookupPledge.RemoteNetworkID)
	require.Equal(t, "", lookupPledge.Recipient)

	// Query for pledge after recording one
	bondAssetPledgeKey := "Pledged_" + defaultAssetType + defaultAssetId
	chaincodeStub.GetStateReturnsForKey(bondAssetPledgeKey, bondAssetPledgeBytes, nil)
	pledgeStatus, err = simpleAsset.GetAssetPledgeStatus(transactionContext, defaultAssetType, defaultAssetId, getLockerECertBase64(), destNetworkID,
		getRecipientECertBase64())
	require.NoError(t, err)
	proto.Unmarshal(pledgeStatus, lookupPledge)
	json.Unmarshal([]byte(lookupPledge.AssetDetails), &lookupPledgeAsset)
	var originalPledgeAsset sa.BondAsset
	json.Unmarshal([]byte(bondAssetPledge.AssetDetails), &originalPledgeAsset)
	require.Equal(t, originalPledgeAsset.Type, lookupPledgeAsset.Type)
	require.Equal(t, originalPledgeAsset.ID, lookupPledgeAsset.ID)
	require.Equal(t, originalPledgeAsset.Owner, lookupPledgeAsset.Owner)
	require.Equal(t, originalPledgeAsset.Issuer, lookupPledgeAsset.Issuer)
	require.Equal(t, bondAssetPledge.LocalNetworkID, lookupPledge.LocalNetworkID)
	require.Equal(t, bondAssetPledge.RemoteNetworkID, lookupPledge.RemoteNetworkID)
	require.Equal(t, bondAssetPledge.Recipient, lookupPledge.Recipient)

	// Query for claim when no asset or claim exists
	chaincodeStub.GetCreatorReturns([]byte(getCreatorInContext("recipient")), nil)
	claimStatusQueried, err := simpleAsset.GetAssetClaimStatus(transactionContext, defaultAssetType, defaultAssetId, getRecipientECertBase64(),
		getLockerECertBase64(), sourceNetworkID, expiry)
	require.NoError(t, err)
	lookupClaim := &common.AssetClaimStatus{}
	proto.Unmarshal(claimStatusQueried, lookupClaim)
	var lookupClaimAsset sa.BondAsset
	json.Unmarshal([]byte(lookupClaim.AssetDetails), &lookupClaimAsset)
	require.Equal(t, "", lookupClaimAsset.Type)
	require.Equal(t, "", lookupClaimAsset.ID)
	require.Equal(t, "", lookupClaimAsset.Owner)
	require.Equal(t, "", lookupClaimAsset.Issuer)
	require.Equal(t, "", lookupClaim.LocalNetworkID)
	require.Equal(t, "", lookupClaim.RemoteNetworkID)
	require.Equal(t, "", lookupClaim.Recipient)
	require.False(t, lookupClaim.ClaimStatus)

	// Query for claim when only asset but no claim exists
	bondAssetKey := defaultAssetType + defaultAssetId
	bondAsset.Owner = getRecipientECertBase64()
	bondAsset.Issuer = getLockerECertBase64()
	bondAssetJSON, _ = json.Marshal(bondAsset)
	chaincodeStub.GetStateReturnsForKey(bondAssetKey, bondAssetJSON, nil)
	chaincodeStub.GetCreatorReturns([]byte(getCreatorInContext("recipient")), nil)
	claimStatusQueried, err = simpleAsset.GetAssetClaimStatus(transactionContext, defaultAssetType, defaultAssetId, getRecipientECertBase64(),
		getLockerECertBase64(), sourceNetworkID, expiry)
	require.NoError(t, err)
	proto.Unmarshal(claimStatusQueried, lookupClaim)
	json.Unmarshal([]byte(lookupClaim.AssetDetails), &lookupClaimAsset)
	require.Equal(t, "", lookupClaimAsset.Type)
	require.Equal(t, "", lookupClaimAsset.ID)
	require.Equal(t, "", lookupClaimAsset.Owner)
	require.Equal(t, "", lookupClaimAsset.Issuer)
	require.Equal(t, "", lookupClaim.LocalNetworkID)
	require.Equal(t, "", lookupClaim.RemoteNetworkID)
	require.Equal(t, "", lookupClaim.Recipient)
	require.False(t, lookupClaim.ClaimStatus)

	// Query for claim after recording both an asset and a claim
	bondAssetClaimKey := "Claimed_" + defaultAssetType + defaultAssetId
	chaincodeStub.GetStateReturnsForKey(bondAssetClaimKey, claimStatusBytes, nil)
	claimStatusQueried, err = simpleAsset.GetAssetClaimStatus(transactionContext, defaultAssetType, defaultAssetId, getRecipientECertBase64(),
		getLockerECertBase64(), sourceNetworkID, expiry)
	require.NoError(t, err)
	proto.Unmarshal(claimStatusQueried, lookupClaim)
	json.Unmarshal([]byte(lookupClaim.AssetDetails), &lookupClaimAsset)
	var originalClaimAsset sa.BondAsset
	json.Unmarshal([]byte(claimStatus.AssetDetails), &originalClaimAsset)
	require.Equal(t, originalClaimAsset.Type, lookupClaimAsset.Type)
	require.Equal(t, originalClaimAsset.ID, lookupClaimAsset.ID)
	require.Equal(t, originalClaimAsset.Owner, lookupClaimAsset.Owner)
	require.Equal(t, originalClaimAsset.Issuer, lookupClaimAsset.Issuer)
	require.Equal(t, claimStatus.LocalNetworkID, lookupClaim.LocalNetworkID)
	require.Equal(t, claimStatus.RemoteNetworkID, lookupClaim.RemoteNetworkID)
	require.Equal(t, claimStatus.Recipient, lookupClaim.Recipient)
	require.Equal(t, claimStatus.ClaimStatus, lookupClaim.ClaimStatus)
}
