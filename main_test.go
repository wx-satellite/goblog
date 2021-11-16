package main

import "testing"

func TestInt64ToString(t *testing.T) {
	res := Int64ToString(12)

	if res == "12" {
		t.Logf("测试通过")
	} else {
		t.Error("测试失败")
	}
}
