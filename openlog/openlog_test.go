package openlog

import (
	"os"
	_ "time"
	"github.com/yanatan16/itertools"
	"testing"
)


func TestTransactionsLogs(t *testing.T) {
	t.Log("Start Logs")
	for i := range TransactionsLogs() {
		switch v := i.(type) {
		case error:
			t.Errorf("error: %s", v)
		case os.FileInfo:
			t.Logf("transaction file: %s", v.Name())
		}
	}
	t.Log("Stop Logs")
}

/**
func TestReadTransactions(t *testing.T) {
	t.Log("Start Transs")
	for i := range ReadTransactions(10) {
		switch v := i.(type) {
		case error:
			t.Errorf("error: %s", v)
		case string:
			t.Logf("transaction data: %s", v)
		}
	}
	t.Log("Stop Transs")
}
/**/

/**
func TestParseUsosLog(t *testing.T) {
	t.Log("Start USO")
        for i := range ParseUsosLog() {
                switch v := i.(type) {
                case error:
                        t.Errorf("error: %s", v)
                case UsoTransporte:
                        t.Logf("uso data: %+v", v)
                }
        }
        t.Log("Stop USO")
}
/**/


func TestCountUsoAfter(t *testing.T) {
        t.Logf("uso count: %v", CountUsosAfter(1523984060009))

	it := ParseUsosLog(10)

	fFilter := func(i interface{}) bool {
		switch v:= i.(type) {
		case UsoTransporte:
			return v.Exitoso
		}
		return false
	}

	itExitosos := itertools.Filter(fFilter, it)

	fMapper := func(i interface{}) interface{} {
		return i.(UsoTransporte).UsoId
	}

	itUsoId := itertools.Map(fMapper, itExitosos)

	list := itertools.List(itUsoId)

	t.Logf("lista: %v\n, count: %v", list, len(list))
}
