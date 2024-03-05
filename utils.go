package utils
import(
  "unsafe"
)

func ForceType[F any,T any](f F,t func(T)) T {
  return *(*T)(unsafe.Pointer(&f))
}

func ForceRawType[T any](f any,t func(T)) T {
  return *(*T)(unsafe.Pointer(&f))
}

func NilOfType[T any](t func(T)) T {
  n := error(nil)
  return *(*T)(unsafe.Pointer(&n))
}

func AllowUnused(a any) {}

func Ptr[T any](value T) *T {
  return &value
}

func AsInterface(i interface{})interface{}{
  return i
}

func AssertType[T any](i interface{},t func(T))T{
  return i.(T)
}

func AssertTypeUnsafe[T any](i interface{},t func(T))T{
  r, ok := i.(T)
  AllowUnused(ok)
  return r
}

func SwitchType[T any](i interface{})T {
  switch v := i.(type) {
  case T:
    return v
  default:
    return ForceType(i,func(t T){})
  }
}
