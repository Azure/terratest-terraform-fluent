// Package ops contains the operations that can be performed on the data extracted from
// the plan, output, or state.	E.g. HasValue, Exists, DoesNotExist, etc.
//
// It is also possible to query the data using gjson queries.
// This will return a ops.Operative type which can be used to compare the data.

package ops
