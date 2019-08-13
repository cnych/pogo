package flags

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	fmt.Println("In Pogo Project")

	var flagInt = flag.Int("flagname", 100, "help message for flagname")
	var flagString = flag.String("flagstr", "abc", "help message for flag string")

	var flagValue int

	flag.IntVar(&flagValue, "flagname2", 1234, "helm message for flagname2")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `pogo version: pogo/1.0.0
Usage: pogo [-h] [-p prefix]

Options:
`)
		flag.PrintDefaults()
	}

	flag.Parse()

	fmt.Println(*flagInt, *flagString, flagValue)

	//flag.Value()

	fmt.Println(flag.Arg(0))
	fmt.Println(flag.Args())
	fmt.Println(flag.NArg())
}
