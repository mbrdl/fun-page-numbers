package main

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strings"
	"time"
)

// todo vast majority of generated terms are no good, fix by delegating to operators for next value
// mostly due to negations -> eliminate duplicate negatives
// todo parallelize

const DigitString = "4812018"
var baseTerms = []int {0, 1, 2, 5, 15, 51, 188, 731, 2950, 12235, 51822, 223191, 974427, 4302645, 19181100, 86211885, 390248055, 1777495635, 8140539950, 37463689775, 173164232965, 803539474345, 3741930523740, 17481709707825, 81912506777200, 384847173838501, 1812610804416698}
var operators = []*Operator{&opAdd, &opMultiply, &opDivide, &opRaise}
const UseNegation = true  // don't use opSubtract in this case!

func main() {
	start := time.Now()
	termEvalCounter := 0
	numbersGeneratedCounter := 0
	resultMap := make(map[int][]string)
	for result := range Generate(DigitString) {
		termEvalCounter++
		reprs := resultMap[result.value]
		if reprs == nil {
			resultMap[result.value] = make([]string, 0)
			numbersGeneratedCounter++
		}
		resultMap[result.value] = append(reprs, result.repr)
	}
	fmt.Printf("took %v\n", time.Since(start))
	fname := fmt.Sprintf("results/%s-complete-%d.txt", DigitString, time.Now().Unix())
	completeWriter := NewMyWriter(fname)
	compactWriter := NewMyWriter(strings.ReplaceAll(fname, "complete", "compact"))
	smallestNumber, gotSmallest := -1, false
	for key := 1; key <= 1000; key++ {
		terms := resultMap[key]
		if terms == nil {
			if !gotSmallest {
				smallestNumber, gotSmallest = key, true
			}
			terms = []string{}
		}
		WriteResult(completeWriter, key, terms, -1)
		WriteResult(compactWriter, key, terms, 3)
	}
	printer := NewNumberPrinter()
	printer.Print("evaluated %d terms\n", termEvalCounter)
	printer.Print("generated %d different numbers\n", numbersGeneratedCounter)
	printer.Print("smallest number not generated: %d\n", smallestNumber)
}

type NumberPrinter struct {
	printer *message.Printer
}

func NewNumberPrinter() *NumberPrinter {
	return &NumberPrinter{message.NewPrinter(language.English)}
}

func (printer *NumberPrinter) Print(template string, number int) {
	_, err := printer.printer.Printf(template, number)
	if err != nil {
		panic(err)
	}
}
