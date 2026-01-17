package keeper

import (
	"fmt"
	"strconv"

	storetypes "cosmossdk.io/store/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/product_certificate/types"
)

// Keeper manages the product_certificate module's state
type Keeper struct {
	cdc        codec.Codec
	storeKey   storetypes.StoreKey
	paramSpace paramtypes.Subspace
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey, paramSpace paramtypes.Subspace) Keeper {
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramSpace: paramSpace,
	}
}

// GetParams returns the current parameters for the product_certificate module
func (k Keeper) GetParams(ctx sdk.Context) (types.Params, error) {
	var params types.Params
	k.paramSpace.GetParamSet(ctx, &params)
	return params, nil
}

// SetParams sets the parameters for the product_certificate module
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	k.paramSpace.SetParamSet(ctx, &params)
	return nil
}

// GetCertificate returns a certificate by ID
func (k Keeper) GetCertificate(ctx sdk.Context, id string) (types.ProductCertificate, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetCertificateKey(id)
	bz := store.Get(key)
	if bz == nil {
		return types.ProductCertificate{}, false
	}

	var cert types.ProductCertificate
	if err := k.cdc.Unmarshal(bz, &cert); err != nil {
		ctx.Logger().Error("Failed to unmarshal certificate",
			"error", err,
			"key", string(types.GetCertificateKey(id)))
		return types.ProductCertificate{}, false
	}
	return cert, true
}

// SetCertificate sets a certificate in the store
func (k Keeper) SetCertificate(ctx sdk.Context, cert types.ProductCertificate) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetCertificateKey(cert.ID)
	bz := k.cdc.MustMarshal(&cert)
	store.Set(key, bz)
	return nil
}

// GetAllCertificates returns all certificates
func (k Keeper) GetAllCertificates(ctx sdk.Context) ([]types.ProductCertificate, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.CertificateKeyPrefix))
	defer iterator.Close()

	var certificates []types.ProductCertificate
	for ; iterator.Valid(); iterator.Next() {
		var cert types.ProductCertificate
		if err := k.cdc.Unmarshal(iterator.Value(), &cert); err != nil {
			ctx.Logger().Error("Failed to unmarshal certificate, skipping",
				"error", err,
				"key", string(iterator.Key()))
			continue
		}
		certificates = append(certificates, cert)
	}

	return certificates, nil
}

// GetVerification returns a verification by certificate ID and verifier
func (k Keeper) GetVerification(ctx sdk.Context, certificateID, verifier string) (types.CertificateVerification, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetVerificationKey(certificateID, verifier)
	bz := store.Get(key)
	if bz == nil {
		return types.CertificateVerification{}, false
	}

	var verification types.CertificateVerification
	if err := k.cdc.Unmarshal(bz, &verification); err != nil {
		ctx.Logger().Error("Failed to unmarshal verification",
			"error", err,
			"key", string(types.GetVerificationKey(certificateID, verifier)))
		return types.CertificateVerification{}, false
	}
	return verification, true
}

// SetVerification sets a verification in the store
func (k Keeper) SetVerification(ctx sdk.Context, verification types.CertificateVerification) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetVerificationKey(verification.CertificateID, verification.Verifier)
	bz := k.cdc.MustMarshal(&verification)
	store.Set(key, bz)
	return nil
}

// GetAllVerifications returns all verifications
func (k Keeper) GetAllVerifications(ctx sdk.Context) ([]types.CertificateVerification, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.VerificationKeyPrefix))
	defer iterator.Close()

	var verifications []types.CertificateVerification
	for ; iterator.Valid(); iterator.Next() {
		var verification types.CertificateVerification
		if err := k.cdc.Unmarshal(iterator.Value(), &verification); err != nil {
			ctx.Logger().Error("Failed to unmarshal verification, skipping",
				"error", err,
				"key", string(iterator.Key()))
			continue
		}
		verifications = append(verifications, verification)
	}

	return verifications, nil
}

// CreateCertificate creates a new certificate
func (k Keeper) CreateCertificate(ctx sdk.Context, cert types.ProductCertificate) error {
	// Check if certificate already exists
	if _, exists := k.GetCertificate(ctx, cert.ID); exists {
		return fmt.Errorf("certificate with ID %s already exists", cert.ID)
	}

	// Validate certificate
	if cert.ID == "" {
		return fmt.Errorf("certificate ID cannot be empty")
	}
	if cert.ProductID == "" {
		return fmt.Errorf("product ID cannot be empty")
	}
	if cert.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	if cert.Status == "" {
		return fmt.Errorf("status cannot be empty")
	}

	// Set certificate
	return k.SetCertificate(ctx, cert)
}

// UpdateCertificate updates an existing certificate
func (k Keeper) UpdateCertificate(ctx sdk.Context, cert types.ProductCertificate) error {
	// Check if certificate exists
	if _, exists := k.GetCertificate(ctx, cert.ID); !exists {
		return fmt.Errorf("certificate with ID %s does not exist", cert.ID)
	}

	// Validate certificate
	if cert.ID == "" {
		return fmt.Errorf("certificate ID cannot be empty")
	}
	if cert.ProductID == "" {
		return fmt.Errorf("product ID cannot be empty")
	}
	if cert.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	if cert.Status == "" {
		return fmt.Errorf("status cannot be empty")
	}

	// Set certificate
	return k.SetCertificate(ctx, cert)
}

// DeleteCertificate deletes a certificate
func (k Keeper) DeleteCertificate(ctx sdk.Context, id string) error {
	// Check if certificate exists
	if _, exists := k.GetCertificate(ctx, id); !exists {
		return fmt.Errorf("certificate with ID %s does not exist", id)
	}

	store := ctx.KVStore(k.storeKey)
	key := types.GetCertificateKey(id)
	store.Delete(key)
	return nil
}

// VerifyCertificate verifies a certificate
func (k Keeper) VerifyCertificate(ctx sdk.Context, certificateID, verifier, status, notes string) error {
	// Check if certificate exists
	if _, exists := k.GetCertificate(ctx, certificateID); !exists {
		return fmt.Errorf("certificate with ID %s does not exist", certificateID)
	}

	// Create verification
	verification := types.NewCertificateVerification(certificateID, verifier, status, notes)
	return k.SetVerification(ctx, verification)
}

// InitGenesis initializes the genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	// Set parameters
	k.SetParams(ctx, genState.Params)

	// Set certificates
	for _, cert := range genState.Certificates {
		k.SetCertificate(ctx, cert)
	}

	// Set verifications
	for _, verification := range genState.Verifications {
		k.SetVerification(ctx, verification)
	}
}

// ExportGenesis exports the genesis state
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	// Get parameters
	params, _ := k.GetParams(ctx)

	// Get all certificates
	certificates, _ := k.GetAllCertificates(ctx)

	// Get all verifications
	verifications, _ := k.GetAllVerifications(ctx)

	return &types.GenesisState{
		Certificates:  certificates,
		Verifications: verifications,
		Params:        params,
	}
}

// BeginBlocker is called at the beginning of every block
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// Log block height
	ctx.Logger().Info(fmt.Sprintf("ProductCertificate BeginBlocker: Block %d", ctx.BlockHeight()))

	// Perform any necessary operations at the beginning of each block
	// This could include:
	// - Certificate expiration checks
	// - Auto-verification processes
	// - Emitting events

	// Example: Emit a block start event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProductCertificateCreated,
			sdk.NewAttribute(types.AttributeKeyCreatedAt, strconv.FormatInt(ctx.BlockTime().Unix(), 10)),
		),
	)
}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) []abci.ValidatorUpdate {
	// Log block height
	ctx.Logger().Info(fmt.Sprintf("ProductCertificate EndBlocker: Block %d", ctx.BlockHeight()))

	// Perform any necessary operations at the end of each block
	// This could include:
	// - Certificate cleanup
	// - Verification processing
	// - Emitting events

	// Example: Emit a block end event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProductCertificateUpdated,
			sdk.NewAttribute(types.AttributeKeyUpdatedAt, strconv.FormatInt(ctx.BlockTime().Unix(), 10)),
		),
	)

	// Return validator updates (if any)
	return []abci.ValidatorUpdate{}
}
