package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type pack struct {
	name      string
	installed string
	updated   string
	updates   int
	removed   string
}

func writePack(f os.File, p pack) {
	f.WriteString(fmt.Sprintf("- Package Name        : %s\n", p.name))
	f.WriteString(fmt.Sprintf("  - Install date      : %s\n", p.installed))
	f.WriteString(fmt.Sprintf("  - Last update date  : %s\n", p.updated))
	f.WriteString(fmt.Sprintf("  - How many updates  : %d\n", p.updates))
	f.WriteString(fmt.Sprintf("  - Removal date      : %s\n", p.removed))
}

func main() {
	fmt.Println("Pacman Log Analyzer")
	if len(os.Args) < 2 {
		fmt.Println("You must send at least one pacman log file to analize")
		fmt.Println("usage: ./pacman_log_analizer <logfile>")
		os.Exit(1)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	packs := make(map[string]*pack)

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
	file.Close()

	instPacks, remPacks, upPacks, currPacks := 0, 0, 0, 0
	for _, eachline := range txtlines {
		s := strings.Split(eachline, " ")
		if len(s) > 4 {
			date := string(eachline[1:17])
			installDate := "-"
			lastUpdate := "-"
			numUpdates := 0
			remDate := "-"
			validLine := true

			name := s[4]
			opType := s[3]

			switch opType {
			case "upgraded":
				upPacks++ // gnrl
				lastUpdate = date
				numUpdates++
			case "installed":
				instPacks++ //gnrl
				installDate = date
			case "removed":
				remPacks++ //gnrl
				remDate = date
			default:
				validLine = false
			}

			if validLine {
				currPacks = instPacks - remPacks //gnrl
				if x, found := packs[name]; found {
					switch opType {
					case "upgraded":
						x.updated = date
						x.updates++
					case "installed":
						x.installed = date
					case "removed":
						x.removed = date
					}
				} else {
					current := pack{name, installDate, lastUpdate, numUpdates, remDate}
					packs[name] = &current
				}
			}

		}

	}

	//writting the final file
	f, err := os.Create("packages_report.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	f.WriteString("Pacman Packages Report\n----------------------")
	f.WriteString(fmt.Sprintf("\n- Installed packages : %d", instPacks))
	f.WriteString(fmt.Sprintf("\n- Removed packages   : %d", remPacks))
	f.WriteString(fmt.Sprintf("\n- Upgraded packages  : %d", upPacks))
	f.WriteString(fmt.Sprintf("\n- Current installed  : %d\n", currPacks))
	f.WriteString("\nList of packages\n----------------\n")

	for _, l := range packs {
		writePack(*f, *l)
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("bytes written successfully")
}
