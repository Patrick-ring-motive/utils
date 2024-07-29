package utils
import(
  "unsafe"
  "reflect"
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

func AllowUnused(a ...any) {}

func Ptr[T any](value T) *T {
  return &value
}

func AsInterface(i interface{})interface{}{
  return i
}

func AssertType[T any](i any,t func(T))T{
  return i.(T)
}

func AssertTypeUnsafe[T any](i interface{},t func(T))T{
  r, ok := i.(T)
  AllowUnused(ok)
  return r
}

func SwitchType[S any,T any](s S,t func(T))T {
  i := AsInterface(s)
  switch v := i.(type) {
  case T:
    return v
  default:
    return ForceType(i,func(t T){})
  }
}

func InvokeAnyFunc(fn interface{}, args interface{}) any {
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