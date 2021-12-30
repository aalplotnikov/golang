package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
	return nil
}

var (
	nameFileFlag     *string
	timeout          *int
	rangeNumbersFlag arrayFlags
)

func init() {
	nameFileFlag = flag.String("file", "text.txt", "name txt file")
	timeout = flag.Int("timeout", 1000, "the program will close after the time has elapsed")
	flag.Var(&rangeNumbersFlag, "range", "range prime numbers")
}

type Writer struct {
	mutex  *sync.Mutex
	Writer *csv.Writer
}

func NewCsvWriter(fileName string) (*Writer, error) {
	File, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(File)
	return &Writer{Writer: w, mutex: &sync.Mutex{}}, nil
}

func (w *Writer) Write(row []string, c1 chan bool) {
	w.mutex.Lock()
	w.Writer.Write(row)
	w.mutex.Unlock()
	c1 <- true
}

func (w *Writer) Flush() {
	w.mutex.Lock()
	w.Writer.Flush()
	w.mutex.Unlock()
}

func main() {
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(*timeout)*time.Millisecond)

	defer cancel()

	w, err := NewCsvWriter(*nameFileFlag)
	if err != nil {
		log.Panic(err)
	}

	c := make(chan string)
	c1 := make(chan bool)

	if len(rangeNumbersFlag) == 0 {
		go rangeNumbers(ctx, 1, 100, c)
		go w.Write([]string{<-c}, c1)
		<-c1
	}

	for i := 0; i < len(rangeNumbersFlag); i++ {
		a, b := rangeToInt(rangeNumbersFlag[i])
		go rangeNumbers(ctx, a, b, c)
	}

	for i := 0; i < len(rangeNumbersFlag); i++ {
		go w.Write([]string{<-c}, c1)
	}

	for i := 0; i < len(rangeNumbersFlag); i++ {
		<-c1
	}
	w.Flush()
}

func rangeNumbers(ctx context.Context, a int, b int, c chan string) {
	prime := "Range " + fmt.Sprint(a) + " : " + fmt.Sprint(b) + "\n"
	isPrime := true
	for i := a; i <= b; i++ {
		for j := 2; j < i; j++ {
			if i%j == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			prime = prime + fmt.Sprint(i) + ", "
		}

		isPrime = true
		select {
		case <-ctx.Done():
			{
				c <- prime
				return
			}
		default:

		}
	}
	c <- prime
}

func rangeToInt(s string) (int, int) {
	str := strings.Split(s, ":")
	if len(str) != 2 {
		return 1, 20
	}

	a, err := strconv.ParseInt(str[0], 10, 32)
	if err != nil {
		a = 1
	}
	b, err := strconv.ParseInt(str[1], 10, 32)
	if err != nil {
		b = 100
	}
	return int(a), int(b)
}
