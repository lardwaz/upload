package option

var (
	defaultENVOptions = &OptsENV{} // prod = false
)

// OptsENV is an implementation of OptionsENV
type OptsENV struct {
	prod bool
}

// IsPROD returns whether if PROD
func (o *OptsENV) IsPROD() bool {
	return o.prod
}

// SetPROD set current ENV to PROD
func (o *OptsENV) SetPROD(b bool) {
	o.prod = b
}
