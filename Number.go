package utils

func AddInt[I int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | uintptr, T int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 |uintptr](i I,t T)I{
  return i + I(t)
}

