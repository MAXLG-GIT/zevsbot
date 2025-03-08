package domain_interfaces

import "image"

type ImageProcessing interface {
	Median(image *image.Image, medianVal int) (*image.Image, error)
}

type Barcode interface {
	ReadImage(img *image.Image) (string, error)
}
