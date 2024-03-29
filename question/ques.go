package question

import (
	"github.com/samtech09/kidoma/internal/util"
)

type Question struct {
	N1        int
	N2        int
	Ans       int
	Remainder int
	Op        string
}

// type QCount struct {
// 	Add  int
// 	Subs int
// 	Div  int
// 	Mul  int
// 	//Total int
// }

type QConfig struct {
	//firstStart bool
	// templateDirectory string
	// templateFileList  []string
	Add                   int
	Subs                  int
	Div                   int
	Mul                   int
	AskNums               bool // if set to true, will ask number of questions to show for each type
	DigitsAdd1            int
	DigitsAdd2            int
	DigitsSub1            int
	DigitsSub2            int
	DigitsDiv1            int
	DigitsDiv2            int
	DigitsMul1            int
	DigitsMul2            int
	SimpleDivision        bool
	Div1Max               int  //maximum number for divident - just to keep division simple with small limited numbers
	Div2Max               int  //maximum number for divisor - just to avoid division by big numbers like 849/89
	Div2Min               int  //minimum number for divisor - to avoid divide by 1 type question
	Mul2Min               int  //minium number for multiplier - to avoid multiply by 1 type question
	AutoIncreaseWrongQues bool // auto add new question of same type if wrong attempted
}

//var QCounter QCount

func (q *QConfig) TotalQuesToAsk() int {
	return q.Add + q.Div + q.Mul + q.Subs
}

func getMax(digit1, digit2 int) (int, int) {
	if digit2 > digit1 {
		panic("Second operand can not be larger than first")
	}
	max1 := 1
	max2 := 1
	for x := 1; x <= digit1; x++ {
		max1 *= 10
	}
	max1--
	for x := 1; x <= digit2; x++ {
		max2 *= 10
	}
	max2--
	return max1, max2
}

func sumQuiz(digit1, digit2 int) Question {
	min := 3
	max1, max2 := getMax(digit1, digit2)

	q := Question{}
	q.N1 = util.GetNum(min, max1)
	q.N2 = util.GetNum(min, max2)
	q.Op = "+"
	q.Ans = q.N1 + q.N2

	return q
}

func subQuiz(digit1, digit2 int) Question {
	min := 3
	max1, max2 := getMax(digit1, digit2)

	q := Question{}
	q.N2 = util.GetNum(min, max2)
	q.N1 = util.GetNum(q.N2, max1)
	q.Op = "-"
	q.Ans = q.N1 - q.N2

	return q
}

func mulQuiz(digit1, digit2, mul2Min int) Question {
	min := mul2Min
	max1, max2 := getMax(digit1, digit2)

	q := Question{}
	q.N1 = util.GetNum(min, max1)
	q.N2 = util.GetNum(min, max2)
	q.Op = "x"
	q.Ans = q.N1 * q.N2

	return q
}

func divQuiz(digit1, digit2 int, simpleDiv bool, div1Max, div2Max, div2Min int) Question {
	min := div2Min
	max1, max2 := getMax(digit1, digit2)

	if div1Max > 0 {
		max1 = div1Max
	}
	if div2Max > 0 {
		max2 = div2Max
	}

	q := Question{}
	q.N2 = util.GetNum(min, max2)
	q.N1 = util.GetNum(q.N2, max1)
	i := q.N1 % q.N2
	if simpleDiv {
		// for simple division, numeratoor must be fully divisible by denominator
		//   so get operands untile mod is 0
		for i = 1; i > 0 && q.N1 != q.N2; i = (q.N1 % q.N2) {
			q.N2 = util.GetNum(min, max2)
			q.N1 = util.GetNum(q.N2, max1)
			// mood := q.N1 % q.N2
			// fmt.Printf("%d mod %d = %d\n", q.N1, q.N2, mood)
		}
	}

	q.Op = "÷"
	q.Ans = q.N1 / q.N2 // integer division, decimals are truncated
	q.Remainder = q.N1 % q.N2

	return q
}

func GetRandomQues(c *QConfig) Question {
	var q Question

	if c.TotalQuesToAsk() > 0 {
		for i := 0; i < 1; i = q.N1 {
			n := util.GetNum(0, 20)
			if n < 6 && c.Add > 0 {
				q = sumQuiz(c.DigitsAdd1, c.DigitsAdd2)
				c.Add--
				break
			}
			if (n > 5 && n < 11) && c.Subs > 0 {
				q = subQuiz(c.DigitsSub1, c.DigitsSub2)
				c.Subs--
				break
			}
			if (n > 10 && n < 16) && c.Mul > 0 {
				q = mulQuiz(c.DigitsMul1, c.DigitsMul2, c.Mul2Min)
				c.Mul--
				break
			}
			if (n > 15 && n < 21) && c.Div > 0 {
				q = divQuiz(c.DigitsDiv1, c.DigitsDiv2, c.SimpleDivision, c.Div1Max, c.Div2Max, c.Div2Min)
				c.Div--
				break
			}
		}
	} else {
		q = Question{}
		q.Op = "."
	}

	return q
}
