package analizador

import (
	"Proyecto1/comandos"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func Analizar(comandoEntero string) {

	fmt.Println(comandoEntero)

	analComando := regexp.MustCompile("^[A-Za-z]+")
	comando := analComando.FindAllString(comandoEntero, 1)
	analBanderas := regexp.MustCompile("(-[A-Za-z]*=([A-Za-z0-9./]*))")
	banderas := analBanderas.FindAllString(comandoEntero, -1)

	//fmt.Println(comando)
	//fmt.Println(banderas)
	if comando != nil {
		ejecutarComando(comando, banderas)
	}

}

func ejecutarComando(comando []string, banderas []string) {

	switch comando[0] {

	case "execute":
		//ejecutar execute
		EjecExecute(banderas)
		break

	case "mkdisk":
		comandos.EjecMkdisk(banderas)
		break

	case "rep":
		//fmt.Println("si llega")
		comandos.EjecRepMkdisk()
		break
	case "fdisk":
		comandos.EjecFdisk(banderas)
		break

	case "exit":
		fmt.Println("cerrando aplicacion")
		os.Exit(0)

	}

}

func EjecExecute(banderas []string) {
	dupla := strings.Split(banderas[0], "=")
	if dupla[0] == "-path" {
		fmt.Println(dupla[1])
		archivo, err := os.Open(dupla[1])

		if err != nil {
			fmt.Println("Error al abrir el archivo: ", err)
			return
		}
		defer archivo.Close()

		scanner := bufio.NewScanner(archivo)

		for scanner.Scan() {
			linea := scanner.Text()
			fmt.Println(linea)

			Analizar(linea)
		}
	}
}
