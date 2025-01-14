package bitsUtil

// And returns x&y.
func And(x, y uint) uint {
	return x & y
}

// Or returns x|y.
func Or(x, y uint) uint {
	return x | y
}

// Xor returns x^y.
func Xor(x, y uint) uint {
	return x ^ y
}

// Not returns ^x.
func Not(x uint) uint {
	return ^x
}

// BitClear (AND NOT) returns x&^y.
func BitClear(x, y uint) uint {
	return x &^ y
}

// LeftShift returns x<<n.
func LeftShift(x uint, n int) uint {
	return x << n
}

// RightShift returns x>>n.
func RightShift(x uint, n int) uint {
	return x >> n
}

// MultiplyBy2n returns x*2^n.
func MultiplyBy2n(x uint, n int) uint {
	return x << n
}

// DivideBy2n returns x/2^n.
func DivideBy2n(x uint, n int) uint {
	return x >> n
}

// IsEven returns true if x is even.
func IsEven(x uint) bool {
	return (x & 1) == 0
}

// IsPowerOfTwo returns true if x is a power of two.
func IsPowerOfTwo(x uint) bool {
	return x != 0 && (x&(x-1)) == 0
}

// IsDivisibleBy8 returns true if x is divisible by 8.
func IsDivisibleBy8(x uint) bool {
	return (x & 7) == 0
}

// IsOppositeSigns returns true if x and y have opposite signs.
func IsOppositeSigns(x, y int) bool {
	return (x ^ y) < 0
}

// SetNth returns x with the nth bit set to 1, and n start from 0 and rightmost.
func SetNth(x uint, n int) uint {
	return x | (1 << n)
}

// UnsetNth returns x with the nth bit set to 0, and n start from 0 and rightmost.
func UnsetNth(x uint, n int) uint {
	return x &^ (1 << n)
}

// IsSetNth returns true if the nth bit of x is set.
func IsSetNth(x uint, n int) bool {
	return (x & (1 << n)) != 0
}

// ToggleNth returns x with the nth bit toggled, and n start from 0 and rightmost.
func ToggleNth(x uint, n int) uint {
	return x ^ (1 << n)
}

// ToggleExceptNth returns x with all bits except the nth bit toggled, and n start from 0 and rightmost.
func ToggleExceptNth(x uint, n int) uint {
	return x ^ (^(1 << n))
}

// ToggleRightN returns x with the rightmost n bits toggled.
func ToggleRightN(x uint, n int) uint {
	return x ^ ((1 << n) - 1)
}
