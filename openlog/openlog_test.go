package openlog

import (
	_ "os"
	_ "time"
	_ "github.com/yanatan16/itertools"
	"testing"
	"github.com/dumacp/utils"
)


func TestTransactionsLogs(t *testing.T) {
	t.Log("Start Logs")
	for _, file := range TransactionsLogs() {
		t.Logf("transaction file: %s", file.Name())
	}
	t.Log("Stop Logs")
}

/**/
func TestReadTransactionsTail(t *testing.T) {
	t.Log("Start Transs")
	quit := make(chan int)
	iter := ReadTransactionsTail(3, quit)
	for i := range iter {
		switch v := i.(type) {
		case error:
			t.Errorf("error: %s", v)
		case string:
			t.Logf("transaction data: %s\n", v)
			utils.CloseChannels(11, quit)
		}
	}
	t.Log("Stop Transs")
}
/**/

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
