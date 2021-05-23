package lanZouApi

import "testing"

func Test(t *testing.T) {
	url := "https://www.lanzoui.com/idYdhoofl6f"
	pwd := "5p6w"
	r := New(url, pwd)
	r.Do()
	t.Log(r.DirectUrl)
}
