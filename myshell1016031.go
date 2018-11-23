package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const flag = "?"
const bit = ":"

func main() {
	stdin := bufio.NewScanner(os.Stdin)
	count := 0
	fmt.Printf("./myshell[00]> ")
	// Loop Scan
	for stdin.Scan() {
		text := stdin.Text()
		arr := strings.Split(text, " ")
		// Exit Command
		if arr[0] == "bye" {
			os.Exit(0)
		}
		// 区切り文字あるか確認
		if checkFlag(arr) {
			// ?より手前の切り出し
			s, flagPoint := searchFlag(arr)
			// ?より後の切り分け
			s2, bitPoint := searchBit(arr, flagPoint)
			args1 := s
			args2 := s2
			args3 := arr[bitPoint+1:]
			// とりあえず最初のコマンド実行
			result, err := doCommand(args1[0], args1)
			if err != nil {
				log.Println(err)
			}
			// Exit出ちゃった時
			if !result {
				_, err := doCommand(args3[0], args3)
				if err != nil {
					log.Println(err)
				}
				// Successした時
			} else {
				_, err := doCommand(args2[0], args2)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			args := arr[0:]
			_, err := doCommand(args[0], args)
			if err != nil {
				log.Println(err)
			}
		}
		count++
		fmt.Printf("./myshell[%02d]> ", count)
	}
}

func doCommand(topCommand string, option []string) (result bool, err error) {
	attr := syscall.ProcAttr{Files: []uintptr{0, 1, 2}}
	cpath, err := exec.LookPath(topCommand)
	if err != nil {
		return
	}
	pid, err := syscall.ForkExec(cpath, option, &attr)
	if err != nil {
		return
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return
	}
	status, err := proc.Wait()
	if err != nil {
		return
	}
	if !status.Success() {
		fmt.Println(status.String())
	}
	result = status.Success()
	return
}

func searchFlag(s []string) (s2 []string, i int) {
	for _, a := range s {
		if a == flag {
			break
		} else {
			s2 = append(s2, a)
		}
		i++
	}
	return
}

func searchBit(s []string, count int) (s2 []string, i int) {
	for _, a := range s {
		if i > count {
			if a == bit {
				break
			} else {
				s2 = append(s2, a)
			}
		}
		i++
	}
	return
}

func checkFlag(s []string) (b bool) {
	b = false
	for _, a := range s {
		if a == flag {
			b = true
			break
		}
	}
	return
}
