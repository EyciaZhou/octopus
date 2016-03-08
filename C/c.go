package COctopus

import (
	"gopkg.in/macaron.v1"
	"github.com/EyciaZhou/msghub-http/C"
	"encoding/json"
	"github.com/EyciaZhou/octopus/M"
)

func RouterGroup(m *macaron.Macaron) {
	m.Get("/getVersion", getVersion)
}

func getVersion(ctx *macaron.Context) {
	version := ctx.Query("version")
	info := ctx.Query("info")

	kvs := MOctopus.KVS{}

	err := json.Unmarshal(([]byte)(info), &kvs)

	if err != nil {
		ctx.JSON(200, C.Error(err))
		return
	}

	v := MOctopus.Operate.GetWithVersion(version, kvs)

	ctx.JSON(200, C.Pack(v))
}
