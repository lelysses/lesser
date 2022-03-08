# lesser


## What is this?

`lesser` defines a type-parameterized interface with one method, `Less`, which returns a boolean for whether the caller, of type `T`, is less than some other instance of `T`. This is blatantly stolen from [Robert Griesemer's talk at Gophercon 2020](https://www.youtube.com/watch?v=TborQFPY2IM) about the type parameters proposal. Probably more controversially, this library also defines a wrapper called `Basic` over the built-in numerical types, exposing the underlying `<` operator through this `Less` method. The reasoning for this follows.

## Why is this?

In the near future, [the `constraints` package](https://pkg.go.dev/golang.org/x/exp/constraints) will be in the standard library. Constraints defines an constraint called `Ordered`, which matches all types that respond to the `<` operator. This is explicitly only suitable for built-in types since user-defined types (structs, etc) may not respond to operators. Conversely, however, built-ins may not respond to methods. As a consequence, any sorting function or ordered collection must specify different versions (`BasicSort` vs `Sort`, `BasicHeap` vs `Heap`, etc) for built-in types and user-defined types. This is an untenable position, and makes the type parameters proposal insufficient and unworkable for a large swathe of the problems it exists to solve (anything involving ordering).

There appear to be three different possible directions to go here:

1. Every time you want to write an ordered collection or a sorting function or anything depending on natural ordering functionality, you have to copy paste it and have one copy use `Lesser` as the constraint and the other use `Ordered` as the constraint. That's not a reasonable request to make. Without this library or something like it ***you are here***. As of the end of last year this was Robert Griesemer's suggested solution, and to the best of my knowledge it still is. It's unclear what the purpose of generics is if we're still going to be forced to use code generation or copy-pasting.

2. Write a library very much like this one, but which defines a series of wrappers for each of the built-in types, each of which has a `Less` method, so then when you want to use `int`s in an ordered collection you'd convert them all to `lesser.Int` first. This is scary and gives me serious Java heebie-jeebies. I do not want Go to turn into a world where we are all using special magic wrappers for every basic data type at all times or something.

3. Exactly what this library is, which is the same as (2) but instead of defining the types manually, you create them at compile-time using type parameters with something like `lesser.Basic[int]` instead of `lesser.Int`. It feels far less likely this way that it'll develop into the scary situation I just described.

4. [SEE ADDENDUM](#addendum) FOR IAN LANCE TAYLOR'S EPIC AND COOL ALTERNATIVE

I don't know if something like this is coming to the standard library, but I'm unwilling to wait until it is. For a while it felt like the `Lesser` interface itself was likely to make it into the `constraints` library, but I'm not sure about that now. All discussion of how the hell to actually use generics has been explicitly pushed until some time after generics are released into the wild. I hate the antichrist.

This doesn't provide a Java-style `compareTo` or the like, but merely exposes a `Less` method that is fit for sorting algorithms and ordered collections. Look elsewhere for something more robust. This is a plug for a hole in the current generics implementation. This is not a "framework"; this is not a "platform".

## How do I use it?

Well first you'll need to get Go 1.18, which has not yet been released. Release candidate 1 was released in February, and you can install it as follows:

```bash
go install golang.org/dl/go1.18rc1@latest
go1.18rc1 download
```

Now you can use the `go1.18rc1` command as an alternative to the `go` command, but with generics enabled.

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
h.Push(lesser.Basic[int]{1})
h.Push(lesser.Basic[int]{2})
h.Push(lesser.Basic[int]{3})
```

To get the value back out of a `Basic`, you just poll the `Val` attribute. It used to be as simple as a cast back to the correct type, but then [this issue](https://github.com/golang/go/issues/45639) happened because we can't have nice things, so `Basic` had to be changed from an alias type to a struct. The ergonomics have suffered as a consequence.

If you want to use it with your own wacky type, `Basic` doesn't apply, since your own type won't implement `Ordered`. Instead, you can just give it a `Less` method.

```go
type Name struct {
	Last  string
	First string
}

func (this Name) Less(other Name) bool {
	if this.Last == other.Last {
		return this.First < other.First
	}
	return this.Last < other.Last
}
```

And now you should be able to directly push some `Name`s onto a heap of `Name`s.

```go
var h Heap[Name]
h.Push(Name{"Marx", "Karl"})
h.Push(Name{"Marx", "Groucho"})
```

## Okay, but really, why?

I don't know if I really think this library is a good idea. I don't know how Go-y it is. The definition of """idiomatic""" Go will probably change in the coming year, and maybe this is the way things are going. I don't know how I feel about that. More than anything I want this library to get other people thinking about how this is actually going to work, because without something like this library it's not possible to make simple ordered collections that operate on both user-defined and built-in data types.

## Addendum

Ian Lance Taylor made (what I think is) a great suggestion [here](https://github.com/golang/go/issues/47632#issuecomment-897168431), arguing for constraining containers to `any` and passing a comparison function to the constructor. This is a good idea and subjectively feels more idiomatic than what I have done here. I think?
