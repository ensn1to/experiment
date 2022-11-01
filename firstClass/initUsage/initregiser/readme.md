### 在init函数中实现“注册模式”
<br/></br>
通过在 init 函数中注册自己的实现的模式，就有效降低了 Go 包对外的直接暴露，尤其是包级变量的暴露，从而避免了外部通过包级变量对包状态的改动 \n
比如通过"lib/pq”库访问pg数据库

```go
import (
    "database/sql"
    _ "github.com/lib/pq"
)

func main() {
    // sql.Open只要传入驱动名称就可以直接拿到数据库句柄
    db, err := sql.Open("postgres", "user=pqgotest dbname=pqgotest sslmode=verify-full")
    if err != nil {
        log.Fatal(err)
    }
    
    age := 21
    rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)
    ...
}
```

查看pq的包，有init， 把postgres注册到了sql包的驱动中，从而避免了pq包的包级变量暴露

```go
func init() {
    sql.Register("postgres", &Driver{})
}
```

其实sql.Open函数就是这个模式中的工厂方法，它根据外部传入的驱动名称“生产”出不同类别的数据库实例句柄。

”注册模式“示例2，读取图片大小
```go
package main

import (
    "fmt"
    "image"
    _ "image/gif" // 以空导入方式注入gif图片格式驱动
    _ "image/jpeg" // 以空导入方式注入jpeg图片格式驱动
    _ "image/png" // 以空导入方式注入png图片格式驱动
    "os"
)

func main() {
    // 支持png, jpeg, gif
    width, height, err := imageSize(os.Args[1]) // 获取传入的图片文件的宽与高
    if err != nil {
        fmt.Println("get image size error:", err)
        return
    }
    fmt.Printf("image size: [%d, %d]\n", width, height)
}

func imageSize(imageFile string) (int, int, error) {
    f, _ := os.Open(imageFile) // 打开图文文件
    defer f.Close()

    img, _, err := image.Decode(f) // 对文件进行解码，得到图片实例
    if err != nil {
        return 0, 0, err
    }

    b := img.Bounds() // 返回图片区域
    return b.Max.X, b.Max.Y, nil
}
```

以上示例支持png,jpeg,gif三种图片格式，原因是因为image/png, image/jpeg, image/gif包都通过init函数把自己”注册“到了image的支持格式列表中了
```go
// 查看image包中的各init，各类型将自己注册到配置中
// $GOROOT/src/image/png/reader.go
func init() {
	image.RegisterFormat("png", pngHeader, Decode, DecodeConfig)
}

// $GOROOT/src/image/jpeg/reader.go
func init() {
	image.RegisterFormat("jpeg", "\xff\xd8", Decode, DecodeConfig)
}

// $GOROOT/src/image/gif/reader.go
func init() {
	image.RegisterFormat("gif", "GIF8?a", Decode, DecodeConfig)
}
```