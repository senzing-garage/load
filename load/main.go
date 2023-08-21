package load

import "context"

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Load interface {
	Load(context.Context) bool
}

// mover is 6601:  https://github.com/Senzing/knowledge-base/blob/main/lists/senzing-product-ids.md
const MessageIdFormat = "senzing-6601%04d"
