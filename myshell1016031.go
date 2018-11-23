package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const ternaryFlag = "?"
const redirectFlag = ">"
const redirectOppFlag = "<"
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
		if checkFlag(arr, ternaryFlag) {
			// ?より手前の切り出し
			s, flagPoint := searchFlag(arr, ternaryFlag)
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
			// > のリダイレクトがあった場合
		} else if checkFlag(arr, redirectFlag) {
			args1, args2 := splitFlag(arr, redirectFlag)
			args1Option := connectOption(args1[1:])
			// 標準入力の取得
			out, err := exec.Command(args1[0], args1Option).Output()
			if err != nil {
				log.Println(err)
			}
			// 書き込みできるようファイルを用意
			fw, err := os.OpenFile(args2[0], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				log.Println(err)
			}
			defer fw.Close()
			// 書き込み
			fmt.Fprintln(fw, string(out))
		} else if checkFlag(arr, redirectOppFlag) {
			s, flagPoint := searchFlag(arr, redirectOppFlag)
			args1 := s
			args2 := args1[flagPoint+1:]
			b, err := ioutil.ReadFile(args2[0])
			if err != nil {
				log.Println(err)
			}
			status := string(b)
			result, status, err := doCommandForRedirect(args1[0], args1)
			if err != nil {
				log.Println(err)
			}
			// 書き込みできるようファイルを用意
			if result {
				fw, err := os.OpenFile(args2[0], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
				if err != nil {
					log.Println(err)
				}
				defer func() {
					if err := fw.Close(); err != nil {
						log.Println(err)
					}
				}()
				_, err = fw.WriteString(status)
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

func doCommandForRedirect(topCommand string, option []string) (result bool, st string, err error) {
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
	st = status.String()
	return
}

func doCommandForOppRedirect(topCommand string, option []string) (result bool, st string, err error) {
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
	st = status.String()
	return
}

func searchFlag(s []string, flag string) (s2 []string, i int) {
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

func checkFlag(s []string, flag string) (b bool) {
	b = false
	for _, a := range s {
		if a == flag {
			b = true
			break
		}
	}
	return
}

func splitFlag(s []string, flag string) (b []string, c []string) {
	splitFlag := false
	for _, a := range s {
		if a == flag {
			splitFlag = true
		} else if splitFlag == false {
			b = append(b, a)
		} else if splitFlag == true {
			c = append(c, a)
		}
	}
	return
}

func connectOption(s []string) (sentence string) {
	for i, a := range s {
		if i == 0 {
			sentence = a
		} else {
			sentence = sentence + " " + a
		}
	}
	return
}
