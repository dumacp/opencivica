package main

import (
	"log"
	"github.com/dumacp/opencivica/openlog"
	"github.com/yanatan16/itertools"
)


func main() {
        //log.Printf("error count: %v\n", openlog.CountUsosAfter(0))
        log.Printf("versions: %q\n", openlog.AppVersions(0))

	/**/
	iterUsosLog := openlog.ParseUsosLog(8)

	it := itertools.Tee(iterUsosLog, 2)

	fFilter1 := func(i interface{}) bool {
		return i.(openlog.UsoTransporte).Exitoso
	}
	fFilter2 := func(i interface{}) bool {
		return !i.(openlog.UsoTransporte).Exitoso
	}

	itUsos := itertools.Filter(fFilter1, it[0])
	itErrors := itertools.Filter(fFilter2, it[1])

	fReducer := func(memo interface{}, element interface{}) interface{} {
		return memo.(int) + 1
	}

	countUsos := itertools.Reduce(itUsos, fReducer, 0)
	countErrors := itertools.Reduce(itErrors, fReducer, 0)

	log.Printf("error count: %v\n", countErrors)
	log.Printf("usos count: %v\n", countUsos)

	/**/
}
