package ion

import (
	"fmt"
	"testing"
)

func TestSharedSymbolTable(t *testing.T) {
	st := NewSharedSymbolTable("test", 2, []string{
		"abc",
		"def",
		"foo'bar",
		"null",
		"def",
		"ghi",
	})

	if st.Name() != "test" {
		t.Errorf("wrong name: %v", st.Name())
	}
	if st.Version() != 2 {
		t.Errorf("wrong version: %v", st.Version())
	}
	if st.MaxID() != 5 {
		t.Errorf("wrong maxid: %v", st.MaxID())
	}

	testFindByName(t, st, "def", 2)
	testFindByName(t, st, "null", 4)
	testFindByName(t, st, "bogus", 0)

	testFindByID(t, st, 0, "")
	testFindByID(t, st, 2, "def")
	testFindByID(t, st, 4, "null")
	testFindByID(t, st, 6, "")

	testString(t, st, `$ion_shared_symbol_table::{name:"test",version:2,symbols:["abc","def","foo'bar","null","ghi"]}`)
}

func TestLocalSymbolTable(t *testing.T) {
	st := NewLocalSymbolTable(nil, []string{"foo", "bar"})

	if st.MaxID() != 2 {
		t.Errorf("wrong maxid: %v", st.MaxID())
	}

	testFindByName(t, st, "$ion", 0)
	testFindByName(t, st, "foo", 1)
	testFindByName(t, st, "bar", 2)
	testFindByName(t, st, "bogus", 0)

	testFindByID(t, st, 0, "")
	testFindByID(t, st, 1, "foo")
	testFindByID(t, st, 2, "bar")
	testFindByID(t, st, 3, "")

	testString(t, st, `$ion_symbol_table::{symbols:["foo","bar"]}`)
}

func TestLocalSymbolTableWithImports(t *testing.T) {
	imports := []SharedSymbolTable{V1SystemSymbolTable}
	st := NewLocalSymbolTable(imports, []string{
		"foo",
		"bar",
	})

	if st.MaxID() != 11 { // 9 from $ion.1, 2 local.
		t.Errorf("wrong maxid: %v", st.MaxID())
	}

	testFindByName(t, st, "$ion", 1)
	testFindByName(t, st, "$ion_shared_symbol_table", 9)
	testFindByName(t, st, "foo", 10)
	testFindByName(t, st, "bogus", 0)

	testFindByID(t, st, 0, "")
	testFindByID(t, st, 1, "$ion")
	testFindByID(t, st, 9, "$ion_shared_symbol_table")
	testFindByID(t, st, 10, "foo")
	testFindByID(t, st, 11, "bar")
	testFindByID(t, st, 12, "")

	testString(t, st, `$ion_symbol_table::{imports:[{name:"$ion",version:1,max_id:9}],symbols:["foo","bar"]}`)
}

func TestSymbolTableBuilder(t *testing.T) {
	b := NewSymbolTableBuilder(V1SystemSymbolTable)

	id, ok := b.Add("name")
	if ok {
		t.Error("Add(name) returned true")
	}
	if id != 4 {
		t.Errorf("Add(name) returned %v", id)
	}

	id, ok = b.Add("foo")
	if !ok {
		t.Error("Add(foo) returned false")
	}
	if id != 10 {
		t.Errorf("Add(foo) returned %v", id)
	}

	id, ok = b.Add("foo")
	if ok {
		t.Error("Second Add(foo) returned true")
	}
	if id != 10 {
		t.Errorf("Second Add(foo) returned %v", id)
	}

	st := b.Build()
	if st.MaxID() != 10 {
		t.Errorf("maxid returned %v", st.MaxID())
	}

	testFindByName(t, st, "$ion", 1)
	testFindByName(t, st, "foo", 10)
	testFindByName(t, st, "bogus", 0)

	testFindByID(t, st, 1, "$ion")
	testFindByID(t, st, 10, "foo")
	testFindByID(t, st, 11, "")
}

func testFindByName(t *testing.T, st SymbolTable, sym string, expected int) {
	t.Run("FindByName("+sym+")", func(t *testing.T) {
		actual, ok := st.FindByName(sym)
		if expected == 0 {
			if ok {
				t.Fatalf("unexpectedly found: %v", actual)
			}
		} else {
			if !ok {
				t.Fatal("unexpectedly not found")
			}
			if actual != expected {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		}
	})
}

func testFindByID(t *testing.T, st SymbolTable, id int, expected string) {
	t.Run(fmt.Sprintf("FindByID(%v)", id), func(t *testing.T) {
		actual, ok := st.FindByID(id)
		if expected == "" {
			if ok {
				t.Fatalf("unexpectedly found: %v", actual)
			}
		} else {
			if !ok {
				t.Fatal("unexpectedly not found")
			}
			if actual != expected {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		}
	})
}

func testString(t *testing.T, st SymbolTable, expected string) {
	t.Run("String()", func(t *testing.T) {
		actual := st.String()
		if actual != expected {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})
}
