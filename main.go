package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
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
	Mbr_tamano         int32
	Mbr_fecha_creacion [19]byte
	Mbr_dsk_signature  int32
	//mbr_partitions     [4]partition
}

func Analizar(comandoEntero string) {

	fmt.Println(comandoEntero)

	analComando := regexp.MustCompile("^[A-Za-z]+")
	comando := analComando.FindAllString(comandoEntero, 1)
	analBanderas := regexp.MustCompile("(-[A-Za-z]*=(.*))")
	banderas := analBanderas.FindAllString(comandoEntero, -1)

	//fmt.Println(comando)
	//fmt.Println(banderas)
	ejecutarComando(comando, banderas)

}

func ejecutarComando(comando []string, banderas []string) {

	switch comando[0] {

	case "execute":
		//ejecutar execute
		ejecExecute(banderas)
		break

	case "mkdisk":
		ejecMkdisk(banderas)
		break

	case "rep":
		//fmt.Println("si llega")
		ejecRep()
		break

	}

}

func ejecExecute(banderas []string) {
	dupla := strings.Split(banderas[0], "=")
	if dupla[0] == "-path" {
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

func ejecMkdisk(banderas []string) {

	// if _, err := os.Stat("Discos"); os.IsNotExist(err) {
	// 	err = os.Mkdir("Discos", 0664)
	// 	if err != nil {
	// 		fmt.Println("Error al crear el directorio Discos: ", err)
	// 		return
	// 	}
	// }

	nombreDisco := "A"

	archivo, err := os.Create(nombreDisco + ".dsk")

	if err != nil {
		fmt.Println("Error al crear el archivo del disco: ", err)
		return
	}
	defer archivo.Close()

	var mbrDisk MBR

	randomNum := rand.Intn(99) + 1
	mbrDisk.Mbr_tamano = 5 * 1024 * 1024
	mbrDisk.Mbr_dsk_signature = int32(randomNum)
	fechaActual := time.Now()
	fecha := fechaActual.Format("2006-01-02 15:04:05")
	copy(mbrDisk.Mbr_fecha_creacion[:], fecha)

	bufer := new(bytes.Buffer)
	for i := 0; i < 1024; i++ {
		bufer.WriteByte(0)
	}

	var totalBytes int = 0
	for totalBytes < int(mbrDisk.Mbr_tamano) {
		c, err := archivo.Write(bufer.Bytes())
		if err != nil {
			fmt.Println("Error al escribir en el archivo: ", err)
			return
		}
		totalBytes += c
	}
	//fmt.Println("Archivo llenado con 0s")

	archivo.Seek(0, 0)
	err = binary.Write(archivo, binary.LittleEndian, &mbrDisk)
	if err != nil {
		fmt.Println("Error al escribir el MBR en el disco: ", err)
		return
	}
	fmt.Println("Disco", nombreDisco, "creado con exito")
}

func ejecRep() {
	archivo, err := os.Open("A.dsk")

	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()
	var disk MBR
	disk.Mbr_dsk_signature = int32(0)
	disk.Mbr_fecha_creacion = [19]byte{}
	disk.Mbr_dsk_signature = int32(0)
	archivo.Seek(int64(0), 0)
	fmt.Println("sss")
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error al leer el MBR del disco: ", err)
		return
	}
	archivo.Close()
	fmt.Println("Tamaño: ", disk.Mbr_tamano)
	fmt.Println("Fecha: ", string(disk.Mbr_fecha_creacion[:]))
	fmt.Println("Signature: ", disk.Mbr_dsk_signature)
	//fmt.Println("Fit: ", string(disk.Dsk_fit[:]))
	//fmt.Println("Partition1: ", string(disk.Mbr_partition1.Part_status[:]))
	//fmt.Println("Partition2: ", string(disk.Mbr_partition2.Part_status[:]))
	//fmt.Println("Partition3: ", string(disk.Mbr_partition3.Part_status[:]))
	//fmt.Println("Partition4: ", string(disk.Mbr_partition4.Part_status[:]))
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

	Analizar(comando)
}
