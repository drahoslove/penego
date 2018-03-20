package storage_test

import (
	"fmt"
	"git.yo2.cz/drahoslav/penego/storage"
	"testing"
)

func TestStorage(test *testing.T) {
	st := storage.New()
	st.Set("bool", true)
	st.Set("string", "ahoj")
	st.Set("int", 42)
	st.Set("float", 3.14)
	// st.Set("list", []interface{}{1,2,3,4})

	if st.Bool("bool") != true {
		test.Errorf("Get bool failed: %v", st.Bool("bool"))
	}
	if st.String("string") != "ahoj" {
		test.Errorf("Get string failed: %v", st.String("string"))
	}
	if st.Int("int") != 42 {
		test.Errorf("Get int failed: %v", st.Int("int"))
	}
	if st.Float("float") != 3.14 {
		test.Errorf("Get float failed: %v", st.Float("float"))
	}
	// if st.List("list").([]interface{})[0] != 1 {
	// 	test.Errorf("Get list item failed: %s", st.Get("list"))
	// }

	fmt.Println(st)
}
