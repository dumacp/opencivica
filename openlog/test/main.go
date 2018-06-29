package main

import (
	"log"
	"time"
	"runtime"
	"flag"
	"github.com/dumacp/opencivica/openlog"
	"github.com/dumacp/utils"
	"github.com/yanatan16/itertools"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")


func main() {
        //log.Printf("versions: %q\n", openlog.AppVersions(0))

flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

	for {
		/**
		log.Printf("goroutines start: %d\n", runtime.NumGoroutine())
		quit1 := make(chan bool)
		itLogs := openlog.ReadAppLogs(0, quit1)
		//iterLogs := itertools.Tee(itLogs, 2)
		log.Printf("goroutines step1: %d\n", runtime.NumGoroutine())

		iterLibs1 := openlog.ParseAppLog(openlog.FuncAppVersionLog, itLogs)
		log.Printf("goroutines step2: %d\n", runtime.NumGoroutine())

		log.Printf("App version: %q\n", (<-iterLibs1).(openlog.AppLogData).Data)
		quit1 <- true
		//close(iterLibs1)
		log.Printf("goroutines step3: %d\n", runtime.NumGoroutine())

		/**
		quit2 := make(chan int)
		itLogs2 := openlog.ReadAppLogs(0, quit2)
		iterLibs2 := openlog.ParseAppLog(openlog.FuncTabVersionLog, itLogs2)
		log.Printf("goroutines step4: %d\n", runtime.NumGoroutine())
		mapVers := make(map[string]string)
		for v:= range iterLibs2 {
			libv := v.(openlog.AppLogData).Data.(map[string]string)
			for key, value := range libv {
				if _, ok := mapVers[key]; !ok {
					mapVers[key] = value
				}
			}
		}
		utils.CloseChannel(quit2, 12)
		log.Printf("Lib versions: %q\n", mapVers)
		/**/


		/**
		quit3 := make(chan int)
                it := openlog.ReadAppLogs(0, quit3)
		//tee1 := itertools.Tee(it, 2)
		//libVersions := openlog.TabVersions(it)
		appVersions := openlog.AppVersions(it)



		//log.Printf("Lib versions: %q\n", libVersions)
		log.Printf("app versions: %q\n", appVersions)
		utils.CloseChannel(quit3, 13)
		utils.FinishChannel(it)
		//utils.FinishChannel(tee1[1])


		/**/

		quit4 := make(chan int)
		iterTranss := openlog.ReadTransactions(0,quit4)
		iterUsosLog := openlog.ParseUsosLog(iterTranss)
		//log.Printf("error count: %v\n", openlog.CountUsosAfter(iterUsosLog,0))


		fFilter0 := func(i interface{}) bool {
			return i.(openlog.UsoTransporte).UsoId > 1530225830461
			//return i.(openlog.UsoTransporte).UsoId > 1523989505114
			//return i.(openlog.UsoTransporte).UsoId > 0
		}

		itLast := itertools.TakeWhile(fFilter0, iterUsosLog)

		it := itertools.Tee(itLast, 2)

		fFilter1 := func(i interface{}) bool {
			switch v := i.(type) {
			case openlog.UsoTransporte:
				return v.Exitoso
			}
			return false
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


		utils.CloseChannel(14, quit4)
		utils.FinishChannel(itErrors, itUsos, it[0], it[1], itLast, iterUsosLog, iterTranss)

		time.Sleep(time.Second * 6)
		num := runtime.NumGoroutine()
		log.Printf("\n\n#############     goroutines end: %d\n\n", num)
		if num > 1 {
			break
		}



	}

if *memprofile != "" {
        f, err := os.Create(*memprofile)
        if err != nil {
            log.Fatal("could not create memory profile: ", err)
        }
        runtime.GC() // get up-to-date statistics
        if err := pprof.WriteHeapProfile(f); err != nil {
            log.Fatal("could not write memory profile: ", err)
        }
        f.Close()
    }

}
