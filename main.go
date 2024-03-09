package main

import (
	"Proyecto1/analizador"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	//comandos.LeerMounts()
	//for {
	//	leerComando()
	//}

	//id := "A" + string(3+48) + "44"
	//fmt.Println(id)

	ruta := "/home/home2/home3/archivo.txt"

	lRuta := strings.Split(ruta[1:], "/")
	lRuta = lRuta[:len(lRuta)-1]
	fmt.Println(ruta[1:])
	fmt.Println(lRuta)
	fmt.Println(len(lRuta))

}

func leerComando() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Ingrese un comando: ")
	comando, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al ingresar el comando: ", err)
		return
	}

	comando = strings.TrimSpace(comando)

	analizador.Analizar(comando)
}
