// lesser defines a type-parameterized interface called "Interface" with one
// method, Less, which returns a boolean for whether the caller, of type T,
// is less than some other instance of T. This is blatantly stolen from
// Robert Griesemer's talk at Gophercon about the type parameters proposal
// in 2020. 
//
// The library then defines a type-parameterized wrapper called "Basic"
// over all Ordered built-in types, allowing these types to implement
// Interface so that Interface can function as a uniform constraint for
// all orderable types, permitting them to all be stored in equivalent
// ordered collections. See the README for rationale.
package lesser

// Interface is an interface that wraps the Less method.
//
// Less compares a caller of type T to some other variable of type T,
// returning true if the caller is the lesser of the two values, and false
// otherwise. If the two values are equal it returns false.
type Interface[T any] interface {
	Less(other T) bool
}

// Basic is a parameterized type that abstracts over the entire class of
// Ordered types (the set of Go built-in types which respond to the <
// operator), and exposes this behavior via a Less method so that they
// fall under the lesser.Interface constraint.
type Basic[N constraints.Ordered] struct{ Val N }

// Less implements Interface[Basic[N]] for Basic[N]. Returns true if the value
// of the caller is less than that of the parameter; otherwise returns
// false.
func (x Basic[N]) Less(y Basic[N]) bool {
	return x.Val < y.Val
}
