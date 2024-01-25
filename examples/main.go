package main

import (
	"fmt"
	"github.com/calvinmclean/babyapi"
	"github.com/cody-cody-wy/babyapiFileUploadParser"
	"github.com/go-chi/render"
	"net/http"
	"os"
)

type TestStruct struct {
	Test string
}

type Types struct {
	// Numbers
	Int    int
	Int8   int8
	Int16  int16
	Int32  int32
	Int64  int64
	Uint   uint
	Uint8  uint8
	Uint16 uint16
	Uint32 uint32
	Uint64 uint64
	Rune   rune
	Byte   byte
	// Uintptr uintptr // why should I deserialize pointers??
	// Floats
	Float32 float32
	Float64 float64
	// Complex not supported by JSON
	// Complex64  complex64
	// Complex128 complex128
	// Others
	Boolean       bool
	String        string
	Struct        TestStruct
	Array         [2]int8
	Array2D       [2][2]int8
	StructArray   [3]TestStruct
	Slice         []string
	Slice2D       [][]float32
	StructSlice   []TestStruct
	SliceArray    [][3]float64
	ArraySlice    [3][]float64
	Image         babyapiFileUploadParser.FileField
	Images        []babyapiFileUploadParser.FileField
	Images2D      [][]babyapiFileUploadParser.FileField
	ImagesArray   [3]babyapiFileUploadParser.FileField
	ImagesArray2D [2][2]babyapiFileUploadParser.FileField
	privateInt    int
}

type Project struct {
	babyapi.DefaultResource

	Name        string `form:"projectName" json:"projectName"`
	Description string
	Test        string
	Image       babyapiFileUploadParser.FileField
	Image2      babyapiFileUploadParser.FileField `form:"OtherImage" json:"OtherImage"`
	Types       Types
}

func main() {
	render.Decode = babyapiFileUploadParser.Decoder

	ProjectApi := babyapi.NewAPI[*Project]("Projects", "/Projects", func() *Project { return &Project{} })

	ProjectApi.SetAfterDelete(func(r *http.Request) *babyapi.ErrResponse {
		id := babyapi.GetIDParam(r, ProjectApi.Name())
		fmt.Println(id)
		err := os.RemoveAll(fmt.Sprintf("./Static/%s/%s", ProjectApi.Name(), id))
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return nil
	})

	ProjectApi.SetOnCreateOrUpdate(func(r *http.Request, project *Project) *babyapi.ErrResponse {
		fmt.Println("ON CREATE")
		babyapiFileUploadParser.WriteAllFileFields("./Static/"+ProjectApi.Name(), project.GetID(), project)
		return nil
	})

	ProjectApi.RunCLI()
}