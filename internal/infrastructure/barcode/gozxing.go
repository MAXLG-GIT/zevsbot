package barcode

import (
	"github.com/disintegration/gift"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
	"github.com/rs/zerolog/log"
	"image"
	"image/png"
	"os"
	"zevsbot/internal/domain/domain_interfaces"
)

type barCodeStr struct {
	imgRedactor domain_interfaces.ImageProcessing
}

func Init(imgRedactor domain_interfaces.ImageProcessing) domain_interfaces.Barcode {
	return &barCodeStr{imgRedactor: imgRedactor}
}

func (bc barCodeStr) ReadImage(img *image.Image) (string, error) {

	// ----------- START barcode decoder

	// 1. Create a new filter list and add some filters.
	g := gift.New(
		gift.Median(3, false),
	)

	// 2. Create a new image of the corresponding size.
	dst := image.NewRGBA(g.Bounds((*img).Bounds()))

	// 3. Use the Draw func to apply the filters to src and store the result in dst.
	g.Draw(dst, *img)

	fileDst, err := os.CreateTemp(os.Getenv("ZEVS_TMP_DIR"), "decode-*.png")
	if err != nil {
		return "", err
	}

	defer func(fileDst *os.File) {
		err = fileDst.Close()
		if err != nil {
			log.Error().Msg(err.Error())
		}
		err = os.Remove(fileDst.Name())
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}(fileDst)

	err = png.Encode(fileDst, dst)
	if err != nil {
		return "", err
	}
	//---------------------------------------

	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(dst)
	// decode image
	hints := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_TRY_HARDER: true,
	}

	barCodeReader := oned.NewCode128Reader()
	result, err := barCodeReader.Decode(bmp, hints)
	if err != nil {
		if err.Error() == "NotFoundException" {
			return "", nil
		}
		return "", err
	}
	// ----------- END barcode decoder
	return result.String(), nil
}
