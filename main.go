package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
)

type CalcType int

const (
	CalcInvalid CalcType = iota
	CalcAnnual
	CalcDiff
	CalcPrincipal
	CalcPeriod
	CalcPayment
)

var (
	payment, principal, interest float64
	periods                      int
	method                       string
)

func init() {
	flag.Float64Var(&payment, "payment", -1, "The payment amount")
	flag.Float64Var(&principal, "principal", -1, "The loan principal")
	flag.IntVar(&periods, "periods", -1, "The number of months needed to repay the loan")
	flag.Float64Var(&interest, "interest", -1, "The annual interest rate")
	flag.StringVar(&method, "type", "", `The type of payment: "annuity" or "diff"`)
	flag.Parse()
}

func main() {
	action, err := getAction()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	switch action {
	case CalcAnnual:
		err = doAnnualCalculations()
	case CalcDiff:
		err = doDiffCalculations()
	}

	if err != nil {
		fmt.Println(err)
	}
}

func incorrectParameters() error {
	return errors.New("Incorrect parameters")
}

func getAction() (CalcType, error) {
	switch method {
	case "annuity":
		return CalcAnnual, nil
	case "diff":
		return CalcDiff, nil
	default:
		return CalcInvalid, incorrectParameters()
	}
}

func doAnnualCalculations() error {
	action, err := getAnnualAction()
	if err != nil {
		return err
	}

	switch action {
	case CalcPeriod:
		periods = calculatePeriod()
		displayPeriods()
	case CalcPrincipal:
		principal = calculatePrincipal()
		displayPrincipal()
	case CalcPayment:
		payment = calculatePayment()
		displayPayment()
	}

	displayOverpayment()

	return nil
}

func getAnnualAction() (CalcType, error) {
	if interest < 0 {
		return CalcInvalid, incorrectParameters()
	}

	switch true {
	case periods < 0 && principal >= 0 && payment >= 0:
		return CalcPeriod, nil
	case periods >= 0 && principal < 0 && payment >= 0:
		return CalcPrincipal, nil
	case periods >= 0 && principal >= 0 && payment < 0:
		return CalcPayment, nil
	default:
		return CalcInvalid, incorrectParameters()
	}
}

func calculatePeriod() int {
	i := getInterestRate()
	n := math.Log(payment/(payment-i*principal)) / math.Log(1+i)

	return int(math.Ceil(n))
}

func calculatePrincipal() float64 {
	i := getInterestRate()
	ni := math.Pow(1+i, float64(periods))
	p := payment * (ni - 1) / (i * ni)

	return math.Floor(p)
}

func calculatePayment() float64 {
	i := getInterestRate()
	ni := math.Pow(1+i, float64(periods))
	a := principal * i * ni / (ni - 1)

	return math.Ceil(a)
}

func displayPeriods() {
	var dates = make([]string, 0, 2)

	years := periods / 12
	months := periods % 12

	if years > 1 {
		dates = append(dates, fmt.Sprintf("%d years", years))
	} else if years == 1 {
		dates = append(dates, "1 year")
	}

	if months > 1 {
		dates = append(dates, fmt.Sprintf("%d months", months))
	} else if months == 1 {
		dates = append(dates, "1 month")
	}

	fmt.Printf("It will take %s to repay this loan!\n", strings.Join(dates, " and "))
}

func displayPrincipal() {
	fmt.Printf("Your loan principal = %d!\n", int(principal))
}

func displayPayment() {
	fmt.Printf("Your annuity payment = %d!\n", int(payment))
}

func displayOverpayment() {
	fmt.Printf("Overpayment = %d\n", int(math.Ceil(payment*float64(periods))-principal))
}

func doDiffCalculations() error {
	// check input values
	if principal < 0 || interest < 0 || periods < 0 {
		return incorrectParameters()
	}

	var total float64

	// do temporary calculations
	pn := principal / float64(periods)
	i := getInterestRate()

	for m := 1; m <= periods; m++ {
		dp := math.Ceil(pn + i*(principal-pn*(float64(m)-1)))
		total += dp

		fmt.Printf("Month %d: payment is %d\n", m, int(dp))
	}

	fmt.Printf("\nOverpayment = %d\n", int(math.Ceil(total-principal)))

	return nil
}

func getInterestRate() float64 {
	return interest / (12 * 100)
}
