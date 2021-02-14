package util

func MapRange(input, inputStart, inputEnd, outputStart, outputEnd float64) float64 {
	return ((input-inputStart)/(inputEnd-inputStart))*(outputEnd-outputStart) + outputStart
}
