package main 

import (
	"log"
	"time"
	"sort"
	"strconv"

	"encoding/json"

	_ "github.com/lib/pq"
    // "github.com/coopernurse/gorp"
    "github.com/go-gorp/gorp"
)


const (
	DUMMY_AMOUNT_MAX = 100000000
	// DUMMY_FEE_MAX = 10000
	NUM_BUCKETS = 100
)


func handleAmounts(amount_channel chan int, dbmap *gorp.DbMap) {
	buckets := initialiseBuckets()
	hour := -1 // time starting flag

	for amount := range amount_channel {
		now := time.Now() //.Format(time.RFC850)
		addToBuckets(amount, &buckets)

		if hour != -1 && now.Hour() != hour {
			saveBuckets(&buckets, dbmap)
			hour = now.Hour()
		}
		if hour == -1 { hour = now.Hour() }
	}
}


func saveBuckets(buckets *map[int]int, dbmap *gorp.DbMap) {
	log.Println("Saving buckets...")

	// first convert keys to strings
	mstring := make(map[string]int)
	for k, v := range *buckets {
		mstring[strconv.Itoa(k)] = v
	}

	// now reset buckets
	*buckets = initialiseBuckets()


	// first convert to json string
	json, err := json.Marshal(mstring)
	if err != nil {
		log.Printf("Error converting buckets to json string: \n %s\n", err.Error())
		return
	}

	// insert into database
	amounts := Amounts{Json: string(json)}
	err = dbmap.Insert(&amounts)
	if err != nil {
		log.Printf("Error inserting into database: \n %s\n", err.Error())
		return
	}
}

// inefficient but simple
func addToBuckets(amount int, buckets *map[int]int) {
	// sort keys
	var keys []int
    for k := range *buckets {
        keys = append(keys, k)
    }
    sort.Ints(keys)

    for _, v := range keys {
    	if amount < v {
    		(*buckets)[v] = (*buckets)[v] + 1
    		return
    	}
    }

    // if new maximum
    (*buckets)[amount] = 1
}

func initialiseBuckets() map[int]int {
	m := make(map[int]int)

	increment := DUMMY_AMOUNT_MAX / NUM_BUCKETS
	for i := 1; i < NUM_BUCKETS; i++ {
		m[i*increment] = 0
	}

	return m
}