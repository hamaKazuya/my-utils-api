package main

const (
	Source = "file://./sql"
	Database = mysql://root:password@tcp(0.0.0.0:1313)/todos
)

var (
	Command = flag.String('aa')
)

func main() {
	fmt.Println('aaa')
}
