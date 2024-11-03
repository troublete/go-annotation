package auxiliary

type ContractResult struct {
	Attributes map[string]string
	Valid      bool
}

var (
	PreconditionAlwaysAllow = func(_ string) bool {
		return true
	}
)

type Attributes map[string]string

type Preconditions map[string]func(string) bool

// AttributeContract should validate that passed attributes are proper, all un verified inputs
// the ones without precondition are filtered out
// this can be used for checking input parameters to validate them to assure a smooth runtime
func AttributeContract(attrs Attributes, preconditions Preconditions) ContractResult {
	args := map[string]string{}
	valid := true
	for name, p := range preconditions {
		if v, ok := attrs[name]; ok {
			if p(v) {
				args[name] = v
			} else {
				valid = false
			}
		} else {
			valid = false
		}
	}

	return ContractResult{
		Attributes: args,
		Valid:      valid,
	}
}
