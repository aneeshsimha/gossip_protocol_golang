package gossip

import "os"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	logf, err := os.Create("logs/logfile")
	check(err)
	defer logf.close()

}
