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

	//fmt.Println(comandoEntero)

	analComando := regexp.MustCompile("^[A-Za-z]+")
	comando := analComando.FindAllString(comandoEntero, 1)
	analBanderas := regexp.MustCompile("(-[A-Za-z0-9]*=([A-Za-z0-9./_]*))")
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

	case "mount":
		comandos.EjecMount(banderas)
		break
	case "unmount":
		comandos.EjecUnMount(banderas)
		break
	case "lmount":
		comandos.EjecLMount()
		break

	case "rep":
		//fmt.Println("si llega")
		EjecRep(banderas)
		break
	case "fdisk":
		comandos.EjecFdisk(banderas)
		break
	case "rmdisk":
		comandos.EjecRmdisk(banderas)
		break
	case "mkfs":
		comandos.EjecMkfs(banderas)
		break
	case "cat":
		comandos.EjecCat(banderas)
		break
	case "logout":
		comandos.EjecLogout()
		break
	case "login":
		comandos.EjecLogin(banderas)
		break
	case "pause":
		comandos.EjecPause()
	case "exit":
		fmt.Println("cerrando aplicacion")
		os.Exit(0)

	}

}

func EjecRep(banderas []string) {

	name := ""
	path := ""
	id := ""
	//ruta := ""

	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-name" {

			name = dupla[1]

		} else if dupla[0] == "-path" {
			path = dupla[1]

		} else if dupla[0] == "-id" {
			id = dupla[1]
		} else if dupla[0] == "-ruta" {
			//ruta = dupla[1]
		} else {
			fmt.Println("Parametro invalido")
		}
	}
	fmt.Println(name)
	fmt.Println(id)

	switch name {
	case "mbr":
		//reporte mbr
		break
	case "disk":
		//reporte disk
		comandos.EjecRepMkdisk(id, path)
		break

	case "inodo":
		//reporte inodo
		index := comandos.VerificarParticionMontada(id)
		if index == -1 {
			fmt.Println("Id no encontrada")
			return
		}
		comandos.EjecRepInodes(index)
		break
	case "journaling":
		//reporte journaling
		break
	case "block":
		//reporte block
		index := comandos.VerificarParticionMontada(id)
		if index == -1 {
			fmt.Println("Id no encontrada")
			return
		}
		comandos.EjecRepBloques(index)
		break
	case "bm_inode":
		//reporte bitmap inodo
		index := comandos.VerificarParticionMontada(id)
		if index == -1 {
			fmt.Println("Id no encontrada")
			return
		}
		comandos.EjecRepBmInodes(index)
		break
	case "bm_block":
		//reporte bitmap block
		index := comandos.VerificarParticionMontada(id)
		if index == -1 {
			fmt.Println("Id no encontrada")
			return
		}
		comandos.EjecRepBmBloques(index)
		break
	case "tree":
		//reporte tree

		comandos.EjecRepTree(id)
		break
	case "sb":
		//reporte sb
		index := comandos.VerificarParticionMontada(id)
		if index == -1 {
			fmt.Println("Id no encontrada")
			return
		}
		comandos.EjecRepSB(index)
		break
	case "file":
		//reporte file
		break
	case "ls":
		//reporte ls
		break

	default:
		fmt.Println("nombre no valido")
		return
	}

}

func EjecExecute(banderas []string) {
	dupla := strings.Split(banderas[0], "=")
	if dupla[0] == "-path" {
		//fmt.Println(dupla[1])
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
