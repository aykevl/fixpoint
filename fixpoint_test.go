package fixpoint

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stretchr/testify/assert"
)

func TestQ24(t *testing.T) {
	two := Q24FromFloat(2)
	for _, f := range []float32{0.25, 1, 10, 0.125} {
		q := Q24FromFloat(f)
		assert.Equal(t, f, q.Float(), "float32 roundtrip failed")
		assert.Equal(t, Q24FromFloat(f*2), q.Mul(two), "multiply by 2")
		square := q.Mul(q)
		assert.Equal(t, Q24FromFloat(f*f), square, "square")
		assert.Equal(t, Q24FromFloat(f), square.Div(q), "div")
	}
}

func TestRotation(t *testing.T) {
	// Create vector to rotate.
	vec1 := mgl32.Vec3{0, 0.8320503, 0.5547002}
	vec2 := Vec3Q24FromFloat(0, 0.8320503, 0.5547002)

	// Create rotation quaternion.
	inc1 := mgl32.Quat{0, mgl32.Vec3{0.07, 0, 0}}
	inc1.W = 1.0 - 0.5*(inc1.X()*inc1.X()+inc1.Y()*inc1.Y()+inc1.Z()*inc1.Z())

	inc2 := QuatQ24{Q24FromInt32(0), Vec3Q24FromFloat(0.07, 0, 0)}
	inc2.W = Q24FromInt32(1).Sub(Q24FromFloat(0.5).Mul(inc2.X().Mul(inc2.X()).Add(inc2.Y().Mul(inc2.Y())).Add(inc2.Z().Mul(inc2.Z()))))

	rotation1 := mgl32.QuatIdent()
	rotation2 := QuatIdent()
	for i := 0; i < 10; i++ {
		// Incrementally add parts of the rotation.
		rotation1 = rotation1.Mul(inc1)
		rotation2 = rotation2.Mul(inc2)

		// Get the vector that belongs to this rotation.
		vec1 := rotation1.Rotate(vec1)
		vec2 := rotation2.Rotate(vec2)
		diffX := vec2.X.Float() - vec1.X()
		diffY := vec2.Y.Float() - vec1.Y()
		diffZ := vec2.Z.Float() - vec1.Z()
		if diffX > 0.00001 || diffY > 0.00001 || diffZ > 0.00001 {
			t.Errorf("difference too big:\nX1: %.8f\nX2: %.8f\nY1: %.8f\nY2: %.8f\nZ1: %.8f\nZ2: %.8f",
				vec1.X(), vec2.X.Float(),
				vec1.Y(), vec2.Y.Float(),
				vec1.Z(), vec2.Z.Float())
		}
	}
}
