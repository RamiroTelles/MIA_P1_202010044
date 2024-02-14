package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type partition struct {
	part_status      [1]byte
	part_type        [1]byte
	part_fit         [1]byte
	part_start       int32
	part_s           int32
	part_name        [16]byte
	part_correlative int32
	part_id          [4]byte
}

type MBR struct {
	mbr_tamano         int32
	mbr_fecha_creacion [10]byte
	mbr_dsk_signature  [1]byte
	mbr_partitions     [4]partition
}

func main() {

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

	analizar(comando)
}

func analizar(comando string) {

}
