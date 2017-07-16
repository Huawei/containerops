package git_test

import (
    "testing"
	"util/git"
	"os"
)



func Test_Clone_HTTP(t *testing.T) {
	testspace := os.Getenv("testspace")

	if err := git.Clone("http://192.168.123.201/yangkghjh/easy-php.git", testspace); err != nil {
		t.Error("Git clone http repository error.")
	} else {
		os.RemoveAll(testspace)
		t.Log("Git clone http repository success.")
	}
}

// func Test_Clone_SSH(t *testing.T) {
// 	testspace := os.Getenv("testspace")

// 	if err := git.Clone("git@192.168.123.201:yangkghjh/easy-php.git", testspace); err != nil {
// 		t.Error("Git clone ssh repository error.")
// 	} else {
// 		os.RemoveAll(testspace)
// 		t.Log("Git clone ssh repository success.")
// 	}
// }