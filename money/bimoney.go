package money

import (
	"fmt"
	"strconv"
	"strings"
)

// Bimoney is an integer representation of money on Binance with up to 8 levels of scale/precision.
type Bimoney int64

func (m Bimoney) PercentageChange(newVal Bimoney) Bimoney {
	if m == 0 {
		return Bimoney(10000000000)
	}

	if strings.HasPrefix(((newVal - m) * 100000000).FormatBimoney(false), "-") {
		diff := newVal - m
		if diff == 0 {
			return Bimoney(0)
		}

		interim := m / 10000
		if interim == 0 {
			return Bimoney(0)
		}

		return (diff / interim) * 1000000
	} else {
		diff := (newVal - m) * 100000000
		if diff == 0 {
			return Bimoney(0)
		}

		interim := diff / m
		return interim * 100
	}
}

func (m Bimoney) AmountFromPercentage(percentage Bimoney) Bimoney {
	if percentage == 0 {
		return Bimoney(0)
	}

	interim := m * (percentage / 100000000)
	if interim == 0 {
		return Bimoney(0)
	}

	return interim / 100
}

func (m Bimoney) PortionOf(value Bimoney) Bimoney {
	if m == 0 || value == 0 {
		return Bimoney(0)
	}

	return m / (value / 10000) * 10000
}

// ParseBimoney converts a string representation of an integer or a float to Bimoney.
// Only dots '.' are considered valid decimal place delimiters, commas ',' are not accepted.
func ParseBimoney(strAmont string) (Bimoney, error) {
	if _, err := strconv.ParseFloat(strAmont, 64); err != nil {
		return 0, fmt.Errorf("cannot convert input string to Bimoney; %w", err)
	}

	if strings.Contains(strAmont, ",") {
		return 0, fmt.Errorf("cannot convert input string containing ',' to Bimoney")
	}

	if !strings.Contains(strAmont, ".") {
		parsedAmount, err := strconv.ParseInt(strAmont, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert input string to Bimoney; %w", err)
		}
		return Bimoney(parsedAmount * 100000000), nil
	}

	pointSplit := strings.Split(strAmont, ".")

	if len(pointSplit) != 2 {
		return 0, fmt.Errorf("cannot convert input string containing multiple '.' to Bimoney")
	}

	if strings.HasSuffix(pointSplit[1], "0") {
		for strings.HasSuffix(pointSplit[1], "0") {
			pointSplit[1] = strings.TrimSuffix(pointSplit[1], "0")
		}

		if pointSplit[1] == "" {
			strAmont = pointSplit[0]
		} else {
			strAmont = strings.Join(pointSplit, ".")
		}
	}

	if !strings.Contains(strAmont, ".") {
		return ParseBimoney(strAmont)
	}

	// prevent removing leading zeros when parsing to int
	var discountLeadingPrefix bool
	if strings.HasPrefix(pointSplit[1], "0") {
		pointSplit[1] = fmt.Sprintf("1%s", pointSplit[1])
		discountLeadingPrefix = true
	}

	afterPointInt, err := strconv.ParseInt(pointSplit[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot convert input string to Bimoney; %w", err)
	}

	// count the number of digits after the decimal point
	digitCounter := 0
	afterPointIntCounter := afterPointInt
	for afterPointIntCounter != 0 {
		afterPointIntCounter /= 10
		digitCounter++
	}

	if discountLeadingPrefix {
		digitCounter -= 1
	}

	if digitCounter == 8 {
		microInt, err := strconv.ParseInt(strings.Replace(strAmont, ".", "", 1), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert input string to Bimoney; %w", err)
		}

		return Bimoney(microInt), nil
	}

	if digitCounter > 8 {
		return 0, fmt.Errorf("cannot convert input string to Bimoney, max 8 places of precision are supported, input requires %d", digitCounter)
	} else {
		for 8-digitCounter != 0 {
			strAmont += "0"
			digitCounter++
		}
	}

	microInt, err := strconv.ParseInt(strings.Replace(strAmont, ".", "", 1), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot convert input string to Bimoney; %w", err)
	}

	return Bimoney(microInt), nil
}

func (b Bimoney) FormatBimoney(fitToLot bool) string {
	bimoneyStr := strconv.FormatInt(int64(b), 10)

	var negative bool
	if strings.HasPrefix(bimoneyStr, "-") {
		bimoneyStr = strings.TrimPrefix(bimoneyStr, "-")
		negative = true
	}

	if len(bimoneyStr) < 8 {
		for len(bimoneyStr) < 8 {
			bimoneyStr = fmt.Sprintf("0%s", bimoneyStr)
		}
		bimoneyStr = fmt.Sprintf("0%s", bimoneyStr)
	} else if len(bimoneyStr) == 8 {
		bimoneyStr = fmt.Sprintf("0%s", bimoneyStr)
	}

	bimoneyFmt := fmt.Sprintf("%s.%s", bimoneyStr[:len(bimoneyStr)-8], bimoneyStr[len(bimoneyStr)-8:])

	if fitToLot {
		pointSplit := strings.Split(bimoneyFmt, ".")
		roundDirectionN, err := strconv.Atoi(string(pointSplit[1][3]))
		if err != nil {
			return bimoneyFmt
		}

		if roundDirectionN >= 5 {
			roundN, err := strconv.Atoi(string(pointSplit[1][2]))
			if err != nil {
				return bimoneyFmt
			}

			roundN++

			bimoneyFmt = fmt.Sprintf("%s.%s", pointSplit[0], fmt.Sprintf("%s%d00000", pointSplit[1][:2], roundN))
		} else if roundDirectionN <= 4 {
			roundN, err := strconv.Atoi(string(pointSplit[1][2]))
			if err != nil {
				return bimoneyFmt
			}

			bimoneyFmt = fmt.Sprintf("%s.%s", pointSplit[0], fmt.Sprintf("%s%d00000", pointSplit[1][:2], roundN))
		}
	}

	if negative {
		bimoneyFmt = fmt.Sprintf("-%s", bimoneyFmt)
	}

	return bimoneyFmt
}
