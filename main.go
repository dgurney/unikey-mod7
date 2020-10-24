package main

/*
   Copyright (C) 2020 Daniel Gurney
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/dgurney/unikey/generator"
	"github.com/dgurney/unikey/validator"
)

const version = "0.0.1"

func main() {
	all := flag.Bool("a", false, "Generate all kinds of keys.")
	bench := flag.Int("bench", 0, "Benchmark generation and validation of N*3 keys.")
	cd := flag.Bool("d", false, "Generate a 10-digit CD key.")
	elevencd := flag.Bool("e", false, "Generate an 11-digit CD key.")
	oem := flag.Bool("o", false, "Generate an OEM key.")
	repeat := flag.Int("r", 1, "Generate n keys.")
	t := flag.Bool("t", false, "Show how long the generation took.")
	validate := flag.String("v", "", "Validate a CD or OEM key.")
	ver := flag.Bool("ver", false, "Show version information and exit")
	flag.Parse()

	if *ver {
		fmt.Printf("unikey-mod7 v%s by Daniel Gurney\n", version)
		return
	}

	if *repeat < 1 {
		*repeat = 1
	}

	if *bench != 0 {
		fmt.Printf("Running key generation benchmark with %d keys of each type...\n", *bench)
		generationBenchmark(*bench)
		return
	}

	var started time.Time
	if *t {
		started = time.Now()
	}

	if *validate != "" {
		k := *validate
		var ki validator.KeyValidator
		vch := make(chan bool)

		switch {
		case len(k) == 12 && k[4:5] == "-":
			ki = validator.Mod7ElevenCD{
				Key: *validate,
			}
		case len(k) == 11 && k[3:4] == "-":
			ki = validator.Mod7CD{
				Key: *validate,
			}
		case len(k) == 23 && k[5:6] == "-" && k[9:10] == "-" && k[17:18] == "-" && len(k[18:]) == 5:
			ki = validator.Mod7OEM{
				Key: *validate,
			}
		default:
			fmt.Println("Could not recognize key type")
			return
		}

		go validator.Validate(ki, vch)
		switch {
		case <-vch:
			fmt.Printf("%s is valid.\n", k)
		default:
			fmt.Printf("%s is invalid.\n", k)
		}

		return
	}

	CDKeych := make(chan string, runtime.NumCPU())
	eCDKeych := make(chan string, runtime.NumCPU())
	OEMKeych := make(chan string, runtime.NumCPU())
	if !*all && !*cd && !*elevencd && !*oem {
		fmt.Println("You must specify what you want to do! Usage:")
		flag.PrintDefaults()
		return
	}
	if *elevencd && *oem && *cd {
		*elevencd, *oem, *cd = false, false, false
		*all = true
	}
	// a and key type are mutually exclusive
	if *elevencd && *all || *oem && *all || *cd && *all {
		*all = false
	}

	oemkey := generator.Mod7OEM{}
	ecdkey := generator.Mod7ElevenCD{}
	cdkey := generator.Mod7CD{}
	for i := 0; i < *repeat; i++ {
		if *all {
			go oemkey.Generate(OEMKeych)
			go cdkey.Generate(CDKeych)
			go ecdkey.Generate(eCDKeych)
			fmt.Println(<-OEMKeych)
			fmt.Println(<-CDKeych)
			fmt.Println(<-eCDKeych)
		}
		if *elevencd {
			go ecdkey.Generate(eCDKeych)
			fmt.Println(<-eCDKeych)
		}
		if *cd {
			go cdkey.Generate(CDKeych)
			fmt.Println(<-CDKeych)
		}
		if *oem {
			go oemkey.Generate(OEMKeych)
			fmt.Println(<-OEMKeych)
		}
	}

	if *t {
		var ended time.Duration
		switch {
		case time.Since(started).Round(time.Second) > 1:
			ended = time.Since(started).Round(time.Millisecond)
		default:
			ended = time.Since(started).Round(time.Microsecond)
		}
		if ended < 1 {
			// Oh Windows...
			fmt.Println("Could not display elapsed time correctly :(")
			return
		}
		mult := 0
		switch {
		default:
			switch {
			case *repeat > 1:
				fmt.Printf("Took %s to generate %d keys.\n", ended, *repeat)
				return
			case *repeat == 1:
				fmt.Printf("Took %s to generate %d key.\n", ended, *repeat)
				return
			}
		case *elevencd && *oem || *elevencd && *cd || *oem && *cd:
			mult = 2
		case *all:
			mult = 3
		}
		fmt.Printf("Took %s to generate %d keys.\n", ended, *repeat*mult)
	}
}
