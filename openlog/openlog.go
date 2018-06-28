/*
Package opencivica contains utilities to get opencivica information from AppTransporte's logs
*/
package openlog


import (
	"os"
	"fmt"
	"time"
	"strings"
	"unicode"
	"regexp"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/net/html"
	"github.com/yanatan16/itertools"
	"github.com/dumacp/utils"
)

const (
	PATH_LOGS = "/SD/OpenCivica_Files/logs/"
	SHORT_FORM = "2018-04-17 11:53:48"
)

type Civica interface{}

//Transportation Use object
type UsoTransporte struct {
	TipoTarjeta	string
	CivicaBefore	Civica
	CivicaAfter	Civica
	PTtvId		int	`json:"p_ttv_id:`
	PLocId		int	`json:"p_loc_id"`
	PEquSerie	int	`json:"p_equ_serie"`
	PPotId		int	`json:"p_pot_id"`
	PCarId		int	`json:"p_car_id"`
	PTnvData	int64	`json:"p_tnv_data"`
	PTnvSeqEquip	int	`json:"p_tnv_seq_equip"`
	PTnvValor	int	`json:"p_tnv_valor"`
	PTnvSaldoPost	int	`json:"p_tnv_saldo_post"`
	SaldoBancos	int	`json:"saldoBancos"`
	PFptId		int	`json:"p_fpt_id"`
	PTnvDataAnt	int64	`json:"p_tnv_data_ant"`
	PLocIdAnt	int	`json:"p_loc_id_ant"`
	PTnvValorAnt	int	`json:"p_tnv_valor_ant"`
	PRotIdAnt	int	`json:"p_rot_id_ant"`
	PRotIdAntAnt	int	`json:"p_rot_id_ant_ant"`
	PTnvContTempo	int	`json:"p_tnv_cont_tempo"`
	PFptIdAnt	int	`json:"p_fpt_id_ant"`
	PPerId		int	`json:"p_per_id"`
	PTnvContarj	int	`json:"p_tnv_contarj"`
	PTnvPasos	int	`json:"p_tnv_pasos"`
	PEquSerieAnt	int	`json:"p_equ_serie_ant(NO_DEFINIDO)"`
	PTnvAalorLiq	int	`json:"p_tnv_valor_liq"`
	PTnvValorCred	int	`json:"p_tnv_valor_cred"`
	PTnvSaldoCred	int	`json:"p_tnv_saldo_cred"`
	PTnvSecUsoEnTrayecto	int	`json:"p_tnv_sec_uso_en_trayecto"`
	PSubRutId	int	`json:"p_subrut_id"`
	UsoId		int64	`json:"usoId"`
	Exitoso		bool	`json:"exitoso"`
	VerTarjBloq	bool	`json:"verTarjBloq"`
	VerListaNegra	bool	`json:"verListaNegra"`
	VerLimitSinUso	bool	`json:"verLimitSinUso"`
	VerPerfilLimiteTiempo	bool	`json:"verPerfilLimiteTiempo"`
	VerFechaValMon	bool	`json:"verFechaValMon"`
	VerFechaValBen	bool	`json:"verFechaValBen"`
	MinFaltantesParaCumplirPerfilLimTiempo	int	`json:"minFaltantesParaCumplirPerfilLimTiempo"`
	TarifNoEncontrada	bool	`json:"tarifNoEncontrada"`
	EstadoEscritura	int	`json:"estadoEscritura"`
	SinSaldo	bool	`json:"sinSaldo"`
	SeHaBloqueado	bool	`json:"seHaBloqueado"`
	SaldoAntesDelUso	int	`json:"saldoAntesDelUso"`
	SaldoBancosAntesDelUso	int	`json:"saldoBancosAntesDelUso"`
	TimeoutTorniqueteExced	bool	`json:"timeoutTorniqueteExced"`
	LecturaMs		int	`json:"lecturaMs"`
	ReglaNegocioMs		int	`json:"reglaNegocioMs"`
	EscrituraMs	int	`json:"escrituraMs"`
}

type AppLogData struct {
	TimeRef		int64
	Data		interface{}
}

//Function to Filter in Parse Iteration
type FuncParseLog func(data string) AppLogData

//Iter of App Log Transsaction (transacciones.html*)
func TransactionsLogs(quit <-chan int) itertools.Iter {

	iterFile := make(itertools.Iter)
	go func() {
		defer func() {
			//fmt.Println("Salida TransactionsLogs")
			close(iterFile)
		}()
		files, err := ioutil.ReadDir(PATH_LOGS)
		if err != nil {
			iterFile <- err
			return
		}
		utils.SortFileInfo(files)
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "transacciones.html") {
				//fmt.Println(file.Name())
				/**/
				select {
				case <-quit:
					//fmt.Println("QUIT0 !!!!!!!! 1")
					return
				case iterFile <- file:
				}
				/**/
			}
		}
		return
	}()
	return iterFile
}
//Iter of APP Logs (log.html*)
func AppLogs(quit <-chan int) itertools.Iter {

	iterFile := make(itertools.Iter)
	go func() {
		defer func() {
			close(iterFile)
		}()
		files, err := ioutil.ReadDir(PATH_LOGS)
		if err != nil {
			iterFile <- err
			return
		}
		utils.SortFileInfo(files)
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "log.html") {
				select {
				case <-quit:
					return
				case iterFile <- file:
				}
			}
		}
		return
	}()
	return iterFile
}

/**/
func messageData(node *html.Node, timeout int) string {
	if node.Type == html.ElementNode && node.Data == "td" {
		for _, a := range node.Attr {
			if a.Key == "title" && a.Val == "Message" {
				return node.FirstChild.Data
			}
		}
		return ""
	}
	return ""
}
/**/

/**/
//Extract MessageData from Log html
func MessageData(node *html.Node, iterData itertools.Iter, timeout int, quit <-chan int) bool {
	defer func() {
		if timeout == 0 {
			close(iterData)
		}
	}()
	if node.Type == html.ElementNode && node.Data == "td" {
		for _, a := range node.Attr {
			if a.Key == "title" && a.Val == "Message" {
				if node.FirstChild.Data != "" {
					select {
					case <-quit:
						return false
					case iterData <- node.FirstChild.Data:
					}
				}
				return true
			}
		}
		return true
	}
	for c := node.LastChild; c != nil; c = c.PrevSibling {
		if !MessageData(c, iterData, timeout+1, quit) {
			return false
		}
	}
	return true
}
/**/

func parseUsoLog(data string) (uso UsoTransporte)  {
	fieldsRaw := strings.Split(data, ";;;")
	if len(fieldsRaw) != 3 {
		return
	}

	f1 := func(r rune) bool {
		return unicode.IsSpace(r) ||  (';' == r)
	}
	fields1 := strings.FieldsFunc(fieldsRaw[0],f1)
	usoSlice := make([]string,0)
	for _, field1 := range fields1 {
		els1 := strings.Split(field1,"=")
		if len(els1) > 1 && els1[0] != "" {
			usoSlice = append(usoSlice,fmt.Sprintf("%q: %v", els1[0], els1[1]))
		}
	}

	usoss := strings.Join(usoSlice,",")
	usoss = "{" + usoss + "}"
	if err := json.Unmarshal([]byte(usoss), &uso); err != nil {
		return
	}

	f2 := func(r rune) bool {
		return unicode.IsSpace(r) ||  (';' == r) || ('=' == r)
	}
	fields2 := strings.FieldsFunc(fieldsRaw[1],f2)
        civicaBeforeSlice := make([]string,0)
        for _, field2 := range fields2 {
		els2 := strings.Split(field2,":")
                if len(els2) > 1 && els2[0] != "" {
                        civicaBeforeSlice = append(civicaBeforeSlice,fmt.Sprintf("%q: %v", els2[0], els2[1]))
                }
        }

	civicaBeforess := strings.Join(civicaBeforeSlice,",")
	civicaBeforess = "{" + civicaBeforess + "}"
	var civicaBefore Civica
        if err := json.Unmarshal([]byte(civicaBeforess), &civicaBefore); err == nil {
                uso.CivicaBefore = civicaBefore
        }

	fields3 := strings.FieldsFunc(fieldsRaw[2],f2)
        civicaAfterSlice := make([]string,0)
        for _, field3 := range fields3 {
		els3 := strings.Split(field3,":")
                if len(els3) > 1 && els3[0] != "" {
                        civicaAfterSlice = append(civicaAfterSlice,fmt.Sprintf("%q: %v", els3[0], els3[1]))
                }
        }

	civicaAfterss := strings.Join(civicaAfterSlice,",")
	civicaAfterss = "{" + civicaAfterss + "}"
	var civicaAfter Civica
        if err := json.Unmarshal([]byte(civicaAfterss), &civicaAfter); err == nil {
                uso.CivicaAfter = civicaAfter
        }

	return
}

//Parse Log Transaction in UsosTranspote
func ParseUsosLog(trs itertools.Iter, quit <-chan int) itertools.Iter {
	/**/
	quit0 := make(chan int)
	quit1 := make(chan int)
	go func() {
		id := <-quit
		go func() {
			select {
			case quit0 <- id:
			default:
				close(quit0)
			}
		}()
		go func() {
			select {
			case quit1 <- id:
			default:
				close(quit1)
			}
		}()
	}()
	/**/

        fMapper := func(i interface{}) interface{} {
		switch v:= i.(type) {
                case string:
			return parseUsoLog(v)
		}
		return parseUsoLog("")
        }

	itMap := itertools.MapQuit(fMapper, trs, quit0)

	fFilter := func(i interface{}) bool {
		switch v := i.(type) {
		case UsoTransporte:
			return v != UsoTransporte{}
		}
		return false
	}
	return itertools.FilterQuit(fFilter, itMap, quit1)

}

//Read log Transactions until quit<- channel
func ReadTransactions(timeout int, quit <-chan int) itertools.Iter {

	quit0 := make(chan int)
	quit1 := make(chan int)

	iterFile := TransactionsLogs(quit0)
	itData := make(itertools.Iter)
	go func() {
		id := 0
		defer func() {
			go func() {
				select {
				case quit0 <- id:
				default:
					close(quit0)
				}
			}()
			go func() {
				select {
				case quit1 <- id:
				default:
					close(quit1)
				}
			}()
			close(itData)
		}()
		for file := range iterFile {
			content, err := os.Open(PATH_LOGS + file.(os.FileInfo).Name())
			if err != nil {
				fmt.Println("ERROR: ",err)
			}
			doc, _ := html.Parse(content)
			iterData := make(itertools.Iter)
			go MessageData(doc, iterData, 0, quit1)
			content.Close()
			doc = nil

			for el := range iterData {
				select {
				case id = <-quit:
					return
				case itData <- el:
				}
			}
		}
	}()
	return itData
}

//Function to Parse App Version info
func FuncAppVersionLog(data string) AppLogData {
	appData := AppLogData{}
	if !strings.Contains(data,"Versiones Librearias") {
		return appData
	}

	re := regexp.MustCompile("\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}")
	timeS := re.FindString(data)
	if timeS != "" {
		loc, _ := time.LoadLocation("America/Bogota")
		t, _ := time.ParseInLocation(SHORT_FORM, timeS, loc)
		appData.TimeRef = t.UnixNano()
	}

	/**/
	fieldsRaw := strings.Split(data, "\n")

	if len(fieldsRaw) < 3 {
		return appData
	}

	versions := make(map[string]string)
	for i:=1; i < len(fieldsRaw)-1; i++ {
		els1 := strings.Split(fieldsRaw[i],":")
		if len(els1) > 1 && els1[0] != "" {
			versions[els1[0]] = strings.Trim(els1[1], " ")
		}
	}

	appData.Data = versions

	return appData
}

//Function to Parse Tables Version info in App
func FuncTabVersionLog(data string) AppLogData {
	appData := AppLogData{}
	if !strings.Contains(data,"Se han actualizado las siguientes tablas") && !strings.Contains(data,"ha iniciado con") {
		return appData
	}

	/**/
	re1 := regexp.MustCompile("\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}")
	timeS := re1.FindString(data)
	if timeS != "" {
		loc, _ := time.LoadLocation("America/Bogota")
		t, _ := time.ParseInLocation(SHORT_FORM, timeS, loc)
		appData.TimeRef = t.UnixNano()
	}

	re2 := regexp.MustCompile("siguientes tablas: +\\[(.+):(.+)\\]")
	fields2 := re2.FindStringSubmatch(data)

	if len(fields2) > 2 && fields2[1] != "" && fields2[2] != "" {
		versions2 := map[string]string{ fields2[1]: fields2[2] }
		appData.Data = versions2
		return appData
	}

	re3 := regexp.MustCompile("ha iniciado con (.+) en version: (\\d+\\.{0,1}\\d{0,4})")
        fields3 := re3.FindStringSubmatch(data)

	if len(fields3) > 2 && fields3[1] != "" && fields3[2] != "" {
		versions3 := map[string]string{ fields3[1]: fields3[2] }
		appData.Data = versions3
		return appData
	}

	return appData
}

//Iter to Parse Function in Iter Data Log
func ParseAppLog(f FuncParseLog, trs itertools.Iter, quit <-chan int) itertools.Iter {
	/**/
	quit0 := make(chan int)
	quit1 := make(chan int)
	go func() {
		id := <-quit
		go func() {
			select {
			case quit0 <- id:
			default:
				close(quit0)
			}
		}()
		go func() {
			select {
			case quit1 <- id:
			default:
				close(quit1)
			}
		}()
	}()
	/**/

        fMapper := func(i interface{}) interface{} {
		switch v:= i.(type) {
                case string:
			return f(v)
                }
		return f("")
        }

	it := itertools.MapQuit(fMapper, trs, quit0)

	fFilter := func(i interface{}) bool {
		switch v:= i.(type) {
		case AppLogData:
			return v != AppLogData{}
		}
		return false
	}

	return itertools.FilterQuit(fFilter, it, quit1)
}

//Iter to Parse App Version info in App
func AppVersions(trs itertools.Iter, quit <-chan int) (versions map[string]string) {
	itVersions := ParseAppLog(FuncAppVersionLog, trs, quit)
	vers := <-itVersions
	if vers == nil {
		return
	}
	return vers.(AppLogData).Data.(map[string]string)
}

//Iter to Parse Tables Version info in App
func TabVersions(trs itertools.Iter, quit <-chan int) (versions map[string]string) {
	itVersions := ParseAppLog(FuncTabVersionLog, trs, quit)
	mapVers := make(map[string]string)
	for v:= range itVersions {
		switch tab := v.(type) {
		case AppLogData:
			switch libv := tab.Data.(type) {
			case map[string]string: 
				for key, value := range libv {
					if _, ok := mapVers[key]; !ok {
					mapVers[key] = value
					}
				}
			}
		}
	}

	return mapVers
}

//Iter to App Logs Data
func ReadAppLogs(timeout int, quit <-chan int) itertools.Iter {

	quit0 := make(chan int)
	quit1 := make(chan int)
	iterFile := AppLogs(quit0)

	/**/
	itData := make(itertools.Iter)
	go func() {
		id := 0
		defer func() {
			go func() {
				select {
				case quit0 <- id:
				default:
					close(quit0)
				}
			}()
			go func() {
				select {
				case quit1 <- id:
				default:
					close(quit1)
				}
			}()
			close(itData)
		}()
		for file := range iterFile {
			content, err := os.Open(PATH_LOGS + file.(os.FileInfo).Name())
			if err != nil {
				fmt.Println("ERROR: ",err)
			}

			doc, _ := html.Parse(content)

			iterData := make(itertools.Iter)

			go MessageData(doc, iterData, 0, quit1)
			content.Close()
			doc = nil
			for el := range iterData {
				select {
				case id = <-quit:
					return
				case itData <- el:
				}
			}
		}
	}()
	/**/
	return itData
}

func CountUsosAfter(iterUsosLog itertools.Iter, timeref int64, quit <-chan int) (int64, int64) {

	fFilter1 := func(i interface{}) bool {
		switch v := i.(type) {
		case UsoTransporte:
			ok := v.Exitoso
			if ok {
				return ok
			}
			return false
		}
		return false
	}
	iterUsos := itertools.FilterQuit(fFilter1,iterUsosLog, quit)

	fFilter2 := func(i interface{}) bool {
		switch v := i.(type) {
		case UsoTransporte:
			if v.UsoId >= timeref {
				return true
			}
		}
		return false
	}

	iterUsosAfter := itertools.TakeWhile(fFilter2,iterUsos)

	fReducer := func(memo interface{}, element interface{}) interface{} {
		switch v:= element.(type) {
		case UsoTransporte:
			if v.UsoId > memo.([]int64)[1] {
				memo.([]int64)[1] = v.UsoId
			}
			memo.([]int64)[0] = memo.([]int64)[0] + 1
		}
		return memo
	}

	memo := []int64{0,0}
	countUsos := itertools.Reduce(iterUsosAfter, fReducer, memo)

	return countUsos.(([]int64))[0], countUsos.(([]int64))[1]
}

func CountErrorsAfter(iterUsosLog itertools.Iter, timeref int64, quit <-chan int) (int64, int64) {
	fFilter1 := func(i interface{}) bool {
		switch v := i.(type) {
		case UsoTransporte:
			return !v.Exitoso && v.UsoId > 0
		}
		return false
	}
	iterUsos := itertools.FilterQuit(fFilter1,iterUsosLog, quit)

	fFilter2 := func(i interface{}) bool {
		switch v := i.(type) {
		case UsoTransporte:
			if v.UsoId >= timeref {
				return true
			}
		}
		return false
	}

	iterUsosAfter := itertools.TakeWhile(fFilter2,iterUsos)

	fReducer := func(memo interface{}, element interface{}) interface{} {
		switch v := element.(type) {
		case UsoTransporte:
			if v.UsoId > memo.([]int64)[1] {
				memo.([]int64)[1] = v.UsoId
			}
			memo.([]int64)[0] = memo.([]int64)[0] + 1
		}
		return memo
	}

	memo := []int64{0,0}
	countUsos := itertools.Reduce(iterUsosAfter, fReducer, memo)

	return countUsos.(([]int64))[0], countUsos.(([]int64))[1]
}


