package types

import (
	"fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// ProductCertificate module constants
const (
	ModuleName   = "product_certificate"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

// Event types
const (
	EventTypeProductCertificateCreated  = "product_certificate_created"
	EventTypeProductCertificateUpdated  = "product_certificate_updated"
	EventTypeProductCertificateDeleted  = "product_certificate_deleted"
	EventTypeProductCertificateVerified = "product_certificate_verified"
)

// Event attribute keys
const (
	AttributeKeyProductID     = "product_id"
	AttributeKeyCertificateID = "certificate_id"
	AttributeKeyOwner         = "owner"
	AttributeKeyStatus        = "status"
	AttributeKeyCreatedAt     = "created_at"
	AttributeKeyUpdatedAt     = "updated_at"
)

// ProductCertificate represents a product certificate
// ROOT FIX: Added protobuf tags to fix unmarshal panic
// COSMOS SDK 0.53.4: Must have protobuf tags if implementing ProtoMessage()
type ProductCertificate struct {
	ID         string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ProductID  string `protobuf:"bytes,2,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Owner      string `protobuf:"bytes,3,opt,name=owner,proto3" json:"owner,omitempty"`
	Status     string `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	Metadata   string `protobuf:"bytes,5,opt,name=metadata,proto3" json:"metadata,omitempty"`
	CreatedAt  int64  `protobuf:"varint,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt  int64  `protobuf:"varint,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	VerifiedAt int64  `protobuf:"varint,8,opt,name=verified_at,json=verifiedAt,proto3" json:"verified_at,omitempty"`
	ExpiresAt  int64  `protobuf:"varint,9,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
}

// ProtoMessage implements proto.Message interface
func (pc *ProductCertificate) ProtoMessage() {}

// Reset implements proto.Message interface
func (pc *ProductCertificate) Reset() {
	*pc = ProductCertificate{}
}

// String implements proto.Message interface
func (pc *ProductCertificate) String() string {
	return fmt.Sprintf("ProductCertificate{ID: %s, ProductID: %s, Owner: %s, Status: %s, Metadata: %s, CreatedAt: %d, UpdatedAt: %d, VerifiedAt: %d, ExpiresAt: %d}",
		pc.ID, pc.ProductID, pc.Owner, pc.Status, pc.Metadata, pc.CreatedAt, pc.UpdatedAt, pc.VerifiedAt, pc.ExpiresAt)
}

// CertificateVerification represents a certificate verification
type CertificateVerification struct {
	CertificateID string `json:"certificate_id"`
	Verifier      string `json:"verifier"`
	Status        string `json:"status"`
	VerifiedAt    int64  `json:"verified_at"`
	Notes         string `json:"notes"`
}

// ProtoMessage implements proto.Message interface
func (cv *CertificateVerification) ProtoMessage() {}

// Reset implements proto.Message interface
func (cv *CertificateVerification) Reset() {
	*cv = CertificateVerification{}
}

// String implements proto.Message interface
func (cv *CertificateVerification) String() string {
	return fmt.Sprintf("CertificateVerification{CertificateID: %s, Verifier: %s, Status: %s, VerifiedAt: %d, Notes: %s}",
		cv.CertificateID, cv.Verifier, cv.Status, cv.VerifiedAt, cv.Notes)
}

// GenesisState represents the genesis state of the product_certificate module
type GenesisState struct {
	Certificates  []ProductCertificate      `json:"certificates"`
	Verifications []CertificateVerification `json:"verifications"`
	Params        Params                    `json:"params"`
}

// ProtoMessage implements proto.Message interface
func (gs *GenesisState) ProtoMessage() {}

// Reset implements proto.Message interface
func (gs *GenesisState) Reset() {
	*gs = GenesisState{}
}

// String implements proto.Message interface
func (gs *GenesisState) String() string {
	return fmt.Sprintf("GenesisState{Certificates: %v, Verifications: %v, Params: %v}",
		gs.Certificates, gs.Verifications, gs.Params)
}

// Params defines the parameters for the product_certificate module
type Params struct {
	MaxCertificates   uint32 `json:"max_certificates"`
	MinMetadataLength uint32 `json:"min_metadata_length"`
	MaxMetadataLength uint32 `json:"max_metadata_length"`
	VerificationFee   string `json:"verification_fee"`
	ExpirationTime    int64  `json:"expiration_time"`
	AutoVerification  bool   `json:"auto_verification"`
}

// DefaultParams returns the default parameters for the product_certificate module
func DefaultParams() Params {
	return Params{
		MaxCertificates:   1000,
		MinMetadataLength: 10,
		MaxMetadataLength: 1000,
		VerificationFee:   "1000000",   // 1 USC
		ExpirationTime:    86400 * 365, // 1 year
		AutoVerification:  false,
	}
}

// Validate validates the parameters
func (p Params) Validate() error {
	if p.MaxCertificates == 0 {
		return fmt.Errorf("max certificates must be positive")
	}
	if p.MinMetadataLength == 0 {
		return fmt.Errorf("min metadata length must be positive")
	}
	if p.MaxMetadataLength < p.MinMetadataLength {
		return fmt.Errorf("max metadata length must be greater than min metadata length")
	}
	if p.VerificationFee == "" {
		return fmt.Errorf("verification fee cannot be empty")
	}
	if p.ExpirationTime <= 0 {
		return fmt.Errorf("expiration time must be positive")
	}
	return nil
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte("MaxCertificates"), &p.MaxCertificates, validateUint32),
		paramtypes.NewParamSetPair([]byte("MinMetadataLength"), &p.MinMetadataLength, validateUint32),
		paramtypes.NewParamSetPair([]byte("MaxMetadataLength"), &p.MaxMetadataLength, validateUint32),
		paramtypes.NewParamSetPair([]byte("VerificationFee"), &p.VerificationFee, validateString),
		paramtypes.NewParamSetPair([]byte("ExpirationTime"), &p.ExpirationTime, validateInt64),
		paramtypes.NewParamSetPair([]byte("AutoVerification"), &p.AutoVerification, validateBool),
	}
}

// Validation functions
func validateString(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateUint32(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateInt64(i interface{}) error {
	_, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

// NewProductCertificate creates a new ProductCertificate
func NewProductCertificate(id, productID, owner, status, metadata string) ProductCertificate {
	now := time.Now().Unix()
	return ProductCertificate{
		ID:         id,
		ProductID:  productID,
		Owner:      owner,
		Status:     status,
		Metadata:   metadata,
		CreatedAt:  now,
		UpdatedAt:  now,
		VerifiedAt: 0,
		ExpiresAt:  now + 86400*365, // 1 year from now
	}
}

// NewCertificateVerification creates a new CertificateVerification
func NewCertificateVerification(certificateID, verifier, status, notes string) CertificateVerification {
	return CertificateVerification{
		CertificateID: certificateID,
		Verifier:      verifier,
		Status:        status,
		VerifiedAt:    time.Now().Unix(),
		Notes:         notes,
	}
}

// DefaultGenesis returns the default genesis state for the product_certificate module
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Certificates:  []ProductCertificate{},
		Verifications: []CertificateVerification{},
		Params:        DefaultParams(),
	}
}

// Validate validates the genesis state
func (gs GenesisState) Validate() error {
	// Validate parameters
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	// Validate certificates
	seenIDs := make(map[string]bool)
	for _, cert := range gs.Certificates {
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

		// Check for duplicate IDs
		if seenIDs[cert.ID] {
			return fmt.Errorf("duplicate certificate ID: %s", cert.ID)
		}
		seenIDs[cert.ID] = true
	}

	// Validate verifications
	for _, verification := range gs.Verifications {
		if verification.CertificateID == "" {
			return fmt.Errorf("verification certificate ID cannot be empty")
		}
		if verification.Verifier == "" {
			return fmt.Errorf("verification verifier cannot be empty")
		}
		if verification.Status == "" {
			return fmt.Errorf("verification status cannot be empty")
		}
		if verification.VerifiedAt <= 0 {
			return fmt.Errorf("verification timestamp must be positive")
		}
	}

	return nil
}
