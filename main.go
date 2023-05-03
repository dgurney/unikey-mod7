package main

/*
   Copyright (C) 2021 Daniel Gurney
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
	"math/rand"
	"time"

	"github.com/dgurney/unikey/generator"
	"github.com/dgurney/unikey/validator"
)

const version = "0.5.0"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	bench := flag.Int("bench", 0, "Benchmark generation and validation of N*3 keys.")
	cd := flag.Bool("d", false, "Generate a 10-digit CD key.")
	elevencd := flag.Bool("e", false, "Generate an 11-digit CD key.")
	oem := flag.Bool("o", false, "Generate an OEM key.")
	repeat := flag.Int("r", 1, "Generate n keys.")
	t := flag.Bool("t", false, "Show how long the generation took.")
	validate := flag.String("v", "", "Validate a CD or OEM key.")
	ver := flag.Bool("ver", false, "Show version information and exit.")
	Is95 := flag.Bool("95", false, "Apply Windows 95 rules to OEM/CD key validation.")
	flag.Parse()

	if *ver {
		fmt.Printf("unikey-mod7 v%s by Daniel Gurney\n", version)
		return
	}

	if *repeat < 1 {
		*repeat = 1
	}

	if *bench != 0 {
		fmt.Printf("Running key generation/validation benchmark with %d keys of each type...\n", *bench)
		k := generationBenchmark(*bench)
		validationBenchmark(k)
		return
	}

	var started time.Time
	if *t {
		started = time.Now()
	}

	if *validate != "" {
		k := *validate
		var ki validator.KeyValidator

		switch {
		case len(k) == 12 && k[4:5] == "-":
			ki = validator.Mod7ElevenCD{
				First:  k[0:4],
				Second: k[5:12],
			}
		case len(k) == 11 && k[3:4] == "-":
			ki = validator.Mod7CD{
				First:  k[0:3],
				Second: k[4:11],
				Is95:   *Is95,
			}
		case len(k) == 11 && *Is95:
			ki = validator.Mod7CD{
				First:  k[0:3],
				Second: k[4:11],
				Is95:   *Is95,
			}
		case len(k) == 23 && k[5:6] == "-" && k[9:10] == "-" && k[17:18] == "-" && len(k[18:]) == 5:
			ki = validator.Mod7OEM{
				First: k[0:5],
				// nice
				Second: k[6:9],
				Third:  k[10:17],
				Fourth: k[18:],
				Is95:   *Is95,
			}
		default:
			fmt.Println("Could not recognize key type")
			return
		}

		err := ki.Validate()
		switch {
		default:
			fmt.Printf("%s is valid\n", k)
		case err != nil:
			fmt.Printf("%s is invalid: %s\n", k, err)
		}

		return
	}

	if !*cd && !*elevencd && !*oem {
		fmt.Println("You must specify what you want to do! Usage:")
		flag.PrintDefaults()
		return
	}

	oemkey := generator.Mod7OEM{}
	ecdkey := generator.Mod7ElevenCD{}
	cdkey := generator.Mod7CD{}
	var k generator.KeyGenerator
	for i := 0; i < *repeat; i++ {
		switch {
		case *elevencd:
			ecdkey.Generate()
			k = &ecdkey
		case *cd:
			cdkey.Generate()
			k = &cdkey
		case *oem:
			oemkey.Generate()
			k = &oemkey
		}
		fmt.Println(k.String())
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
		switch {
		case *repeat > 1:
			fmt.Printf("Took %s to generate %d keys.\n", ended, *repeat)
			return
		case *repeat == 1:
			fmt.Printf("Took %s to generate %d key.\n", ended, *repeat)
			return
		}
	}
}
