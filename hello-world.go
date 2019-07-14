package main

/*
solver for this kinda game - SKYLINE SAVIOUR / HIMMELSRETTER - https://www.puzzlescript.net/editor.html?hack=d2d091063e73a44796628e2b12d09989 (not final link)
*/

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type punkt struct {
	x int
	y int
}

const (
	vPunkt   uint8 = 1 << 0
	vR       uint8 = 1 << 1
	vL       uint8 = 1 << 2
	vU       uint8 = 1 << 3
	vO       uint8 = 1 << 4
	statisch uint8 = 1 << 5
)

var rahmenzeichnen = [...]string{
	"□", //0x00000
	"╶", //0x00001
	"╴", //0x00010
	"─", //0x00011
	"╷", //0x00100
	"┌", //0x00101
	"┐", //0x00110
	"┬", //0x00111
	"╵", //0x01000
	"└", //0x01001
	"┘", //0x01010
	"┴", //0x01011
	"│", //0x01100
	"├", //0x01101
	"┤", //0x01110
	"┼", //0x01111
	"▩", //0x10000
	"╺", //0x10001
	"╸", //0x10010
	"━", //0x10011
	"╻", //0x10100
	"┏", //0x10101
	"┓", //0x10110
	"┳", //0x10111
	"╹", //0x11000
	"┗", //0x11001
	"┛", //0x11010
	"┻", //0x11011
	"┃", //0x11100
	"┣", //0x11101
	"┫", //0x11110
	"╋"} //0x11111

func randrange(minI int, maxE int) int {
	return rand.Intn(maxE-minI) + minI
}

func randrangeI(minI int, maxI int) int {
	return rand.Intn(maxI+1-minI) + minI
}

func karteZuString(b int, h int, raster [][]uint8) string {
	result := ""
	for j := 0; j < h; j++ {
		for i := 0; i < b; i++ {
			if raster[i][j] == 0 {
				result += "."
			} else {
				v := raster[i][j] >> 1
				result += rahmenzeichnen[v]
				// result += fmt.Sprintf("%x", raster[i][j])
			}
		}
		if j+1 < h {
			result += "\n"
		}
	}
	return result
}

func löscheZeile(b int, h int, raster [][]uint8, y int) {

	for i := 0; i < b; i++ {
		raster[i][y] = 0
	}
	if y < h-1 {
		//lösche nach oben zeigende Pfeilen von der unteren Zeile
		for i := 0; i < b; i++ {
			raster[i][y+1] = raster[i][y+1] &^ vO
		}
	}
	if y > 0 {
		//lösche nach oben zeigende Pfeilen von der unteren Zeile
		for i := 0; i < b; i++ {
			raster[i][y-1] = raster[i][y-1] &^ vU
		}
	}
}

func vollFallenLassen(b int, h int, raster [][]uint8) {
	for {
		if lassFallen(b, h, raster) == false {
			return
		}
	}
}

//gibt true zurück wenn etw gemacht war
func lassFallen(b int, h int, raster [][]uint8) bool {

	// lösche metadaten
	for i := 0; i < b; i++ {
		for j := 0; j < h; j++ {
			raster[i][j] = raster[i][j] & ^statisch
		}
	}

	for i := 0; i < b; i++ {
		if raster[i][h-1] != 0 {
			raster[i][h-1] = raster[i][h-1] | statisch
		}
	}

	result := false
	modified := true
	for modified {
		modified = false
		for i := 0; i < b; i++ {
			for j := h - 1; j >= 0; j-- {
				if raster[i][j] == 0 || raster[i][j]&statisch != 0 {
					continue
				}
				//propagate upwards without problems, horizontally only through links
				machKram := false
				if i > 0 { //guck links
					l := raster[i][j]&vL == vL
					if l && raster[i-1][j]&statisch == statisch {
						machKram = true
					}
				}
				if i < b-1 { //guck rechts
					r := raster[i][j]&vR == vR
					if r && raster[i+1][j]&statisch == statisch {
						machKram = true
					}
				}
				if j > 0 { //guck oben
					o := raster[i][j]&vO == vO
					if o && raster[i][j-1]&statisch == statisch {
						machKram = true
					}
				}
				if j < h-1 { //guck unten (ohe Prüfen)
					o := true //raster[i][j]&vU == vU
					if o && raster[i][j+1]&statisch == statisch {
						machKram = true
					}
				}
				if machKram {
					raster[i][j] = raster[i][j] | statisch
					modified = true
				}
			}
		}
	}

	for i := 0; i < b; i++ {
		for j := h - 1; j >= 1; j-- {
			if raster[i][j]&statisch == statisch {
				continue
			}
			obenesStück := raster[i][j-1]
			if obenesStück&statisch == statisch {
				continue
			}

			if obenesStück != 0 {
				raster[i][j] = obenesStück
				result = true
			} else {
				raster[i][j] = 0
			}
		}
	}

	// erste reihe löschen
	for i := 0; i < b; i++ {
		if (raster[i][0] & statisch) == 0 {
			raster[i][0] = 0
		}
	}

	return result
}

func bauRaster(b int, h int, größe int) [][]uint8 {
	if größe > b*h {
		größe = b * h
	}
	raster := make([][]uint8, b)
	for i := range raster {
		raster[i] = make([]uint8, h)
	}

	mitte := b / 2

	raster[mitte][h-1] = 1

	punkteZahl := 1

	minx := mitte - 1
	maxx := mitte + 1
	miny := h - 2

	for punkteZahl < größe {
		var rx int
		if randrangeI(0, 1) == 0 {
			rx = randrangeI(minx, mitte)
		} else {
			rx = randrangeI(mitte, maxx)
		}
		ry := randrangeI(miny, h-1)

		if raster[rx][ry] != 0 {
			continue
		}

		// if ry == h-1 {
		// 	//gut
		// } else
		{
			if (ry > 0 && raster[rx][ry-1] != 0) ||
				(ry < h-1 && raster[rx][ry+1] != 0) ||
				(rx > 0 && raster[rx-1][ry] != 0) ||
				(rx < b-1 && raster[rx+1][ry] != 0) {
				//gut
			} else {
				continue
			}
		}

		raster[rx][ry] = 1
		punkteZahl++
		if rx > 0 && (rx-1) < minx {
			minx = rx - 1
		}
		if rx < (b-1) && (rx+1) > maxx {
			maxx = rx + 1
		}
		if ry > 0 && ry-1 < miny {
			miny = ry - 1
		}

	}

	for i := minx + 1; i <= maxx-1; i++ {
		for j := miny + 1; j <= h-2; j++ {
			if raster[i][j] == 0 {
				continue
			}
			n := raster[i][j-1]
			s := raster[i][j+1]
			e := raster[i+1][j]
			w := raster[i-1][j]
			nw := raster[i-1][j-1]
			ne := raster[i+1][j-1]
			sw := raster[i-1][j+1]
			se := raster[i+1][j+1]

			if n+s+e+w+nw+ne+sw+se >= 7 {
				raster[i][j] = 0
			}

		}

	}

	for i := 0; i < b; i++ {
		for j := 0; j < h; j++ {
			if raster[i][j] == 0 {
				continue
			}
			raster[i][j] = vPunkt

			if i > 0 {
				if raster[i-1][j] != 0 {
					raster[i][j] = raster[i][j] | vL
				}
			}
			if i < b-1 {
				if raster[i+1][j] != 0 {
					raster[i][j] = raster[i][j] | vR
				}
			}
			if j > 0 {
				if raster[i][j-1] != 0 {
					raster[i][j] = raster[i][j] | vO
				}
			}
			if j < h-1 {
				if raster[i][j+1] != 0 {
					raster[i][j] = raster[i][j] | vU
				}
			}
		}
	}

	return raster
}

func obersteZeile(b int, h int, raster [][]uint8) int {
	for j := 0; j < h; j++ {
		for i := 0; i < b; i++ {
			if raster[i][j] != 0 {
				return j
			}
		}
	}
	return h
}

func kopiereRaster(b int, h int, raster [][]uint8) [][]uint8 {

	duplicate := make([][]uint8, b)
	data := make([]uint8, b*h)
	for i := range raster {
		start := i * h
		end := start + h
		duplicate[i] = data[start:end:end]
		copy(duplicate[i], raster[i])
	}
	return duplicate
}

func losRaster(b int, h int, raster [][]uint8, bleibendezüge int, scores []int) {
	gipfel := obersteZeile(b, h, raster)
	for j := gipfel; j < h; j++ {
		kopie := kopiereRaster(b, h, raster)
		löscheZeile(b, h, kopie, j)
		vollFallenLassen(b, h, kopie)

		if bleibendezüge == 1 {
			var h = obersteZeile(b, h, kopie)
			scores[h]++
		} else {
			losRaster(b, h, kopie, bleibendezüge-1, scores)
		}
	}
}

func nichtnullkleinste(ar []int) int {
	min := ar[0]
	for i := 1; i < len(ar); i++ {
		if ar[i] != 0 && ar[i] < min || min == 0 {
			min = ar[i]
		}
	}
	return min
}

func summe(ar []int) int {
	result := ar[0]
	for i := 1; i < len(ar); i++ {
		result += ar[i]
	}
	return result
}

func ratio(ar []int) float64 {
	k := nichtnullkleinste(ar)
	s := summe(ar)

	return float64(k) / float64(s)
}

func rechtste(ar []int) int {
	for i := len(ar) - 1; i >= 0; i-- {
		if ar[i] != 0 {
			return i
		}
	}
	return -1
}

func prettyPrint(b int, h int, raster [][]uint8, explosionen int, himmelhöhe int) string {

	var sb = b + 4

	s := "~"
	for i := 1; i < sb; i++ {
		s += "f"
	}
	s += "\n"

	playersetzt := false

	for j := 0; j < h; j++ {
		s += "f"

		d := h - 2 - j
		if d%2 == 0 && d >= 0 && (d/2) < explosionen {
			s += "b"
		} else {
			s += "i"
		}

		s += "f"
		for i := 0; i < b; i++ {
			if raster[i][j] == 0 {
				if playersetzt == false {
					playersetzt = true
					s += "p"
				} else if himmelhöhe == j {
					s += "g"
				} else {
					s += "."
				}
			} else {
				s += "#"
			}
		}

		s += "f\n"
	}

	for i := 0; i < sb; i++ {
		s += "f"
	}
	s += "\n"

	return s
}

type lösungsstruct struct {
	b           int
	h           int
	raster      [][]uint8
	scores      []int
	explosionen int
	himmelhöhe  int
	ratio       float64
}

func machwas(b int, h int, g int, explosionen int) lösungsstruct {

	raster := bauRaster(b, h, g)

	// fmt.Println()
	// fmt.Println(karteZuString(b, h, raster))
	// fmt.Println()

	scores := make([]int, h+1)
	losRaster(b, h, raster, explosionen, scores)

	r := ratio(scores)
	himmelhöhe := rechtste(scores)

	return lösungsstruct{
		b:           b,
		h:           h,
		raster:      raster,
		scores:      scores,
		himmelhöhe:  himmelhöhe,
		explosionen: explosionen,
		ratio:       r,
	}
}

func prettyPrintVersuch(ls lösungsstruct) {

	fmt.Println()
	fmt.Println(ls.ratio)
	fmt.Println()
	fmt.Println(ls.scores)
	fmt.Println()

	output := prettyPrint(ls.b, ls.h, ls.raster, ls.explosionen, ls.himmelhöhe)

	fmt.Println(output)
}

func main() {

	rand.Seed(time.Now().UnixNano())

	b := flag.Int("Breite", 20, "Gitterbreite")
	h := flag.Int("Höhe", 30, "Gitterhöhe")
	g := flag.Int("Große", 50, "Strukturgröße")
	jederp := flag.Int("jeder", 1000, "jeder")
	explosionen := flag.Int("Explosionen", 3, "Explosionen")
	versuche := flag.Int("Versuche", 100, "Versuchen")
	flag.Parse()

	lösungen := make([]lösungsstruct, 0)

	fmt.Printf("Breite=%d Höhe=%d Große=%d Explosionen=%d Versuche=%d\n", *b, *h, *g, *explosionen, *versuche)

	jeder := *jederp
	var max = *versuche
	for i := 0; i < max; i++ {
		cand := machwas(*b, *h, *g, *explosionen)
		lösungen = append(lösungen, cand)

		if i > 0 && i%jeder == 0 {
			fmt.Printf("processed %d\n", i)
		}
	}

	sort.Slice(lösungen, func(i, j int) bool {
		return lösungen[i].ratio > lösungen[j].ratio
	})

	lmin := len(lösungen) - 6
	lmax := len(lösungen) - 1
	if lmin < 0 {
		lmin = 0
	}

	for i := lmin; i <= lmax; i++ {
		lösung := lösungen[i]
		prettyPrintVersuch(lösung)
	}
}
