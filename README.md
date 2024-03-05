
# Utils Package Documentation

## Overview
The `utils` package offers a collection of utility functions designed to facilitate advanced type manipulation, pointer generation, and development conveniences in Go. These functions leverage unsafe operations to provide flexibility beyond the Go type system's safety constraints. Users should apply these utilities with understanding of their underlying behavior and potential side effects.

### ForceType
`ForceType` retrieves the original type of an object previously converted into an 'any' type. This function is particularly useful for regaining type-specific operations on a value that was stored or passed as an 'any' type, ensuring the original type's memory representation is correctly interpreted.

```go
var originalType float64 = 42.0
var genericType any = originalType // originalType stored as 'any'
// Retrieve the original type from genericType
retrievedType := utils.ForceType[any, float64](genericType, func(float64) {})
fmt.Println(retrievedType) // Outputs: 42.0
```

### ForceRawType
`ForceRawType` is an earlier version of `ForceType` that might exhibit unexpected behavior, such as returning a pointer instead of a value, or vice versa. It is retained for backward compatibility or specific scenarios where its unique side effects are desirable.

```go
// Example use case intentionally omitted due to potential side effects and need for specific understanding.
```

### NilOfType
`NilOfType` generates a typed nil value without initializing an actual instance of the type. This utility is beneficial for declaring variables of a certain type for interfaces or when needing a placeholder value without an initial concrete value.

```go
var typedNil interface{} = utils.NilOfType[interface{}](func(interface{}) {})
fmt.Println(typedNil == nil) // Outputs: true
```

### AllowUnused
`AllowUnused` explicitly marks variables as used to circumvent compiler warnings about unused variables. This function supports development stages, allowing for the definition of objects for future use or objects whose instantiation has side effects beneficial for application behavior.

```go
var unusedVar = "This is not used yet"
utils.AllowUnused(unusedVar) // Prevents compiler warnings about 'unusedVar' being unused.
```

### Ptr
`Ptr` generates a pointer to a given value, including interfaces and types where the '&' operator cannot be directly applied. This function facilitates the creation of pointers to any value or object, expanding the usability of pointers in Go.

```go
var value interface{} = "example"
var pointerToValue = utils.Ptr(value)
fmt.Println(*pointerToValue) // Outputs: example
```

### `AsInterface`

- **Purpose**: Returns the input value as an `interface{}` type, effectively serving as a no-operation for interface conversion.
- **Signature**: `func AsInterface(i interface{}) interface{}`

### `AssertType`

- **Purpose**: Performs a type assertion on an interface value to convert it to a specific type `T`, panicking if the assertion fails.
- **Signature**: `func AssertType[T any](i interface{}, t func(T)) T`

### `AssertTypeUnsafe`

- **Purpose**: Similar to `AssertType`, but it ignores the success of the type assertion, primarily used when the developer is certain of the underlying type and wishes to avoid checking the assertion result.
- **Signature**: `func AssertTypeUnsafe[T any](i interface{}, t func(T)) T`

### `SwitchType`

- **Purpose**: Uses a type switch to attempt conversion of an interface value to type `T`, falling back to an unsafe conversion if the direct type assertion fails.
- **Signature**: `func SwitchType[T any](i interface{}, t func(T)) T`
