package email

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"reflect"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/dtos/utildtos"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/commons"
	"github.com/gofiber/fiber/v2"
)

type Emailservice interface{
	//Crea el mail para poder enviar
	CreateEmailService(c *fiber.Ctx) (*bytes.Buffer,string,string,error)
	SendEmailService(body *bytes.Buffer,contentType,tempDirectoryPath string)error
}

type emailService struct{
	commons        commons.Commons
}

func NewService(c commons.Commons) Emailservice {
	email:= &emailService{
		commons: c,
	}
	return email
}


func (e *emailService) CreateEmailService(c *fiber.Ctx )(*bytes.Buffer,string,string,error){
	//Verifico los datos que recibo 
	multipartFormData, erro:=c.MultipartForm()
	//Obtengo los archivos
	if erro !=nil{
		return nil,"","", errors.New("error con el archivo recibido")
	}
	//Obtengo los correos y los archivos
	emailRecibido:= c.FormValue("Email")
	files:=multipartFormData.File["Archivos"]
	//Controlo que mande un archivo y un email
	if len(emailRecibido)==0 || len(files)==0{
		return nil,"","", errors.New("debe enviar un correo y un archivo",)
	}
	if len(files)>1{
		return nil,"","", errors.New("puede enviar solo un archivo")
	}
	//Creo la carpeta para almacenar temporalmente los archivos
	tempDirectoryPath:=	fmt.Sprintf(config.DIR_BASE + config.DOC_CL + "/emailDocs")
	if _, err := os.Stat(tempDirectoryPath); os.IsNotExist(err) {
		err = os.MkdirAll(tempDirectoryPath, 0755)
		if err != nil {
			return nil,"","", errors.New("ocurrio un error al crear la carpeta para almacenar los archivos temporales")
		}
	}
	file:=multipartFormData.File["Archivos"][0]
	//Creo un nombre unico para cada archivo
	fechaActual := time.Now()
	fechaFormato := fmt.Sprintf("%02d%02d%d%02d%02d%02d%02d",
		fechaActual.Day(), 
		fechaActual.Month(), 
		fechaActual.Year(),
		fechaActual.Hour(), 
		fechaActual.Minute(), 
		fechaActual.Second(),
		fechaActual.Nanosecond(),
	)
	fileName:=fmt.Sprintf("%s-%s",fechaFormato,file.Filename)
	filePath:=fmt.Sprintf("%s/%s-%s",tempDirectoryPath,fechaFormato,file.Filename)
	error:=c.SaveFile(file,filePath)
	if error!=nil{
		fmt.Println("Error al guardar el archivo",error)
		return nil,"","", error
	}

	//Creamos el Content-Type ya que recibimos un "application/pdf" y si no enviamos un "text/{extension del archivo}, nos genera un archivo extra sin informacion"
	extension:= path.Ext(filePath)
	if len(extension)>0 {
		extension= extension[1:]
	}
	//Qui
	//Defino mis datos a enviar
	emailDatos:=utildtos.RequestDatosMail{
		Asunto: 			"Envio de archivos",
		From: 				"Wee.ar",
		Nombre: 			"Wee.ar",
		Mensaje: 			"Los archivos enviados se encuentran adjuntados",
		CamposReemplazar: 	[]string{""},
		AdjuntarEstado: 	true,
		TipoEmail:			"adjunto",
		Attachment: utildtos.Attachment{
			Name: fileName,
			ContentType: fmt.Sprintf("text/%s",extension),
			WithFile: true,
		},
	}
	
	//Creo el writer para crear el fordata para archivos
	body:= &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//Obtengo el tipo y el valor de emailDatos
	typeEmail := reflect.TypeOf(emailDatos)
	valueEmail := reflect.ValueOf(emailDatos)
	//Recorro emailDatos y voy agregandolo segun clave-valor al form-data
	for i:=0;i<typeEmail.NumField();i++{
		//Obtengo los campos clave-valor
		key:=typeEmail.Field(i)
		value:=valueEmail.Field(i)
		//Obtengo sus valores respectivos en string
		valueName:=fmt.Sprintf("%v",value.Interface())
		keyName:=key.Name

		//El attachment al ser otra estructura, la procesamos diferente
		if keyName=="Attachment"{
			//Formateo el attachment	
			jsonAttachment, err := json.Marshal(value.Interface())
			if err != nil {
				log.Fatal(err)	
				return nil,"","", errors.New("ocurrio un error al crear los datos del email")

			}
			err =writer.WriteField(keyName,  string(jsonAttachment))
			if err != nil {
				log.Fatal(err)
				return nil,"","", errors.New("ocurrio un error al crear los datos del email")
			}
		}
		//Si es del tipo email lo proceso luego del for
		if keyName=="Email"{
			continue
		}
		//Los agrego al form-data
		err:= writer.WriteField(keyName,valueName)
		if err != nil {
			log.Fatal(err)
			return nil,"","", errors.New("ocurrio un error al crear los datos del email")
		}
	}
	//El Email lo agrego a parte por que tiene un formato especial

	err :=writer.WriteField("Email",  string(emailRecibido))
	if err != nil {
		log.Fatal(err)
		return nil,"","", errors.New("ocurrio un error al crear los datos del email")
	}
	
	//Se agrega el campo extra InformarPago que no se encuentra en la estructura, para no crear otra estructura por 1 solo campo
	err =writer.WriteField("InformarPago", "false")
	if err != nil {
		log.Fatal(err)
		return nil,"","", errors.New("ocurrio un error al crear los datos del email")
	}

	//Recorro los path de los archivos para poder adjuntarlos al formdata
		//Abro el archivo
		f, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
			return nil,"","", errors.New("ocurrio un error al adjuntar el archivo")
		}

		//Cierro el archivo luego de utilizarlo
		//Creo el form-data
		fileWriter, err := writer.CreateFormFile("Archivos", filePath)
		if err != nil {
			return nil,"","", errors.New("ocurrio un error al adjuntar el archivo")
		}
		//Copio el archivo en el form-data
		_, err = io.Copy(fileWriter, f)
		if err != nil {
			log.Fatal(err)
			return nil,"","", errors.New("ocurrio un error al adjuntar el archivo")
		}
		f.Close()

		err = writer.Close()
		contentType:=writer.FormDataContentType()
		if err != nil {
			log.Fatal(err)
			return nil,"","", errors.New("ocurrio un error al adjuntar el archivo")
		}

	return body, contentType,tempDirectoryPath,nil
}



func (e *emailService) SendEmailService(body *bytes.Buffer,contentType,tempDirectoryPath string)error{

	url:= config.SNS_EMAIL+"/emails/enviar-email"
	// Creo la peticion HTTP con el form-data
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("ApiKey", config.SNS_EMAIL_APIKEY)
	// Envio la peticion http
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("ocurrio un error al enviar el correo")
	}
	defer resp.Body.Close()
	erro:= os.RemoveAll(tempDirectoryPath)
	if erro!=nil {
		fmt.Println("Error al borrar los archivos temporales")
	}
	if resp.StatusCode!=200{
		//Leo la respuesta que recibo
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error al leer el body")
			return errors.New("error en servicio de email")
		}
		//Decodifico para poder enviar el mensaje
		var response map[string]string
		err= json.Unmarshal([]byte(string(responseBody)), &response)
		if err != nil {
			fmt.Println("Error decodificar el mensaje del emails service", err)
			return errors.New("error en servicio de email")
		}
		//Si existe el mensaje lo envio
		message, ok := response["message"]
		if !ok {
			fmt.Println("No existe mensaje a devolver")
			return errors.New("error en servicio de email")
		}
	
		return errors.New(message)
	}
	
	return nil
}