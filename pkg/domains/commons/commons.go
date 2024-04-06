package commons

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/commonsdtos"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Commons interface {
	NewUUID() string
	IsValidUUID(u string) (bool, error)
	/*
		BuscarArchivos permite obtener un slice con la informacion de los archivos que se encuentra en un directorio.
		recibe como parametro la ruta del directorio que se desea observar
	*/
	LeerDirectorio(rutaFTP string) ([]fs.FileInfo, error)
	// /*
	// 	MoverArchivos permite mover un archivo en otra ubicación se le debe pasar los siguientes parametros:
	// 	- ruta de origen
	// 	- ruta de destino
	// 	- nombre del archivo
	// */
	// MoverArchivos(rutaOrigen, rutaDestino, nombreArchivo string) error

	/*
		BorrarARchivo permite borrar un archivo en un directorio.
		recibe como paramentro la ruta del directorio y el nombre del archivo
	*/
	BorrarArchivo(rutaFTP, nombreArchivo string) error
	/*
			 BorrarDirectorio permite borrar  el directorio temporal creado para alojar los archivos de cierre de lote.
		 	recibe como paramentro la ruta del directorio temporal
	*/
	BorrarDirectorio(ruta string) error
	//Crea un nuevo archivo se debe expecificar ruta completa con el nombre
	CreateFile(ruta string) (archivo *os.File, erro error)
	//Remove un archivo se debe especificar la ruta completa con el nombre
	RemoveFile(ruta string) (erro error)
	//Crea la ruta completa con el nombre del archivo para que se use al crear
	CreateFileName(file commonsdtos.FileName) string
	//Elimina las tildes y caracteres especiales de un string y lo transforma en mayuscula
	NormalizeStrings(str string) (string, error)
	//Se utiliza para comprimir uno o mas archivos
	ZipFiles(request commonsdtos.ZipFilesRequest) (erro error)
	//Se utiliza para guardar un archivo pdf
	SaveFiberPdf(file *multipart.FileHeader, ruta string, fiber *fiber.Ctx) (erro error)
	/*
		Función para crear el mensage que se enviará por correo
		to - email destinatario
		from - email remitente
		value - corpo ya en formato html del mensaje
	*/
	CreateMessage(to []string, from, value string, Subject string) string

	FormatFecha() (fechaI time.Time, fechaF time.Time, erro error)

	LeerArchivo(patch string) (archivo *os.File, erro error)
	EscribirArchivo(datos string, file *os.File) (erro error)
	GuardarCambios(file *os.File) (erro error)

	ConvertirFormatoFecha(fecha string) string

	ConvertirFecha(fecha string) string

	RemoveAccents(valor string) (string, error)

	// retorna la fecha con HH:mm:ss limites de esa fecha en formato string
	GetDateLastMoment(fecha time.Time) (fechaISO string)
	// retorna la fecha con HH:mm:ss iniciales de esa fecha en formato string
	GetDateFirstMoment(fecha time.Time) (fechaISO string)

	ConvertirFechaToDDMMYYYY(fecha string) string

	ToBase64(b []byte) string
	// retorna la fecha inicial con su momento inicial del mes, y la fecha actual con su momento final
	// Se utiliza para calculo  de un mes corrientes
	GetFechaInicioActualMes() (FechaInicio, FechaFin string, erro error)
	// recibe slice de numeros y los pasa a slice de string
	NumberSliceToString(input []uint) string
	// recibe una fecha string con formato "yyyy-mm-dd" y retorna string el ultimo momento con fechaISO de tipo 2023-09-01T23:59:59Z
	DateYMDtoDateLastMoment(fechaIN string) (fechaISO string, erro error)
	// recibe una fecha string con formato "yyyy-mm-dd" y retorna string el primer momento con fechaISO de tipo 2023-09-01T00:00:00Z
	DateYMDtoDateFirstMoment(fechaIN string) (fechaISO string, erro error)
	// recibe Time. retorna string first moment
	DateStringToTimeFirstMoment(fechaString string) (fechaTime time.Time, erro error)
	// recibe Time. retorna string last moment
	DateStringToTimeLastMoment(fechaString string) (fechaTime time.Time, erro error)
	// recibe dos fechas time y retorna string en formato YYYYMM. Sirve para periodos de fechas
	DateTimeToYYYYMM(fechaInicio, fechaFin time.Time) (fecha string, erro error)
}

type commons struct {
	fileRepository FileRepository
}

func NewCommons(fl FileRepository) Commons {
	return &commons{
		fileRepository: fl,
	}
}

func (c commons) NewUUID() string {
	return uuid.NewV4().String()
}

func (c commons) IsValidUUID(u string) (bool, error) {
	_, err := uuid.FromString(u)
	if err != nil {
		return false, fmt.Errorf(ERROR_UUID)
	}
	return true, nil
}

func (c commons) CreateMessage(to []string, from, value string, Subject string) string {

	body := value
	header := make(map[string]string)
	header["From"] = from
	header["To"] = to[0]
	header["Subject"] = Subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""

	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	return message

}

func (c commons) CreateFileName(file commonsdtos.FileName) string {
	if len(file.FechaEspecifica) != 0 {
		return fmt.Sprintf("%s%s_%v.%s", file.RutaBase, file.Nombre, file.FechaEspecifica, file.Extension)
	}
	if file.UsaFecha {
		fechaActual := time.Now()
		fechaFormato := fmt.Sprintf("%02d%02d%d%02d%02d%02d",

			fechaActual.Day(), fechaActual.Month(), fechaActual.Year(),
			fechaActual.Hour(), fechaActual.Minute(), fechaActual.Second())

		return fmt.Sprintf("%s%s_%v.%s", file.RutaBase, file.Nombre, fechaFormato, file.Extension)

	} else {
		return fmt.Sprintf("%s%s.%s", file.RutaBase, file.Nombre, file.Extension)
	}

}

func (c commons) CreateFile(ruta string) (archivo *os.File, erro error) {

	var nombreValido = regexp.MustCompile(`([a-zA-Z0-9\s_\\.\-\(\):])+(.doc|.docx|.pdf|.txt|.csv|.xlsx|.xml|.zip)$`)

	if nombreValido.MatchString(ruta) {

		// Verifica si el archivo existe
		// si no existe lo crear
		if c.fileRepository.ExisteArchivo(ruta) {
			var file, err = c.fileRepository.CrearArchivo(ruta)

			if err != nil {
				logs.Error(err.Error())
				erro = fmt.Errorf(ERROR_FILE_CREATE)
			} else {
				archivo = file
			}
		} else {
			erro = fmt.Errorf(ERROR_FILE_EXIST)
		}

	} else {
		erro = fmt.Errorf(ERROR_FILE_NAME)
	}

	return
}

func (c commons) RemoveFile(ruta string) (erro error) {

	erro = c.fileRepository.EliminarArchivo(ruta)

	return

}

func (c commons) LeerDirectorio(rutaFTP string) ([]fs.FileInfo, error) {
	// lee el contenido del directorio que se le pasa por parammetro
	archivos, erro := ioutil.ReadDir(rutaFTP)
	//fmt.Printf("%T", archivos)
	if erro != nil {
		logs.Error(ERROR_READ_ARCHIVO + erro.Error())
		return nil, errors.New(ERROR_READ_ARCHIVO)
	}
	return archivos, nil
}

func (c commons) BorrarArchivo(rutaFTP, nombreArchivo string) error {
	err := os.Remove(rutaFTP + "/" + nombreArchivo)
	if err != nil {
		logs.Error(ERROR_REMOVER_ARCHIVO + err.Error())
		return errors.New(ERROR_REMOVER_ARCHIVO)
	}
	return nil
}
func (c commons) BorrarDirectorio(ruta string) error {
	err := os.RemoveAll(ruta) //+ "/" + nombreArchivo
	if err != nil {
		logs.Error(ERROR_REMOVER_DIRECTORIO + err.Error())
		return errors.New(ERROR_REMOVER_DIRECTORIO)
	}
	return nil
}
func (c commons) NormalizeStrings(str string) (string, error) {
	var normalizer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, err := transform.String(normalizer, str)
	if err != nil {
		logs.Error(err.Error())
		return "", fmt.Errorf(ERROR_NORMALIZAR)
	}
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		logs.Error(err.Error())
		return "", fmt.Errorf(ERROR_NORMALIZAR)
	}
	s = re.ReplaceAllString(s, " ")

	return strings.ToUpper(s), nil
}

func (c *commons) ZipFiles(request commonsdtos.ZipFilesRequest) (erro error) {

	erro = request.IsValid()

	if erro != nil {
		return
	}

	archivo, erro := c.CreateFile(request.NombreArchivo)

	if erro != nil {
		return
	}

	defer archivo.Close()

	zipWriter := zip.NewWriter(archivo)

	defer zipWriter.Close()

	// Agrega los archivos a al zip
	for _, file := range request.Rutas {
		if erro = c._addFileToZip(zipWriter, file); erro != nil {
			return erro
		}
	}
	return
}

func (c *commons) _addFileToZip(zipWriter *zip.Writer, infoFile commonsdtos.InfoFile) error {

	fileToZip, err := c.fileRepository.AbrirArchivo(infoFile.RutaCompleta)

	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Obtiene la información del archivo
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	//Acá se pone la ruta completa hay que ver si no se debería
	//guardar el nombre no mas.
	header.Name = infoFile.NombreArchivo

	//Esto puede ser que no sea necesario porque sirve para comprimir
	//de forma mas eficiente
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func (c *commons) SaveFiberPdf(file *multipart.FileHeader, ruta string, fiber *fiber.Ctx) (erro error) {

	fiber.SaveFile(file, ruta)

	return
}

func (c *commons) FormatFecha() (fechaI time.Time, fechaF time.Time, erro error) {
	startTime := time.Now()
	fechaConvert := startTime.Format("2006-01-02") //YYYY.MM.DD
	fec := strings.Split(fechaConvert, "-")

	dia, err := strconv.Atoi(fec[len(fec)-1])
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}

	mes, err := strconv.Atoi(fec[1])
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}

	anio, err := strconv.Atoi(fec[0])
	if err != nil {
		erro = errors.New(ERROR_CONVERSION_DATO)
		return
	}

	fechaI = time.Date(anio, time.Month(mes), dia, 0, 0, 0, 0, time.UTC)
	fechaF = time.Date(anio, time.Month(mes), dia, 23, 59, 59, 0, time.UTC)

	return
}

func (c commons) LeerArchivo(path string) (archivo *os.File, erro error) {
	// lee el contenido del directorio que se le pasa por parammetro
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	return file, nil
}

func (c commons) EscribirArchivo(datos string, file *os.File) (erro error) {
	_, erro = file.WriteString(datos)
	if erro != nil {
		return
	}
	return
}

func (c commons) GuardarCambios(file *os.File) (erro error) {
	erro = file.Sync()
	if erro != nil {
		return
	}
	defer file.Close()
	return
}

func (c commons) ConvertirFormatoFecha(fecha string) string {
	total := 10
	resultado := fecha[0:4] + "-" + fecha[5:7] + "-" + fecha[8:total]
	return resultado
}

// formato recibido 10-10-2023 rsultado 2023-10-10
func (c commons) ConvertirFecha(fecha string) string {
	total := 10
	resultado := fecha[6:total] + "-" + fecha[3:5] + "-" + fecha[0:2]
	return resultado
}

// de tipo yyyy-mm-dd a dd-mm-yyyy
func (c commons) ConvertirFechaToDDMMYYYY(fecha string) string {
	total := 10
	anio := fecha[0:4]
	mes := fecha[5:7]
	dia := fecha[8:total]
	formatDMY := dia + "-" + mes + "-" + anio
	return formatDMY
}

func (c commons) GetFechaInicioActualMes() (FechaInicio, FechaFin string, erro error) {
	//fecha desde: el inicio del corriente mes
	t := time.Now()
	fecha_inicio := fmt.Sprintf("%d-%02d-01", t.Year(), int(t.Month()))
	//fecha fin: la fecha actual
	fecha_fin := t.Format("2006-01-02")

	// parse date
	fechaInicioTime, erro := time.Parse("2006-01-02", fecha_inicio)
	if erro != nil {
		return
	}

	// parse date
	fechaFinTime, erro := time.Parse("2006-01-02", fecha_fin)
	if erro != nil {
		return
	}

	FechaInicio = c.GetDateFirstMoment(fechaInicioTime)

	FechaFin = c.GetDateLastMoment(fechaFinTime)

	return
}

// recibe una fecha string en formato dd-mm-yyyy y la devuelve en formato yyyy-mm-dd
func ConvertirFechaYYYYMMDD(fecha string) string {
	total := 10
	anio := fecha[6:total]
	mes := fecha[3:5]
	dia := fecha[0:2]
	formatYMD := anio + "-" + mes + "-" + dia
	return formatYMD
}

var normalizer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

func normalize(str string) (string, error) {
	s, _, err := transform.String(normalizer, str)
	if err != nil {
		return "", err
	}
	return s, err
}

func (c commons) RemoveAccents(valor string) (string, error) {
	result, _ := normalize(valor)

	return result, nil
}

// devuelve una fecha string en formato ISO 8601 con la HH:mm:ss finales del dia.
// Uso: comparar limites de fechas
func (c commons) GetDateLastMoment(fecha time.Time) (fechaISO string) {
	year, month, day := fecha.Date()
	t := time.Date(year, month, day, 23, 59, 59, 999, fecha.Location())
	return t.Format(time.RFC3339)
}

// devuelve una fehca string en formato ISO 8601 con la HH:mm:ss finales del dia.
// Uso: comparar limites de fechas
func (c commons) GetDateFirstMoment(fecha time.Time) (fechaISO string) {
	year, month, day := fecha.Date()
	t := time.Date(year, month, day, 00, 00, 00, 1, fecha.Location())
	return t.Format(time.RFC3339)
}

// func (c commons) MoverArchivos(rutaOrigen, rutaDestino, nombreArchivo string) error {
// 	origen, erro := os.Open(rutaOrigen + "/" + nombreArchivo)
// 	if erro != nil {
// 		//msgOpen := ERROR_OPEN_ARCHIVO // + erro.Error()
// 		logs.Error(ERROR_OPEN_ARCHIVO + erro.Error())
// 		return errors.New(ERROR_OPEN_ARCHIVO)
// 	}
// 	defer origen.Close()
// 	destino, erro := os.Create(rutaDestino + "/" + nombreArchivo)
// 	if erro != nil {
// 		//msgCreate := ERROR_CREATE_ARCHIVO //+ erro.Error()
// 		logs.Error(ERROR_CREATE_ARCHIVO + erro.Error())
// 		return errors.New(ERROR_CREATE_ARCHIVO)
// 	}
// 	defer destino.Close()
// 	_, erro1 := io.Copy(destino, origen)
// 	if erro1 != nil {
// 		//msgCreate := ERROR_CREATE_ARCHIVO + erro1.Error()
// 		logs.Error(ERROR_CREATE_ARCHIVO + erro1.Error())
// 		return errors.New(ERROR_CREATE_ARCHIVO)
// 	}
// 	return nil
// }

func (c *commons) ToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func (c *commons) NumberSliceToString(input []uint) string {
	stringSlice := make([]string, len(input))
	for i, num := range input {
		stringSlice[i] = fmt.Sprint(num)
	}
	return strings.Join(stringSlice, ",")
}

func StringToUintSliceNumber(input string) ([]uint, error) {
	numbersStrSlice := strings.Split(input, ",")
	var numbers []uint
	for _, numStr := range numbersStrSlice {
		num, err := strconv.ParseUint(numStr, 10, 32)
		if err != nil {
			error := fmt.Errorf("error al convertir en la funcion StringToUintSliceNumber: %s", err.Error())
			return []uint{}, error
		}
		numbers = append(numbers, uint(num))
	}

	return numbers, nil
}

// recibe una fecha string con formato "yyyy-mm-dd" y retorna el ultimo momento con fechaISO de tipo 2023-09-01T23:59:59Z
func (c commons) DateYMDtoDateLastMoment(fechaIN string) (fechaISO string, erro error) {
	fecha, erro := time.Parse("2006-01-02", fechaIN)
	if erro == nil {
		fechaISO = c.GetDateLastMoment(fecha)
		return
	}
	// Si no se puede analizar con el formato "2006-01-02", intentar con el formato RFC3339
	fecha, erro = time.Parse(time.RFC3339, fechaIN)
	if erro == nil {
		fechaISO = c.GetDateLastMoment(fecha)
		return
	}
	return
}

// recibe una fecha string con formato "yyyy-mm-dd" y retorna el primer momento con fechaISO de tipo 2023-09-01T00:00:00Z
func (c commons) DateYMDtoDateFirstMoment(fechaIN string) (fechaISO string, erro error) {
	fecha, erro := time.Parse("2006-01-02", fechaIN)
	if erro == nil {
		fechaISO = c.GetDateFirstMoment(fecha)
		return
	}
	// Si no se puede analizar con el formato "2006-01-02", intentar con el formato RFC3339
	fecha, erro = time.Parse(time.RFC3339, fechaIN)
	if erro == nil {
		fechaISO = c.GetDateFirstMoment(fecha)
		return
	}

	return
}

func (c commons) DateStringToTimeFirstMoment(fechaString string) (fechaTime time.Time, erro error) {
	fechaTime, erro = time.Parse("2006-01-02", fechaString)
	if erro == nil {
		return
	}

	// Si no se puede analizar con el formato "2006-01-02", intentar con el formato completo
	fechaTime, erro = time.Parse(time.RFC3339, fechaString)
	if erro == nil {
		return
	}

	return
}

func (c commons) DateStringToTimeLastMoment(fechaString string) (fechaTime time.Time, erro error) {
	fecha, erro := time.Parse("2006-01-02", fechaString)
	if erro == nil {
		fechaTime = GetDateLastMomentTime(fecha)
		return
	}
	// Si no se puede analizar con el formato "2006-01-02", intentar con el formato completo
	fechaTime, erro = time.Parse(time.RFC3339, fechaString)
	if erro == nil {
		return
	}
	fechaTime = GetDateLastMomentTime(fecha)
	return
}

func (c commons) DateTimeToYYYYMM(fechaInicio, fechaFin time.Time) (fecha string, erro error) {

	// Obtener el año y el mes de ambas fechas
	year1, month1, _ := fechaInicio.Date()
	year2, month2, _ := fechaFin.Date()

	// Concatenar el año y el mes en el formato deseado
	result1 := fmt.Sprintf("%04d%02d", year1, int(month1))
	result2 := fmt.Sprintf("%04d%02d", year2, int(month2))

	// Comparar los resultados en formato YYYYMM
	if result1 != result2 {
		erro = fmt.Errorf("el periodo no corresponde al mismo mes")
		return
	}

	fecha = result1
	return
}
