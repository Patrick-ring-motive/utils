
# Utils Package Documentation

The infamous utils package.
Here there be dragons.
Use with extreme caution.
This package has many functions that allow you to completely sidestep the go type system.
They attempt to do so in a way that is as safe as possible.
If used correctly, they are very powerful.
If you are considering using these tools then you have probably missed the correct way to do  what you want.
If you have to use this, please consider pairing it with the SafeOps package utilities to wrap your code in a safe way.

Some notes on some unusual patterns that I employ.

`func(T){}` is a function that takes a single typed parameter and does nothing.
This is my way of passing a type reference around without having to instantiate it.
Mostly it is abstracted away but you will see it show in parameters sometimes.
Utils has a convenient function to do this called TypeRef which you can use to generate these references.
You can generate them from an abstract type like so: `TypeRef[int]()` or from a concrete type like so: `TypeRef(0)`.

`*[1]T` is a pointer to an array of length 1 of type T. This ensures that values are passed by reference and not
copied and facilitates passing values back from a defer/recover block.

That brings me to the unconventional error handling pattern.

```
func Example(input inputType)outputType{
  var z outputType
  carrier := *[1]outputType{z}
  example(carrier,input)
  return carrier[0]
}
func example(carrier *[1]outputType,input inputType){
	defer func() {
		if r := recover(); r != nil {
			carrier[0] = fallbackValue
		}
	}()
	carrier[0] = attemptSomething(input)
	if(carrier[0] == nil){
		carrier[0] = fallbackValue
	}
}
```

This is a pattern than handles errors by "returning" a fallback value
on panic or nil. This pattern is difficult to abstract out because go generics dont handle various function types well and the defer needs to happen one function call deeper than where we intend to recover a panic.