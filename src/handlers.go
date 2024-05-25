package main

import (
	"fmt"
	"strconv"
)

func HandleMulti(args []string, from *int, to *int) {
	if len(args) > 0 && len(args) < 2 {
		tempFrom, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error converting first argument to integer:", err)
			return
		}
		if tempFrom < *from || tempFrom > *to {
			fmt.Println("The first argument must be greater than ", *from, " and less than ", *to)
			return
		} else {
			*from = tempFrom
			return
		}
	} else if len(args) == 2 {
		tempFrom, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error converting first argument to integer:", err)
			return
		}
		if tempFrom < *from || tempFrom > *to {
			fmt.Println("The first argument must be greater than ", *from, " and less than ", *to)
			return
		} else {
			*from = tempFrom
		}
		tempTo, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error converting second argument to integer:", err)
			return
		}
		if tempTo > *to || tempTo < *from {
			fmt.Println("The second argument must be less than ", *to, " and greater than ", *from)
			return
		} else {
			*to = tempTo
			return
		}
	} else if len(args) > 2 {
		fmt.Println("Please provide 2 or less arguments for scraping multiple chapters")
		return
	}
}

func HandleSingle(single int, from int, to int) {
	if single < from || single > to {
		fmt.Println("The chapter must be between ", from, " and ", to)
		fmt.Println("You searched for: ", single)
		return
	}
}
