package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
	"unicode"

	"github.com/fatih/color"
)

const LineSep = "_____________________________________________\n"

const Banner = `
 █▄▀ ▄▀█ █   █▄▄ █▀█ ▀█▀
 █ █ █▀█ █   █▄█ █▄█  █

 Code by: https://github.com/Nux-xader
 Version: 5.1.0
 ` + LineSep

const (
	lowerAbjad = "abcdefghijklmnopqrstuvwxyz"
	upperAbjad = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var MaleFirstNames = []string{
	"hafizhan", "gandhi", "aldio", "benny", "vicky", "jufianto",
	"abdur", "abdurrahman", "bakti", "daniel", "dayu", "edi",
	"fadil", "fahmi", "fairuzi", "gustian", "hermawan",
	"ibnuyohanzah", "adnil", "nurgivo", "rahmadi",
	"rahmat", "said", "tarikhul", "vido", "wahyu", "aditya",
	"afrian", "debby", "deny", "edmund", "fajar", "fathiya",
	"fauzar", "habzer", "ikbal", "ikhsan", "jayus", "narendra",
	"naufal", "fadana", "rangga", "rangga", "saiful", "taufik",
	"teddy", "vigo", "wahyu", "yofaldi", "agus", "agustiando",
	"aidil", "alfajri", "bayu", "desnando", "desri", "ilham",
	"ilham", "jukhri", "mardiyyat", "mukhtar", "nazarudin",
	"rano", "raynaldi", "teresno", "wendi", "zukri", "alfi",
	"andre", "andre", "benni", "deswanto", "ferry", "firdaus",
	"hamdani", "hijrah", "indra", "indri", "kemal",
	"miftahur", "nofan", "nurudin", "rudi", "thovanni", "tian",
	"yudiatma", "zukri", "andre", "andrianto", "andryan", "angga",
	"firdaus", "benny", "harika", "khairul", "reza", "rudi",
	"tio", "boby", "dias", "dicky", "dika", "husnul", "irpandi",
	"kurniawan", "padli", "rian", "riandi", "rianto", "ruwadi",
	"sugeng", "tomi", "usthalay", "winggo", "yunaldi", "arif",
	"arie", "asri", "astri", "boby", "dimas", "doni", "hendri",
	"irwan", "isa", "novri", "pajar", "ridho", "ridwan", "supriadi",
	"tommy", "aulia", "agung", "ahmad", "arie", "azwar", "azwir",
	"budi", "dwiki", "dwiza", "frans", "fuad", "herman", "ichsyan",
	"pikril", "rikal", "riyan", "rizki", "syaiful", "sumirah",
	"ade", "kasnanto", "teguh", "sakiyo", "dwi", "laily", "sigit",
	"linda", "untung", "sarminah", "wiwik", "saputra",
}

var MaleLastNames = []string{
	"basyir", "setiawan", "kusuma", "indri", "syah", "hanafi", "anggara",
	"aziz", "setyono", "thoha", "sutedi", "maulana", "amin", "purba",
	"saputri", "aziz", "raharjo", "fahrizal", "gazalba", "pasha", "adhy",
	"sahputra", "triono", "fiyandana", "tona", "santoso", "susanto",
	"baskoro", "marhamah", "shodiq", "rizki", "gunti", "ade", "rizqi",
	"firdaus", "syahputra", "rahman", "zuhdi", "novrian", "putra",
	"fauzi", "ardani", "ahmad", "dava", "ali", "suyoto", "zainuri",
	"tiono", "wijaya", "sumono", "alwi", "arnolis", "suryawan", "djovanka",
	"suhada", "mufarokah", "aida", "gustriansyah", "supartono", "supangat",
	"dio", "widodo", "sudibyo", "sujono", "pratama", "andriano", "mustopo",
	"saleh", "syuhada", "doli", "kasmini", "rifki", "ari", "fauzi", "nurrohim",
	"asnawi", "zakaria", "farhan", "hasan", "arif", "junaidi", "shidqi",
	"apsyarin", "rahmat", "surya", "jalil", "kuswanto", "malik", "hafitz",
	"ganda", "dian", "azmeer", "ramadhan", "mahendra", "efendi", "hasan",
}

var FemaleFirstNames = []string{
	"yulia", "afni", "lisa", "sumilah", "listiani", "anindita", "shallyn",
	"anggi", "atik", "ngatini", "giarsih", "ririn", "alisa", "aprillia",
	"sadra", "lilik", "baniyah", "olivvia", "febi", "neni", "nailul",
	"sehliya", "rita", "nurhayati", "asmunah", "aqra", "endah", "pagiriana",
	"damiyem", "yona", "misiyani", "masinah", "nuning", "mardiana", "cia",
	"sugini", "dwi", "suprihatin", "bonem", "azizah", "rohana", "reno",
	"mawarni", "sumidah", "meilany", "kamiah", "fauziyah", "misti", "diah",
	"resi", "eni", "indriani", "rahmadani", "sadiah", "ulfi", "fitra",
	"anin", "kusningsih", "rindi", "shecylia", "nuriatik", "umi", "vega",
	"lailatul", "amalia", "mutiara", "jaemah", "katmini", "gisyella", "tumi",
	"sutinem", "isti", "ervi", "evi", "leginem", "susi", "leny", "ratmini",
	"elly", "refi", "aprilia", "ngadini", "anisya", "erna", "azza", "adella",
	"era", "nadhirotul", "mukaromah", "wartinah", "rusmiyati", "vina", "tukini",
	"fatmawati", "surni", "qomariyah", "nurdiana", "imania", "puji", "lely",
	"khofia", "meila", "mesinem", "sukini", "bintiah", "mutingah", "aulia",
	"hatini", "donna", "amal", "samiyah", "sitti", "purwati", "surti",
	"lailatun", "aan", "fairi", "suyatun", "maslina", "atun", "agustina",
	"zakiah", "viki", "resa", "mira", "suwarsih", "suriyani", "feny", "felly",
	"shela", "khusnul", "yovita", "nunit", "aliyatul", "sugeng", "munawaroh",
	"kokom", "sudarmi", "salwa", "mahara", "indriyani", "sunarti", "anik",
	"yatini", "ilda", "tumirah", "naimatul", "legiyem", "rubini", "tria",
	"dea", "nurisya", "rozakiah", "madha", "ike", "saskia", "sulistiyani",
	"nora", "yuliani", "nanda", "sulastri", "mahera", "sutinah", "mariyamah",
	"ida", "sri", "intan", "dahlia", "ayu", "retno", "lolita", "sukati",
	"syafa", "nuriyati", "sutarsih", "wiji", "suminten", "neza", "yuli",
	"rati", "aam", "ika", "susanti", "suyani", "dean", "sundari", "ratu",
	"neneng", "hadmini", "rubiyati", "elisa", "indah", "handani", "hadista",
	"murjiatun", "apri", "misbahatussudury", "khuzaimah", "yance",
	"elisnawati", "een", "wahyuni", "inggar", "binar", "ngatina", "riya",
	"supiyah", "siti", "fitria", "juminah", "nuriah", "welas", "azzah",
	"lentati", "nuraini", "revianita", "miswen", "saminah", "alfazola",
	"wulan", "shafa", "halima", "silvia", "fauziah", "ganti", "meli",
	"nurfatrianti", "wiwit", "tessa", "nikawati", "ayla",
}

var FemaleLastNames = []string{
	"rahayu", "agustia", "dwi", "yulia", "hapti", "juniah", "listiani", "meldiani",
	"nasukah", "liyanti", "rahma", "arta", "faridhotul", "andriyani", "nurhayati",
	"dzakiroh", "mulyani", "dhita", "permata", "syifa", "rejeki", "nisa", "sriyati",
	"dinata", "naviah", "julita", "murtini", "azizah", "puji", "saputri",
	"oktavianti", "mawarni", "aulia", "nabilla", "khamidah", "leli", "irul", "mutia",
	"pertiwi", "nida", "rahmadani", "hanifah", "aidah", "masruroh", "witdiawati",
	"zaenab", "yati", "anisa", "nailatul", "puspita", "umayah", "salamah", "badrina",
	"cita", "hayati", "atika", "zulaiha", "amelia", "alifatul", "cantika",
	"meisyabila", "ariyani", "kharisma", "yuli", "masithoh", "madona", "aprilia",
	"yuliastuti", "untari", "mulyati", "rismiani", "cahaya", "fatmawati", "anjani",
	"novitasari", "saritsha", "widya", "miswardi", "dwi", "azzahroh", "ambarwati",
	"nusaybah", "daryani", "khayati", "savitri", "naila", "zulaikha", "cahyati",
	"ain", "ulva", "kotimah", "eria", "purwati", "azijah", "aan", "nurhidayati",
	"ayu", "atun", "elfina", "dah", "wulan", "royani", "hotijah", "fauziah",
	"nada", "lestari", "pauzia", "eviyana", "rohmah", "rahmawati", "khusniah",
	"indrayani", "maslihatun", "indriyani", "sunarti", "mariani", "izzati",
	"yani", "fitri", "fitry", "gustiana", "fatonah", "wardani", "yuliani",
	"afriani", "sulastri", "romlah", "wahyuti", "hanifatur", "desi", "widia",
	"fatimah", "aeni", "febriyanti", "latifah", "devy", "sabila", "uswatun",
	"mubarokah", "tri", "nurviani", "astuti", "susanti", "afitri", "marpoah",
	"kholifah", "hanzuneira", "afrianti", "rafiani", "sundari", "artika", "naca",
}

func RandString(n int, lower, upper, digit, sym bool) string {
	var chars = ""
	if lower {
		chars += lowerAbjad
	}
	if upper {
		chars += upperAbjad
	}
	if digit {
		chars += "0123456789"
	}
	if sym {
		chars += "!@#$%^&*()-_=+[]{}|;:,.<>?/~"
	}

	if len(chars) == 0 {
		chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	}

	var letterRunes = []rune(chars)
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func OriPrint(msg string) {
	color.White(fmt.Sprintf(" [ ] github.com/Nux-xader [%v] %v", time.Now().Format("02 15:04:05"), msg))
}

func WarningPrint(msg string) {
	color.Yellow(fmt.Sprintf(" [#] github.com/Nux-xader [%v] %v", time.Now().Format("02 15:04:05"), msg))
}

func DangerPrint(msg string) {
	color.Red(fmt.Sprintf(" [!] github.com/Nux-xader [%v] %v", time.Now().Format("02 15:04:05"), msg))
}

func SuccessPrint(msg string) {
	color.Green(fmt.Sprintf(" [+] [%v] %v", time.Now().Format("02 15:04:05"), msg))
}

func IsValidName(name string) bool {
	re := regexp.MustCompile(`^[a-zA-Z ]+$`)
	return re.MatchString(name)
}

func IsValidNIK(nik string) bool {
	re := regexp.MustCompile(`^\d{16}$`)
	return re.MatchString(nik)
}

func IsDigit(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}
