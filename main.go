package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/leekchan/accounting"
	"github.com/olekukonko/tablewriter"

	"github.com/mndrix/rand"
)

var AvailableCases map[int]float64

func putCashinCases() {
	cashSlices := []float64{1, 5, 10, 15, 25, 50, 75, 100, 200, 300, 400, 500, 750, 1000, 5000, 10000, 25000, 50000, 75000, 100000, 200000, 300000, 400000, 500000, 750000, 1000000}
	caseSlices := []int{}
	for bc := 1; bc <= 26; bc++ {
		caseSlices = append(caseSlices, bc)
	}

	// shuffle the cash

	for i := range cashSlices {
		j := rand.Intn(i + 1)
		cashSlices[i], cashSlices[j] = cashSlices[j], cashSlices[i]
	}

	// create a map of cash and cases

	//fmt.Printf("Shuffled Cash: %v \n", cashSlices)

	for i, d := range cashSlices {
		AvailableCases[caseSlices[i]] = d
	}

}

func removeGuesses(CaseID int) {

	delete(AvailableCases, CaseID)

}

func whatisleft() ([]int, []float64) {

	cl := []int{}
	dl := []float64{}

	for c, d := range AvailableCases {
		cl = append(cl, c)
		dl = append(dl, d)
	}

	sort.Ints(cl)
	sort.Float64s(dl)

	return cl, dl

}

func dealerOffer(round int) (offer float64) {
	// weighted offers, if there are many low cases left the offer will be higher than median to encourage further gambling
	// if there are more high value cases the offer will be higher to encourage players to quit early

	_, dollarsleft := whatisleft()
	round = round - 1

	// figure out what the total dollar value left is
	// figure out what the total dollar value is for the lowest and highest set of values
	var grandTotal float64 = 0
	var lowTotal float64 = 0
	var highTotal float64 = 0
	for i, value := range dollarsleft {
		grandTotal += value
		if i <= len(dollarsleft)/2 {
			lowTotal += value
		}
		if i >= len(dollarsleft)/2 {
			highTotal += value
		}

	}

	// weighing how much the low and high totals matter
	lowFactor := []float64{0.07, 0.09, 0.13, 0.17, 0.2, 0.25, 0.33, 0.5, 0.5, 0.5, 0.5}
	highFactor := []float64{0.11, 0.16, 0.21, 0.26, 0.31, 0.32, 0.33, 0.34, 0.34, 0.34, 0.34}

	offer = (lowTotal * lowFactor[round]) + (highTotal * lowFactor[round] * highFactor[round])

	if offer > HighestOffer {
		HighestOffer = offer
		BestRound = round
	}

	return offer

}

func UserGuess(round int) {
	ac := accounting.Accounting{Symbol: "$", Precision: 0, Thousand: ",", FormatNegative: "%s(%v)"}
	mf := ac.FormatMoney

	requiredSelections := map[int]int{
		1:  5,
		2:  5,
		3:  4,
		4:  3,
		5:  2,
		6:  1,
		7:  1,
		8:  1,
		9:  1,
		10: 1,
		11: 1,
	}

	for selection := 1; selection <= requiredSelections[round]; selection++ {
		casesleft, dollarsleft := whatisleft()
		data := [][]string{}

		for i, v := range dollarsleft {
			data = append(data, []string{strconv.Itoa(casesleft[i]), mf(v)})
		}
		clearscreen()
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Cases Left In Play", "Money Left In Play"})
		table.SetFooter([]string{"Round: " + strconv.Itoa(round), "Required Picks: " + strconv.Itoa(selection) + "/" + strconv.Itoa(requiredSelections[round])})
		table.SetAutoMergeCells(true)
		table.SetRowLine(true)
		table.AppendBulk(data)
		table.Render()

		fmt.Print("\n\n\n")

		fmt.Print("Enter Case Number: ")
		var playerselected int
		fmt.Scanln(&playerselected)
		if _, check := AvailableCases[playerselected]; check == false {
			selection = selection - 1
			fmt.Printf("Sorry %v is not a valid case number \n", playerselected)
			time.Sleep(2 * time.Second)
			continue
		}
		fmt.Printf("Case # %v contained:  %v \n", playerselected, mf(AvailableCases[playerselected]))
		removeGuesses(playerselected)
	}
	if round < 11 {
		fmt.Printf("Next round you will need to pick %v cases", requiredSelections[round+1])
	}
	time.Sleep(2 * time.Second)

}

func simrun() {

	//This is just a test function for me to play without interaction... not required

	ac := accounting.Accounting{Symbol: "$", Precision: 0, Thousand: ",", FormatNegative: "%s(%v)"}
	mf := ac.FormatMoney

	// guessing 21 is the big winner
	exampleGuess := map[int][]int{
		1:  []int{11, 7, 4, 25, 19},
		2:  []int{12, 8, 5, 26, 20},
		3:  []int{10, 6, 3, 24},
		4:  []int{9, 5, 2},
		5:  []int{18, 23},
		6:  []int{1},
		7:  []int{13},
		8:  []int{14},
		9:  []int{15},
		10: []int{16},
		11: []int{17},
	}

	fmt.Printf("\n\nUser is going to run to the end with Case 21.   It's current value is %v ... Starting Sim \n\n", mf(AvailableCases[21]))

	highestOffer := float64(0)
	bestRound := 0
	for round := 1; round <= 11; round++ {

		for _, v := range exampleGuess[round] {
			removeGuesses(v)
		}

		casesleft, dollarsleft := whatisleft()
		//fmt.Printf("Cases Left: %v \n", casesleft)
		//dlstring := "Dollars Left:  "
		data := [][]string{}

		for i, v := range dollarsleft {
			data = append(data, []string{mf(v), strconv.Itoa(casesleft[i])})
			//dlstring = dlstring + mf(v) + " | "
		}
		//fmt.Println(dlstring)go

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Money Left In Play", "Cases Left In Play"})
		table.SetFooter([]string{"Round: " + strconv.Itoa(round), "Required Picks: X "})
		table.SetAutoMergeCells(true)
		table.SetRowLine(true)
		table.AppendBulk(data)
		table.Render()

		offer := dealerOffer(round)
		if offer > highestOffer {
			highestOffer = offer
			bestRound = round
		}

		var grandTotal float64 = 0
		for _, value := range dollarsleft {
			grandTotal += value
		}
		median := grandTotal / float64(len(dollarsleft))

		fmt.Printf("Round %v - Weighted Offer: %v  |  Median: %v   |  M&O Delta: %v  | Final Delta:  %v\n", round, mf(offer), mf(median), mf(offer-median), mf(AvailableCases[21]-offer))

	}

	fmt.Printf("User ended up with %v, the best offer was %v which was offered in Round %v", mf(AvailableCases[21]), mf(highestOffer), bestRound)

}

func clearscreen() {
	// currently windows only
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func dealorno(round int) bool {

	ac := accounting.Accounting{Symbol: "$", Precision: 0, Thousand: ",", FormatNegative: "%s(%v)"}
	mf := ac.FormatMoney
	var accepted bool

	offer := dealerOffer(round)
	casesleft, dollarsleft := whatisleft()
	data := [][]string{}
	for i, v := range dollarsleft {
		data = append(data, []string{strconv.Itoa(casesleft[i]), mf(v)})
	}
	clearscreen()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Cases Left In Play", "Money Left In Play"})
	table.SetFooter([]string{"Round: " + strconv.Itoa(round), "Dealer Offer: " + mf(offer)})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()

	fmt.Print("\n\n\n")
	fmt.Printf("The Dealer just called, he is authorizing me to offer you %v to stop playing.  Do you accept?\n", mf(offer))
	valid := 0
	for valid == 0 {
		validresponse := map[string]bool{
			"yes": true,
			"no":  true,
		}
		var playerselected string
		fmt.Print("YES or NO: ")
		fmt.Scanln(&playerselected)
		if _, check := validresponse[strings.ToLower(playerselected)]; check == false {
			fmt.Printf("Sorry %v is not a valid choice \n", playerselected)
			continue
		}

		if strings.ToLower(playerselected) == "yes" {
			accepted = true
			valid = 1
			fmt.Printf("All Right!  You end this game with %v ... but let's play on and see what you could have gotten \n", mf(offer))
			TakeHome = offer
			ExitRound = round
			time.Sleep(2 * time.Second)
			return accepted

		}

		if strings.ToLower(playerselected) == "no" {
			accepted = false
			valid = 1
			fmt.Println("Play On!")
			time.Sleep(1 * time.Second)
			return accepted

		}

	}
	return accepted
}

var HighestOffer float64 = 0
var BestRound int = 0
var TakeHome float64 = 0
var DealTaken bool = false
var ExitRound int = 0

func main() {

	AvailableCases = make(map[int]float64)
	ac := accounting.Accounting{Symbol: "$", Precision: 0, Thousand: ",", FormatNegative: "%s(%v)"}
	mf := ac.FormatMoney

	putCashinCases()
	//fmt.Printf("Starting Cases: %v \n", AvailableCases)
	//simrun()

	for round := 1; round <= 11; round++ {
		UserGuess(round)
		if DealTaken == true {
			if round < 11 {
				offer := dealerOffer(round)
				fmt.Println("You would have been offered ", mf(offer))
				time.Sleep(1 * time.Second)
			}
		}
		if DealTaken == false {
			DealTaken = dealorno(round)
		}

	}

	_, dollarsleft := whatisleft()
	fmt.Println("The Final Case Contains! ", mf(dollarsleft[0]))

	if DealTaken == true {
		fmt.Printf("You walk out of here with %v in round %v \n", mf(TakeHome), ExitRound)
		fmt.Printf("The best offer was %v which was offered in Round %v", mf(HighestOffer), BestRound)
	}

	if DealTaken == false {
		fmt.Println("Congrats!")
		if HighestOffer > dollarsleft[0] {
			fmt.Printf("The best offer was %v which was offered in Round %v   \n\n\n", mf(HighestOffer), BestRound)
		}
	}

	var playerselected string
	fmt.Print("\n\n\n Press Enter to Exit the Game")
	fmt.Scanln(&playerselected)

}
