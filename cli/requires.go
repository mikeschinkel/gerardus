package cli

// ExistenceInDB will return only one of the three existence modes so
// it can be checked with case statement or simple equals (=).
func ExistenceInDB(requires ArgRequires) ArgRequires {
	return requires & (IgnoreCheck | MustPassCheck | MustFailCheck)
}

// ArgEmptiness will return only one of the three emtiness modes so
// it can be checked with case statement or simple equals (=).
func ArgEmptiness(requires ArgRequires) ArgRequires {
	return requires & (EmptyOk | MustBeEmpty | NotEmpty)
}

func AndRequires(requires ...ArgRequires) (ar ArgRequires) {
	for _, r := range requires {
		ar &= r
	}
	return ar
}
