package noise

import (
	"math"
	"math/rand"
)

type Perlin2D struct {
	RandomNoise [][][]float64
}

func (P *Perlin2D) Init(seed int64, SizeX, SizeY, octavs int) {
	rand.Seed(seed)
	n := 1
	P.RandomNoise = make([][][]float64, octavs)
	for i := 0; i < octavs; i++ {
		P.RandomNoise[i] = make([][]float64, SizeX*n)
		for j := 0; j < SizeX*n; j++ {
			P.RandomNoise[i][j] = make([]float64, SizeY*n)
			for k := 0; k < SizeY*n; k++ {
				P.RandomNoise[i][j][k] = rand.Float64() * 2
			}
		}
		n *= 2
	}
}

func (P *Perlin2D) GetCos(x float64) float64 {
	k := x*0.9 + x*x*(3-2*x)*0.1
	return 1 - 2*k*k*(3-2*k)
}

func (P *Perlin2D) GetPseudoRandomGradientVector(i, j, octav int) [2]float64 {
	var cos, sin float64
	x := P.RandomNoise[octav][i][j]
	if x > 1 {
		//cos := x * x * (3 - 2*x)
		cos = P.GetCos(x - 1)
		sin = -math.Sqrt(1 - cos*cos)
		//cos := math.Cos(x*math.Pi)
	} else {
		//cos := x * x * (3 - 2*x)
		cos = P.GetCos(x)
		sin = math.Sqrt(1 - cos*cos)
		//cos := math.Cos(x*math.Pi)
	}
	return [2]float64{cos, -sin}
}

func (P *Perlin2D) QunticCurve(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func (P *Perlin2D) Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func (P *Perlin2D) Dot(a, b [2]float64) float64 {
	return a[0]*b[0] + a[1]*b[1]
}

func (P *Perlin2D) GetNoise(fx, fy float64, octav int) float64 {
	left := int(math.Floor(fx))
	top := int(math.Floor(fy))
	pointInQuadX := fx - float64(left)
	pointInQuadY := fy - float64(top)

	topLeftGradient := P.GetPseudoRandomGradientVector(left, top, octav)
	topRightGradient := P.GetPseudoRandomGradientVector(left+1, top, octav)
	bottomLeftGradient := P.GetPseudoRandomGradientVector(left, top+1, octav)
	bottomRightGradient := P.GetPseudoRandomGradientVector(left+1, top+1, octav)

	distanceToTopLeft := [2]float64{pointInQuadX, pointInQuadY}
	distanceToTopRight := [2]float64{pointInQuadX - 1, pointInQuadY}
	distanceToBottomLeft := [2]float64{pointInQuadX, pointInQuadY - 1}
	distanceToBottomRight := [2]float64{pointInQuadX - 1, pointInQuadY - 1}

	tx1 := P.Dot(distanceToTopLeft, topLeftGradient)
	tx2 := P.Dot(distanceToTopRight, topRightGradient)
	bx1 := P.Dot(distanceToBottomLeft, bottomLeftGradient)
	bx2 := P.Dot(distanceToBottomRight, bottomRightGradient)

	pointInQuadX = P.QunticCurve(pointInQuadX)
	pointInQuadY = P.QunticCurve(pointInQuadY)

	tx := P.Lerp(tx1, tx2, pointInQuadX)
	bx := P.Lerp(bx1, bx2, pointInQuadX)
	tb := P.Lerp(tx, bx, pointInQuadY)

	//return tb
	return tb
}

func (P *Perlin2D) Noise(fx, fy float64, octaves int, persistence float64) float64 {
	amplitude := 1.
	max := -1.
	result := 0.
	octav := 0
	for octav < octaves {
		max += amplitude
		result += P.GetNoise(fx, fy, octav) * amplitude
		amplitude *= persistence
		fx *= 2
		fy *= 2
		octav++
	}
	return result / max
}
