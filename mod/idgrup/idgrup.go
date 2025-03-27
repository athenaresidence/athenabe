package idgrup

import (
	"github.com/gocroot/lite/model"
)

func IDGroup(Pesan model.WAMessage) (reply string) {

	return "Hai.. hai.. ini dia id group kaka :\n" + Pesan.Group_id + "\nsimpan dan catat baik baik ya.. \nmakasih"
}
