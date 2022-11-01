### 重置包级的变量值

以flag包为例

包级变量CommandLine，通过NewFlagSet初始化
```go
var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	f := &FlagSet{
		name:          name,
		errorHandling: errorHandling,
	}
	f.Usage = f.defaultUsage
	return f
}

```
Usage被默认为f.defaultUsage，且没有提供其他New方法可以自定义Usage

如果想自定义Usage，该如何做？flag通过init函数修改了CommandLine变量的值
```go
func init() {
	CommandLine.Usage = commandLineUsage // 重置
}


// commandLineUsage是不能导出的，但包含了可导出Usage
func commandLineUsage() {
	Usage()
}

// 用户自定义Usage既可
var Usage = func() {
	fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	PrintDefaults()
}
```

