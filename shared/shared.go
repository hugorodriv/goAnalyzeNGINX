package shared

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var databaseLoc = "./data/dbip-country-lite-2024-04.csv"

type IpDatabase struct {
	arr   []IpRange
	count int
}
type IpRange struct {
	begin   uint32
	end     uint32
	country string
}

func IpToInt(ip string) uint32 {
	split := strings.Split(ip, ".")
	if len(split) != 4 {
		fmt.Println("Error: Invalid IP")
		fmt.Println("IP: ", ip)
		return 0
	}
	part1, _ := strconv.Atoi(split[0])
	part2, _ := strconv.Atoi(split[1])
	part3, _ := strconv.Atoi(split[2])
	part4, _ := strconv.Atoi(split[3])

	return uint32((part1 << 24) + (part2 << 16) + (part3 << 8) + part4)
}

func FindCountry(ip string, database IpDatabase) string {
	ipInt := IpToInt(ip)

	// binary search
	l := 0
	r := database.count - 1
	for l <= r {
		m := l + (r-l)/2
		if database.arr[m].begin <= ipInt && database.arr[m].end >= ipInt {
			return database.arr[m].country
		} else if database.arr[m].begin > ipInt {
			r = m - 1
		} else {
			l = m + 1
		}
	}
	return "Unknown"
}

func ParseDatabase() (IpDatabase, error) {
	file, err := os.Open(databaseLoc)
	if err != nil {
		fmt.Println("Error: ", err)
		return IpDatabase{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		fmt.Println("Error: ", err)
		return IpDatabase{}, err
	}

	var tempArr []IpRange
	var count int

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		if strings.Contains(line[0], ":") {
			break
		}
		begin := IpToInt(line[0])
		end := IpToInt(line[1])
		country := line[2]

		tempIpR := IpRange{begin, end, country}
		tempArr = append(tempArr, tempIpR)
		count++
	}
	return IpDatabase{tempArr, count}, nil
}
