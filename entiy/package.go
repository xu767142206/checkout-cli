package entiy

import (
	"bytes"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/liushuochen/gotable"
	"github.com/liushuochen/gotable/table"
	"github.com/xu767142206/checkout-cli/tools"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	LINUX = "linux"
	MACOS = "macos"
	WIN   = "cygwin"
)

var OsMaps = map[string]string{
	"darwin":  MACOS,
	"linux":   LINUX,
	"windows": WIN,
}

// Package swoole_cli包
type Package struct {
	Url          string    `json:"url"`
	Filename     string    `json:"filename"`
	Etag         string    `json:"etag"`
	LastModified string    `json:"last_modified"`
	Size         string    `json:"size"`
	Name         string    `json:"-"`
	Version      string    `json:"-"`
	Date         time.Time `json:"-"`
}

type netWorkLogic struct {
	OsType string
}

var Meteorologic = new(netWorkLogic)
var client = resty.New()

const DownloadsPath = "./resources"

func init() {
	Meteorologic.OsType = OsMaps[runtime.GOOS]
	err := tools.CreateDir(DownloadsPath)
	if err != nil {
		log.Fatalln(err)
	}
	client.SetOutputDirectory(DownloadsPath)
}

func GetNetWorkLogic() *netWorkLogic {
	return Meteorologic
}

func (workLogic *netWorkLogic) GetSwooleCliList() []Package {

	resp, err := client.R().
		SetResult(make([]Package, 0)).
		SetHeader("Accept", "application/json").
		Get("https://www.swoole.com/download?out=json")
	if err != nil {
		log.Fatalln(err)
	}
	packages := make([]Package, 0)
	if result, ok := resp.Result().(*[]Package); ok {

		for _, v := range *result {
			if strings.Contains(v.Filename, workLogic.OsType) {
				fileName := path.Base(v.Filename)

				v.Name = fileName[0 : len(fileName)-len(path.Ext(fileName))]
				fiedstrs := strings.Split(fileName, "-")

				for _, str := range fiedstrs {
					if len(str) > 0 && str[0] == 'v' {
						v.Version = str
						break
					}
				}
				v.Date, _ = time.Parse("2006-01-02T15:04:05.000Z", v.LastModified)

				packages = append(packages, v)
			}

		}

	}

	return packages
}

// Totable
func (workLogic *netWorkLogic) Totable(packages []Package) *table.Table {

	table, _ := gotable.Create("No.", "date", "version", "name")

	for i, v := range packages {

		table.AddRow(map[string]string{
			"No.":     strconv.Itoa(i + 1),
			"date":    v.Date.Format("2006-01-02 15:04:05"),
			"version": v.Version,
			"name":    v.Name,
		})
	}
	return table
}

// Serach
func (workLogic *netWorkLogic) Serach(list []Package, filed string) []Package {

	packages := make([]Package, 0)

	for _, v := range list {
		if strings.Contains(v.Name, filed) {
			packages = append(packages, v)
		}
	}

	return packages
}

func (workLogic *netWorkLogic) Download(zipUrl string) io.Reader {

	resp, err := client.R().Get(zipUrl)

	if err != nil {
		log.Fatalln(err)
	}
	return bytes.NewReader(resp.Body())
}

func (workLogic *netWorkLogic) GetVersionPackge(list []Package, filed string) (Package, error) {

	for _, v := range list {
		if v.Name == strings.TrimSpace(filed) {
			return v, nil
		}
	}
	return Package{}, fmt.Errorf("未找到相关的版本")

}

func (workLogic *netWorkLogic) ReadDir() ([]fs.FileInfo, error) {
	return ioutil.ReadDir(DownloadsPath)
}
