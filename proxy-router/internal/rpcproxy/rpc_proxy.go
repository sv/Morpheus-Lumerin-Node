package rpcproxy

import (
	"context"
	"math/big"

	constants "github.com/MorpheusAIs/Morpheus-Lumerin-Node/proxy-router/internal/internal"
	"github.com/MorpheusAIs/Morpheus-Lumerin-Node/proxy-router/internal/internal/interfaces"
	"github.com/MorpheusAIs/Morpheus-Lumerin-Node/proxy-router/internal/internal/lib"
	"github.com/MorpheusAIs/Morpheus-Lumerin-Node/proxy-router/internal/internal/repositories/registries"
	"github.com/MorpheusAIs/Morpheus-Lumerin-Node/proxy-router/internal/internal/rpcproxy/structs"
	"github.com/gin-gonic/gin"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RpcProxy struct {
	rpcClient        *ethclient.Client
	providerRegistry *registries.ProviderRegistry
	modelRegistry    *registries.ModelRegistry
	marketplace      *registries.Marketplace
	sessionRouter    *registries.SessionRouter

	legacyTx   bool
	privateKey string
}

func NewRpcProxy(rpcClient *ethclient.Client, diamonContractAddr common.Address, privateKey string, log interfaces.ILogger, legacyTx bool) *RpcProxy {
	providerRegistry := registries.NewProviderRegistry(diamonContractAddr, rpcClient, log)
	modelRegistry := registries.NewModelRegistry(diamonContractAddr, rpcClient, log)
	marketplace := registries.NewMarketplace(diamonContractAddr, rpcClient, log)
	sessionRouter := registries.NewSessionRouter(diamonContractAddr, rpcClient, log)
	return &RpcProxy{
		rpcClient:        rpcClient,
		providerRegistry: providerRegistry,
		modelRegistry:    modelRegistry,
		marketplace:      marketplace,
		sessionRouter:    sessionRouter,
		legacyTx:         legacyTx,
		privateKey:       privateKey,
	}
}

func (rpcProxy *RpcProxy) GetLatestBlock(ctx context.Context) (uint64, error) {
	return rpcProxy.rpcClient.BlockNumber(ctx)
}

func (rpcProxy *RpcProxy) GetAllProviders(ctx context.Context) (int, gin.H) {
	addrs, providers, err := rpcProxy.providerRegistry.GetAllProviders(ctx)
	if err != nil {
		return constants.HTTP_STATUS_BAD_REQUEST, gin.H{"error": err.Error()}
	}

	result := make([]*structs.Provider, len(addrs))
	for i, value := range providers {
		result[i] = &structs.Provider{
			Address:   addrs[i],
			Endpoint:  value.Endpoint,
			Stake:     value.Stake,
			IsDeleted: value.IsDeleted,
			Timestamp: value.Timestamp,
		}
	}

	return constants.HTTP_STATUS_OK, gin.H{"providers": result}
}

func (rpcProxy *RpcProxy) GetAllModels(ctx context.Context) (int, gin.H) {
	ids, models, err := rpcProxy.modelRegistry.GetAllModels(ctx)
	if err != nil {
		return constants.HTTP_STATUS_BAD_REQUEST, gin.H{"error": err.Error()}
	}

	result := make([]*structs.Model, len(ids))
	for i, value := range models {
		result[i] = &structs.Model{
			Id:        lib.BytesToString(ids[i][:]),
			IpfsCID:   lib.BytesToString(value.IpfsCID[:]),
			Fee:       value.Fee,
			Stake:     value.Stake,
			Owner:     value.Owner,
			Name:      value.Name,
			Tags:      value.Tags,
			Timestamp: value.Timestamp,
			IsDeleted: value.IsDeleted,
		}
	}

	return constants.HTTP_STATUS_OK, gin.H{"models": result}
}

func (rpcProxy *RpcProxy) GetBidsByProvider(ctx context.Context, providerAddr common.Address, offset *big.Int, limit uint8) (int, gin.H) {
	ids, bids, err := rpcProxy.marketplace.GetBidsByProvider(ctx, providerAddr, offset, limit)
	if err != nil {
		return constants.HTTP_STATUS_BAD_REQUEST, gin.H{"error": err.Error()}
	}

	result := make([]*structs.Bid, len(ids))
	for i, value := range bids {
		result[i] = &structs.Bid{
			Id:             lib.BytesToString(ids[i][:]),
			ModelAgentId:   lib.BytesToString(value.ModelAgentId[:]),
			Provider:       value.Provider,
			Nonce:          value.Nonce,
			CreatedAt:      value.CreatedAt,
			DeletedAt:      value.DeletedAt,
			PricePerSecond: value.PricePerSecond,
		}
	}
	return constants.HTTP_STATUS_OK, gin.H{"bids": result}
}

func (rpcProxy *RpcProxy) GetBidsByModelAgent(ctx context.Context, modelId [32]byte, offset *big.Int, limit uint8) (int, gin.H) {
	ids, bids, err := rpcProxy.marketplace.GetBidsByModelAgent(ctx, modelId, offset, limit)
	if err != nil {
		return constants.HTTP_STATUS_BAD_REQUEST, gin.H{"error": err.Error()}
	}

	result := make([]*structs.Bid, len(ids))
	for i, value := range bids {
		result[i] = &structs.Bid{
			Id:             lib.BytesToString(ids[i][:]),
			ModelAgentId:   lib.BytesToString(value.ModelAgentId[:]),
			Provider:       value.Provider,
			Nonce:          value.Nonce,
			CreatedAt:      value.CreatedAt,
			DeletedAt:      value.DeletedAt,
			PricePerSecond: value.PricePerSecond,
		}
	}

	return constants.HTTP_STATUS_OK, gin.H{"bids": result}
}

func (rpcProxy *RpcProxy) OpenSession(ctx *gin.Context) (int, gin.H) {
	var reqPayload map[string]interface{}
	if err := ctx.ShouldBindJSON(&reqPayload); err != nil {
		return constants.HTTP_STATUS_BAD_REQUEST, gin.H{"error": err.Error()}
	}

	bidIdStr, ok := reqPayload["bidId"].(string)
	if !ok {
		return constants.HTTP_STATUS_BAD_REQUEST, gin.H{"error": "bidId is required"}
	}
	stakeStr, ok := reqPayload["stake"].(string)
	if !ok {
		print(ok)
		return constants.HTTP_STATUS_BAD_REQUEST, gin.H{"error": "stake is required"}
	}

	stake, ok := new(big.Int).SetString(stakeStr, 10)
	if !ok {
		return constants.HTTP_STATUS_BAD_REQUEST, gin.H{"error": "stake is invalid"}
	}

	bidId := common.FromHex(bidIdStr)

	transactOpt, err := rpcProxy.getTransactOpts(ctx, rpcProxy.privateKey)
	if err != nil {
		return constants.HTTP_STATUS_BAD_REQUEST, gin.H{"error": err.Error()}
	}

	sessionId, err := rpcProxy.sessionRouter.OpenSession(transactOpt, [32]byte(bidId), stake)
	if err != nil {
		return constants.HTTP_STATUS_BAD_REQUEST, gin.H{"error": err.Error()}
	}

	return constants.HTTP_STATUS_OK, gin.H{"sessionId": sessionId}
}

func (rpcProxy *RpcProxy) GetSession(ctx *gin.Context, sessionId string) (int, gin.H) {
	session, err := rpcProxy.sessionRouter.GetSession(ctx, sessionId)
	if err != nil {
		return constants.HTTP_STATUS_BAD_REQUEST, gin.H{"error": err.Error()}
	}
	return constants.HTTP_STATUS_OK, gin.H{"session": session}
}

func (rpcProxy *RpcProxy) getTransactOpts(ctx context.Context, privKey string) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}

	chainId, err := rpcProxy.rpcClient.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, err
	}

	// TODO: deal with likely gasPrice issue so our transaction processes before another pending nonce.
	if rpcProxy.legacyTx {
		gasPrice, err := rpcProxy.rpcClient.SuggestGasPrice(ctx)
		if err != nil {
			return nil, err
		}
		transactOpts.GasPrice = gasPrice
	}

	transactOpts.Value = big.NewInt(0)
	transactOpts.Context = ctx

	return transactOpts, nil
}
