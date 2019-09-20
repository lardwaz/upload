package box

import "os"

// MockAsset is a mock implementation of sdk.AssetBoxer
type MockAsset struct{}

// NewMockAsset returns a new MockAsset
func NewMockAsset() *MockAsset {
	return &MockAsset{}
}

// Open implements the sdk.AssetBoxer
func (m *MockAsset) Open(name string) (*os.File, error) {
	return os.Open(name)
}
