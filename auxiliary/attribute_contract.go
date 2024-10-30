package auxiliary

type ContractResult struct {
	Attributes map[string]string
	Valid      bool
}

type Preconditions map[string]func(string) bool

func AttributeContract(attrs map[string]string, preconditions Preconditions) ContractResult {
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
