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
		if check_flag(arr) {
			s, flagPoint := search_flag(arr)
			s2, bitPoint := search_bit(arr, flagPoint)
			args1 := s
			args2 := s2
			args3 := arr[bitPoint+1:]
			attr := syscall.ProcAttr{Files: []uintptr{0, 1, 2}}
			cpath, err := exec.LookPath(args1[0])
			if err != nil {
				log.Println(err)
			} else {
				pid, err := syscall.ForkExec(cpath, args1, &attr)
				if err != nil {
					log.Println(err)
				}
				proc, err := os.FindProcess(pid)
				status, err := proc.Wait()
				if err != nil {
					log.Println(err)
				}
				if !status.Success() {
					fmt.Println(status.String())
					cpath, err := exec.LookPath(args3[0])
					if err != nil {
						log.Println(err)
					}
					pid, err := syscall.ForkExec(cpath, args3, &attr)
					proc, err := os.FindProcess(pid)
					status, err := proc.Wait()
					if err != nil {
						log.Println(err)
					}
					if !status.Success() {
						fmt.Println(status.String())
					}
				} else {
					cpath, err := exec.LookPath(args2[0])
					if err != nil {
						log.Println(err)
					}
					pid, err := syscall.ForkExec(cpath, args2, &attr)
					proc, err := os.FindProcess(pid)
					status, err := proc.Wait()
					if err != nil {
						log.Println(err)
					}
					if !status.Success() {
						fmt.Println(status.String())
					}
				}
			}
		} else {
			args := arr[0:]
			attr := syscall.ProcAttr{Files: []uintptr{0, 1, 2}}
			cpath, err := exec.LookPath(args[0])
			if err != nil {
				log.Println(err)
			} else {
				pid, err := syscall.ForkExec(cpath, args, &attr)
				if err != nil {
					log.Println(err)
				}
				proc, err := os.FindProcess(pid)
				status, err := proc.Wait()
				if err != nil {
					log.Println(err)
				}
				if !status.Success() {
					fmt.Println(status.String())
				}
			}
		}
		count++
		fmt.Printf("./myshell[%02d]> ", count)
	}
}

func search_flag(s []string) (s2 []string, i int) {
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

func search_bit(s []string, count int) (s2 []string, i int) {
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

func check_flag(s []string) (b bool) {
	b = false
	for _, a := range s {
		if a == flag {
			b = true
			break
		}
	}
	return
}
