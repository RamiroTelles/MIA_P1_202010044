package main

import (
	"Proyecto1/analizador"
	"Proyecto1/comandos"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	comandos.LeerMounts()

	for {
		leerComando()
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
