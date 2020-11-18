package tool_mongo

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// 写入文件
func WriteFile(name string, tpl string, file string, data interface{}, funcs ...template.FuncMap) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	t := template.New(name)

	if len(funcs) > 0 {
		t.Funcs(funcs[0])
	}

	t.Delims("{{{", "}}}")
	t, err = t.Parse(tpl)
	if err != nil {
		return err
	}

	if err := t.Execute(f, data); err != nil {
		return err
	}

	return GoFmt(file)
}

// 使用go fmt 格式化文件
func GoFmt(file string) error {
	_, err := ExecCmd("", "gofmt", "-w", file)
	return err
}

// 下划线转驼峰
func ToTitle(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	return strings.Replace(strings.Title(str), " ", "", -1)
}

// 下划线转驼峰(首字母小写)
func ToCamel(str string) string {
	arr := strings.Split(str, "_")
	if len(arr) <= 1 {
		return str
	}
	first := strings.ToLower(arr[0])
	arr = arr[1:]
	str = strings.Join(arr, " ")
	return first + strings.Replace(strings.Title(str), " ", "", -1)
}

func ExecCmd(dir string, name string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String(), fmt.Errorf("%s\n%s", err.Error(), stderr.String())
	}
	return stdout.String(), nil
}
