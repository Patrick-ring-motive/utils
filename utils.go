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
import(
  "unsafe"
  "reflect"
)


func TypeOf[T any](t ...T)func(T){
  return func(T){}
}

func TypeRef[T any](t ...T)func(T){
  return func(T){}
}

func ConvertType[I any,T any](i I,t ...func(T))T{
  switch v := any(i).(type) {
  case T:
    return T(v)
  default:
    return T(AssertType[I,T](i))
  }
}

func AssertType[I any,T any](i I,t ...func(T))T{
  r, ok := any(i).(T)
  if(!ok){
    r = SwitchType[I,T](i)
  }
  return r
}

func AssertTypeUnsafe[T any](i interface{},t ...func(T))T{
  r, ok := i.(T)
  AllowUnused(ok)
  return r
}

func SwitchType[S any,T any](s S,t ...func(T))T {
  i := AsInterface(s)
  switch v := i.(type) {
  case T:
    return v
  default:
    return ForceType[S,T](s)
  }
}

func ForceType[F any,T any](f F,t ...func(T)) T {
  var z T
  a := &[1]T{z}
  forceType(a, f)
  return a[0]
}
func forceType[F any,T any](a *[1]T,f F) {
  defer func() {
    if r := recover(); r != nil {
      a[0] = ForceRawType(f,func(T){})
    }
  }()
  a[0]=*(*T)(unsafe.Pointer(&f))
}

func ForceRawType[T any](f any,t ...func(T)) T {
  var z T
  a := &[1]T{z}
  forceType(a, f)
  return a[0]
}
func forceRawType[T any](a *[1]T,f any) {
  defer func() {
    if r := recover(); r != nil {
      a[0] = ZeroOfType(func(T){})
    }
  }()
  a[0]=*(*T)(unsafe.Pointer(&f))
}

func ZeroOfType[T any](t ...func(T)) T{
  var x T
  return x
}

func NilOfType[T any](t ...func(T)) T {
  n := any(nil)
  return *(*T)(unsafe.Pointer(&n))
}

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

func AllowUnused(a ...any) {}

func Ptr[T any](value T) *T {
  return &value
}

func AsInterface(i interface{})interface{}{
  return i
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