package img_processing

import (
	"github.com/disintegration/gift"
	"image"
	"zevsbot/internal/domain/domain_interfaces"
)

type giftStruct struct {
	gift *gift.GIFT
}

func Init() domain_interfaces.ImageProcessing {
	return &giftStruct{gift: gift.New()}
}

func (g *giftStruct) Median(inImage *image.Image, medianVal int) (*image.Image, error) {
	var outImg image.Image
	g.gift.Add(gift.Median(medianVal, false))

	dst := image.NewRGBA(g.gift.Bounds((*inImage).Bounds()))

	g.gift.Draw(dst, outImg)

	g.gift.Empty()
	return &outImg, nil
}
