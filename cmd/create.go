package cmd

import (
	"fmt"
	box_lib "github.com/gh-zhangpeng/box-lib"
	"github.com/spf13/cobra"
	"os"
	"text/template"
)

var module string
var controllerMethod string
var serviceMethod string

type File struct {
	PackageName string
	Imports     []string
}

type Controller struct {
	FuncName string
}

type Service struct {
	FuncName string
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create controller/service",
	Long:  "create controller/service command",
	Run: func(cmd *cobra.Command, args []string) {
		if len(module) == 0 {
			fmt.Println("请输入模块名字")
		}
		output := "./cli-output"
		//模块路径
		controllerDirPath := output + "/controller/" + module
		//判断模块是否存在
		exists, err := box_lib.Exists(controllerDirPath)
		if err != nil {
			fmt.Printf("检查 controller 是否存在失败，err: %s\n", err.Error())
			return
		}
		if !exists {
			fmt.Println("controller 文件夹不存在，创建 controller 文件夹")
			err := os.MkdirAll(controllerDirPath, 0777)
			if err != nil {
				fmt.Printf("创建 controller 文件夹失败，err: %s\n", err.Error())
				return
			}
		}

		if len(controllerMethod) == 0 {
			fmt.Printf("请输入 controller 方法名")
		}

		controllerPath := controllerDirPath + "/" + module + ".go"
		err = createFile(controllerPath, "./tpls/controller.tpl", File{
			PackageName: "controller",
			Imports:     []string{"github.com/gin-gonic/gin"},
		}, Controller{
			FuncName: controllerMethod,
		})
		if err != nil {
			fmt.Printf("创建 controller 失败，err: %s\n", err.Error())
			return
		}

		if len(serviceMethod) == 0 {
			serviceMethod = controllerMethod
		}

		serviceDirPath := output + "/service/" + module
		servicePath := serviceDirPath + "/" + module + ".go"
		//判断 controller 是否存在
		exists, err = box_lib.Exists(serviceDirPath)
		if err != nil {
			fmt.Printf("err: %s\n", err.Error())
			return
		}
		if !exists {
			fmt.Println("service 文件夹不存在，创建 service 文件夹")
			err := os.MkdirAll(serviceDirPath, 0777)
			if err != nil {
				fmt.Printf("创建 service 文件夹失败，err: %s\n", err.Error())
				return
			}
		}

		err = createFile(servicePath, "./tpls/service.tpl", File{
			PackageName: "service",
			Imports:     []string{"github.com/gin-gonic/gin"},
		}, Service{
			FuncName: serviceMethod,
		})
		if err != nil {
			fmt.Printf("创建 service 失败，err: %s\n", err.Error())
			return
		}

		fmt.Printf("创建完成，请前往 %s 查看创建的文件。", output)
		//writer := bytes.NewBufferString("")
		//err = targetTpl.Execute(writer, data)
		//fmt.Println(writer)
		//if err != nil {
		//	fmt.Printf("生成 controller 失败，err: %s\n", err.Error())
		//	return
		//}
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

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&module, "module", "m", "module", "create module")
	createCmd.Flags().StringVarP(&controllerMethod, "controller", "c", "", "create controller")
	createCmd.Flags().StringVarP(&serviceMethod, "service", "s", "", "create service")
}
