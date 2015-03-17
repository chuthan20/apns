package apns

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	token, err := hex.DecodeString("F389410AE1B57972DBBF6EB0C05C2626AB69EDE88F523D7EED49FA6E63A6C266")
	if err != nil {
		t.Fatal(err)
	}
	config, err := LoadConfig("config.json")
	if err != nil {
		t.Fatal(err)
	}
	conn, err := config.Connect()
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	total := 500000
	wg.Add(total)
	start := time.Now()
	for y := 0; y < 2; y++ {
		go func(y int) {
			for i := 0; i < total/2; i++ {
				msg := &Message{Payload: map[string]interface{}{
					"aps": map[string]interface{}{
						"alert": fmt.Sprintf("Test message %d-%d", y+1, i+1),
						"badge": i,
					},
					"time":   time.Now().Format(time.RFC3339Nano),
					"uint32": rand.Uint32(),
					"inf64":  rand.Int63(),
					"float":  rand.Float64(),
				}}
				conn.Send(msg, token)
				// time.Sleep(50 * time.Millisecond)
				// if i%(rand.Intn(500)+1) == 0 {
				// 	time.Sleep(time.Duration(rand.Intn(150)) * time.Millisecond)
				// }
				wg.Done()
			}
		}(y)
	}
	wg.Wait()
	fmt.Println("time", time.Since(start).String(), "to send", total, "messages")
	time.Sleep(1 * time.Second)
	fmt.Println("Count:", conn.counter)
}
