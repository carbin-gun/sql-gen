package writers

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/carbin-gun/sql-gen/model"
	"github.com/pkg/errors"
)

var writers map[string]Writer

//Writer write-way
type Writer interface {
	Write(map[string]*model.TableMeta) error
}

func init() {
	writers = map[string]Writer{
		"console": ConsoleWriter{},
		"file":    FileWriter{},
	}
}

//Lookup to query key-specifed writer
func Lookup(w string) (Writer, bool) {
	writer, ok := writers[w]
	return writer, ok
}

//ConsoleWriter write to console
type ConsoleWriter struct{}

func (c ConsoleWriter) Write(data map[string]*model.TableMeta) error {

	fmt.Println("<==================================>")
	for k, v := range data {
		fmt.Printf("table: <%s>\n", k)
		fmt.Printf("sql: < %s >\n\n", v.JointColumns())
	}
	fmt.Println("<==================================>")

	return nil

}

//FileWriter writes data to go file
type FileWriter struct{}

var (
	//if call method, omit parentheses is ok
	tpl = `
package tables
{{range .data}}
type {{.StructName}} struct{}
func (t *{{.StructName}}) Column() string {
	return "{{.JointColumns}}"
}
{{end}}	
	`
	metaTpl = `
package tables
type Table interface{
	Column() string
}	

	`
)

const (
	DirName      = "database/tables/"
	FileName     = "tables.go"
	metaFileName = "meta.go"
)

func (FileWriter) MakeDir() {
	if err := os.MkdirAll(filepath.Dir(DirName), 0755); err != nil {
		panic("error create database/tables directory")
	}
}

func (w FileWriter) WriteMeta(workDir string) {
	filename := filepath.Join(DirName, metaFileName)
	file, err := os.OpenFile(filepath.Join(workDir, filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	t, err := template.New("FILE_WRITER").Parse(metaTpl)
	if err != nil {
		panic(fmt.Errorf("meta template create error:%v", err))
	}
	if err := t.Execute(file, nil); err != nil {
		panic(fmt.Sprintf("meta render error:%v", err))
	}
}
func (w FileWriter) WriteData(workDir string, data map[string]*model.TableMeta) {
	filename := filepath.Join(DirName, FileName)
	file, err := os.OpenFile(filepath.Join(workDir, filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(fmt.Sprintf("open data file error:%v", err))
	}

	defer file.Close()
	t, err := template.New("FILE_WRITER").Parse(tpl)
	if err != nil {
		panic(fmt.Sprintf("parse writer data tempalte error:%v", err))
	}
	params := map[string]interface{}{
		"data": data,
	}
	err = t.Execute(file, params)
	if err != nil {
		panic(fmt.Sprintf("write to data file [tables.go] error:%v", err))
	}
	err = format(filename)
	if err != nil {
		panic(fmt.Sprintf("go format data file [tables.go] error:%v", err))
	}
}

func (w FileWriter) Write(data map[string]*model.TableMeta) error {
	w.MakeDir()
	workDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "os.Getwd() error")
	}
	w.WriteMeta(workDir)
	w.WriteData(workDir, data)
	fmt.Println("write OK!!!")
	return nil
}

func format(file string) error {
	cmd := exec.Command("gofmt", "-w", file)
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "gofmt format code error")
	}
	return nil
}
