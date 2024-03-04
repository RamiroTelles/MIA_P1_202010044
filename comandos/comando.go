package comandos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type partition struct {
	Part_status      [1]byte
	Part_type        [1]byte
	Part_fit         [1]byte
	Part_start       int32
	Part_s           int32
	Part_name        [16]byte
	Part_correlative int32
	Part_id          [4]byte
}

type MBR struct {
	Mbr_tamano         int32
	Mbr_fecha_creacion [19]byte
	Mbr_dsk_signature  int32
	MBR_dsk_fit        [1]byte
	Mbr_partitions     [4]partition
}

type EBR struct {
	Part_mount [1]byte
	Part_fit   [1]byte
	Part_start int32
	Part_s     int32
	Part_next  int32
	Part_name  [16]byte
}

func LeerMounts() {

	archivos, err := os.ReadDir("MIA/P1")
	if err != nil {
		fmt.Println("Error al leer el directorio: ", err)
		return
	}
	//Declarar las letras del abecedario
	//letras := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//Nombre del disco a partir de la cantidad de discos, por ejemplo A=1, B=2, C=3
	//nombreDisco := string(letras[len(archivos)])

	for _, ruta := range archivos {
		archivo, err := os.OpenFile("MIA/P1/"+ruta.Name(), os.O_RDWR, 0777)
		fmt.Println(ruta.Name()[:1])
		if err != nil {
			fmt.Println("Error al abrir el disco: ", err)
			return
		}
		defer archivo.Close()

		var disk MBR
		archivo.Seek(int64(0), 0)
		err = binary.Read(archivo, binary.LittleEndian, &disk)
		if err != nil {
			fmt.Println("Error al leer el MBR del disco: ", err)
			return
		}
		for i, part := range disk.Mbr_partitions {
			if part.Part_status == [1]byte{'1'} {
				var newMount Mount
				newMount.Id = string(part.Part_id[:])

				newMount.LetterValor = ruta.Name()[:1]
				newMount.Name = string(part.Part_name[:])
				newMount.partNum = int32(i)
				newMount.Part_type = part.Part_type
				newMount.Size = part.Part_s
				newMount.Start = part.Part_start

				particionesMontadas = append(particionesMontadas, newMount)

			}

		}

	}

}

func EjecFdisk(banderas []string) {
	unit := "k"
	size := -1
	driveLetter := ""
	name := ""
	type1 := "p"
	fit := "w"
	//delete := 0

	despTemp := binary.Size(MBR{}) + 1

	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-size" {

			size, _ = strconv.Atoi(dupla[1])

		} else if dupla[0] == "-unit" {
			unit = dupla[1]
		} else if dupla[0] == "-driveletter" {
			driveLetter = dupla[1]
		} else if dupla[0] == "-name" {
			name = dupla[1]

		} else if dupla[0] == "-type" {
			type1 = dupla[1]
		} else if dupla[0] == "-fit" {
			if dupla[1] == "bf" {
				fit = "b"
			} else if dupla[1] == "ff" {
				fit = "f"
			} else if dupla[1] == "wf" {
				fit = "w"
			} else {
				fmt.Println("Valor fit no valido")
				return
			}
		} else if dupla[0] == "-delete" {
			if dupla[1] == "full" {
				//delete = 1
			} else {
				fmt.Println("parametro invalido para delete")
			}
		} else {
			fmt.Println("Parametro invalido")
		}
	}

	if name == "" {
		fmt.Println("nombre invalido")
		return
	}

	if size < 0 {
		fmt.Println("tamano no valido")
		return
	}

	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	}

	//fmt.Println(unit)
	//fmt.Println(size)
	//fmt.Println(driveLetter)
	//fmt.Println(name)
	//fmt.Println(type1)

	archivo, err := os.OpenFile("MIA/P1/"+driveLetter+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var disk MBR
	archivo.Seek(int64(0), 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error al leer el MBR del disco: ", err)
		return
	}

	numPart := -1

	for i, particion := range disk.Mbr_partitions {

		if particion.Part_s == int32(0) {
			numPart = i
			break

		} else {
			despTemp += int(particion.Part_s) + 1
		}

	}

	var partExtend partition
	for _, part := range disk.Mbr_partitions {
		if part.Part_type == [1]byte{'e'} {
			if type1 == "e" {
				fmt.Println("Ya existe una particion extendida")
				return
			}
			partExtend = part

		}
	}

	if type1 == "l" {
		if partExtend.Part_type != [1]byte{'e'} {
			fmt.Println("No existe particion extendida")
			return
		}
		var ebr EBR
		despTemp = int(partExtend.Part_start)

		for {

			archivo.Seek(int64(despTemp), 0)
			binary.Read(archivo, binary.LittleEndian, &ebr)
			if ebr.Part_s != 0 {
				if strings.Contains(string(ebr.Part_name[:]), name) {
					fmt.Println("Error: El nombre de la particion ya existe")
					return
				}
				despTemp += int(ebr.Part_s) + 1 + binary.Size(EBR{})

			} else {
				break

			}

		}

		if int32(despTemp)+int32(binary.Size(EBR{}))+int32(size)+1 > partExtend.Part_start+partExtend.Part_s {
			fmt.Println("Error: No hay espacio para crear la particion")
			return
		}
		//Crear el nuevo EBR
		var nuevoEBR EBR
		nuevoEBR.Part_mount = [1]byte{'0'}
		nuevoEBR.Part_fit = [1]byte{fit[0]}
		nuevoEBR.Part_start = int32(despTemp) + 1 + int32(binary.Size(EBR{}))
		nuevoEBR.Part_s = int32(size)
		nuevoEBR.Part_next = int32(-1)

		copy(nuevoEBR.Part_name[:], name)
		//Escribir el nuevo EBR
		archivo.Seek(int64(despTemp), 0)
		binary.Write(archivo, binary.LittleEndian, &nuevoEBR)
		archivo.Close()
		fmt.Println("Particion logica creada con exito")
		return

	} else {
		var nuevaPar partition

		nuevaPar.Part_status = [1]byte{'0'}

		if type1 == "p" || type1 == "e" {

			nuevaPar.Part_type = [1]byte{type1[0]}
		} else {
			fmt.Println("Tipo de particion no valida")
		}

		nuevaPar.Part_fit = [1]byte{fit[0]}

		if numPart < 0 {
			fmt.Println("No hay particiones disponibles")
			return
		}

		nuevaPar.Part_start = int32(despTemp)

		nuevaPar.Part_s = int32(size)

		nuevaPar.Part_name = [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		copy(nuevaPar.Part_name[:], name)

		if despTemp+int(nuevaPar.Part_s)+1 > int(disk.Mbr_tamano) {
			fmt.Println("tamano insuficiente para la particion")
			return
		}
		nuevaPar.Part_correlative = 0

		disk.Mbr_partitions[numPart] = nuevaPar

		archivo.Seek(0, 0)
		binary.Write(archivo, binary.LittleEndian, &disk)
		archivo.Close()

		fmt.Println("particion creada con exito")

	}

}

func EjecMkdisk(banderas []string) {

	// if _, err := os.Stat("Discos"); os.IsNotExist(err) {
	// 	err = os.Mkdir("Discos", 0664)
	// 	if err != nil {
	// 		fmt.Println("Error al crear el directorio Discos: ", err)
	// 		return
	// 	}
	// }
	unit := "m"
	size := -1

	fit := "f"

	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-size" {

			size, _ = strconv.Atoi(dupla[1])

		} else if dupla[0] == "-unit" {
			unit = dupla[1]

		} else if dupla[0] == "-fit" {
			if dupla[1] == "bf" {
				fit = "b"
			} else if dupla[1] == "ff" {
				fit = "f"
			} else if dupla[1] == "wf" {
				fit = "w"
			}
		} else {
			fmt.Println("Parametro invalido")
		}
	}

	archivos, err := os.ReadDir("MIA/P1")
	if err != nil {
		fmt.Println("Error al leer el directorio: ", err)
		return
	}
	//Declarar las letras del abecedario
	letras := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//Nombre del disco a partir de la cantidad de discos, por ejemplo A=1, B=2, C=3
	nombreDisco := string(letras[len(archivos)])

	archivo, err := os.Create("MIA/P1/" + nombreDisco + ".dsk")

	if err != nil {
		fmt.Println("Error al crear el archivo del disco: ", err)
		return
	}
	defer archivo.Close()

	var mbrDisk MBR

	randomNum := rand.Intn(99) + 1

	if size < 0 {

		fmt.Println("Valor invalido para el parametro -size")
		return

	}

	if unit == "k" {
		mbrDisk.Mbr_tamano = int32(size) * 1024
	} else if unit == "m" {
		mbrDisk.Mbr_tamano = int32(size) * 1024 * 1024
	} else {
		fmt.Println("El valor del parametro -unit no es valido")
		return
	}

	mbrDisk.Mbr_dsk_signature = int32(randomNum)
	fechaActual := time.Now()
	fecha := fechaActual.Format("2006-01-02 15:04:05")
	copy(mbrDisk.Mbr_fecha_creacion[:], fecha)

	fitBytes := []byte(fit)

	mbrDisk.MBR_dsk_fit = [1]byte(fitBytes)

	mbrDisk.Mbr_partitions[0].Part_status = [1]byte{'0'}
	mbrDisk.Mbr_partitions[1].Part_status = [1]byte{'0'}
	mbrDisk.Mbr_partitions[2].Part_status = [1]byte{'0'}
	mbrDisk.Mbr_partitions[3].Part_status = [1]byte{'0'}

	mbrDisk.Mbr_partitions[0].Part_type = [1]byte{'0'}
	mbrDisk.Mbr_partitions[1].Part_type = [1]byte{'0'}
	mbrDisk.Mbr_partitions[2].Part_type = [1]byte{'0'}
	mbrDisk.Mbr_partitions[3].Part_type = [1]byte{'0'}

	mbrDisk.Mbr_partitions[0].Part_fit = [1]byte{'0'}
	mbrDisk.Mbr_partitions[1].Part_fit = [1]byte{'0'}
	mbrDisk.Mbr_partitions[2].Part_fit = [1]byte{'0'}
	mbrDisk.Mbr_partitions[3].Part_fit = [1]byte{'0'}

	mbrDisk.Mbr_partitions[0].Part_start = 0
	mbrDisk.Mbr_partitions[1].Part_start = 0
	mbrDisk.Mbr_partitions[2].Part_start = 0
	mbrDisk.Mbr_partitions[3].Part_start = 0

	mbrDisk.Mbr_partitions[0].Part_s = 0
	mbrDisk.Mbr_partitions[1].Part_s = 0
	mbrDisk.Mbr_partitions[2].Part_s = 0
	mbrDisk.Mbr_partitions[3].Part_s = 0

	mbrDisk.Mbr_partitions[0].Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	mbrDisk.Mbr_partitions[1].Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	mbrDisk.Mbr_partitions[2].Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	mbrDisk.Mbr_partitions[3].Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}

	mbrDisk.Mbr_partitions[0].Part_correlative = 0
	mbrDisk.Mbr_partitions[1].Part_correlative = 0
	mbrDisk.Mbr_partitions[2].Part_correlative = 0
	mbrDisk.Mbr_partitions[3].Part_correlative = 0

	mbrDisk.Mbr_partitions[0].Part_id = [4]byte{'0', '0', '0', '0'}
	mbrDisk.Mbr_partitions[1].Part_id = [4]byte{'0', '0', '0', '0'}
	mbrDisk.Mbr_partitions[2].Part_id = [4]byte{'0', '0', '0', '0'}
	mbrDisk.Mbr_partitions[3].Part_id = [4]byte{'0', '0', '0', '0'}

	bufer := new(bytes.Buffer)
	for i := 0; i < 1024; i++ {
		bufer.WriteByte(0)
	}

	var totalBytes int = 0
	fmt.Println(mbrDisk.Mbr_tamano)
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

func EjecRmdisk(banderas []string) {
	driveletter := ""

	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-driveletter" {

			driveletter = dupla[1]

		} else {
			fmt.Println("Parametro invalido")
		}
	}

	err := os.Remove("MIA/P1/" + driveletter + ".dsk")

	if err != nil {
		fmt.Println("Error al crear el archivo: ")
		fmt.Println(err)
		return
	}
	fmt.Println("Disco eliminado con exito")

}

func EjecMount(banderas []string) {

	driveletter := ""
	name := ""

	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-driveletter" {

			driveletter = dupla[1]

		} else if dupla[0] == "-name" {
			name = dupla[1]

		} else {
			fmt.Println("Parametro invalido")
		}
	}

	archivo, err := os.OpenFile("MIA/P1/"+driveletter+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var disk MBR
	archivo.Seek(int64(0), 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error al leer el MBR del disco: ", err)
		return
	}
	numPart := -1
	for i, part := range disk.Mbr_partitions {

		if strings.Contains(string(part.Part_name[:]), name) {
			if string(part.Part_type[:]) == "p" && part.Part_status == [1]byte{'0'} {
				numPart = i
				break
			} else {
				fmt.Println("Error al montar la particion")
				return
			}

		}

	}

	if numPart == -1 {
		fmt.Println("particion no encontrada")
		return
	}

	contador := 1
	for i := 0; i < len(particionesMontadas); i++ {
		if particionesMontadas[i].LetterValor == driveletter {
			contador++
		}
	}

	id := driveletter + string(contador+48) + "44"

	for _, mounts := range particionesMontadas {
		if strings.Contains(id, mounts.Id) {
			fmt.Println("Particion ya montada")
		}
	}

	idB := []byte(id)

	disk.Mbr_partitions[numPart].Part_status = [1]byte{'1'}
	disk.Mbr_partitions[numPart].Part_id = [4]byte(idB)
	disk.Mbr_partitions[numPart].Part_correlative = int32(contador)
	var newMount Mount
	newMount.LetterValor = driveletter
	newMount.Id = id
	newMount.Name = name
	newMount.Part_type = disk.Mbr_partitions[numPart].Part_type
	newMount.Start = disk.Mbr_partitions[numPart].Part_start
	newMount.Size = disk.Mbr_partitions[numPart].Part_s
	newMount.partNum = int32(numPart)

	particionesMontadas = append(particionesMontadas, newMount)

	archivo.Seek(0, 0)
	binary.Write(archivo, binary.LittleEndian, &disk)
	archivo.Close()

	fmt.Println("Particion montada con exito")

}

func EjecUnMount(banderas []string) {
	id := ""

	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-id" {
			id = dupla[1]

		} else {
			fmt.Println("Parametro invalido")
		}
	}

	index := VerificarParticionMontada(id)

	if index == -1 {
		fmt.Println("Id de la particion no encontrada")
		return
	}

	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var disk MBR
	archivo.Seek(int64(0), 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error al leer el MBR del disco: ", err)
		return
	}

	disk.Mbr_partitions[particionesMontadas[index].partNum].Part_status = [1]byte{'0'}
	disk.Mbr_partitions[particionesMontadas[index].partNum].Part_id = [4]byte{0, 0, 0, 0}
	disk.Mbr_partitions[particionesMontadas[index].partNum].Part_correlative = 0

	particionesMontadas = append(particionesMontadas[:index], particionesMontadas[index+1:]...)

	archivo.Seek(0, 0)
	binary.Write(archivo, binary.LittleEndian, &disk)
	archivo.Close()

	fmt.Println("Particion desmontada con exito")
}

func EjecLMount() {

	if len(particionesMontadas) == 0 {
		fmt.Println("no hay particiones montadas")
	}

	for _, mounts := range particionesMontadas {
		fmt.Print("Id: ")
		fmt.Println(mounts.Id)
		fmt.Print("Disco: ")
		fmt.Println(mounts.LetterValor)
		fmt.Print("Nombre Particion: ")
		fmt.Println(mounts.Name)
		fmt.Print("Tipo: ")
		fmt.Println(string(mounts.Part_type[:]))
		fmt.Print("Inicio: ")
		fmt.Println(mounts.Start)
		fmt.Print("Tamano: ")
		fmt.Println(mounts.Size)
	}
}

func EjecRepMBR() {
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

func EjecRepMkdisk(id string, path string) {
	archivo, err := os.Open("MIA/P1/" + string(id[0]) + ".dsk")

	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var disk MBR
	archivo.Seek(int64(0), 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error al leer el MBR del disco: ", err)
		return
	}
	sizeMBR := int(disk.Mbr_tamano)
	libre := int(disk.Mbr_tamano)

	Dot := "digraph grid {bgcolor=\"slategrey\" label=\" Reporte Disk \"layout=dot "
	Dot += "labelloc = \"t\"edge [weigth=1000 style=dashed color=red4 dir = \"both\" arrowtail=\"open\" arrowhead=\"open\"]"
	Dot += "node[shape=record, color=lightgrey]a0[label=\"MBR"

	for _, part := range disk.Mbr_partitions {
		if part.Part_s != 0 {
			libre -= int(part.Part_s)
			Dot += "|"
			if part.Part_type == [1]byte{'e'} {
				Dot += "{Extendida"
				libreExtendida := int(part.Part_s)

				var ebr EBR
				desp := int(part.Part_start)
				archivo.Seek(int64(desp), 0)
				err := binary.Read(archivo, binary.LittleEndian, &ebr)
				if err != nil {
					fmt.Println("Error al leer el EBR: ", err)
					return
				}

				if ebr.Part_s != 0 {
					Dot += "|{EBR"

					Dot += "|Logica"
					porcentaje := (float64(ebr.Part_s) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					//libre -= int(ebr.Part_s)

					desp += int(ebr.Part_s) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(desp), 0)
					binary.Read(archivo, binary.LittleEndian, &ebr)
					for {

						if ebr.Part_s == 0 {
							break
						}
						Dot += "|EBR"
						Dot += "|Logica"
						porcentaje := (float64(ebr.Part_s) * float64(100)) / float64(sizeMBR)
						Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
						libre -= int(ebr.Part_s)

						desp += int(ebr.Part_s) + 1 + binary.Size(EBR{})
						archivo.Seek(int64(desp), 0)
						binary.Read(archivo, binary.LittleEndian, &ebr)

					}
					if libreExtendida > 0 {
						Dot += "|Libre"
						porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
						Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					}
					Dot += "}}"

				} else {
					Dot += "|Libre"
					porcentaje := (float64(part.Part_s) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					Dot += "}"
				}
				//libre -= int(ebr.Part_s)

			} else {
				Dot += "Primaria"
				porcentaje := (float64(part.Part_s) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"

			}

		}
	}

	if libre > 0 {
		Dot += "|Libre"
		porcentaje := (float64(libre) * float64(100)) / float64(sizeMBR)
		Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
	}
	Dot += "\"];\n}"

	//Crear el archivo .dot
	DotName := "ReporteDisk.dot"
	archivoDot, err := os.Create(DotName)
	if err != nil {
		fmt.Println("Error al crear el archivo .dot: ", err)
		return
	}
	defer archivoDot.Close()
	_, err = archivoDot.WriteString(Dot)
	if err != nil {
		fmt.Println("Error al escribir el archivo .dot: ", err)
		return
	}
	//Generar la imagen
	cmd := exec.Command("dot", "-T", "png", DotName, "-o", "ReporteDisk.png")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar la imagen: ", err)
		return
	}

	fmt.Println("Reporte generado con exito")
}
