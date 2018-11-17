// Package fixpoint implements fixed-point arithmetic and vector operations. It
// has been inspired on and is partially copied from the github.com/go-gl/mathgl
// package.
package fixpoint

// Useful link:
// https://spin.atomicobject.com/2012/03/15/simple-fixed-point-math/

// Q24 is a Q7.24 fixed point integer type that has 24 bits of precision to the
// right of the fixed point. It is designed to be used as a more efficient
// replacement for unit vectors with some extra room to avoid overflow.
type Q24 struct {
	N int32
}

// Q24FromFloat converts a float32 to the same number in fixed point format.
// Inverse of .Float().
func Q24FromFloat(x float32) Q24 {
	return Q24{int32(x * (1 << 24))}
}

// Q24FromInt32 returns a fixed point integer with all decimals set to zero.
func Q24FromInt32(x int32) Q24 {
	return Q24{x << 24}
}

// Float returns the floating point version of this fixed point number. Inverse
// of Q24FromFloat.
func (q Q24) Float() float32 {
	return float32(q.N) / (1 << 24)
}

// Int32Scaled returns the underlying fixed point number multiplied by scale.
func (q Q24) Int32Scaled(scale int32) int32 {
	return q.N / (1 << 24 / scale)
}

// Add returns the argument plus this number.
func (q1 Q24) Add(q2 Q24) Q24 {
	return Q24{q1.N + q2.N}
}

// Sub returns the argument minus this number.
func (q1 Q24) Sub(q2 Q24) Q24 {
	return Q24{q1.N - q2.N}
}

// Neg returns the inverse of this number.
func (q1 Q24) Neg() Q24 {
	return Q24{-q1.N}
}

// Mul returns this number multiplied by the argument.
func (q1 Q24) Mul(q2 Q24) Q24 {
	return Q24{int32((int64(q1.N) * int64(q2.N)) >> 24)}
}

// Div returns this number divided by the argument.
func (q1 Q24) Div(q2 Q24) Q24 {
	return Q24{int32((int64(q1.N) << 24) / int64(q2.N))}
}

// Vec3Q24 is a 3-dimensional vector with Q24 fixed point elements.
type Vec3Q24 struct {
	X Q24
	Y Q24
	Z Q24
}

// Vec3Q24FromFloat returns the fixed-point vector of the given 3 floats.
func Vec3Q24FromFloat(x, y, z float32) Vec3Q24 {
	return Vec3Q24{Q24FromFloat(x), Q24FromFloat(y), Q24FromFloat(z)}
}

// Add returns this vector added to the argument.
func (v1 Vec3Q24) Add(v2 Vec3Q24) Vec3Q24 {
	// Copied from go-gl/mathgl and modified.
	return Vec3Q24{v1.X.Add(v2.X), v1.Y.Add(v2.Y), v1.Z.Add(v2.Z)}
}

// Mul returns this vector multiplied by the argument.
func (v1 Vec3Q24) Mul(c Q24) Vec3Q24 {
	// Copied from go-gl/mathgl and modified.
	return Vec3Q24{v1.X.Mul(c), v1.Y.Mul(c), v1.Z.Mul(c)}
}

// Dot returns the dot product between this vector and the argument.
func (v1 Vec3Q24) Dot(v2 Vec3Q24) Q24 {
	// Copied from go-gl/mathgl and modified.
	return v1.X.Mul(v2.X).Add(v1.Y.Mul(v2.Y)).Add(v1.Z.Mul(v2.Z))
}

// Cross returns the cross product between this vector and the argument.
func (v1 Vec3Q24) Cross(v2 Vec3Q24) Vec3Q24 {
	// Copied from go-gl/mathgl and modified.
	return Vec3Q24{v1.Y.Mul(v2.Z).Sub(v1.Z.Mul(v2.Y)), v1.Z.Mul(v2.X).Sub(v1.X.Mul(v2.Z)), v1.X.Mul(v2.Y).Sub(v1.Y.Mul(v2.X))}
}

// QuatQ24 is a quaternion with Q24 fixed point elements.
type QuatQ24 struct {
	W Q24
	V Vec3Q24
}

// QuatIdent returns the identity quaternion.
func QuatIdent() QuatQ24 {
	return QuatQ24{Q24FromInt32(1), Vec3Q24{}}
}

// X returns the X part of this quaternion.
func (q QuatQ24) X() Q24 {
	return q.V.X
}

// Y returns the Y part of this quaternion.
func (q QuatQ24) Y() Q24 {
	return q.V.Y
}

// Z returns the Z part of this quaternion.
func (q QuatQ24) Z() Q24 {
	return q.V.Z
}

// Mul returns this quaternion multiplied by the argument.
func (q1 QuatQ24) Mul(q2 QuatQ24) QuatQ24 {
	// Copied from go-gl/mathgl and modified.
	return QuatQ24{q1.W.Mul(q2.W).Sub(q1.V.Dot(q2.V)), q1.V.Cross(q2.V).Add(q2.V.Mul(q1.W)).Add(q1.V.Mul(q2.W))}
}

// Rotate returns the vector from the argument rotated by the rotation this
// quaternion represents.
func (q1 QuatQ24) Rotate(v Vec3Q24) Vec3Q24 {
	// Copied from go-gl/mathgl and modified.
	cross := q1.V.Cross(v)
	// v + 2q_w * (q_v x v) + 2q_v x (q_v x v)
	return v.Add(cross.Mul(Q24FromInt32(2).Mul(q1.W))).Add(q1.V.Mul(Q24FromInt32(2)).Cross(cross))
}
