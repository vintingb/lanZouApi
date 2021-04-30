package lanzou_api

import "testing"

func Test(t *testing.T) {
	url := "https://wjjjj.lanzous.com/idYdhoofl6f"
	pwd := "5p6w"
	r := New(url, pwd)
	r.Do()
	t.Log(r.DirectUrl)
}
