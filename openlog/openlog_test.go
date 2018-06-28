package openlog

import (
	"os"
	_ "time"
	_ "github.com/yanatan16/itertools"
	"testing"
	"github.com/dumacp/utils"
)


func TestTransactionsLogs(t *testing.T) {
	t.Log("Start Logs")
	quit1 := make(chan int)
	for i := range TransactionsLogs(quit1) {
		switch v := i.(type) {
		case error:
			t.Errorf("error: %s", v)
		case os.FileInfo:
			t.Logf("transaction file: %s", v.Name())
		}
	}
	utils.CloseChannels(11, quit1)
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
	quit1 := make(chan int)
	quit2 := make(chan int)
	quit3 := make(chan int)
	trs := ReadTransactions(0, quit1)
	usosLog := ParseUsosLog(trs, quit2)
	usos, tRef := CountUsosAfter(usosLog, 0, quit3)
	//usos, tRef := CountUsosAfter(usosLog, 1523984060009, quit3)

	utils.CloseChannels(11, quit1, quit2, quit3)
        t.Logf("uso count: %v (%v)", usos, tRef)
}
