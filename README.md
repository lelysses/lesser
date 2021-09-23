# lesser

## What is this?

`lesser` defines a type-parameterized interface with one method, `Less`, which returns a boolean for whether the caller, of type `T`, is less than some other instance of `T`. This is blatantly stolen from [Robert Griesemer's talk at Gophercon 2020](https://www.youtube.com/watch?v=TborQFPY2IM) about the type parameters proposal. Probably more controversially, this library also defines a wrapper called `Basic` over the built-in numerical types, exposing the underlying `<` operator through this `Less` method. The reasoning for this follows.

## Why is this?

In the near future the standard library will, in a `constraints` package, have something like this (from [the type parameters proposal](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#operations-based-on-type-sets)):

```go
// Ordered is a type constraint that matches any ordered type.
// An ordered type is one that supports the <, <=, >, and >= operators.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}
```

This is explicitly only suitable for built-in types since user-defined types (structs etc) may not respond to operators. Conversely, however, built-ins may not respond to methods. As a consequence, any sorting function or ordered collection must specify different versions (`BasicSort` vs `Sort`, `BasicHeap` vs `Heap`, etc) for built-in types and user-defined types. This is an untenable position, and makes the type parameters proposal insufficient and unworkable for a large swathe of the problems it exists to solve (anything involving ordering).

There appear to be three different possible directions to go here:

1. Every time you want to write an ordered collection or a sorting function or anything depending on natural ordering functionality, you have to copy paste it and have one copy use `Lesser` as the constraint and the other use `Ordered` as the constraint. That's not a reasonable request to make. Without this library or something like it ***you are here***. As of the end of last year this was Robert Griesemer's suggested solution, and to the best of my knowledge it still is.

2. Write a library very much like this one, but which defines a series of wrappers for each of the built-in types, each of which has a `Less` method, so then when you want to use `int`s in an ordered collection you'd convert them all to `lesser.Int` first. This is scary and gives me serious Java heebie-jeebies. I do not want Go to turn into a world where we are all using special magic wrappers for every basic data type at all times or something.

3. Exactly what this library is, which is the same as the former but instead of defining the types manually, you create them at compile-time using type parameters with something like `lesser.Basic[int]` instead of `lesser.Int`. It feels far less likely this way that it'll develop into the scary situation I just described.

I don't know if something like this is coming to the standard library, but I'm unwilling to wait until it is. I have a sneaking suspicion the `Lesser` interface itself is likely to make it into the `constraints` library alongside `Ordered`, at which point I'll likely remove each of them here. This doesn't provide a Java-style `compareTo` or the like, but merely exposes a `Less` method that is fit for sorting algorithms and ordered collections. Look elsewhere for something more robust. This is a plug for a hole in the current generics implementation. This is not a "framework"; this is not a "platform".

## How do I use it?

Well first you'll need to get Go 1.18, which might seem hard since it hasn't been released yet. You can get what amounts to the dev branch of the Go compiler using [`gotip`](https://pkg.go.dev/golang.org/dl/gotip):

```bash
go install golang.org/dl/gotip@latest
gotip download
```

Now you can use the `gotip` command as an alternative to the `go` command, but with generics enabled.

If you want to build this into your project, use `gotip` when updating or initializing your `go.mod`, or otherwise make sure that your `go.mod` notes that we're using `go 1.18`, and make sure you build or run with `gotip`, not `go`, or you'll get compiler syntax errors.

## Okay, but how do I *use* it?

If you want to build an ordered collection, say a binary heap, do something like this:

```go
package foo

import "github.com/lelysses/lesser"

// a cool heap
type Heap[T lesser.Interface[T]] []T

// push an item onto it
func (h *Heap[T]) Push(val T) {
	// let's pretend ...
}
```

Now your heap works for user-defined types which expose a `Less[T] bool` method, but it also works for built-in types, although the mechanism for using built-in types is somewhat more complicated.

If you want to initialize the `Heap` we defined above for the built-in concrete type `int`, and then push the numbers 1, 2, and 3 onto it, you'd do the following:

```go
var h Heap[lesser.Basic[int]]
h.Push(lesser.Basic[int](1))
h.Push(lesser.Basic[int](2))
h.Push(lesser.Basic[int](3))
```

If you want to use it with your own wacky type, you won't need to (and can't) use `Basic`, since your own type won't implement `Ordered`. Instead, you can just do something like this:

```go
type Book struct {
	Title string
	ISBN uint64
}

func (b Book) Less(other Book) bool {
	return p.ISBN < other.ISBN
}
```

And now you should be able to directly push some `Book`s onto a heap of `Book`s.

```go
var h Heap[Book]
h.Push(Book{"Cool Book", 1234567890123})
h.Push(Book{"Another Cool Book", 2345678901234})
h.Push(Book{"And Still Another Cool Book", 3456789012345})
```

## Okay, but really, why?

I don't know if I really think this library is a good idea. I don't know how Go-y it is. The definition of """idiomatic""" Go will probably change in the coming year, and maybe this is the way things are going. I don't know how I feel about that. More than anything I want this library to get other people thinking about how this is actually going to work, because without something like this library it's not possible to make simple ordered collections that operate on both user-defined and built-in data types.
