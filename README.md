# Fixed point vector library

This is a Go library for calculations on 3D objects. All calculations are done
with fixed-point math so it is fast on hardware without floating point unit.

## Performance

This library can multiply two quaternions and rotate 12 vectors by this
quaternion in about 1 millisecond on a Cortex-M0 running at 16MHz when using the
[TinyGo](https://github.com/aykevl/tinygo) compiler.

## License

This library is licensed under a 3-clause BSD license.

Some code has been copied from the [mathgl](https://github.com/go-gl/mathgl)
library and has been modified to use fixed point arithmetic. Both libraries use
the same license.
