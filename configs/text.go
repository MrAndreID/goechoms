package configs

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
)

func LoadVersion() {
	figure.NewFigure("MrAndreID", "standard", true).Print()

	fmt.Println("====================================================================== v1.0.3")

	fmt.Println()
}
