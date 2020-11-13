package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

//Creamos clase satelite
type Satellite struct {
	Nrosatellite    int     `json:"nrosatellite"`
	X               float32 `json:"x"`
	Y               float32 `json:"y"`
	Nombre          string  `json:"nombre"`
	Distanciaemisor float32 `json:"distanciaemisor"`
}

func actualizarDistancia(paramnombre string, paramdistancia float32) {

	//cambiamos el parametro distancia de cada satellite
	var n string = satellites[0].Nombre
	fmt.Println(n)
	for idx := 0; idx < len(satellites); idx++ {
		if strings.ToUpper(satellites[idx].Nombre) == strings.ToUpper(paramnombre) {

			satellites[idx].Distanciaemisor = paramdistancia
		}
		idx += idx
	}
}

//CREAMOS LOS OBJETOS SATELLITES
type allsatellites []Satellite

var satellites = allsatellites{
	{
		Nrosatellite:    1,
		X:               -500,
		Y:               -200,
		Nombre:          "Kenobi",
		Distanciaemisor: 0,
	},
	{
		Nrosatellite:    2,
		X:               100,
		Y:               -100,
		Nombre:          "SkyWalker",
		Distanciaemisor: 0,
	},
	{
		Nrosatellite:    3,
		X:               500,
		Y:               100,
		Nombre:          "Sato",
		Distanciaemisor: 0,
	},
}

////////////////////////////////////////// clase nave
type Nave struct {
	Nave   int
	X      float32
	Y      float32
	Nombre string
}

//OBTENER COORDENADA DE NAVE
func GetLocation(distances ...float32) (x, y float32) {
	x = 2
	y = 2
	return x, y
}

//DESCRIFRAR MENSAJE DE NAVE----------------------------------------------------------------
func GetMessage(messages ...[]string) (msg string) {
	//averiguar por si los strings vienen con distinta cantidad de palabras
	//var max int
	var frase string
	var referencia []string = messages[0]
	var izq string
	var der string

	//llenamos los ""
	for i := 0; i < len(referencia); i++ {
		izq = ""
		der = ""

		//tomo el primer valor que no sea ""
		if referencia[i] != "" {
			//ponemos el valor izquierdo del indice i
			if i-1 >= 0 && referencia[i-1] != "" {

				izq = referencia[i-1]
			}
			//ponemos el valor izq al indice i
			if i+1 < len(referencia) && referencia[i+1] != "" {
				der = referencia[i+1]
			}

			//si encuentro otro valor en el slice que esta en blanco a la izq o derecha lo relleno
			if der != "" || izq != "" {
				for j := 0; j < len(referencia); j++ {
					if referencia[i] == referencia[j] {
						if j+1 < len(referencia) {
							if referencia[j+1] == "" && der != "" {
								referencia[j+1] = der
							}
						}

						if j-1 >= 0 && referencia[j-1] == "" && izq != "" {
							referencia[j-1] = izq
						}
					}
				}
			}

		}

	}

	type Valorcantidad struct {
		Palabra  string
		Cantidad int
	}

	type Valorescantidades []Valorcantidad

	var palabras Valorcantidad
	var v Valorescantidades

	//recorremos la cadena completa para contar las cantidades
	for f := 0; f < len(referencia); f++ {
		palabras.Palabra = ""
		palabras.Cantidad = 0
		if referencia[f] != "" {
			palabras.Palabra = referencia[f]
			palabras.Cantidad = palabras.Cantidad + 1
			for g := f + 1; g < len(referencia); g++ {
				if referencia[f] == referencia[g] {
					palabras.Cantidad = palabras.Cantidad + 1

				}
			}
			palabras.Cantidad = palabras.Cantidad / 3
			v = append(v, palabras)
		}

	}

	//ARMAMOS FRASE FINAL
	var i = 0
	for j := i; j < len(v); j++ {
		if v[0].Cantidad == 0 && v[0].Palabra == v[j].Palabra {
			break
		}
		frase = frase + v[j].Palabra + " "

		v[j].Cantidad = v[j].Cantidad - 1
	}
	msg = frase
	return msg
}

// /////////////////////////////////////////clase mensaje

type Mensaje struct {
	Distance float32  `json:"distance"`
	Message  []string `json:"message"`
	Name     string   `json:"name"`
}

type Mensajes struct {
	Satellites []Mensaje `json:"satellites"`
}

//----------------------------------------------------------
func getLocation(w http.ResponseWriter, r *http.Request) {

	//RECIBIR JSON
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "error en datos enviados")
		return
	}

	var mensajes Mensajes
	json.Unmarshal(reqBody, &mensajes)
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusCreated)
	if err != nil {
		log.Fatal(err)
	}

	//EN BASE A LOS MENSAJES QUE LLEGAN ACTUALIZO LA DISTANCIA AL SATELLITE EN CADA SATELLITE
	//TAMBIEN ENVIO A DESCRIFRAR EL
	a := []float32{}

	var mens []string

	//var coordX, coordY float32

	for i := 0; i < len(mensajes.Satellites); i++ {
		//persistir datos en atributo de clase Satellite
		go actualizarDistancia(mensajes.Satellites[i].Name, mensajes.Satellites[i].Distance)
		//a contiene las distancias para pasar a GetLocation
		a = append(a, mensajes.Satellites[i].Distance)

		//dejo los mensajes de menor longitud como el de mas longitud
		var maxlong int = len(mensajes.Satellites[0].Message)
		for i := 0; i < len(mensajes.Satellites); i++ {
			if maxlong < len(mensajes.Satellites[i].Message) {
				maxlong = len(mensajes.Satellites[i].Message)

			}
			if len(mensajes.Satellites[i].Message) < maxlong {
				for j := len(mensajes.Satellites[i].Message); j < maxlong; j++ {
					mensajes.Satellites[i].Message = append(mensajes.Satellites[i].Message, "")
					var frase1 []string
					frase1 = append(frase1, "")
					for k := 0; k < len(mensajes.Satellites[i].Message); k++ {
						frase1 = append(frase1, mensajes.Satellites[i].Message[k])
					}
					mensajes.Satellites[i].Message = frase1
				}
			}
		}

		//mens para obtner []string para enviar a GetMensajes
		for j := 0; j < len(mensajes.Satellites[i].Message); j++ {
			mens = append(mens, mensajes.Satellites[i].Message[j])
		}
	}

	//coordX, coordY = GetLocation(a...)

	//DESCIFRAR MENSAJE DE NAVE
	//devolver a w mensajedefinitivo
	var msgdef string = GetMessage(mens)
	json.NewEncoder(w).Encode(msgdef)
	return
}

func topsecret_split(w http.ResponseWriter, req *http.Request) {

}

func create_topsecret_split(w http.ResponseWriter, req *http.Request) {

}

/////////////////////////////////////////main
func main() {

	router := mux.NewRouter().StrictSlash(true)

	//endpoints / rutas
	router.HandleFunc("/topsecret", getLocation).Methods("POST")
	//router.HandleFunc("/topsecret_split/{satellite_name}", topsecret_split).Methods("GET")
	//router.HandleFunc("/topsecret_split/{satellite_name}", create_topsecret_split).Methods("POST")

	//log
	log.Fatal(http.ListenAndServe(":8080", router))

}
