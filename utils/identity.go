package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var provinceData = map[int]map[int][]int{
	11: {
		1:  {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18},
		2:  {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		3:  {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24},
		4:  {1, 2, 3, 7, 8, 10, 11, 12, 13, 17, 18, 19, 20, 21},
		5:  {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		6:  {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
		7:  {3, 4, 5, 6, 7, 8, 9, 11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 22, 24, 25, 27, 29, 31},
		8:  {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27},
		9:  {1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		10: {1, 2, 4, 6, 9, 10, 11, 12, 13, 14, 16},
		11: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
		12: {1, 2, 3, 4, 5, 6, 7, 8, 9},
		13: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		14: {1, 2, 3, 4, 5, 6, 7, 8, 9},
		15: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		16: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		17: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		18: {1, 2, 3, 4, 5, 6, 7, 8},
		71: {1, 2, 3, 4, 5, 6, 7, 8, 9},
		72: {1, 2},
		73: {1, 2, 3, 4},
		74: {1, 2, 3, 4, 5},
		75: {1, 2, 3, 4, 5},
	},
	12: {
		1:  {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		2:  {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		3:  {1, 2, 3, 4, 5, 6, 7, 14, 20, 21, 22, 29, 30, 31},
		4:  {5, 6, 10, 11, 20, 21, 27, 28, 29, 35},
		5:  {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
		6:  {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
		7:  {1, 2, 3, 4, 5, 6, 7, 8, 9, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 31, 32, 33},
		8:  {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
		9:  {8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
		10: {1, 2, 7, 8, 9, 14, 18, 19, 20},
		11: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		12: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 19, 20, 21, 22, 23, 24},
		13: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
		14: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 25, 26, 27, 28, 29, 30},
		15: {1, 2, 3, 4, 5, 6, 7, 8},
		16: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		17: {1, 2, 3, 4, 5, 6, 7, 8, 9},
		18: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
		19: {1, 2, 3, 4, 5, 6, 7},
		20: {1, 2, 3, 4, 5, 6, 7, 8, 9},
		21: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		22: {1, 2, 3, 4, 5},
		23: {1, 2, 3, 4, 5, 6, 7, 8},
		24: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		25: {1, 2, 3, 4, 5, 6, 7, 8},
		71: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21},
		72: {1, 2, 3, 4, 5, 6, 7, 8},
		73: {1, 2, 3, 4},
		74: {1, 2, 3, 4, 5, 6},
		75: {1, 2, 3, 4, 5},
		76: {1, 2, 3, 4, 5},
		77: {1, 2, 3, 4, 5, 6},
	},
}

type OrdererIdentity struct {
	Title   string
	Name    string
	Id      string
	Phone   string
	Email   string
	Address string
}

func fakeNik() string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	// Randomly select province
	provinsiKeys := make([]int, 0, len(provinceData))
	for k := range provinceData {
		provinsiKeys = append(provinsiKeys, k)
	}
	provinsiCode := provinsiKeys[r.Intn(len(provinsiKeys))]
	kabupatenKotaMap := provinceData[provinsiCode]

	// Randomly select kabupaten/kota
	kabupatenKeys := make([]int, 0, len(kabupatenKotaMap))
	for k := range kabupatenKotaMap {
		kabupatenKeys = append(kabupatenKeys, k)
	}
	kabupatenCode := kabupatenKeys[r.Intn(len(kabupatenKeys))]
	kecamatanList := kabupatenKotaMap[kabupatenCode]

	// Randomly select kecamatan
	kecamatanCode := kecamatanList[r.Intn(len(kecamatanList))]

	// Randomly select date of birth
	year := r.Intn(2007-1990+1) + 1990
	month := r.Intn(12) + 1
	day := r.Intn(28) + 1 // Simplified to 28 to avoid month/leap year complexities

	// Randomly select gender and adjust day if female
	gender := r.Intn(2) // 0 for male, 1 for female
	if gender == 1 {
		day += 40
	}

	// Generate random sequence
	sequence := r.Intn(9999-1+1) + 1

	nik := fmt.Sprintf("%02d%02d%02d%02d%02d%s%04d",
		provinsiCode,
		kabupatenCode,
		kecamatanCode,
		day,
		month,
		strconv.Itoa(year)[2:],
		sequence,
	)

	return nik
}

func FakeIdentity() (identity OrdererIdentity) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	prefixes := []string{
		"0812", "0813", "0814", "0815", "0816", "0851", "0852",
		"0821", "0822", "0823", "0831", "0832", "0833", "0838",
	}
	identity.Phone = fmt.Sprintf("%s%d", prefixes[r.Intn(len(prefixes))], r.Intn(90000000)+10000000)

	switch r.Intn(2) {
	case 1:
		identity.Title = "Mr."
		identity.Name = MaleFirstNames[r.Intn(len(MaleFirstNames)-1)]
		identity.Name += " " + MaleLastNames[r.Intn(len(MaleLastNames)-1)]
	default:
		identity.Title = "Mrs."
		identity.Name = FemaleFirstNames[r.Intn(len(FemaleFirstNames)-1)]
		identity.Name += " " + FemaleLastNames[r.Intn(len(FemaleLastNames)-1)]
	}

	identity.Id = fakeNik()
	identity.Email = RandString(r.Intn(10)+5, true, false, true, false) + "@gmail.com"

	identity.Address = "Jl. " + RandString(r.Intn(15)+5, true, false, false, false)
	return
}

func increment(id, max int) int {
	if id >= max {
		return 1
	}
	return id + 1
}

func NikRotator(nik string) string {
	id, _ := strconv.Atoi(nik[len(nik)-2:])
	return fmt.Sprintf("%s%02d", nik[:len(nik)-2], increment(id, 99))
}
