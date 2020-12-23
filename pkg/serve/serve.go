package serve

import (
	"fmt"
	"net/http"

	build "github.com/DhruvBisla/ogolb/pkg/build"
)

func Serve() {
	build.Build()
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	fmt.Println("Listening on :8000...")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}
