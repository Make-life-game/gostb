package main

import (
	"encoding/json"
	"flag"
	"github/Make-life-game/gostn"
	"os"
)

var path = flag.String("p", "", "Path to save, parent directory should be created before.")
var w = flag.Uint("w", 320, "Path to save, parent directory should be created before.")
var h = flag.Uint("h", 0, "Path to save, parent directory should be created before.")
var v = flag.Bool("v", true, "Path to save, parent directory should be created before.")

func main() {
	flag.Parse()
	encoder := json.NewEncoder(os.Stdout)

	if *path == "" {
		r := map[string]interface{}{
			"path": path,
			"code": 3,
			"msg":  "-p args required",
		}
		encoder.Encode(r)
		return
	}

	img := gostn.GetFullScreenShot(*w, *h, *v)
	code, msg := gostn.UpdateScreenshotInfo(*path, img)
	r := map[string]interface{}{
		"path": path,
		"code": code,
		"msg":  msg,
	}

	encoder.Encode(r)

}
