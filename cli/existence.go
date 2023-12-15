package cli

// Existence will set ArgRequires to be only one of the three existence modes so
// it can be checked with case statement or simple equals (=).
func Existence(requires ArgRequires) ArgRequires {
	return requires & (OkToExist | MustExist | MustNotExist)
}
