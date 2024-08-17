package utils

/*
The infamous utils package.
Here there be dragons.
Use with extreme caution.
This package has many functions that allow you to completely sidestep the go type system.
They attempt to do so in a way that is as safe as possible.
If used correctly, they are very powerful.
If you are considering using these tools then you have probably missed the correct way to do  what you want.
If you have to use this, please consider pairing it with the SafeOps package utilities to wrap your code in a safe way.
*/
import (
	"reflect"
	"unsafe"
)

/*
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
*/

/*Type reference via `func(T)`*/
func TypeRef[T any](t ...T) func(T) {
	return func(T) {}
}

/*
ConvertType is the start of a chain of type conversion methods that get increasingly more aggressive in their attempts to convert the type.
This attenpts basic type conversion `type(i)` and falls back to AsserType on failure.
*/
func ConvertType[I any, T any](i I, t ...func(T)) T {
	switch v := any(i).(type) {
	case T:
		return T(v)
	default:
		return T(AssertType[I, T](i))
	}
}

/*Shorthand for ConvertType but takes only one parameter*/
func Convert[I any, T any](i I) T {
	switch v := any(i).(type) {
	case T:
		return T(v)
	default:
		return T(AssertType[I, T](i))
	}
}

/*
AssertType attempts the basic type assertiont `i.(type)`.
It falls back SwitchType on failure.
*/
func AssertType[I any, T any](i I, t ...func(T)) T {
	r, ok := any(i).(T)
	if !ok {
		r = SwitchType[I, T](i)
	}
	return r
}

/*Shorthand for AssertType but takes only one parameter*/
func Assert[F any, T any](f F) T {
	return AssertType[F, T](f)
}

/*AssertTypeUnsafe attempts type assertion and returns whatever the result even if it fails*/
func AssertTypeUnsafe[T any](i interface{}, t ...func(T)) T {
	r, ok := i.(T)
	AllowUnused(ok)
	return r
}

/*Shorthand for AssertTypeUnsafe but takes only one parameter*/
func AssertUnsafe[T any](f any) T {
	return AssertTypeUnsafe[T](f)
}

/*SwitchType uses a type switch to try and convert types. It falls back to ForceType on failure to match type.*/
func SwitchType[S any, T any](s S, t ...func(T)) T {
	i := AsInterface(s)
	switch v := i.(type) {
	case T:
		return v
	default:
		return ForceType[S, T](s)
	}
}

/*Shorthand for SwitchType but takes only one parameter*/
func Switch[F any, T any](f F) T {
	return SwitchType[F, T](f)
}

/*
ForceType uses unsafe.Pointer with generics to forcibly convert between types.
It falls back to ForceRawType on panic
*/
func ForceType[F any, T any](f F, t ...func(T)) T {
	var z T
	a := &[1]T{z}
	forceType(a, f)
	return a[0]
}
func forceType[F any, T any](a *[1]T, f F) {
	defer func() {
		if r := recover(); r != nil {
			a[0] = ForceRawType(f, func(T) {})
		}
	}()
	a[0] = *(*T)(unsafe.Pointer(&f))
}

/*Shorthand for ForceType but takes only one parameter*/
func Force[F any, T any](f F) T {
	return ForceType[F, T](f)
}

/*
ForceRawType uses unsafe.Pointer and the `any` type to forcibly convert types.
On panic it falls back to the zero value of the target type which ends the chain
*/
func ForceRawType[T any](f any, t ...func(T)) T {
	var z T
	a := &[1]T{z}
	forceType(a, f)
	return a[0]
}
func forceRawType[T any](a *[1]T, f any) {
	defer func() {
		if r := recover(); r != nil {
			a[0] = ZeroOfType(func(T) {})
		}
	}()
	a[0] = *(*T)(unsafe.Pointer(&f))
}

/*Alias for ForceRawType but only 1 parameter*/
func Coerce[T any](f any) T {
	return ForceRawType[T](f)
}

/*
ZeroOfType returns the zero value of the target type.
You can do `ZeroOfType(func(T){})` or `ZeroOfType[T]()` if your type is inferable
*/
func ZeroOfType[T any](t ...func(T)) T {
	var x T
	return x
}

/*Shorthand version for ZeroOfType but accepts no parameters*/
func ZeroOf[T any]() T {
	var x T
	return x
}

/*
NilOfType takes an any interface of nil and coerces it into the target type
You can do `NilOfType(func(T){})` or `NilOfType[T]()` if your type is inferable
*/
func NilOfType[T any](t ...func(T)) T {
	n := any(nil)
	return *(*T)(unsafe.Pointer(&n))
}

/*Shorthand version for NilOfType but accepts no parameters*/
func NilOf[T any]() T {
	n := any(nil)
	return *(*T)(unsafe.Pointer(&n))
}

/*
This does not work yet. Don't use it.

InitOfType uses reflection to detect if the target type is one that has a nil zero value.
If it does then an empty shell ov the target type is returned.
Otherwise the zero value is returned.
For example if the target type is a slice, then an empty slice is returned
*/
func InitOfType[T any](t ...func(T)) T {
	var instance T
	typ := reflect.TypeOf(instance)

	switch typ.Kind() {
	case reflect.Slice:
		val := reflect.MakeSlice(typ, 0, 0)
		return val.Interface().(T)
	case reflect.Map:
		val := reflect.MakeMap(typ)
		return val.Interface().(T)
	case reflect.Ptr:
		val := reflect.New(typ.Elem())
		return val.Interface().(T)
	case reflect.Chan:
		val := reflect.MakeChan(typ, 0)
		return val.Interface().(T)
	case reflect.Func:
		val := reflect.MakeFunc(typ, func(args []reflect.Value) (results []reflect.Value) {
			returns := make([]reflect.Value, typ.NumOut())
			for i := range returns {
				returns[i] = reflect.Zero(typ.Out(i))
			}
			return returns
		})
		return val.Interface().(T)
	case reflect.Interface:
		if typ.NumMethod() == 0 {
			return reflect.Zero(typ).Interface().(T)
		}
	default:
		return instance
	}
	return instance
}

/*Shorthand version for InitOfType but accepts no parameters*/
func InitOf[T any]() T {
	return InitOfType[T]()
}

/*
AllowUnused is a function that allows you to ignore variables.
No action is performed on the values passed in.
The compiler just stops complaining.
Yes you can use this to ignore errors. Probably not a good idea.
*/
func AllowUnused(a ...any) {}

/*
Returns a pointer of the given value.
Convenient for inlining pointer creation.
*/
func Ptr[T any](value T) *T {
	return &value
}

/*
AsInterface returns the given input as an interface.
Helps with type hints or restrictions
*/
func AsInterface(i interface{}) interface{} {
	return i
}

/*
AsAny returns the given input as an interface.
Helps with type hints or restrictions
*/
func AsAny(i any) any {
	return i
}

/*
AsAny returns the given input as an generic interface.
Helps with type hints or restrictions
*/
func AsGeneric[generic any](g generic) generic {
	return g
}

/*
Invoke takes an interface of a function and an interface of a list of input parapeters and attempts to execute the function using those parameters.
This is enabled using reflection.
This function is highly unstable and allows for all kinds of strange possibilities.
*/
func Invoke(fn interface{}, args interface{}) any {
	fnVal := reflect.ValueOf(fn)
	fnType := fnVal.Type()
	numIn := fnType.NumIn()
	in := make([]reflect.Value, numIn)
	for i, arg := range args.([]any) {
		argVal := reflect.ValueOf(arg)
		in[i] = argVal
	}
	out := fnVal.Call(in)
	result := make([]interface{}, len(out))
	for i, o := range out {
		result[i] = o.Interface()
	}
	return result[0]
}
