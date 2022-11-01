## Go初始化依赖关系
<br>
</br>
### 验证内容：
- main包依赖pkg1、pkg2
- pkg1、pkg2依赖pkg3

<br>
</br>
### 结果
```shell
➜  firstClass git:(master) ✗ go run main.go 
# pkg3只会被初始化一次
pkg3: const c1 has been initialized
pkg3: const c2 has been initialized
pkg3: var v1 has been initialized
pkg3: var v2 has been initialized
pkg3: first init func invoked
pkg3: second init func invoked

pkg1: const c1 has been initialized
pkg1: const c2 has been initialized
pkg1: var v1 has been initialized
pkg1: var v2 has been initialized
pkg1: first init func invoked
pkg1: second init func invoked

pkg2: const c1 has been initialized
pkg2: const c2 has been initialized
pkg2: var v1 has been initialized
pkg2: var v2 has been initialized
pkg2: first init func invoked
pkg2: second init func invoked

main: const c1 has been initialized
main: const c2 has been initialized
main: var v1 has been initialized
main: var v2 has been initialized
main: first init func invoked
main: second init func invoked
```
- pkg3先被初始化，且只被初始化一次
- 每个包内按照常量 -> 变量 -> init()顺序初始化

<br>
</br>
### 总结
- 依赖包按“**深度优先”的次序**进行初始化；
- 每个包内按以“**常量 -> 变量 -> init 函数**”的顺序进行初始化；
- 包内的多个 init 函数按出现次序进行自动调用