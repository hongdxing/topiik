package test

import (
	"fmt"
	"testing"
	"topiik/internal/util"
)

func TestSplitCommandLine(t *testing.T) {
	strs := util.SplitCommandLine(`"aaa" bbb "this, for sanity, should not be parts"`)
	fmt.Println(strs)
	if len(strs) != 3 {
		t.Fatal()
	}

}

func TestSplitCommandLineWithEscape(t *testing.T) {
	strs := util.SplitCommandLine(`"aaa" bbb "this, for sanity,\" should not be parts'"`)
	for _, s := range strs {
		fmt.Println(s)
	}
	if len(strs) != 3 {
		t.Fatal()
	}
}

func TestSplitCommandLineWithUnpairedQuote(t *testing.T) {
	strs := util.SplitCommandLine(`"aaa bbb "this, for sanity,\" should not be parts'"`)
	fmt.Println(strs)
	for _, s := range strs {
		fmt.Println(s)
	}
	if len(strs) != 3 {
		t.Fatal()
	}
}
