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

func EjecFdisk(banderas []string) {
	unit := "k"
	size := -1
	driveLetter := ""
	name := ""
	type1 := "p"

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
		} else {
			fmt.Println("Parametro invalido")
		}
	}

	//fmt.Println(unit)
	//fmt.Println(size)
	//fmt.Println(driveLetter)
	//fmt.Println(name)
	//fmt.Println(type1)

	archivo, err := os.OpenFile(driveLetter+".dsk", os.O_RDWR, 0777)
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

	//fmt.Println(disk)

	numPart := -1
	for i, particion := range disk.Mbr_partitions {

		if particion.Part_s == int32(0) {
			numPart = i
			break

		} else {
			despTemp += int(particion.Part_s) + 1
		}

	}

	var nuevaPar partition

	nuevaPar.Part_status = [1]byte{'1'}

	if type1 == "p" || type1 == "e" || type1 == "l" {

		nuevaPar.Part_type = [1]byte{type1[0]}
	} else {
		fmt.Println("Tipo de particion no valida")
	}

	nuevaPar.Part_fit = [1]byte{'w'}

	if numPart < 0 {
		fmt.Println("No hay particiones disponibles")
		return
	}

	nuevaPar.Part_start = int32(despTemp)

	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	}

	if size < 0 {
		fmt.Println("tamano no valido")
		return
	}

	nuevaPar.Part_s = int32(size)

	if name == "" {
		fmt.Println("nombre invalido")
		return
	}

	nuevaPar.Part_name = [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	copy(nuevaPar.Part_name[:], name)

	if despTemp+int(nuevaPar.Part_s)+1 > int(disk.Mbr_tamano) {
		fmt.Println("tamano insuficiente para la particion")
		return
	}

	disk.Mbr_partitions[numPart] = nuevaPar

	archivo.Seek(0, 0)
	binary.Write(archivo, binary.LittleEndian, &disk)
	archivo.Close()

	fmt.Println("particion creada con exito")

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

	nombreDisco := "A"

	archivo, err := os.Create(nombreDisco + ".dsk")

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

func EjecRep() {
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

func EjecRepMkdisk() {
	archivo, err := os.Open("A.dsk")

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

			Dot += "Primaria"
			porcentaje := (float64(part.Part_s) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
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
