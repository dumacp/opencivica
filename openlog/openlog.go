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
	"io/ioutil"
	"encoding/json"
	"golang.org/x/net/html"
	"github.com/yanatan16/itertools"
	"github.com/dumacp/utils"
)

const (
	PATH_LOGS = "/SD/OpenCivica_Files/logs/"
)

type Civica interface{}

type UsoTransporte struct {
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

type FuncParseLog func(data string) interface{}

func TransactionsLogs() itertools.Iter {

	iterFile := make(itertools.Iter)
	go func() {
		defer close(iterFile)
		files, err := ioutil.ReadDir(PATH_LOGS)
		if err != nil {
			iterFile <- err
			return
		}
		utils.SortFileInfo(files)
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "transacciones.html") {
				//fmt.Println(file.Name())
				iterFile <- file
			}
		}
		return
	}()

	return iterFile
}

func AppLogs() itertools.Iter {

	iterFile := make(itertools.Iter)
	go func() {
		defer close(iterFile)
		files, err := ioutil.ReadDir(PATH_LOGS)
		if err != nil {
			iterFile <- err
			return
		}
		utils.SortFileInfo(files)
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "log.html") {
				//fmt.Println(file.Name())
				iterFile <- file
			}
		}
		return
	}()

	return iterFile
}

/**/
func MessageData(node *html.Node, iterData itertools.Iter, timeout int) {
	defer close(iterData)
	//fmt.Printf("Node: %v\n", node)
	if node.Type == html.ElementNode && node.Data == "td" {
		for _, a := range node.Attr {
			if a.Key == "title" && a.Val == "Message" {
				//fmt.Printf("node.Next %v: %s\n", node.Data, node.FirstChild.Data)
				iterData <- node.FirstChild.Data
				return
			}
		}
		return
	}
	itSlice := make([]itertools.Iter,0)
	for c := node.LastChild; c != nil; c = c.PrevSibling {
		it := make(itertools.Iter)
		itSlice = append(itSlice, it)
		go MessageData(c, it, timeout)
	}

	for _, ch := range itSlice {
		for v := range ch {
			iterData <- v
		}
	}
	if timeout > 0 {
		time.Sleep(time.Millisecond * time.Duration(timeout))
	}
}
/**/

/**
func MessageData(z *html.Tokenizer, iterData itertools.Iter, count int) {
	//fmt.Printf("z.Token %v: %v\n", count, z)
	defer close(iterData)
	for {
		tt := z.Next()
		switch tt {
		case html.StartTagToken:
			t := z.Token()
			if t.Data == "td" {
				inner := z.Next()
				if inner == html.TextToken {
					iterData <- string(z.Text())
				}
			}
		case html.ErrorToken:
			fmt.Printf("ErrorToken: %v\n", z.Err())
			fmt.Printf("ErrorData: %s\n", z.Text())
			return
		}
	}
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

func ParseUsosLog(timeout int) itertools.Iter {
	trs := ReadTransactions(timeout)
        fMapper := func(i interface{}) interface{} {
                var usoi interface{}
                usoi = parseUsoLog(i.(string))
                return usoi
        }

        return  itertools.Map(fMapper, trs)
}


func ReadTransactions(timeout int) itertools.Iter {

	iterFile := TransactionsLogs()

	itSlice := make([]itertools.Iter,0)
	for file := range iterFile {
		content, err := os.Open(PATH_LOGS + file.(os.FileInfo).Name())
		if err != nil {
			fmt.Println("ERROR: ",err)
		}

		doc, _ := html.Parse(content)
		//doc := html.NewTokenizer(strings.NewReader(string(content)))

		iterData := make(itertools.Iter)
		itSlice = append(itSlice, iterData)

		go MessageData(doc, iterData, timeout)
	}

	return itertools.Chain(itSlice...)
}

func FuncAppVersionLog(data string) interface{} {
	if !strings.Contains(data,"Versiones Librearias") {
		return nil
	}

	/**/
	fieldsRaw := strings.Split(data, "\n")

	if len(fieldsRaw) < 3 {
		return nil
	}

	versions := make(map[string]string)
	for i:=1; i < len(fieldsRaw)-1; i++ {
		els1 := strings.Split(fieldsRaw[i],":")
		if len(els1) > 1 && els1[0] != "" {
			versions[els1[0]] = strings.Trim(els1[1], " ")
		}
	}

	return versions
}

func ParseAppLog(f FuncParseLog, timeout int) itertools.Iter {
	trs := ReadAppLogs(timeout)
        fMapper := func(i interface{}) interface{} {
		return f(i.(string))
        }

	return itertools.Map(fMapper, trs)
}

func AppVersions(timeout int) (versions map[string]string) {
	itVersions := ParseAppLog(FuncAppVersionLog, timeout)
        fFilter := func(i interface{}) bool {
		if i != nil {
			return i.(map[string]string) != nil
		}
		return false
        }

	itVers := itertools.Filter(fFilter, itVersions)
	vers := <-itVers
	return vers.(map[string]string)
}

func ReadAppLogs(timeout int) itertools.Iter {

	iterFile := AppLogs()

	itSlice := make([]itertools.Iter,0)
	for file := range iterFile {
		content, err := os.Open(PATH_LOGS + file.(os.FileInfo).Name())
		if err != nil {
			fmt.Println("ERROR: ",err)
		}

		doc, _ := html.Parse(content)
		//doc := html.NewTokenizer(strings.NewReader(string(content)))

		iterData := make(itertools.Iter)
		itSlice = append(itSlice, iterData)

		go MessageData(doc, iterData, timeout)
	}

	return itertools.Chain(itSlice...)
}

func CountUsosAfter(iterUsosLog itertools.Iter, timeref int64) int {

	fFilter1 := func(i interface{}) bool {
		switch v := i.(type) {
		case UsoTransporte:
			return v.Exitoso
		}
		return false
	}
	iterUsos := itertools.Filter(fFilter1,iterUsosLog)

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
		switch element.(type) {
		case UsoTransporte:
			return memo.(int) + 1
		}
		return memo.(int)
	}

	countUsos := itertools.Reduce(iterUsosAfter, fReducer, 0)

	return countUsos.(int)
}

func CountErrorsAfter(iterUsosLog itertools.Iter, timeref int64) int {

	fFilter1 := func(i interface{}) bool {
		switch v := i.(type) {
		case UsoTransporte:
			return !v.Exitoso && v.UsoId > 0
		}
		return false
	}
	iterUsos := itertools.Filter(fFilter1,iterUsosLog)

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
		switch element.(type) {
		case UsoTransporte:
			return memo.(int) + 1
		}
		return memo.(int)
	}

	countUsos := itertools.Reduce(iterUsosAfter, fReducer, 0)

	return countUsos.(int)
}


