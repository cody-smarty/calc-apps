package handlers

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"strconv"
)

type CSVHandler struct {
	logger      *log.Logger
	input       *csv.Reader
	output      *csv.Writer
	calculators map[string]Calculator
}

func NewCSVHandler(logger *log.Logger, input io.Reader, output io.Writer, calculators map[string]Calculator) *CSVHandler {
	return &CSVHandler{
		logger:      logger,
		input:       csv.NewReader(input),
		output:      csv.NewWriter(output),
		calculators: calculators,
	}
}

func (this *CSVHandler) Handle() error {
	defer this.output.Flush()
	var a, b, c int
	var calculator Calculator
	this.input.FieldsPerRecord = 3
	for {
		record, err := this.input.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else if errors.Is(err, csv.ErrFieldCount) {
				this.logger.Printf("%v: %v", err, len(record))
				continue
			} else {
				panic(err)
				// TODO -- handle other csv errors
			}
		}
		if a, err = strconv.Atoi(record[0]); err != nil {
			this.logger.Println("invalid arg:", record[0])
			continue
		}
		var ok bool
		if calculator, ok = this.calculators[record[1]]; !ok {
			this.logger.Println("unsupported operator:", record[1])
			continue
		}
		if b, err = strconv.Atoi(record[2]); err != nil {
			this.logger.Println("invalid arg:", record[2])
			continue
		}
		c = calculator.Calculate(a, b)
		result := append(record, strconv.Itoa(c))
		this.logger.Println(result)
		if err = this.output.Write(result); err != nil {
			panic(err)
			// TODO -- handle Write error
		}
	}
	return this.output.Error()
}
