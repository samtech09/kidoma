package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/samtech09/kidoma/internal/util"
	"github.com/samtech09/kidoma/question"
)

var (
	currentQues question.Question
	config      = question.QConfig{}
	qdata       QuesData
)

type QuesData struct {
	Q            question.Question
	IsCorrect    bool
	TotalCorrect int
	TotalWrong   int
	IsResult     bool
	LastAns      string
	BgImage      string
	Qcounter     int
	SoundFile    string
	GifFile      string
	QAAdd        QAnalysis
	QASubs       QAnalysis
	QADiv        QAnalysis
	QAMul        QAnalysis
	HasMore      bool
	TotalQues    int
}

type QAnalysis struct {
	Qtype   string
	Correct int
	Wrong   int
}

func main() {
	port := os.Getenv("KIDOMA_PORT")
	if port == "" {
		port = "8081"
	}

	initConfig()

	if config.AskNums {
		// ask number of questions to user
		fmt.Print("Additions: ")
		fmt.Scanf("%d", &config.Add)
		fmt.Print("Substraction: ")
		fmt.Scanf("%d", &config.Subs)
		fmt.Print("Multiplication: ")
		fmt.Scanf("%d", &config.Mul)
		fmt.Print("division: ")
		fmt.Scanf("%d", &config.Div)
	}
	qdata.TotalQues = config.TotalQuesToAsk()

	r := mux.NewRouter()
	r.HandleFunc("/", redirectHandler)
	r.HandleFunc("/quiz", quizHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./assets"))))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// run static server
	fmt.Println("Listening on http://localhost:", port)
	server.ListenAndServe()
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "static/index.html", http.StatusPermanentRedirect)
}

func quizHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	if r.Method == "POST" {
		qdata.IsCorrect = false
		qdata.IsResult = true
		currentQues = question.Question{}
		ans := r.FormValue("ans")
		// remove leading zeros if any
		ans = strings.TrimPrefix(ans, "0")
		qdata.LastAns = ans
		if strconv.Itoa(qdata.Q.Ans) == ans {
			qdata.IsCorrect = true
			qdata.TotalCorrect++
			//fmt.Printf("Q: %d %s %d = %d,   UserAns: %s\n", qdata.Q.N1, qdata.Q.Op, qdata.Q.N2, qdata.Q.Ans, ans)
		} else {
			qdata.TotalWrong++
			fmt.Printf("Q-%d: %d %s %d = %d,   UserAns: %s  [wrong]\n", qdata.Qcounter, qdata.Q.N1, qdata.Q.Op, qdata.Q.N2, qdata.Q.Ans, ans)

			if config.AutoIncreaseWrongQues {
				if qdata.Q.Op == "+" {
					config.Add += 1
				} else if qdata.Q.Op == "-" {
					config.Subs += 1
				} else if qdata.Q.Op == "÷" {
					config.Div += 1
				} else if qdata.Q.Op == "x" {
					config.Mul += 1
				}
				qdata.TotalQues += 1
			}
		}

		x := util.GetNum(1, 4)
		if qdata.IsCorrect {
			// get random sucess gif and sound
			qdata.GifFile = "media/s" + strconv.Itoa(x) + ".gif"
			qdata.SoundFile = "media/s" + strconv.Itoa(x) + ".mp3"
		} else {
			// get random wrong gif and sound
			qdata.GifFile = "media/w" + strconv.Itoa(x) + ".gif"
			qdata.SoundFile = "media/w" + strconv.Itoa(x) + ".mp3"
		}

		addQA(qdata.Q.Op, qdata.IsCorrect)

	} else {
		// get new question
		qdata.IsResult = false
		currentQues = question.GetRandomQues(&config)
		qdata.Q = currentQues

		//fmt.Print(currentQues)

		if (qdata.Qcounter % 4) == 0 {
			qdata.BgImage = "img/" + strconv.Itoa(util.GetNum(1, 11)) + ".jpg"
		}
		qdata.HasMore = (config.TotalQuesToAsk() > 0)
		//fmt.Printf("More: %d  ", question.QCounter.Total())
		qdata.Qcounter++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(qdata)
}

func addQA(op string, isCorrect bool) {
	switch op {
	case "+":
		if isCorrect {
			qdata.QAAdd.Correct++
		} else {
			qdata.QAAdd.Wrong++
		}
	case "-":
		if isCorrect {
			qdata.QASubs.Correct++
		} else {
			qdata.QASubs.Wrong++
		}
	case "÷":
		if isCorrect {
			qdata.QADiv.Correct++
		} else {
			qdata.QADiv.Wrong++
		}
	case "x":
		if isCorrect {
			qdata.QAMul.Correct++
		} else {
			qdata.QAMul.Wrong++
		}
	}
}

//Initconfig read config file and initialize AppConfig struct
func initConfig() {
	conffile := "config.json"

	file, err := os.Open(conffile)
	if err != nil {
		log.Fatal("Missing config file.\n", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("Cannot parse config.\n", err)
	}
}

// func getInput(qtype string) bool {
// askinput:
// 	fmt.Printf("%s [y/n] : ", qtype)
// 	// read user input
// 	var userAns string
// 	_, err := fmt.Scanf("%s", &userAns)
// 	if err != nil {
// 		fmt.Println("Invalid input. Try again.")
// 		goto askinput
// 	}
// 	if (userAns == "y") || (userAns == "Y") {
// 		return true
// 	}
// 	return false
// }
