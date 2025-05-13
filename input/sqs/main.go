package sqs

import "errors"

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// load is 6201:  https://github.com/senzing-garage/knowledge-base/blob/main/lists/senzing-product-ids.md
const ComponentID = 6201

// Log message prefix.
const (
	Prefix           = "load: "
	OptionCallerSkip = 4
)

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Error level ranges and usage:
// Level 	Range 		Use 							Comments.
var IDMessages = map[int]string{
	// TRACE 	0000-0999 	Entry/Exit tracing 				May contain sensitive data.
	// DEBUG 	1000-1999 	Values seen during processing 	May contain sensitive data.
	// INFO 	2000-2999 	Process steps achieved
	2004: Prefix + "Unable to retrieve the config: %v",
	2005: Prefix + "Unable to reach Sz: %v",
	2006: Prefix + "Unable to initialize Sz: %v",
	2999: Prefix + "So long and thanks for all the fish.",
	// WARN 	3000-3999 	Unexpected situations, but processing was successful
	// ERROR 	4000-4999 	Unexpected situations, processing was not successful
	// FATAL 	5000-5999 	The process needs to shutdown
	5000: Prefix + "Fatal error, there was an unexpected issue; please report this as a bug: %v",
	// PANIC 	6000-6999 	The underlying system is at issue
	//
	//			8000-8999 	Reserved for observer messages
}

// Status strings for specific messages.
var IDStatuses = map[int]string{}

var errForPackage = errors.New("input.sqs")
