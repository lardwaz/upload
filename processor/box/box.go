package box

import sdk "go.lsl.digital/lardwaz/sdk/upload"

// Asset represents an assetbox
var Asset sdk.AssetBoxer

// Set sets the asset box to retrieve static assets
func Set(assetBox sdk.AssetBoxer) {
	Asset = assetBox
}
