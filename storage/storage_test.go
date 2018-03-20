package storage_test

import (
	"testing"
	"fmt"
	"git.yo2.cz/drahoslav/penego/storage"
)

func TestStorage(test *testing.T) {
	st := storage.New()
	st.Set("bool", true)
	st.Set("string", "ahoj")
	st.Set("int", 42)
	st.Set("list", []interface{}{1,2,3,4})

	if st.Get("bool") != true {
		test.Errorf("Get bool failed: %s", st.Get("bool"))
	}
	if st.Get("string") != "ahoj" {
		test.Errorf("Get string failed: %s", st.Get("string"))
	}
	if st.Get("int") != 42 {
		test.Errorf("Get int failed: %s", st.Get("int"))
	}
	if st.Get("list").([]interface{})[0] != 1 {
		test.Errorf("Get list item failed: %s", st.Get("list"))
	}

	fmt.Println(st)
}