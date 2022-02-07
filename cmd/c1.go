/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	box_lib "github.com/gh-zhangpeng/box-lib"
	"github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"html/template"
	"os"
	"strings"
)

type Api struct {
	Path string `json:"path"`
	Comment string `json:"comment"`
	Module string `json:"module"`
	ControllerFuncName string `json:"controllerFuncName"`
	ServiceFuncName string `json:"serviceFuncName"`
}

type File struct {
	PackageName string
	Imports     []string
}

// c1Cmd represents the c1 command
var c1Cmd = &cobra.Command{
	Use:   "c1",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("c1 called")
		controllerOutput := "./cli-output/controllers"
		serviceOutput := "./cli-output/services"
		if temp := viper.GetString("create.output.controller"); len(temp) > 0 {
			controllerOutput = temp
		}
		if temp := viper.GetString("create.output.service"); len(temp) > 0 {
			serviceOutput = temp
		}
		fmt.Printf("crontroller 生成路径为：%+v\n", controllerOutput)
		fmt.Printf("service 生成路径为：%+v\n", serviceOutput)

		apis := getApis()
		if len(apis) == 0 {
			fmt.Println("apis is empty")
			return
		}
		for _, v := range apis {
			valid, err := checkApiValid(v)
			if err != nil {
				return 
			} else if !valid {
				return
			}
			module := v.Module
			controllerFuncName := v.ControllerFuncName
			serviceFuncName := v.ServiceFuncName
			if len(module) == 0 {
				pathItems := strings.Split(v.Path, "/")
				module = pathItems[len(pathItems) - 2]
			}
			if len(controllerFuncName) == 0 {
				pathItems := strings.Split(v.Path, "/")
				controllerFuncName = box_lib.FirstUpper(pathItems[len(pathItems) - 1])
			}
			if len(serviceFuncName) == 0 {
				serviceFuncName = controllerFuncName
			}
			fmt.Printf("module: %s, controllerFuncName: %s, serviceFuncName: %s\n", module, controllerFuncName, serviceFuncName)
			//模块路径
			moduleControllerPath := controllerOutput + "/" + module
			fmt.Printf("%+v\n", moduleControllerPath)
			err = createDir(moduleControllerPath)
			if err != nil {
				fmt.Printf("createDir module controller dir path fail, err: %s\n", err.Error())
				return
			}
			//模块路径
			moduleServicePath := serviceOutput + "/" + module
			fmt.Printf("%+v\n", moduleServicePath)
			err = createDir(moduleServicePath)
			if err != nil {
				fmt.Printf("createDir module service dir path fail, err: %s\n", err.Error())
				return
			}
		}
	},
}

func createFile(filePath string, tplPath string, fileData File, tplData interface{}) error {
	exists, err := box_lib.Exists(filePath)
	if err != nil {
		fmt.Printf("检查文件是否存在失败，err: %s\n", err.Error())
		return err
	}
	flag := os.O_RDWR | os.O_APPEND
	if !exists {
		fmt.Println("文件不存在，创建文件")
		flag = os.O_CREATE | flag
	}
	//打开或创建 controller 文件
	controller, err := os.OpenFile(filePath, flag, 0777)
	defer controller.Close()
	if err != nil {
		fmt.Printf("打开文件失败，err: %s\n", err.Error())
		return err
	}

	if !exists {
		//读取文件模板
		fileTpl, err := template.ParseFiles("./tpls/file.tpl")
		if err != nil {
			fmt.Printf("解析 file 模板失败，err: %s\n", err.Error())
			return err
		}
		err = fileTpl.Execute(controller, fileData)
		if err != nil {
			fmt.Printf("渲染 file 模版失败，err: %s\n", err.Error())
			return err
		}
	}

	targetTpl, err := template.ParseFiles(tplPath)
	if err != nil {
		fmt.Printf("解析目标模板失败，err: %s\n", err.Error())
		return err
	}
	err = targetTpl.Execute(controller, tplData)
	if err != nil {
		fmt.Printf("渲染目标模版失败，err: %s\n", err.Error())
		return err
	}
	return nil
}

func createDir(path string) error {
	exists, err := box_lib.Exists(path)
	if err != nil {
		fmt.Printf("check dir path fail, err: %s\n", err.Error())
		return err
	}
	if exists {
		//fmt.Println("dir path already exists，next step...")
		return nil
	}
	fmt.Println("module dir path does not exist, creating...")
	return os.MkdirAll(path, 0777)
}

func checkApiValid(api Api) (bool, error) {
	if len(api.Path) == 0 {
		fmt.Println("please check path, the path cannot be empty")
		return false, nil
	}
	match, err := box_lib.Match("^(/[A-Za-z0-9]+)+$", api.Path)
	if err != nil {
		fmt.Printf("match path fail, err: %s\n", err.Error())
		return false, nil
	}
	if !match {
		fmt.Println("please check path, the rule of the path is ^(/[A-Za-z0-9])+$")
		return false, nil
	}
	return true, nil
}

func getApis() []Api {
	apisConfig := viper.Get("create.apis")
	apisString, err := jsoniter.MarshalToString(apisConfig)
	if err != nil {
		fmt.Printf("marshal apis to string fail, err: %s\n", err.Error())
		return nil
	}
	var apis []Api
	err = jsoniter.UnmarshalFromString(apisString, &apis)
	if err != nil {
		fmt.Printf("unmarshal apis string fail, err: %s\n", err.Error())
		return nil
	}
	return apis
}

func init() {
	rootCmd.AddCommand(c1Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// c1Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// c1Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
