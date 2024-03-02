package utils
import(
  "unsafe"
)

func ForceType[T any](obj any,t func(T)) T {
  return *(*T)(unsafe.Pointer(&obj))
}

func NilOfType[T any](t func(T)) T {
  n := error(nil)
  return *(*T)(unsafe.Pointer(&n))
}

func Pass(a any) {}

func Ptr[T any](value T) *T {
  return &value
}

