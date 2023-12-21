package cli

// ArgExistence will return only one of the three existence modes, so it can be
// checked with case statement or simple equals (=).
func ArgExistence(requires ArgRequires) ArgRequires {
	return requires & (IgnoreExists | MustExist | NotExist)
}

// ArgEmptiness will return only one of the three emptiness modes, so it can be
// checked with case statement or simple equals (=).
func ArgEmptiness(requires ArgRequires) ArgRequires {
	return requires & (EmptyOk | MustBeEmpty | NotEmpty)
}

// ArgValidation ArgEmptiness will return MustValidate if set, otherwise 0, so it
// can be checked with simple equals (=).
func ArgValidation(requires ArgRequires) ArgRequires {
	return requires & MustValidate
}

func AndRequires(requires ...ArgRequires) (ar ArgRequires) {
	for _, r := range requires {
		ar &= r
	}
	return ar
}
