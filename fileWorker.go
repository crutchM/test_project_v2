package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

var c chan string

func PerformLinks() {
	tmp := readFile()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		WriteInResultFile()
		wg.Done()
	}()
	c = make(chan string, 10000)
	// 1. Все функции стремятся взять передаваемые параметры по значению
	// поэтому url := v это как минимум странное решение, поскольку строки копируются
	// 2. (70 000 грутин + 3) достижимо для сервера, но смысла от простаивающих горутин мало будет
	routines := make(chan struct{}, 10000) //буффер для работающих горутин(надеюсь вариант избавления от простаивающих рутин приемлем)
	for _, v := range tmp {
		url := v
		routines <- struct{}{}
		wg.Add(1)
		go func(url string) {
			c <- Ping(url)
			<-routines
			wg.Done()
		}(url)
	}
	wg.Wait()
	close(c)
}

// Функция чтения файла превратилась в функцию обработки данных из файла
// и запуск механизма записи в другой файл.
//постарался выделить обработку и запуск записи в файл более логично
func readFile() []string {
	data, err := os.Open("sites")
	if err != nil {
		panic(err)
	}
	sc := bufio.NewScanner(data)
	// --------------------------
	// За этот блок большой минус
	// Причны:
	// 1. Не понятно для чего нужен массив, предполагаю для того чтобы временно переложить данные
	//если не создавать некий "буфер" то при сканировании ввода происходит блокировка
	// 2. Для чего вызывать panic если ты ТУТ ЖЕ её отлавилваешь?
	// неправильно понял суть отлавливания ошибок, вроде разобрался, поменял
	var tmp []string
	for sc.Scan() {
		tmp = append(tmp, sc.Text())
	}
	defer data.Close()
	defer recovery()
	return tmp
}

func WriteInResultFile() {
	defer recovery()
	f, err := os.OpenFile("result", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {

		panic(err)
	}
	// А для чего нужна дата в файле?
	//для некой "отчетности", дабы понимать актуальность информации, сделано больше для собственного удобства
	f.WriteString(time.Now().Format("02-01 15:04") + "\n")
	for {
		val, ok := <-c
		if !ok {
			break
		} else {
			if _, err = f.WriteString(val); err != nil {
				panic(err)
			}
		}

	}
	f.WriteString("\n")
	f.Close()
}

func recovery() {
	if msg := recover(); msg != nil {
		fmt.Println(msg)
	}
}
