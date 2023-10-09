package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const AccessHeader = "superuser"
const ErrHeader = "ваш доступ ограничен к запрашиваему ресурсу"
const MathematicSymbols = "+-"

func getRoot(w http.ResponseWriter, r *http.Request) {
	log.Print("got / request")

	//ua := r.Header.Get("User-Access")
	// API должно быть доступно только если в HTTP Header есть “User-Access”
	//со значением “superuser”. В случае отказа в доступе, нужно вывести
	//сообщение в консоли и отправить соответствующий ответ клиенту.
	//if ua != AccessHeader {
	//	log.Printf("неверный заголовок входа User-Access:%s", ua)
	//	http.Error(w, ErrHeader, http.StatusForbidden)
	//	return
	//}
	//только операции сложения и вычитания.
	//Например, с Frontend-а к тебе приходит строка “2+2-3-5+1”.
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("не удается прочитать тело сообщения err:%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(bytes) == 0 {
		bytes = []byte(r.URL.RawQuery)
	}
	// поскольку мы имеем дело только с цифрами и знаками, то можно итерироваться по байтам
	var nums []int = make([]int, 1)
	var plus []string
	for pos, char := range bytes {
		charStr := string(char)
		log.Printf("character %s starts at byte position %d", charStr, pos)
		if strings.Contains(MathematicSymbols, charStr) {
			// математический знак
			plus = append(plus, charStr)
			nums = append(nums, 0)
		} else {
			intVar, err := strconv.Atoi(charStr)
			if err != nil {
				textErr := fmt.Sprintf("ошибка чтения числа err:%s", err)
				log.Printf(textErr)
				http.Error(w, textErr, http.StatusInternalServerError)
				return
			}
			
			fullInt := nums[len(nums)-1]*10 + intVar // сдвигаем и прибавляем
			//fullInt, err := strconv.Atoi(fmt.Sprintf("%d%d", nums[len(nums)-1], intVar)) // или подставляем и преобразуем
			
			nums[len(nums)-1] = fullInt
		}
	}
	// посчитаем
	var summ int
	for i, num := range nums {
		if i == 0 {
			summ = num
			continue
		}
		var sznak string
		if len(plus) > 0 {
			sznak = plus[0]
			plus = plus[1:]
		}
		switch sznak {
		case "+":
			summ += num
		case "-":
			summ -= num
		}
	}
	w.Write([]byte(fmt.Sprint(summ)))
}

func main() {
	fmt.Println("")
	http.HandleFunc("/", getRoot)
	fmt.Println("http://localhost:3333/?2+2-3-5+1\n")
	fmt.Println("http://localhost:3333/?852+7656-133\n")
	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed ok\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
