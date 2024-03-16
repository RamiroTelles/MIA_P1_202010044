package comandos

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Content_J struct {
	Operation [10]byte
	Path      [100]byte
	Content   [100]byte
	Date      [19]byte
}

type Journaling struct {
	Size      int32
	Ultimo    int32
	Contenido [50]Content_J
}

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

type superBloque struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_blocks_count int32
	S_free_inodes_count int32
	S_mtime             [19]byte
	S_umtime            [19]byte
	S_mnt_count         int32
	S_magic             int32
	S_inode_s           int32
	S_block_s           int32
	S_firts_ino         int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
}

type inodo struct {
	I_uid   int32
	I_gid   int32
	I_s     int32
	I_atime [19]byte
	I_ctime [19]byte
	I_mtime [19]byte
	I_block [15]int32
	I_type  [1]byte
	I_perm  [3]byte
}

type Inodo struct {
	I_uid   int32
	I_gid   int32
	I_s     int32
	I_atime [19]byte
	I_ctime [19]byte
	I_mtime [19]byte
	I_block [15]int32
	I_type  [1]byte
	I_perm  [3]byte
}

type b_content struct {
	B_name  [12]byte
	B_inodo int32
}

type bloqueCarpeta struct {
	B_content [4]b_content
}

type bloqueArchivos struct {
	B_content [64]byte
}

type bloqueApuntadores struct {
	B_pointers [16]int32
}

func LeerMounts() {
	uId = -1
	gId = -1

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
		//fmt.Println(ruta.Name()[:1])
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
			return
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
	//fmt.Println(mbrDisk.Mbr_tamano)
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
			return
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
			return
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

func EjecMkfs(banderas []string) {
	id := ""
	typeVar := "full"
	fs := "2fs"

	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-id" {
			id = dupla[1]

		} else if dupla[0] == "-type" {
			typeVar = dupla[1]

		} else if dupla[0] == "-fs" {
			if dupla[1] == "2fs" || dupla[1] == "3fs" {
				typeVar = dupla[1]
			} else {
				fmt.Println("No se acepta el tipo de Fs")
				return
			}

		} else {
			fmt.Println("Parametro invalido")
		}
	}

	index := VerificarParticionMontada(id)
	if index == -1 {
		fmt.Println("Particion no montada")
		return
	}

	if typeVar != "full" {
		fmt.Println("Tipo de formateo invalido")
		return
	}

	var n int
	if fs == "2fs" {
		n = int(math.Floor(float64(int(particionesMontadas[index].Size)-int(binary.Size(superBloque{}))) / float64(4+int(binary.Size(inodo{}))+3*int(binary.Size(bloqueArchivos{})))))

		//crear el ext2

		crearExt2(index, n, particionesMontadas[index])
	} else {
		//formatear primero la particion
		//crear el ext3
		n = int(math.Floor(float64(int(particionesMontadas[index].Size)-int(binary.Size(superBloque{}))) / float64(4+int(binary.Size(inodo{}))+3*int(binary.Size(bloqueArchivos{}))+binary.Size(Journaling{}))))

		crearExt3(index, n, particionesMontadas[index])
		//fmt.Println("Crear ext3", n)
	}

}

func crearExt2(index int, n int, mountActual Mount) {
	//var newSuperBloque superBloque
	//fmt.Println(newSuperBloque)

	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()
	archivo.Seek(int64(mountActual.Start), 0)

	for i := int32(0); i < mountActual.Size; i++ {

		err := binary.Write(archivo, binary.LittleEndian, [1]byte{0})
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	var sbloque superBloque

	sbloque.S_filesystem_type = 2
	sbloque.S_bm_inode_start = int32(mountActual.Size) + int32(binary.Size(superBloque{}))
	sbloque.S_bm_block_start = sbloque.S_bm_inode_start + int32(n)
	sbloque.S_inode_start = sbloque.S_bm_block_start + int32(3*n)
	sbloque.S_block_start = sbloque.S_inode_start + int32(n*int(binary.Size(inodo{})))

	sbloque.S_inodes_count = int32(n)
	sbloque.S_blocks_count = int32(3 * n)

	sbloque.S_free_inodes_count = int32(n)
	sbloque.S_free_blocks_count = int32(3 * n)
	fechaActual := time.Now()
	fechaF := fechaActual.Format("2006-01-02 15:04:05")
	copy(sbloque.S_mtime[:], []byte(fechaF))
	copy(sbloque.S_umtime[:], []byte(fechaF))
	sbloque.S_mnt_count = 1
	sbloque.S_magic = 61267
	sbloque.S_inode_s = int32(binary.Size(inodo{}))
	sbloque.S_block_s = int32(binary.Size(bloqueArchivos{}))
	sbloque.S_firts_ino = 0
	sbloque.S_first_blo = 0

	var newInodo inodo

	var newblock bloqueCarpeta

	newInodo.I_uid = 1
	newInodo.I_gid = 1
	newInodo.I_s = 0
	copy(newInodo.I_atime[:], []byte(fechaF))
	copy(newInodo.I_ctime[:], []byte(fechaF))
	copy(newInodo.I_mtime[:], []byte(fechaF))

	for i := int32(0); i < 15; i++ {
		newInodo.I_block[i] = -1
	}

	newInodo.I_block[0] = 0
	newInodo.I_type = [1]byte{'0'}
	newInodo.I_perm = [3]byte{'6', '6', '4'}

	copy(newblock.B_content[0].B_name[:], ".")
	newblock.B_content[0].B_inodo = 0
	copy(newblock.B_content[1].B_name[:], "..")
	newblock.B_content[1].B_inodo = 0
	newblock.B_content[2].B_inodo = -1
	newblock.B_content[3].B_inodo = -1

	archivo.Seek(int64(sbloque.S_inode_start), 0)
	err = binary.Write(archivo, binary.LittleEndian, &newInodo)
	if err != nil {
		fmt.Println("Error al escribir el inodo 0: ", err)
		return
	}
	archivo.Seek(int64(sbloque.S_block_start), 0)
	err = binary.Write(archivo, binary.LittleEndian, &newblock)
	if err != nil {
		fmt.Println("Error al escribir el bloque 0: ", err)
		return
	}

	sbloque.S_free_blocks_count--
	sbloque.S_free_inodes_count--
	sbloque.S_firts_ino++
	sbloque.S_first_blo++

	archivo.Seek(int64(sbloque.S_bm_block_start), 0)
	err = binary.Write(archivo, binary.LittleEndian, [1]byte{1})
	if err != nil {
		fmt.Println("Error al escribir el bitmap inodo 0: ", err)
		return
	}

	archivo.Seek(int64(sbloque.S_bm_inode_start), 0)
	err = binary.Write(archivo, binary.LittleEndian, [1]byte{1})
	if err != nil {
		fmt.Println("Error al escribir el bitmap bloque 0: ", err)
		return
	}

	archivo.Seek(int64(mountActual.Start), 0)
	err = binary.Write(archivo, binary.LittleEndian, &sbloque)
	if err != nil {
		fmt.Println("Error al escribir el superbloque: ", err)
		return
	}

	archivo.Close()

	uId = 1
	gId = 1

	CrearArchivo("/users.txt", "1,G,root\n1,U,root,root,123\n", false, index)

	uId = -1
	gId = -1

	fmt.Println("EXT2 creado con exito")

}

func crearExt3(index int, n int, mountActual Mount) {
	//var newSuperBloque superBloque
	//fmt.Println(newSuperBloque)
	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()
	archivo.Seek(int64(mountActual.Size), 0)

	for i := int32(0); i < mountActual.Size; i++ {

		err := binary.Write(archivo, binary.LittleEndian, [1]byte{0})
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	var sbloque superBloque

	sbloque.S_filesystem_type = 3
	sbloque.S_bm_inode_start = int32(mountActual.Size) + int32(binary.Size(superBloque{})) + int32(binary.Size(Journaling{}))
	sbloque.S_bm_block_start = sbloque.S_bm_inode_start + int32(n)
	sbloque.S_inode_start = sbloque.S_bm_block_start + int32(3*n)
	sbloque.S_block_start = sbloque.S_inode_start + int32(n*int(binary.Size(inodo{})))

	sbloque.S_inodes_count = int32(n)
	sbloque.S_blocks_count = int32(3 * n)

	sbloque.S_free_inodes_count = int32(n)
	sbloque.S_free_blocks_count = int32(3 * n)
	fechaActual := time.Now()
	fechaF := fechaActual.Format("2006-01-02 15:04:05")
	copy(sbloque.S_mtime[:], []byte(fechaF))
	copy(sbloque.S_umtime[:], []byte(fechaF))
	sbloque.S_mnt_count = 1
	sbloque.S_magic = 61267
	sbloque.S_inode_s = int32(binary.Size(inodo{}))
	sbloque.S_block_s = int32(binary.Size(bloqueArchivos{}))
	sbloque.S_firts_ino = 0
	sbloque.S_first_blo = 0

	var newInodo inodo

	var newblock bloqueCarpeta

	newInodo.I_uid = 1
	newInodo.I_gid = 1
	newInodo.I_s = 0
	copy(newInodo.I_atime[:], []byte(fechaF))
	copy(newInodo.I_ctime[:], []byte(fechaF))
	copy(newInodo.I_mtime[:], []byte(fechaF))

	for i := int32(0); i < 15; i++ {
		newInodo.I_block[i] = -1
	}

	newInodo.I_block[0] = 0
	newInodo.I_type = [1]byte{'0'}
	newInodo.I_perm = [3]byte{'6', '6', '4'}

	copy(newblock.B_content[0].B_name[:], ".")
	newblock.B_content[0].B_inodo = 0
	copy(newblock.B_content[1].B_name[:], "..")
	newblock.B_content[1].B_inodo = 0
	newblock.B_content[2].B_inodo = -1
	newblock.B_content[3].B_inodo = -1

	var journal Journaling

	var tempContentJ Content_J

	copy(tempContentJ.Operation[:], "mkdir")
	copy(tempContentJ.Path[:], "/")
	copy(tempContentJ.Content[:], "")
	copy(tempContentJ.Date[:], []byte(fechaF))

	journal.Contenido[0] = tempContentJ
	journal.Size = 1
	journal.Ultimo = 0

	archivo.Seek(int64(sbloque.S_inode_start), 0)
	err = binary.Write(archivo, binary.LittleEndian, &newInodo)
	if err != nil {
		fmt.Println("Error al escribir el inodo 0: ", err)
		return
	}
	archivo.Seek(int64(sbloque.S_block_start), 0)
	err = binary.Write(archivo, binary.LittleEndian, &newblock)
	if err != nil {
		fmt.Println("Error al escribir el bloque 0: ", err)
		return
	}

	sbloque.S_free_blocks_count--
	sbloque.S_free_inodes_count--
	sbloque.S_firts_ino++
	sbloque.S_first_blo++

	archivo.Seek(int64(sbloque.S_bm_block_start), 0)
	err = binary.Write(archivo, binary.LittleEndian, [1]byte{1})
	if err != nil {
		fmt.Println("Error al escribir el bitmap Inodo 0: ", err)
		return
	}

	archivo.Seek(int64(sbloque.S_bm_inode_start), 0)
	err = binary.Write(archivo, binary.LittleEndian, [1]byte{1})
	if err != nil {
		fmt.Println("Error al escribir el bitmap Bloque 0: ", err)
		return
	}

	archivo.Seek(int64(mountActual.Start), 0)
	err = binary.Write(archivo, binary.LittleEndian, &sbloque)
	if err != nil {
		fmt.Println("Error al escribir el superbloque: ", err)
		return
	}
	err = binary.Write(archivo, binary.LittleEndian, &journal)
	if err != nil {
		fmt.Println("Error al escribir el journaling: ", err)
		return
	}

	archivo.Close()

	uId = 1
	gId = 1
	//crear el archivo usuarios con el mkfile
	CrearArchivo("/users.txt", "1,G,root\n1,U,root,root,123\n", false, index)

	uId = -1
	gId = -1

	fmt.Println("EXT3 creado con exito")

}

func EjecEdit(banderas []string) {
	ruta := ""

	cont := ""
	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-path" {
			ruta = dupla[1]

		} else if dupla[0] == "-cont" {
			cont = dupla[1]
		} else {
			fmt.Println("Parametro invalido")
		}
	}
	if cont != "" {
		cont = obtenerArchivo(cont)
	} else {
		fmt.Println("No se ingreso el -cont")
		return
	}

	if ruta == "" {
		fmt.Println("No se ingreso el -path")
		return
	}

	index := VerificarParticionMontada(actualIdMount)

	editarArchivo(index, ruta, cont)

}

func EjecMkfile(banderas []string) {
	ruta := ""
	r := false
	size := 0
	cont := ""
	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-path" {
			ruta = dupla[1]

		} else if dupla[0] == "-r" {
			r = true
		} else if dupla[0] == "-size" {
			size1, err := strconv.Atoi(dupla[1])
			size = size1
			if err != nil {
				fmt.Println("No se pudo leer ell -size")
				return
			}
		} else if dupla[0] == "-cont" {
			cont = dupla[1]
		} else {
			fmt.Println("Parametro invalido")
		}
	}
	//fmt.Println(r)
	fmt.Println(size)

	if cont != "" {
		cont = obtenerArchivo(cont)
	}

	if ruta == "" {
		fmt.Println("No se ingreso el -path")
		return
	}
	index := VerificarParticionMontada(actualIdMount)

	CrearArchivo(ruta, cont, r, index)
}

func obtenerArchivo(ruta string) string {
	archivo, err := os.Open(ruta)

	if err != nil {
		fmt.Println("Error al abrir el archivo: ", err)
		return ""
	}
	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)
	lineas := ""
	for scanner.Scan() {
		lineas += scanner.Text()

	}

	return lineas
}

func EjecMkdir(banderas []string) {
	ruta := ""
	r := false
	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-path" {
			ruta = dupla[1]

		} else if dupla[0] == "-r" {
			r = true
		} else {
			fmt.Println("Parametro invalido")
		}
	}
	fmt.Println(r)
	index := VerificarParticionMontada(actualIdMount)
	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	var numInodo int
	if ruta == "/" {
		numInodo = 0
	} else {
		lRuta := strings.Split(ruta[1:], "/")
		lRuta = lRuta[:len(lRuta)-1]
		numInodo = obtenerNumInodo(lRuta, archivo, sblock)
	}

	archivo.Seek(int64(sblock.S_inode_start+int32(numInodo)*int32(binary.Size(inodo{}))), 0)
	var inodoTemp inodo
	nombre := strings.Split(ruta[1:], "/")
	nombre = nombre[len(nombre)-1:]
	err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		return
	}

	var punteroTemp int
	for i, ptr := range inodoTemp.I_block {

		punteroTemp = rellenarBloques(int(ptr), i, nombre[0], &sblock, archivo)

		if punteroTemp != -1 {
			inodoTemp.I_block[i] = int32(punteroTemp)

			break

		}

	}

	var newInodo inodo

	newInodo.I_uid = int32(uId)
	newInodo.I_gid = int32(gId)
	newInodo.I_s = 0
	inodoTemp.I_s = inodoTemp.I_s + newInodo.I_s
	fechaActual := time.Now()
	fechaF := fechaActual.Format("2006-01-02 15:04:05")

	copy(newInodo.I_atime[:], []byte(fechaF))
	copy(newInodo.I_ctime[:], []byte(fechaF))
	copy(newInodo.I_mtime[:], []byte(fechaF))

	for i := int32(0); i < 15; i++ {
		newInodo.I_block[i] = -1
	}

	newInodo.I_type = [1]byte{'0'}
	newInodo.I_perm = [3]byte{'6', '6', '4'}
	newInodo.I_block[0] = sblock.S_first_blo

	var newBlockC bloqueCarpeta
	newBlockC.B_content[0].B_inodo = sblock.S_firts_ino
	copy(newBlockC.B_content[0].B_name[:], ".")
	newBlockC.B_content[1].B_inodo = int32(numInodo)
	copy(newBlockC.B_content[1].B_name[:], "..")
	newBlockC.B_content[2].B_inodo = -1
	newBlockC.B_content[3].B_inodo = -1

	archivo.Seek(int64(sblock.S_inode_start+int32(numInodo)*int32(binary.Size(inodo{}))), 0)
	err = binary.Write(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al escribir el inodo : ", err)
		return
	}

	archivo.Seek(int64(sblock.S_inode_start+sblock.S_firts_ino*int32(binary.Size(inodo{}))), 0)
	err = binary.Write(archivo, binary.LittleEndian, &newInodo)
	if err != nil {
		fmt.Println("Error al escribir el inodo archivos: ", err)
		return
	}

	archivo.Seek(int64(sblock.S_bm_inode_start+sblock.S_firts_ino), 0)
	err = binary.Write(archivo, binary.LittleEndian, [1]byte{1})
	if err != nil {
		fmt.Println("Error al escribir el bitmap Inodos: ", err)
		return
	}

	sblock.S_firts_ino = encontrarInodoLibre(&sblock, archivo)
	sblock.S_free_inodes_count--

	archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueCarpeta{})*int(sblock.S_first_blo))), 0)
	err = binary.Write(archivo, binary.LittleEndian, &newBlockC)
	if err != nil {
		fmt.Println("Error al escribir el bloque Carpetas: ", err)
		return
	}

	archivo.Seek(int64(sblock.S_bm_block_start+sblock.S_first_blo), 0)
	err = binary.Write(archivo, binary.LittleEndian, [1]byte{1})
	if err != nil {
		fmt.Println("Error al escribir el bloque Carpetas: ", err)
		return
	}

	sblock.S_first_blo = encontrarBloqueLibre(&sblock, archivo)
	sblock.S_free_blocks_count--

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Write(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al escribir el superbloque: ", err)
		return
	}
	if sblock.S_filesystem_type == 3 {
		var journal Journaling
		archivo.Seek(int64(particionesMontadas[index].Start+int32(binary.Size(superBloque{}))), 0)
		binary.Read(archivo, binary.LittleEndian, &journal)
		copy(journal.Contenido[journal.Size].Operation[:], "mkdir")
		copy(journal.Contenido[journal.Size].Path[:], ruta)
		copy(journal.Contenido[journal.Size].Content[:], "")
		copy(journal.Contenido[journal.Size].Date[:], []byte(fechaF))
		archivo.Seek(int64(particionesMontadas[index].Start+int32(binary.Size(superBloque{}))), 0)
		err = binary.Write(archivo, binary.LittleEndian, &journal)
		if err != nil {
			fmt.Println("Error al escribir el journaling: ", err)
			return
		}
	}

	//cerrar archivo
	archivo.Close()

}

func leerArchivo(inodoTemp inodo, archivo *os.File, sblock superBloque) string {

	result := ""
	var bArchivo bloqueArchivos

	for _, ptr := range inodoTemp.I_block {

		if ptr != -1 {
			archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueArchivos{}))*ptr), 0)

			err := binary.Read(archivo, binary.LittleEndian, &bArchivo)
			if err != nil {
				fmt.Println("Error al leer el bloque archivos: ", err)
				return ""
			}
			result += strings.TrimRight(string(bArchivo.B_content[:]), string(rune(0)))

		}
	}
	return result
}

func CrearArchivo(ruta string, cont string, r bool, index int) {

	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	var numInodo int
	if ruta == "/" {
		numInodo = 0
	} else {
		lRuta := strings.Split(ruta[1:], "/")
		lRuta = lRuta[:len(lRuta)-1]
		numInodo = obtenerNumInodo(lRuta, archivo, sblock)
	}

	archivo.Seek(int64(sblock.S_inode_start+int32(numInodo)*int32(binary.Size(inodo{}))), 0)
	var inodoTemp inodo
	nombre := strings.Split(ruta[1:], "/")
	nombre = nombre[len(nombre)-1:]
	err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		return
	}
	//var despTemp int
	//var band bool
	var punteroTemp int
	for i, ptr := range inodoTemp.I_block {

		punteroTemp = rellenarBloques(int(ptr), i, nombre[0], &sblock, archivo)

		if punteroTemp != -1 {
			inodoTemp.I_block[i] = int32(punteroTemp)

			break

		}

	}
	//fmt.Println(punteroTemp)
	//escribir el inodo

	//crear el nuevo inodo para el archivo
	var newInodo inodo

	newInodo.I_uid = int32(uId)
	newInodo.I_gid = int32(gId)
	newInodo.I_s = int32(len(cont))
	inodoTemp.I_s = inodoTemp.I_s + newInodo.I_s
	fechaActual := time.Now()
	fechaF := fechaActual.Format("2006-01-02 15:04:05")

	copy(newInodo.I_atime[:], []byte(fechaF))
	copy(newInodo.I_ctime[:], []byte(fechaF))
	copy(newInodo.I_mtime[:], []byte(fechaF))

	for i := int32(0); i < 15; i++ {
		newInodo.I_block[i] = -1
	}

	newInodo.I_type = [1]byte{'1'}
	newInodo.I_perm = [3]byte{'6', '6', '4'}

	//crear los bloques correspondientes para el archivo
	if escribirBloquesArchivo(&newInodo, cont, &sblock, archivo) {
		fmt.Println("Error al escribir el archivo")
		return
	}

	//escribir ese inodo y esos bloques
	archivo.Seek(int64(sblock.S_inode_start+int32(numInodo)*int32(binary.Size(inodo{}))), 0)
	err = binary.Write(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al escribir el inodo : ", err)
		return
	}

	archivo.Seek(int64(sblock.S_inode_start+sblock.S_firts_ino*int32(binary.Size(inodo{}))), 0)
	err = binary.Write(archivo, binary.LittleEndian, &newInodo)
	if err != nil {
		fmt.Println("Error al escribir el inodo archivos: ", err)
		return
	}

	archivo.Seek(int64(sblock.S_bm_inode_start+sblock.S_firts_ino), 0)
	err = binary.Write(archivo, binary.LittleEndian, [1]byte{1})
	if err != nil {
		fmt.Println("Error al escribir el bitmap Inodos: ", err)
		return
	}

	sblock.S_firts_ino = encontrarInodoLibre(&sblock, archivo)
	sblock.S_free_inodes_count--

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Write(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al escribir el superbloque: ", err)
		return
	}
	if sblock.S_filesystem_type == 3 {
		var journal Journaling
		archivo.Seek(int64(particionesMontadas[index].Start+int32(binary.Size(superBloque{}))), 0)
		binary.Read(archivo, binary.LittleEndian, &journal)

		copy(journal.Contenido[journal.Size].Operation[:], "mkfile")
		copy(journal.Contenido[journal.Size].Path[:], ruta)
		copy(journal.Contenido[journal.Size].Content[:], cont)
		copy(journal.Contenido[journal.Size].Date[:], []byte(fechaF))
		archivo.Seek(int64(particionesMontadas[index].Start+int32(binary.Size(superBloque{}))), 0)
		err = binary.Write(archivo, binary.LittleEndian, &journal)
		if err != nil {
			fmt.Println("Error al escribir el journaling: ", err)
			return
		}
	}

	//cerrar archivo
	archivo.Close()
}

func EjecCat(banderas []string) {

	var paths []string

	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if strings.Contains(dupla[0], "-file") {
			paths = append(paths, dupla[1])

		} else {
			fmt.Println("Parametro invalido")
		}
	}
	//id
	//actualIdMount = "A144"
	if actualIdMount == "" {
		fmt.Println("No hay una sesion iniciada")
		return
	}

	index := VerificarParticionMontada(actualIdMount)

	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	var numInodo int
	var inodoTemp inodo
	for i, ruta := range paths {

		if ruta == "/" {
			numInodo = 0
		} else {
			lRuta := strings.Split(ruta[1:], "/")

			numInodo = obtenerNumInodo(lRuta, archivo, sblock)
		}
		//fmt.Println(numInodo)
		archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(inodo{}))*int32(numInodo)), 0)
		err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
		if err != nil {
			fmt.Println("Error al leer el inodo: ", err)
			return
		}
		if numInodo != -1 {
			fmt.Println(leerArchivo(inodoTemp, archivo, sblock))
		} else {
			fmt.Println("No se encontro la ruta del archivo: " + strconv.Itoa(i+1))
		}

	}

}

func EjecMkGrp(banderas []string) {
	name := ""
	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-name" {
			name = dupla[1]

		} else {
			fmt.Println("Parametro invalido")
		}
	}

	if name == "" {
		fmt.Println("Ingrese el campo -name")
		return
	}
	if uId != 1 {
		fmt.Println("Solo el usuario root puede crear grupos")
	}

	index := VerificarParticionMontada(actualIdMount)

	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	var inodoTemp inodo

	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(inodo{}))), 0)
	err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		return
	}

	txt := leerArchivo(inodoTemp, archivo, sblock)
	nexId := nextIdGroup(txt)

	txt += strconv.Itoa(nexId) + ",G," + name + "\n"

	archivo.Close()
	editarArchivo(index, "/users.txt", txt)

	fmt.Println("Grupo creado con exito")

}

func EjecRmGrp(banderas []string) {
	name := ""
	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-name" {
			name = dupla[1]

		} else {
			fmt.Println("Parametro invalido")
		}
	}

	if name == "" {
		fmt.Println("Ingrese el campo -name")
		return
	}

	if uId != 1 {
		fmt.Println("Solo el usuario root puede eliminar grupos")
	}

	index := VerificarParticionMontada(actualIdMount)

	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	var inodoTemp inodo

	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(inodo{}))), 0)
	err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		return
	}

	txt := leerArchivo(inodoTemp, archivo, sblock)

	lineas := strings.Split(txt, "\n")
	lineas = lineas[:len(lineas)-1]

	band := false
	for i, linea := range lineas {

		if linea[2] == 'G' && linea[0] != '0' {
			campos := strings.Split(linea, ",")
			if campos[2] == name {
				nuevaLinea := []byte(linea)
				nuevaLinea[0] = '0'
				lineas[i] = string(nuevaLinea)
				band = true
				break

			}
		}

	}
	newTxt := ""
	for _, linea := range lineas {

		newTxt += linea + "\n"

	}

	if !band {
		fmt.Println("No se encontro el grupo")
		return
	}
	archivo.Close()
	editarArchivo(index, "/users.txt", newTxt)

	fmt.Println("Grupo eliminado con exito con exito")

}

func EjecMkUsr(banderas []string) {
	name := ""
	pass := ""
	group := ""
	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-user" {
			name = dupla[1]

		} else if dupla[0] == "-pass" {
			pass = dupla[1]

		} else if dupla[0] == "-grp" {
			group = dupla[1]

		} else {
			fmt.Println("Parametro invalido")
		}
	}

	if name == "" {
		fmt.Println("Ingrese el campo -name")
		return
	}

	if pass == "" {
		fmt.Println("Ingrese el campo -pass")
		return
	}

	if group == "" {
		fmt.Println("Ingrese el campo -grp")
		return
	}

	if uId != 1 {
		fmt.Println("Solo el usuario root puede crear usuarios")
	}

	index := VerificarParticionMontada(actualIdMount)

	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	var inodoTemp inodo

	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(inodo{}))), 0)
	err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		return
	}

	txt := leerArchivo(inodoTemp, archivo, sblock)
	nexId := nextIdUser(txt)

	txt += strconv.Itoa(nexId) + ",U," + group + "," + name + "," + pass + "\n"

	archivo.Close()
	editarArchivo(index, "/users.txt", txt)

	fmt.Println("Grupo creado con exito")

}

func EjecRmUsr(banderas []string) {
	name := ""
	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-user" {
			name = dupla[1]

		} else {
			fmt.Println("Parametro invalido")
		}
	}

	if name == "" {
		fmt.Println("Ingrese el campo -name")
		return
	}

	if uId != 1 {
		fmt.Println("Solo el usuario root puede eliminar usuarios")
	}

	index := VerificarParticionMontada(actualIdMount)

	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	var inodoTemp inodo

	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(inodo{}))), 0)
	err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		return
	}

	txt := leerArchivo(inodoTemp, archivo, sblock)

	lineas := strings.Split(txt, "\n")
	lineas = lineas[:len(lineas)-1]

	band := false
	for i, linea := range lineas {

		if linea[2] == 'U' && linea[0] != '0' {
			campos := strings.Split(linea, ",")
			if campos[3] == name {
				nuevaLinea := []byte(linea)
				nuevaLinea[0] = '0'
				lineas[i] = string(nuevaLinea)
				band = true
				break

			}
		}

	}

	newTxt := ""
	for _, linea := range lineas {

		newTxt += linea + "\n"

	}

	if !band {
		fmt.Println("No se encontro el Usuario")
		return
	}
	archivo.Close()
	editarArchivo(index, "/users.txt", newTxt)

	fmt.Println("Grupo eliminado con exito con exito")

}

func editarArchivo(index int, ruta string, cont string) {
	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	var numInodo int
	if ruta == "/" {
		numInodo = 0
	} else {
		lRuta := strings.Split(ruta[1:], "/")

		numInodo = obtenerNumInodo(lRuta, archivo, sblock)
	}

	eliminarBloquesInodo(numInodo, &sblock, archivo)

	var inodoTemp inodo
	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(inodo{}))*int32(numInodo)), 0)
	err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		return
	}

	//crear los bloques correspondientes para el archivo
	if escribirBloquesArchivo(&inodoTemp, cont, &sblock, archivo) {
		fmt.Println("Error al escribir el archivo")
		return
	}
	inodoTemp.I_s = int32(len(cont))

	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(inodo{}))*int32(numInodo)), 0)
	err = binary.Write(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al escribir el inodo: ", err)
		return
	}

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Write(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al escribir el superbloque: ", err)
		return
	}

	archivo.Close()

}

func eliminarBloquesInodo(numInodo int, sblock *superBloque, archivo *os.File) {
	var inodoTemp inodo
	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(inodo{}))*int32(numInodo)), 0)
	err := binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		return
	}

	for i, ptr := range inodoTemp.I_block {
		if ptr != -1 {
			eliminarBloque(int(ptr), sblock, archivo)
			inodoTemp.I_block[i] = -1
		}
	}

	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(inodo{}))*int32(numInodo)), 0)
	err = binary.Write(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al escribir el inodo: ", err)
		return
	}

}

func eliminarBloque(ptr int, sblock *superBloque, archivo *os.File) {

	bufer := make([]byte, 64)

	archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueArchivos{}))*int32(ptr)), 0)
	err := binary.Write(archivo, binary.LittleEndian, &bufer)
	if err != nil {
		fmt.Println("Error al eliminar el bloque: ", err)
		return
	}

	archivo.Seek(int64(sblock.S_bm_block_start+int32(ptr)), 0)
	err = binary.Write(archivo, binary.LittleEndian, [1]byte{0})
	if err != nil {
		fmt.Println("Error al eliminar el bitmap Bloque: ", err)
		return
	}
	sblock.S_free_blocks_count++

	sblock.S_first_blo = encontrarBloqueLibre(sblock, archivo)

}

func nextIdUser(txt string) int {

	lineas := strings.Split(txt, "\n")
	id := 1

	lineas = lineas[:len(lineas)-1]

	var numTemp int
	for _, linea := range lineas {
		if linea[2] == 'U' {
			numTemp, _ = strconv.Atoi(string(linea[0]))
			if numTemp > id {
				id = numTemp
			}

		}
	}

	return id + 1

}

func nextIdGroup(txt string) int {

	lineas := strings.Split(txt, "\n")
	id := 1

	lineas = lineas[:len(lineas)-1]

	var numTemp int
	for _, linea := range lineas {
		if linea[2] == 'G' {
			numTemp, _ = strconv.Atoi(string(linea[0]))
			if numTemp > id {
				id = numTemp
			}

		}
	}

	return id + 1

}

func EjecLogin(banderas []string) {

	user := ""
	pass := ""
	id := ""

	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-user" {
			user = dupla[1]

		} else if dupla[0] == "-pass" {
			pass = dupla[1]

		} else if dupla[0] == "-id" {
			id = dupla[1]

		} else {
			fmt.Println("Parametro invalido")
		}
	}

	if user == "" {
		fmt.Println("Ingrese el campo -user")
		return
	}

	if pass == "" {
		fmt.Println("Ingrese el campo -pass")
		return
	}

	if id == "" {
		fmt.Println("Ingrese el campo -id")
		return
	}

	if uId != -1 {
		fmt.Println("Ya hay una sesion iniciada")
		return

	}

	index := VerificarParticionMontada(id)

	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var inodoTemp inodo
	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	if !(sblock.S_filesystem_type == 2 || sblock.S_filesystem_type == 3) {
		println("El sistema de archivos no es 2fs ni 3fs o no está formateado")
		return
	}

	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(inodo{}))*1), 0)
	err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	usersTxt := leerArchivo(inodoTemp, archivo, sblock)

	lineas := strings.Split(usersTxt, "\n")
	lineas = lineas[:len(lineas)-1]
	nombreGrupo := ""
	for _, linea := range lineas {

		if linea[2] == 'U' && linea[0] != '0' {
			campos := strings.Split(linea, ",")
			if campos[3] == user && campos[4] == pass {
				uId, _ = strconv.Atoi(campos[0])
				nombreGrupo = campos[2]
				break

			}
		}

	}
	if nombreGrupo == "" {
		fmt.Println("Usuario no encontrado")
		return
	}

	for _, linea := range lineas {

		if linea[2] == 'G' {
			campos := strings.Split(linea, ",")
			if campos[2] == nombreGrupo {
				gId, _ = strconv.Atoi(campos[0])
				if linea[0] == '0' {
					fmt.Println("Grupo no encontrado")
					uId = -1
					gId = -1
					return
				}
				break

			}
		}

	}

	fmt.Println("Usuario Logueado con exito")
	fmt.Println("Usuario: " + strconv.Itoa(uId))
	fmt.Println("Grupo: " + strconv.Itoa(gId))
	actualIdMount = id
}

func EjecLogout() {
	if uId == -1 {
		fmt.Println("Aun no hay sesion iniciada")
		return
	}
	uId = -1
	gId = -1
	actualIdMount = ""
	fmt.Println("Sesion cerrada con Exito")
}

func EjecPause() {
	fmt.Println("Presiona Enter para continuar...")
	esperarEnter()
	//fmt.Println("Continuando...")
}

func esperarEnter() {
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}

func encontrarBloqueLibre(sblock *superBloque, archivo *os.File) int32 {

	sizeBitmap := int((sblock.S_bm_block_start - sblock.S_bm_inode_start) * 3)
	var tempByte [1]byte
	archivo.Seek(int64(sblock.S_bm_block_start), 0)
	for i := 0; i < sizeBitmap; i++ {

		err := binary.Read(archivo, binary.LittleEndian, &tempByte)
		if err != nil {
			fmt.Println("Error al leer el bitmap Bloques: ", err)
			return -1
		}
		if tempByte == [1]byte{0} {
			return int32(i)
		}
	}
	return -1
}

func encontrarInodoLibre(sblock *superBloque, archivo *os.File) int32 {

	sizeBitmap := int((sblock.S_bm_block_start - sblock.S_bm_inode_start))
	var tempByte [1]byte
	archivo.Seek(int64(sblock.S_bm_inode_start), 0)
	for i := 0; i < sizeBitmap; i++ {

		err := binary.Read(archivo, binary.LittleEndian, &tempByte)
		if err != nil {
			fmt.Println("Error al leer el bitmap Inodos: ", err)
			return -1
		}
		if tempByte == [1]byte{0} {
			return int32(i)
		}
	}
	return -1
}

// ---------------------Funcion Rellenar bloques controla bloques de apuntadores

/* func rellenarBloques(ptr int, tipo int, name string, sblock *superBloque, archivo *os.File) int {
	var punteroTemp int
	var bloquePtr bloqueApuntadores
	if tipo == 14 {
		// trile indirecto

		if ptr == -1 {
			//crear bloque doble indirecto
			punteroTemp = crearBloqueCarpetas(name, sblock, archivo)
			punteroTemp = crearBloquePtr1(punteroTemp, sblock, archivo)
			punteroTemp = crearBloquePtr1(punteroTemp, sblock, archivo)
			return crearBloquePtr1(punteroTemp, sblock, archivo)

		} else {
			//Leer el bloque de punteros
			archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueApuntadores{})*ptr)), 0)
			binary.Read(archivo, binary.LittleEndian, &bloquePtr)
			for i, punteroBlock := range bloquePtr.b_pointers {
				punteroTemp = rellenarBloques(int(punteroBlock), tipo-1, name, sblock, archivo)
				if punteroTemp != -1 {
					bloquePtr.b_pointers[i] = int32(punteroTemp)
					archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueApuntadores{})*ptr)), 0)
					binary.Write(archivo, binary.LittleEndian, &bloquePtr)
					return ptr
				}
			}
		}

	} else if tipo == 13 {
		//doble indirecto
		if ptr == -1 {
			//crear bloque doble indirecto
			punteroTemp = crearBloqueCarpetas(name, sblock, archivo)
			punteroTemp = crearBloquePtr1(punteroTemp, sblock, archivo)
			return crearBloquePtr1(punteroTemp, sblock, archivo)

		} else {
			//Leer el bloque de punteros
			archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueApuntadores{})*ptr)), 0)
			binary.Read(archivo, binary.LittleEndian, &bloquePtr)
			for i, punteroBlock := range bloquePtr.b_pointers {
				punteroTemp = rellenarBloques(int(punteroBlock), tipo-1, name, sblock, archivo)
				if punteroTemp != -1 {
					bloquePtr.b_pointers[i] = int32(punteroTemp)
					archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueApuntadores{})*ptr)), 0)
					binary.Write(archivo, binary.LittleEndian, &bloquePtr)
					return ptr
				}
			}
		}

	} else if tipo == 12 {
		//simple indirecto
		if ptr == -1 {
			//crear bloque simple indirecto
			punteroTemp = crearBloqueCarpetas(name, sblock, archivo)
			return crearBloquePtr1(punteroTemp, sblock, archivo)
		} else {
			//Leer el bloque de punteros
			archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueApuntadores{})*ptr)), 0)
			binary.Read(archivo, binary.LittleEndian, &bloquePtr)
			for i, punteroBlock := range bloquePtr.b_pointers {
				punteroTemp = rellenarBloques(int(punteroBlock), tipo-1, name, sblock, archivo)
				if punteroTemp != -1 {
					bloquePtr.b_pointers[i] = int32(punteroTemp)
					archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueApuntadores{})*ptr)), 0)
					binary.Write(archivo, binary.LittleEndian, &bloquePtr)
					return ptr
				}
			}
		}

	} else {
		//directo
		if ptr == -1 {
			//crear bloque Carpetas
			return crearBloqueCarpetas(name, sblock, archivo)
		} else {
			if rellenarBloqueCarpetas(name, ptr, sblock, archivo) {
				return ptr
			}
		}
	}

	return -1

} */

func escribirBloquesArchivo(newInodo *inodo, cont string, sblock *superBloque, archivo *os.File) bool {

	cantBloques := len(cont) / 64
	fmt.Println(cantBloques)
	//fmt.Println("Pasa por aqui")
	if cantBloques > 15 {
		fmt.Println("no se pudo escribir el archivo,archivo muy grande")
		return true
	} else {

		for i := 0; i <= cantBloques; i++ {

			if len(cont) > 64 {
				var newblock bloqueArchivos
				copy(newblock.B_content[:], cont[:64])
				cont = cont[64:]

				newInodo.I_block[i] = sblock.S_first_blo

				archivo.Seek(int64(sblock.S_block_start+sblock.S_first_blo*int32(binary.Size(bloqueArchivos{}))), 0)
				err := binary.Write(archivo, binary.LittleEndian, &newblock)
				if err != nil {
					fmt.Println("Error al escribir el bloque archivos: ", err)
					return true
				}

				archivo.Seek(int64(sblock.S_bm_block_start+sblock.S_first_blo), 0)
				err = binary.Write(archivo, binary.LittleEndian, [1]byte{1})
				if err != nil {
					fmt.Println("Error al escribir el bitmap bloques: ", err)
					return true
				}

				sblock.S_first_blo = encontrarBloqueLibre(sblock, archivo)
				sblock.S_free_blocks_count--

			} else {
				var newblock bloqueArchivos
				copy(newblock.B_content[:], cont)

				newInodo.I_block[i] = sblock.S_first_blo

				archivo.Seek(int64(sblock.S_block_start+sblock.S_first_blo*int32(binary.Size(bloqueArchivos{}))), 0)
				err := binary.Write(archivo, binary.LittleEndian, &newblock)
				if err != nil {
					fmt.Println("Error al escribir el bloque archivos: ", err)
					return true
				}

				archivo.Seek(int64(sblock.S_bm_block_start+sblock.S_first_blo), 0)
				err = binary.Write(archivo, binary.LittleEndian, [1]byte{1})
				if err != nil {
					fmt.Println("Error al escribir el bitmap bloques: ", err)
					return true
				}

				sblock.S_first_blo = encontrarBloqueLibre(sblock, archivo)
				sblock.S_free_blocks_count--

			}

		}
	}
	return false
}

// -------------Funcion Rellenar bloques no controla bloques de apuntadores
func rellenarBloques(ptr int, tipo int, name string, sblock *superBloque, archivo *os.File) int {

	//directo
	if ptr == -1 {
		//crear bloque Carpetas
		return crearBloqueCarpetas(name, sblock, archivo)
	} else {
		if rellenarBloqueCarpetas(name, ptr, sblock, archivo) {
			return ptr
		}
	}

	return -1

}

func rellenarBloqueCarpetas(name string, ptr int, sblock *superBloque, archivo *os.File) bool {

	var blockTemp bloqueCarpeta

	archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueCarpeta{}))*int32(ptr)), 0)

	err := binary.Read(archivo, binary.LittleEndian, &blockTemp)
	if err != nil {
		fmt.Println("Error al leer el bloque: ", err)
		return false
	}

	for i, cont := range blockTemp.B_content {
		if cont.B_inodo == -1 {
			copy(blockTemp.B_content[i].B_name[:], name)
			blockTemp.B_content[i].B_inodo = sblock.S_firts_ino

			archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueCarpeta{}))*int32(ptr)), 0)
			err = binary.Write(archivo, binary.LittleEndian, &blockTemp)
			if err != nil {
				fmt.Println("Error al escribir el bloque carpetas: ", err)
				return false
			}
			return true
		}
	}

	return false
}

/* func rellenarBloquePtr(ptrBloqueEscribir int, ptrBloqueLeer int, sblock *superBloque, archivo *os.File) bool {

	var blockTemp bloqueApuntadores

	archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueCarpeta{}))*int32(ptrBloqueLeer)), 0)

	binary.Read(archivo, binary.LittleEndian, &blockTemp)

	for i, cont := range blockTemp.b_pointers {
		if cont == -1 {
			blockTemp.b_pointers[i] = int32(ptrBloqueEscribir)

			archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueCarpeta{}))*int32(ptrBloqueLeer)), 0)
			binary.Write(archivo, binary.LittleEndian, &blockTemp)
			return true
		}
	}

	return false
} */

func crearBloqueCarpetas(name string, sblock *superBloque, archivo *os.File) int {
	//ptrblock := sblock.S_first_blo

	var newblock bloqueCarpeta

	copy(newblock.B_content[0].B_name[:], []byte(name))
	newblock.B_content[0].B_inodo = sblock.S_firts_ino
	newblock.B_content[1].B_inodo = -1
	newblock.B_content[2].B_inodo = -1
	newblock.B_content[3].B_inodo = -1

	archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueArchivos{}))*sblock.S_first_blo), 0)
	err := binary.Write(archivo, binary.LittleEndian, &newblock)
	if err != nil {
		fmt.Println("Error al escribir el bloque de carpetas: ", err)
		return -1
	}

	archivo.Seek(int64(sblock.S_bm_block_start+sblock.S_first_blo), 0)
	err = binary.Write(archivo, binary.LittleEndian, &[1]byte{1})
	if err != nil {
		fmt.Println("Error al escribir bitmap Bloques: ", err)
		return -1
	}
	result := sblock.S_first_blo
	sblock.S_first_blo = encontrarBloqueLibre(sblock, archivo)
	sblock.S_free_blocks_count--

	return int(result)
}

func crearBloquePtr1(ptr int, sblock *superBloque, archivo *os.File) int {
	//ptrblock := sblock.S_first_blo

	var newblock bloqueApuntadores

	newblock.B_pointers[0] = int32(ptr)

	for i := 1; i < len(newblock.B_pointers); i++ {
		newblock.B_pointers[i] = -1

	}

	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(bloqueArchivos{}))*sblock.S_first_blo), 0)
	binary.Write(archivo, binary.LittleEndian, &newblock)

	archivo.Seek(int64(sblock.S_bm_block_start+sblock.S_first_blo), 0)
	binary.Write(archivo, binary.LittleEndian, [1]byte{1})

	sblock.S_first_blo++
	sblock.S_free_blocks_count--

	return int(sblock.S_first_blo - 1)
}

func obtenerNumInodo(ruta []string, archivo *os.File, sblock superBloque) int {
	numInodo := 0

	for _, nombreCarpeta := range ruta {
		numInodo = buscarInodo(nombreCarpeta, numInodo, int(sblock.S_inode_start), archivo, int(sblock.S_block_start))
		if numInodo == -1 {
			//no lo encontro, verificar el R para crear la nueva carpeta e numInodo= ptr nueva carpeta
			return numInodo
		}
	}

	return numInodo
}

func buscarInodo(ruta string, numInodo int, inicioBytesInodos int, archivo *os.File, inicioBytesBlock int) int {

	despTemp := inicioBytesInodos + numInodo*(binary.Size(inodo{}))

	var inodoTemp inodo
	var bloqueCarpetaTemp bloqueCarpeta
	archivo.Seek(int64(despTemp), 0)
	err := binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		return -1
	}

	for _, ptr := range inodoTemp.I_block {

		if ptr != -1 {
			despTemp = inicioBytesBlock + int(ptr*int32(binary.Size(bloqueCarpeta{})))
			archivo.Seek(int64(despTemp), 0)
			err = binary.Read(archivo, binary.LittleEndian, &bloqueCarpetaTemp)
			if err != nil {
				fmt.Println("Error al leer el bloqueCarpeta: ", err)
				return -1
			}
			for _, cont := range bloqueCarpetaTemp.B_content {
				if strings.Contains(string(cont.B_name[:]), ruta) {
					return int(cont.B_inodo)
				}
			}

		}
	}
	return -1

}

func EjecRepMBR(id string) {
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

	Dot := "digraph grid {bgcolor=\"slategrey\" label=\" Reporte MBR \"layout=dot "
	Dot += "labelloc = \"t\"edge [weigth=1000 style=dashed color=red4 dir = \"both\" arrowtail=\"open\" arrowhead=\"open\"]"
	Dot += "a0[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">MBR</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">mbr_tamano</TD><TD>" + strconv.Itoa(int(disk.Mbr_tamano)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">mbr_fecha_creacion</TD><TD>" + string(disk.Mbr_fecha_creacion[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">mbr_disk_signature</TD><TD>" + strconv.Itoa(int(disk.Mbr_dsk_signature)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">dsk_fit</TD><TD>" + string(disk.MBR_dsk_fit[:]) + "</TD></TR>\n"
	var ebrTemp EBR
	var despTemp int
	for _, part := range disk.Mbr_partitions {
		if part.Part_type == [1]byte{'e'} {
			name := strings.TrimRight(string(part.Part_name[:]), string(rune(0)))
			Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">Particion</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_status</TD><TD>" + string(part.Part_status[:]) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_type</TD><TD>p</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_fit</TD><TD>" + string(part.Part_fit[:]) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_start</TD><TD>" + strconv.Itoa(int(part.Part_start)) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_s</TD><TD>" + strconv.Itoa(int(part.Part_s)) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_name</TD><TD>" + name + "</TD></TR>\n"
			despTemp = int(part.Part_start)
			archivo.Seek(int64(despTemp), 0)
			err = binary.Read(archivo, binary.LittleEndian, &ebrTemp)
			if err != nil {
				fmt.Println("Error al leer el EBR: ", err)
				return
			}
			for ebrTemp.Part_s > 0 {
				name = strings.TrimRight(string(ebrTemp.Part_name[:]), string(rune(0)))
				Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">Particion Logica</TD></TR>\n"
				Dot += "<TR><TD bgcolor=\"lightgrey\">part_mount</TD><TD>" + string(ebrTemp.Part_mount[:]) + "</TD></TR>\n"

				Dot += "<TR><TD bgcolor=\"lightgrey\">part_fit</TD><TD>" + string(ebrTemp.Part_fit[:]) + "</TD></TR>\n"
				Dot += "<TR><TD bgcolor=\"lightgrey\">part_start</TD><TD>" + strconv.Itoa(int(ebrTemp.Part_start)) + "</TD></TR>\n"
				Dot += "<TR><TD bgcolor=\"lightgrey\">part_s</TD><TD>" + strconv.Itoa(int(ebrTemp.Part_s)) + "</TD></TR>\n"
				Dot += "<TR><TD bgcolor=\"lightgrey\">part_s</TD><TD>" + strconv.Itoa(int(ebrTemp.Part_s)) + "</TD></TR>\n"
				Dot += "<TR><TD bgcolor=\"lightgrey\">part_next</TD><TD>" + strconv.Itoa(int(ebrTemp.Part_next)) + "</TD></TR>\n"
				Dot += "<TR><TD bgcolor=\"lightgrey\">part_name</TD><TD>" + name + "</TD></TR>\n"

				despTemp += int(ebrTemp.Part_s) + 1 + binary.Size(EBR{})
				archivo.Seek(int64(despTemp), 0)
				err = binary.Read(archivo, binary.LittleEndian, &ebrTemp)
				if err != nil {
					fmt.Println("Error al leer el EBR: ", err)
					return
				}

			}

		} else if part.Part_type == [1]byte{'p'} {
			name := strings.TrimRight(string(part.Part_name[:]), string(rune(0)))
			Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">Particion</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_status</TD><TD>" + string(part.Part_status[:]) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_type</TD><TD>p</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_fit</TD><TD>" + string(part.Part_fit[:]) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_start</TD><TD>" + strconv.Itoa(int(part.Part_start)) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_s</TD><TD>" + strconv.Itoa(int(part.Part_s)) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_name</TD><TD>" + name + "</TD></TR>\n"

		}
	}

	Dot += "</TABLE>>];\n}"

	//Crear el archivo .dot
	DotName := "ReporteMbr.dot"
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
	cmd := exec.Command("dot", "-T", "png", DotName, "-o", "ReporteMbr.png")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar la imagen: ", err)
		return
	}

	fmt.Println("Reporte generado con exito")
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
					err = binary.Read(archivo, binary.LittleEndian, &ebr)
					if err != nil {
						fmt.Println("Error al leer el ebr: ", err)
						return
					}
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
						err = binary.Read(archivo, binary.LittleEndian, &ebr)
						if err != nil {
							fmt.Println("Error al leer el ebr: ", err)
							return
						}

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

func EjecRepBmInodes(index int) {
	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	bitmap := ""
	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}
	size := int(sblock.S_bm_block_start - sblock.S_bm_inode_start)
	var temp [1]byte
	archivo.Seek(int64(sblock.S_bm_inode_start), 0)
	//fmt.Println(size)
	for i := 0; i < size; i++ {
		err = binary.Read(archivo, binary.LittleEndian, &temp)
		if err != nil {
			fmt.Println("Error al leer el bitmap inodos: ", err)
			return
		}
		if temp == [1]byte{1} {
			bitmap += "1 "
		} else {
			bitmap += "0 "
		}
		if i%20 == 0 {
			bitmap += "\n"
		}

	}
	archivo.Close()

	bitmapFile, err := os.Create("RepBitmapsInodos.txt")
	if err != nil {
		fmt.Println("Error al crear el archivo txt")
	}
	defer bitmapFile.Close()
	bitmapFile.Write([]byte(bitmap))

	bitmapFile.Close()
	fmt.Println("Reporte generado con exito")

}

func EjecRepBmBloques(index int) {
	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	bitmap := ""
	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}
	size := int(sblock.S_bm_block_start-sblock.S_bm_inode_start) * 3

	var temp [1]byte
	archivo.Seek(int64(sblock.S_bm_block_start), 0)
	for i := 0; i < size; i++ {
		err = binary.Read(archivo, binary.LittleEndian, &temp)
		if err != nil {
			fmt.Println("Error al leer el bitmap bloques: ", err)
			return
		}
		if temp == [1]byte{1} {
			bitmap += "1 "
		} else {
			bitmap += "0 "
		}
		if i%20 == 0 && i != 0 {
			bitmap += "\n"
		}

	}
	archivo.Close()

	bitmapFile, err := os.Create("RepBitmapsBloques.txt")
	if err != nil {
		fmt.Println("Error al crear el archivo txt")
	}
	defer bitmapFile.Close()
	bitmapFile.Write([]byte(bitmap))

	bitmapFile.Close()
	fmt.Println("Reporte generado con exito")

}

func EjecRepSB(index int) {

	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var sblock superBloque
	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	Dot := "digraph grid {bgcolor=\"slategrey\" label=\" Reporte SuperBlock \"layout=dot "
	Dot += "labelloc = \"t\"edge [weigth=1000 style=dashed color=red4 dir = \"both\" arrowtail=\"open\" arrowhead=\"open\"]"
	Dot += "a0[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">SuperBlock</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_filesystem_type</TD><TD>" + strconv.Itoa(int(sblock.S_filesystem_type)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_inodes_count</TD><TD>" + strconv.Itoa(int(sblock.S_inodes_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_blocks_count</TD><TD>" + strconv.Itoa(int(sblock.S_blocks_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_free_blocks_count</TD><TD>" + strconv.Itoa(int(sblock.S_free_blocks_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_free_inodes_count</TD><TD>" + strconv.Itoa(int(sblock.S_free_inodes_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_mtime</TD><TD>" + string(sblock.S_mtime[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_umtime</TD><TD>" + string(sblock.S_umtime[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_mnt_count</TD><TD>" + strconv.Itoa(int(sblock.S_mnt_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_magic</TD><TD>" + strconv.Itoa(int(sblock.S_magic)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_inode_size</TD><TD>" + strconv.Itoa(int(sblock.S_inode_s)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_block_size</TD><TD>" + strconv.Itoa(int(sblock.S_block_s)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_first_ino</TD><TD>" + strconv.Itoa(int(sblock.S_firts_ino)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_first_blo</TD><TD>" + strconv.Itoa(int(sblock.S_first_blo)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_bm_inode_start</TD><TD>" + strconv.Itoa(int(sblock.S_bm_inode_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_bm_block_start</TD><TD>" + strconv.Itoa(int(sblock.S_bm_block_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_inode_start</TD><TD>" + strconv.Itoa(int(sblock.S_inode_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_block_start</TD><TD>" + strconv.Itoa(int(sblock.S_block_start)) + "</TD></TR>\n"
	Dot += "</TABLE>>];\n}"

	archivoDot, err := os.Create("reporteSB.dot")
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

	cmd := exec.Command("dot", "-T", "png", "reporteSB.dot", "-o", "reporteSB.png")

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar la imagen: ", err)
		return
	}

	fmt.Println("Reporte generado con exito")
}

func EjecRepInodes(index int) {

	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var sblock superBloque
	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}
	Dot := "digraph grid {\nbgcolor=\"slategrey\";\n label=\" Reporte Inodos \";\n layout=dot;\n "
	Dot += "labelloc = \"t\"; \n edge [weight=1000 style=dashed color=red4 dir = \"both\" arrowtail=open arrowhead=open];\n"
	var inodoTemp inodo
	archivo.Seek(int64(sblock.S_inode_start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		return
	}

	Dot += "inodo"
	Dot += strconv.Itoa(0)
	Dot += "[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">Inodo " + strconv.Itoa(0) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">I_uid</TD><TD>" + strconv.Itoa(int(inodoTemp.I_uid)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">I_gid</TD><TD>" + strconv.Itoa(int(inodoTemp.I_uid)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">I_s</TD><TD>" + strconv.Itoa(int(inodoTemp.I_s)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">I_atime</TD><TD>" + string(inodoTemp.I_atime[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">I_ctime</TD><TD>" + string(inodoTemp.I_ctime[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">I_mtime</TD><TD>" + string(inodoTemp.I_mtime[:]) + "</TD></TR>\n"
	for i, ptr := range inodoTemp.I_block {
		Dot += "<TR><TD bgcolor=\"lightgrey\">I_block[" + strconv.Itoa(i) + "]</TD><TD>" + strconv.Itoa(int(ptr)) + "</TD></TR>\n"
	}

	Dot += "<TR><TD bgcolor=\"lightgrey\">I_type</TD><TD>" + string(inodoTemp.I_type[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">I_perm</TD><TD>" + string(inodoTemp.I_perm[:]) + "</TD></TR>\n"

	Dot += "</TABLE>>];\n"
	var byteTemp [1]byte
	for i := 1; i < int(sblock.S_inodes_count); i++ {
		archivo.Seek(int64(sblock.S_bm_inode_start+int32(i)), 0)
		err = binary.Read(archivo, binary.LittleEndian, &byteTemp)
		if err != nil {
			fmt.Println("Error al leer el bitmap inodos: ", err)
			return
		}
		if byteTemp == [1]byte{1} {
			archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(inodo{}))*int32(i)), 0)
			err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)

			if err != nil {
				fmt.Println("Error al leer el inodo: ", err)
				return
			}
			Dot += "inodo"
			Dot += strconv.Itoa(i)
			Dot += "[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">Inodo " + strconv.Itoa(i) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">I_uid</TD><TD>" + strconv.Itoa(int(inodoTemp.I_uid)) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">I_gid</TD><TD>" + strconv.Itoa(int(inodoTemp.I_uid)) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">I_s</TD><TD>" + strconv.Itoa(int(inodoTemp.I_s)) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">I_atime</TD><TD>" + string(inodoTemp.I_atime[:]) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">I_ctime</TD><TD>" + string(inodoTemp.I_ctime[:]) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">I_mtime</TD><TD>" + string(inodoTemp.I_mtime[:]) + "</TD></TR>\n"
			for i, ptr := range inodoTemp.I_block {
				Dot += "<TR><TD bgcolor=\"lightgrey\">I_block[" + strconv.Itoa(i) + "]</TD><TD>" + strconv.Itoa(int(ptr)) + "</TD></TR>\n"
			}

			Dot += "<TR><TD bgcolor=\"lightgrey\">I_type</TD><TD>" + string(inodoTemp.I_type[:]) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">I_perm</TD><TD>" + string(inodoTemp.I_perm[:]) + "</TD></TR>\n"

			Dot += "</TABLE>>];\n"
			Dot += "inodo" + strconv.Itoa(i-1) + " -> inodo" + strconv.Itoa(i) + ";\n"
		}

	}
	Dot += "}"
	archivoDot, err := os.Create("reporteInode.dot")
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

	cmd := exec.Command("dot", "-T", "png", "reporteInode.dot", "-o", "reporteInode.png")

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar la imagen: ", err)
		return
	}

	fmt.Println("Reporte generado con exito")
}

func EjecRepBloques(index int) {

	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var sblock superBloque
	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}
	Dot := "digraph grid {\n bgcolor=\"slategrey\";\n label=\" Reporte Bloques \";\n layout=dot;\n "
	Dot += "labelloc = \"t\";\n edge [weight=1000 style=dashed color=red4 dir = \"both\" arrowtail=open arrowhead=open];\n"
	var bCarpeta bloqueCarpeta
	archivo.Seek(int64(sblock.S_block_start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &bCarpeta)
	if err != nil {
		fmt.Println("Error al leer el bloque carpetas: ", err)
		return
	}

	Dot += "bloque"
	Dot += strconv.Itoa(0)
	Dot += "[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">bloque " + strconv.Itoa(0) + "</TD></TR>\n"

	Dot += "<TR><TD bgcolor=\"lightgrey\">b_name</TD><TD>b_inodo</TD></TR>\n"
	for _, cont := range bCarpeta.B_content {
		nam := strings.TrimRight(string(cont.B_name[:]), string(rune(0)))
		Dot += "<TR><TD bgcolor=\"lightgrey\">" + nam + "</TD><TD>" + strconv.Itoa(int(cont.B_inodo)) + "</TD></TR>\n"

	}

	Dot += "</TABLE>>];\n"
	var byteTemp [1]byte
	var bArchivo bloqueArchivos
	for i := 1; i < int(sblock.S_inodes_count); i++ {
		archivo.Seek(int64(sblock.S_bm_block_start+int32(i)), 0)
		err = binary.Read(archivo, binary.LittleEndian, &byteTemp)
		if err != nil {
			fmt.Println("Error al leer el bitmap bloques: ", err)
			return
		}
		if byteTemp == [1]byte{1} {
			archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueArchivos{}))*int32(i)), 0)
			err = binary.Read(archivo, binary.LittleEndian, &bArchivo)

			if err != nil {
				fmt.Println("Error al leer el bloque archivos: ", err)
				return
			}
			Dot += "bloque"
			Dot += strconv.Itoa(i)
			Dot += "[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">bloque " + strconv.Itoa(i) + "</TD></TR>\n"
			cont := strings.TrimRight(string(bArchivo.B_content[:]), string(rune(0)))
			Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">" + cont + "</TD></TR>\n"

			Dot += "</TABLE>>];\n"
			Dot += "bloque" + strconv.Itoa(i-1) + " -> bloque" + strconv.Itoa(i) + ";\n"
		}

	}

	Dot += "}"

	archivoDot, err := os.Create("reporteBloque.dot")
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

	cmd := exec.Command("dot", "-T", "png", "reporteBloque.dot", "-o", "reporteBloque.png")

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar la imagen: ", err)
		return
	}

	fmt.Println("Reporte generado con exito")
}

func EjecRepTree(id string) {

	index := VerificarParticionMontada(id)
	if index == -1 {
		fmt.Println("Id no encontrada")
		return
	}

	//Abrir el disco
	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	//Leer el superbloque
	var sb superBloque
	err = binary.Read(archivo, binary.LittleEndian, &sb)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	//Buscar el inodo raiz
	var raiz inodo
	archivo.Seek(int64(sb.S_inode_start), 0)
	binary.Read(archivo, binary.LittleEndian, &raiz)
	Dot := "digraph H {\n"
	Dot += "node [pad=\"0.5\", nodesep=\"0.5\", ranksep=\"1\"];\n"
	Dot += "node [shape=plaintext];\n"
	Dot += "graph [bb=\"0,0,352,154\"];\n"
	Dot += "rankdir=LR;\n"
	Dot += crearDotNodoTree(0, archivo, sb)
	Dot += "}"

	archivoDot, err := os.Create("reporteTree.dot")
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

	cmd := exec.Command("dot", "-T", "png", "reporteTree.dot", "-o", "reporteTree.png")

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar la imagen: ", err)
		return
	}

	fmt.Println("Reporte generado con exito")

}

func crearDotNodoTree(numInodo int, archivo *os.File, sblock superBloque) string {
	var inodoTemp inodo
	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(inodo{}))*int32(numInodo)), 0)

	err := binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo")
		return ""
	}

	Dot := "inodo" + strconv.Itoa(numInodo) + "[label = <\n"
	Dot += "<TABLE border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
	Dot += "<tr><td bgcolor=\"lightgrey\" colspan=\"2\">Inodo" + strconv.Itoa(numInodo) + "</td></tr>\n"
	Dot += "<tr><td>i_uid</td><td>" + strconv.Itoa(int(inodoTemp.I_uid)) + "</td></tr>\n"
	Dot += "<tr><td>i_gid</td><td>" + strconv.Itoa(int(inodoTemp.I_gid)) + "</td></tr>\n"
	Dot += "<tr><td>i_size</td><td>" + strconv.Itoa(int(inodoTemp.I_s)) + "</td></tr>\n"
	Dot += "<tr><td>i_atime</td><td>" + string(inodoTemp.I_atime[:]) + "</td></tr>\n"
	Dot += "<tr><td>i_ctime</td><td>" + string(inodoTemp.I_ctime[:]) + "</td></tr>\n"
	Dot += "<tr><td>i_mtime</td><td>" + string(inodoTemp.I_mtime[:]) + "</td></tr>\n"
	Dot += "<tr><td>i_type</td><td>" + string(inodoTemp.I_type[:]) + "</td></tr>\n"
	Dot += "<tr><td>i_perm</td><td>" + string(inodoTemp.I_perm[:]) + "</td></tr>\n"
	nodosDot := ""
	enlacesDot := ""
	for i, ptr := range inodoTemp.I_block {
		Dot += "<TR><TD bgcolor=\"lightgrey\">I_block[" + strconv.Itoa(i) + "]</TD><TD port='" + strconv.Itoa(i) + "'>" + strconv.Itoa(int(ptr)) + "</TD></TR>\n"
		if ptr != -1 {

			enlacesDot += "inodo" + strconv.Itoa(numInodo) + ":" + strconv.Itoa(i) + " -> bloque" + strconv.Itoa(int(ptr)) + ";\n"

			nodosDot += crearDotBloqueTree(int(ptr), string(inodoTemp.I_type[:]), archivo, sblock)
		}

	}
	Dot += "</TABLE>>];\n"
	Dot += nodosDot
	Dot += enlacesDot

	return Dot

}

func crearDotBloqueTree(ptr int, tipo string, archivo *os.File, sblock superBloque) string {
	Dot := ""
	if strings.Contains(tipo, "1") {
		var bloqueT bloqueArchivos
		archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueArchivos{}))*int32(ptr)), 0)

		err := binary.Read(archivo, binary.LittleEndian, &bloqueT)

		if err != nil {
			fmt.Println("Error al leer el inodo")
			return ""
		}

		Dot += "bloque"
		Dot += strconv.Itoa(ptr)
		Dot += "[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
		Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">bloque " + strconv.Itoa(ptr) + "</TD></TR>\n"
		cont := strings.TrimRight(string(bloqueT.B_content[:]), string(rune(0)))
		Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">" + cont + "</TD></TR>\n"

		Dot += "</TABLE>>];\n"

	} else {
		var bloqueT bloqueCarpeta
		archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueArchivos{}))*int32(ptr)), 0)

		err := binary.Read(archivo, binary.LittleEndian, &bloqueT)

		if err != nil {
			fmt.Println("Error al leer el inodo")
			return ""
		}

		Dot += "bloque"
		Dot += strconv.Itoa(ptr)
		Dot += "[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
		Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">bloque " + strconv.Itoa(ptr) + "</TD></TR>\n"

		Dot += "<TR><TD bgcolor=\"lightgrey\">b_name</TD><TD>b_inodo</TD></TR>\n"
		enlacesDot := ""
		nodosDot := ""
		for i, cont := range bloqueT.B_content {

			nam := strings.TrimRight(string(cont.B_name[:]), string(rune(0)))
			Dot += "<TR><TD bgcolor=\"lightgrey\">" + nam + "</TD><TD port= '" + strconv.Itoa(i) + "'>" + strconv.Itoa(int(cont.B_inodo)) + "</TD></TR>\n"

			if cont.B_inodo != -1 {
				if nam != "." && nam != ".." {
					nodosDot += crearDotNodoTree(int(cont.B_inodo), archivo, sblock)
					enlacesDot += "bloque" + strconv.Itoa(ptr) + ":" + strconv.Itoa(i) + " -> inodo" + strconv.Itoa(1) + ";\n"
				}
			}

		}
		Dot += "</TABLE>>];\n"
		Dot += nodosDot
		Dot += enlacesDot
	}

	return Dot
}

func EjecRepFile(id string, ruta string) {

	index := VerificarParticionMontada(id)
	if index == -1 {
		fmt.Println("id no encontrado")
		return
	}

	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	var numInodo int
	var inodoTemp inodo
	txt := ""

	if ruta == "/" {
		fmt.Println("No se encontro la ruta del archivo ")
		return
	} else {
		lRuta := strings.Split(ruta[1:], "/")

		numInodo = obtenerNumInodo(lRuta, archivo, sblock)
	}
	//fmt.Println(numInodo)
	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(inodo{}))*int32(numInodo)), 0)
	err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		return
	}
	if numInodo != -1 {
		txt = leerArchivo(inodoTemp, archivo, sblock)
	} else {
		fmt.Println("No se encontro la ruta del archivo ")
		return
	}

	archivo.Close()

	archivoFile, err := os.Create("RepFile.txt")
	if err != nil {
		fmt.Println("Error al crear el archivo txt")
	}
	defer archivoFile.Close()
	archivoFile.Write([]byte(txt))

	archivoFile.Close()
	fmt.Println("Reporte generado con exito")
}

func EjecRepJournaling(id string) {

	index := VerificarParticionMontada(id)

	if index == -1 {
		fmt.Println("id no encontrado")
		return
	}

	archivo, err := os.OpenFile("MIA/P1/"+particionesMontadas[index].LetterValor+".dsk", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}
	if sblock.S_filesystem_type == 2 {
		fmt.Println("No se puede generar el reporte journal a un EXT2")
		return
	}

	var journalTemp Journaling

	archivo.Seek(int64(particionesMontadas[index].Start+int32(binary.Size(superBloque{}))), 0)
	binary.Read(archivo, binary.LittleEndian, &journalTemp)
	Dot := "digraph grid {\nbgcolor=\"slategrey\";\n label=\" Reporte Inodos \";\n layout=dot;\n "
	Dot += "labelloc = \"t\"; \n edge [weight=1000 style=dashed color=red4 dir = \"both\" arrowtail=open arrowhead=open];\n"
	Dot += "a0[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">Operacion</TD><TD>Path</TD><TD>Contenido</TD><TD>Fecha</TD></TR>\n"
	op := ""
	path := ""
	cont := ""
	fecha := ""

	for i := 0; i < int(journalTemp.Size); i++ {
		op = strings.TrimRight(string(journalTemp.Contenido[i].Operation[:]), string(rune(0)))
		path = strings.TrimRight(string(journalTemp.Contenido[i].Path[:]), string(rune(0)))
		cont = strings.TrimRight(string(journalTemp.Contenido[i].Content[:]), string(rune(0)))
		fecha = strings.TrimRight(string(journalTemp.Contenido[i].Date[:]), string(rune(0)))

		Dot += "<TR><TD bgcolor=\"lightgrey\">" + op + "</TD><TD>" + path + "</TD> <TD>" + cont + "</TD> <TD>" + fecha + "</TD></TR>\n"

	}

	Dot += "</TABLE>>];\n}"

}
