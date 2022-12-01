package main

import (
	"testing"
)

// 是要对MaleCount方法编写单元测试代码。对于这种依赖外部数据库操作的方法，
// 我们的惯例是使用“伪对象（fake object）”来冒充真实的Stmt接口实现
// 就是Stmt接口类型的方法集合中有四个方法，而MaleCount函数只使用了Stmt接口的一个方法Exec。
// 如果我们针对每个测试用例所用的伪对象都实现这四个方法，那么这个工作量有些大

// 方法：
// 建立了一个fakeStmtForMaleCount的伪对象类型,嵌入了Stmt接口类型
// 这样fakeStmtForMaleCount就实现了Stmt接口
// 我们也实现了快速建立伪对象的目的。接下来我们只需要为fakeStmtForMaleCount实现MaleCount所需的Exec方法，
// 就可以满足这个测试的要求了。
type fakeStmtForMaleCount struct {
	Stmt
}

func (fakeStmtForMaleCount) Exec(stmt string, args ...string) (Result, error) {
	return Result{Count: 5}, nil
}

func TestEmployeeMaleCount(t *testing.T) {
	f := fakeStmtForMaleCount{}
	c, _ := MaleCount(f)
	if c != 5 {
		t.Errorf("want: %d, actual: %d", 5, c)
		return
	}
}

