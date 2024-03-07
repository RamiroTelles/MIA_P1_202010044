package main

import (
	"Proyecto1/analizador"
	"bufio"
	"encoding/binary"
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

	archivo, err := os.OpenFile("pruebas.txt", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()
	archivo.Seek(1, 0)
	for i := int32(0); i < 17; i++ {

		err := binary.Write(archivo, binary.LittleEndian, [1]byte{0})
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

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
