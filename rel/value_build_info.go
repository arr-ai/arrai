package rel

// BuildInfoTuple represents arrai build information.
type BuildInfoTuple struct {
	Tuple
}

// String returns a string which is prettified by readable format.
func (t BuildInfoTuple) String() string {
	str, err := PrettifyString(t, 0)
	if err != nil {
		return err.Error()
	}

	return str
}

// NewBuildInfoTuple returns a new BuildInfoTuple.
func NewBuildInfoTuple(info Tuple) BuildInfoTuple {
	return BuildInfoTuple{info}
}
