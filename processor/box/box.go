package box

import "go.lsl.digital/lardwaz/upload"

// Asset represents an assetbox
var Asset upload.AssetBoxer

// Set sets the asset box to retrieve static assets
func Set(assetBox upload.AssetBoxer) {
	Asset = assetBox
}
