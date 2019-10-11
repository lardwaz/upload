package upload

import "os"

// AssetBoxer represents an asset box
type AssetBoxer interface {
	Open(string) (*os.File, error)
}
