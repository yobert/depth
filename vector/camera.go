package vector

type Camera struct {
	// actual screen resolution (for unprojecting!)
	Width  float64
	Height float64

	// Seed data for view frustom
	YFov Degree
	Near float64
	Far  float64

	// View frustum
	View Frustum

	Projection       M44
	ModelView        M44
	ModelViewInverse M44

	// these will get used in Update() to generate
	// the modelview matrix
	Position V3
	RotAxis  Euler

	// This is just for caching
	ModelViewProjection M44
}

// Before calling this, set cam Width, Height, YFov, Near, and Far
func (cam *Camera) SetupViewProjection() {
	x_ratio := cam.Width / cam.Height
	cam.View = PerspectiveFrustum(cam.YFov, x_ratio, cam.Near, cam.Far)
	cam.Projection = cam.View.M44()
}

func (cam *Camera) SetupModelView() {

	cam.ModelViewInverse = IdentityM44()

	cam.ModelViewInverse = cam.ModelViewInverse.Mult(TranslateM44(cam.Position))
	cam.ModelViewInverse = cam.ModelViewInverse.Mult(RotateAxisM33(V3{0.0, 0.0, 1.0}, cam.RotAxis.Z).M44())
	cam.ModelViewInverse = cam.ModelViewInverse.Mult(RotateAxisM33(V3{0.0, 1.0, 0.0}, cam.RotAxis.Y).M44())
	cam.ModelViewInverse = cam.ModelViewInverse.Mult(RotateAxisM33(V3{1.0, 0.0, 0.0}, cam.RotAxis.X).M44())

	cam.ModelView = cam.ModelViewInverse.Inverse()
}

func (cam *Camera) Unproject(ix, iy float64) [2]V3 {
	near := V4{
		2.0*ix/cam.Width - 1.0,
		2.0*(cam.Height-iy)/cam.Height - 1.0,
		-1, 1}

	far := V4{near.X, near.Y, 1, 1}

	modelview := cam.ModelView
	projection := cam.Projection

	m := modelview.MultX(projection).Inverse()
	//m := projection.Mult(modelview).Inverse()

	return [2]V3{
		m.MultV4(near).HomogeneousToCartesian(),
		m.MultV4(far).HomogeneousToCartesian()}
}

func Ortho(left, right, bottom, top, near, far float64) M44 {
	tx := (right + left) / (right - left)
	ty := (top + bottom) / (top - bottom)
	tz := (far + near) / (far - near)

	return M44{
		2 / (right - left), 0, 0, -tx,
		0, 2 / (top - bottom), 0, -ty,
		0, 0, -2 / (far - near), -tz,
		0, 0, 0, 1,
		//		2 / (right - left), 0, 0, 0,
		//		0, 2 / (top - bottom), 0, 0,
		//		0, 0, -2 / (far - near), 0,
		//		-tx, -ty, -tz, 1,
	}
}
