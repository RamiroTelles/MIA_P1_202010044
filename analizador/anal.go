package analizador

import (
	"fmt"
	"regexp"
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

func Analizar(comando string) {

	fmt.Println(comando)

	analLex := regexp.MustCompile("[A-Za-z]+((\\s)*(-[A-Za-z]*=.*))*")
	encontrado := analLex.FindAllString(comando, -1)

	fmt.Println(encontrado)

}
