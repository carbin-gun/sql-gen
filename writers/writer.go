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
	package store
	const(
		{{range .data}}
		{{.StructName}}Columns = "{{.JointColumns}}"
		{{end}}
	)
	`
)

func (FileWriter) Write(data map[string]*model.TableMeta) error {
	workDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "os.Getwd() error")
	}
	filename := "columns.go"
	file, err := os.OpenFile(filepath.Join(workDir, filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrap(err, "write to file error")
	}

	defer file.Close()
	t, err := template.New("FILE_WRITER").Parse(tpl)
	if err != nil {
		return errors.Wrap(err, "parse writer tempalte error")
	}
	params := map[string]interface{}{
		"data": data,
	}
	err = t.Execute(file, params)
	if err != nil {
		return errors.Wrap(err, "write to  columns.go file error")
	}
	err = format(filename)
	if err != nil {
		return err
	}
	fmt.Println("write to columns.go OK")
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
