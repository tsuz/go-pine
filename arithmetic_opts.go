package pine

// ArithmeticOpts defines handling of special cases on arithmetic operations
type ArithmeticOpts struct {
	NilHandlInst NilHandlInst
}

// NilHandlInst defines how to handle if any of arithmetic values are nil
type NilHandlInst int

const (
	// NilValueReturnNil returns nil if any of arithmetic values are nil
	NilValueReturnNil NilHandlInst = iota
	// NilValueReturnZero returns zero if any of arithmetic values are nil
	NilValueReturnZero
)
