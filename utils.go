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



