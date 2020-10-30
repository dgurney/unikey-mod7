package main

import (
	"fmt"
	"time"

	"github.com/dgurney/unikey/generator"
	"github.com/dgurney/unikey/validator"
)

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

// generationBenchmark generates amount * 3 keys and times it
func generationBenchmark(amount int) []string {
	oem := generator.Mod7OEM{}
	cd := generator.Mod7CD{}
	ecd := generator.Mod7ElevenCD{}
	kch := make(chan generator.KeyGenerator)
	keys := make([]string, 0)
	started := time.Now()
	for i := 0; i < amount; i++ {
		go generator.Generate(oem, kch)
		k := <-kch
		keys = append(keys, k.String())
	}
	for i := 0; i < amount; i++ {
		go generator.Generate(cd, kch)
		k := <-kch
		keys = append(keys, k.String())
	}
	for i := 0; i < amount; i++ {
		go generator.Generate(ecd, kch)
		k := <-kch
		keys = append(keys, k.String())
	}

	var ended time.Duration
	switch {
	case time.Since(started).Round(time.Second) > 1:
		ended = time.Since(started).Round(time.Millisecond)
	default:
		ended = time.Since(started).Round(time.Microsecond)
	}
	fmt.Printf("Took %s to generate %d keys.\n", ended, len(keys))

	return keys
}

// generationBenchmark validates N keys and times it.
func validationBenchmark(keys []string) {
	vch := make(chan bool)
	var ki validator.KeyValidator
	started := time.Now()
	for _, k := range keys {
		switch {
		case len(k) == 12 && k[4:5] == "-":
			ki = validator.Mod7ElevenCD{
				Key: k,
			}
		case len(k) == 11 && k[3:4] == "-":
			ki = validator.Mod7CD{
				Key: k,
			}
		case len(k) == 23 && k[5:6] == "-" && k[9:10] == "-" && k[17:18] == "-" && len(k[18:]) == 5:
			ki = validator.Mod7OEM{
				Key: k,
			}
		}
		go validator.Validate(ki, vch)
		<-vch
	}

	var ended time.Duration
	switch {
	case time.Since(started).Round(time.Second) > 1:
		ended = time.Since(started).Round(time.Millisecond)
	default:
		ended = time.Since(started).Round(time.Microsecond)
	}

	fmt.Printf("Took %s to validate %d keys.\n", ended, len(keys))
}
