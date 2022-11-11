package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/xu767142206/checkout-cli/entiy"
	"github.com/xu767142206/checkout-cli/tools"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {

	fmt.Println()

	logic := entiy.GetNetWorkLogic()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:      "search",
				Usage:     "列出所有适合系统的cli版本列表",
				Aliases:   []string{"s"},
				ArgsUsage: "example: search v4.8.12",
				Action: func(c *cli.Context) error {

					version := c.Args().First()
					//查看列表
					list := logic.GetSwooleCliList()

					if version != "" {
						list = logic.Serach(list, version)
					}
					totable := logic.Totable(list)
					fmt.Println(totable)
					return nil
				},
			},
			{
				Name:    "install",
				Usage:   "安装swoole-cli相应的版本",
				Aliases: []string{"i"},
				Action: func(c *cli.Context) error {

					version := c.Args().First()
					if version == "" {
						log.Fatalln("请输入安装的package！")
					}

					list := logic.GetSwooleCliList()

					packge, err := logic.GetVersionPackge(list, version)
					if err != nil {
						log.Fatalln(err)
					}
					log.Printf("解析:%s...\n", packge.Url)
					//下载
					log.Println("下载ing...")
					err = tools.DownloadFile(path.Join(entiy.DownloadsPath, packge.Filename), packge.Url)
					if err != nil {
						log.Fatalln(err)
					}

					os.RemoveAll(path.Join(entiy.DownloadsPath, packge.Name))
					fmt.Println(path.Join(entiy.DownloadsPath, packge.Name))
					log.Println("解压ing...")
					//解压
					err = tools.Unpack(path.Join(entiy.DownloadsPath, packge.Filename), entiy.DownloadsPath)
					if err != nil {
						log.Fatalln(err)
					}
					os.Remove(path.Join(entiy.DownloadsPath, packge.Filename))
					log.Println("安装完成!")
					return nil
				},
			},
			{
				Name:    "uninstall",
				Usage:   "卸载swoole-cli相应的版本",
				Aliases: []string{"u"},
				Action: func(c *cli.Context) error {

					name := c.Args().First()

					//查看列表
					filePath, err := logic.ReadDir()
					if err != nil {
						log.Fatalln(err)
					}
					for _, v := range filePath {
						if v.Name() == strings.TrimSpace(name) {
							log.Println("找到安装包：", v.Name(), "准备删除...")
							os.RemoveAll(path.Join(entiy.DownloadsPath, v.Name()))
							log.Println(v.Name(), "已经删除!")
							break
						}
					}
					log.Println("卸载完成!")
					return nil

				},
			},

			{
				Name:    "use",
				Aliases: []string{"c"},
				Usage:   "切换swoole-cli版本",
				Action: func(c *cli.Context) error {

					name := c.Args().First()

					//查看列表
					filePath, err := logic.ReadDir()
					if err != nil {
						log.Fatalln(err)
					}

					version := ""
					for _, v := range filePath {
						if v.Name() == strings.TrimSpace(name) {
							version = v.Name()
							break
						}
					}
					if version == "" {
						log.Fatalln("未找到该包")
					}

					absPath, _ := filepath.Abs(entiy.DownloadsPath)
					if logic.OsType == entiy.WIN {
						ioutil.WriteFile(
							"./swoole-cli.bat",
							[]byte(`@`+filepath.Join(absPath, version, "bin", "swoole-cli.exe")+` %*`),
							0777,
						)
					}
					log.Fatalln("切换成功!")
					return nil
				},
			},

			{
				Name:    "list",
				Usage:   "本地安装的swoole-cli版本",
				Aliases: []string{"l"},
				Action: func(c *cli.Context) error {
					//查看列表
					filePath, err := logic.ReadDir()
					if err != nil {
						log.Fatalln(err)
					}
					for i, v := range filePath {
						fmt.Println(i+1, ":", v.Name())
					}
					return nil
				},
			},

			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "获取工具版本号",
				Action: func(c *cli.Context) error {
					fmt.Println("checkout-swoole_cli is simple swoole-cli version manager on Windows.\nCopyright (C) 2022 vance <xu767142206@163.com>")
					fmt.Println("Current version v1.0.0 64 bit.")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
