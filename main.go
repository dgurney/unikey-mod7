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
	"time"

	"github.com/dgurney/unikey/generator"
	"github.com/dgurney/unikey/validator"
)

const version = "0.1.3"

func main() {
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
			fmt.Printf("%s is valid\n", k)
		default:
			fmt.Printf("%s is invalid\n", k)
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
	gch := make(chan generator.KeyGenerator)
	for i := 0; i < *repeat; i++ {
		switch {
		case *elevencd:
			go generator.Generate(ecdkey, gch)
		case *cd:
			go generator.Generate(cdkey, gch)
		case *oem:
			go generator.Generate(oemkey, gch)
		}
		k := <-gch
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
