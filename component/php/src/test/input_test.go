package input_test

import (
    "testing"
	"util/input"
)

func Test_HandleInput(t *testing.T) {
	keys := []string{
		"a",
		"b",
		"c",
		"d",
	}
	CO_DATA := "a=a_value b=1 c="
	result := map[string]string{}
	result_correct := map[string]string{
		"a": "a_value",
		"b": "1",
		"c": "",
		"d": "",
	}
	
	if err := input.HandleInput(CO_DATA, keys, result); err != nil {
		t.Error("Handle input error.")
	} else if compareMap(result, result_correct) == false {
		t.Errorf("Decode result error: %v.Should be %v", result, result_correct)
	} else {
		t.Log("Handle input success.")
	}
}

func compareMap(map1 map[string]string, map2 map[string]string) bool {
	for k, v := range map2 {
		if v != map1[k] {
			return false
		}
	}

	return true
}