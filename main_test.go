/*
Created & Written by: Puneeth Reddy
*/

package main

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"testing"
	"unicode"
)


//TODO:Need to go through some edge test cases as well. works for most of them, need to think more to find out if there are any cases that fail tests.


func TestService_Install(t *testing.T) {
	input:=[]string{"NETCARD","TELNET","DNS","BROWSER"}
	//input1:=[]string{"TELNET","DNS","BROWSER","TELNET"}
	//input2:=[]string{"TCPIP","DNS","foo"}

	newservice:=NewService()

	newservice.PackageInfo= map[string]packageinfo{"TELNET":packageinfo{0,[]string{"TCPIP","NETCARD"},false,false},
		"TCPIP":packageinfo{0,[]string{"NETCARD"},false,false},
		"DNS":packageinfo{0,[]string{"TCPIP","NETCARD"},false,false},
		"BROWSER":packageinfo{0,[]string{"TCPIP","HTML"},false,false},}

	for _,val:=range input{
		newservice.Install(val,true)
	}
	//fmt.Println(newservice)

	expectedservice:=map[string]packageinfo{"TELNET":packageinfo{0,[]string{"TCPIP","NETCARD"},true,true},
		"TCPIP":packageinfo{3,[]string{"NETCARD"},true,false},
		"DNS":packageinfo{0,[]string{"TCPIP","NETCARD"},true,true},
		"BROWSER":packageinfo{0,[]string{"TCPIP","HTML"},true,true},
		"NETCARD":packageinfo{3,nil,true,true},"HTML":packageinfo{1,nil,true,false}}

	if !reflect.DeepEqual(newservice.PackageInfo,expectedservice){
		t.Errorf("Expected %v but got %v",expectedservice,newservice.PackageInfo)
	}

}

func TestService_ProcessDependency(t *testing.T) {
	input:=[]string{"TELNET","TCPIP","NETCARD"}
	input1:=[]string{"TCPIP","NETCARD"}
	input2:=[]string{"DNS","TCPIP","NETCARD"}
	input3:=[]string{"BROWSER","TCPIP","HTML"}

	newservice:=NewService()

	newservice.ProcessDependency(input,3)
	newservice.ProcessDependency(input1,2)
	newservice.ProcessDependency(input2,3)
	newservice.ProcessDependency(input3,3)

	ExpectedOutput:=map[string]packageinfo{"TELNET":packageinfo{0,[]string{"TCPIP","NETCARD"},false,false},
		"TCPIP":packageinfo{0,[]string{"NETCARD"},false,false},
		"DNS":packageinfo{0,[]string{"TCPIP","NETCARD"},false,false},
		"BROWSER":packageinfo{0,[]string{"TCPIP","HTML"},false,false},}

	if !reflect.DeepEqual(newservice.PackageInfo,ExpectedOutput){
		t.Errorf("Expected %v but got %v",ExpectedOutput,newservice.PackageInfo)
	}

}

func TestService_ListInstalledPackages(t *testing.T) {
	//input:=[]string{"NETCARD","TELNET","DNS","BROWSER"}
	input1 := []string{"TELNET", "DNS", "BROWSER", "TELNET"}
	//input2:=[]string{"TCPIP","DNS","foo"}

	newservice := NewService()

	newservice.PackageInfo = map[string]packageinfo{"TELNET": packageinfo{0, []string{"TCPIP", "NETCARD"}, false,false},
		"TCPIP":   packageinfo{0, []string{"NETCARD"}, false,false},
		"DNS":     packageinfo{0, []string{"TCPIP", "NETCARD"}, false,false},
		"BROWSER": packageinfo{0, []string{"TCPIP", "HTML"}, false,false},}

	for _, val := range input1 {
		newservice.Install(val,true)
	}
	output:=newservice.ListInstalledPackages()
	sort.Strings(output)
	expectedoutput:=[]string{"TELNET","DNS","TCPIP","NETCARD","BROWSER","HTML"}
	sort.Strings(expectedoutput)
	if !reflect.DeepEqual(output,expectedoutput){
		t.Errorf("Expected %v but got %v",expectedoutput,output)
	}
}


func TestService_Remove(t *testing.T) {
	input:=[]string{"NETCARD","TELNET","DNS","BROWSER"}
	//input1 := []string{"TELNET", "DNS", "BROWSER", "TELNET"}
	//input2:=[]string{"TCPIP","DNS","foo"}

	newservice := NewService()

	newservice.PackageInfo = map[string]packageinfo{"TELNET": packageinfo{0, []string{"TCPIP", "NETCARD"}, false,false},
		"TCPIP":   packageinfo{0, []string{"NETCARD"}, false,false},
		"DNS":     packageinfo{0, []string{"TCPIP", "NETCARD"}, false,false},
		"BROWSER": packageinfo{0, []string{"TCPIP", "HTML"}, false,false},}

	for _,val:=range input{
		newservice.Install(val,true)
	}
	newservice.Install("foo",true)
	newservice.Install("hello",true)

	for _, val := range input {
		newservice.Remove(val,true)
	}
	output:=newservice.ListInstalledPackages()
	sort.Strings(output)
	expectedoutput:=[]string{"foo","hello"}
	if !reflect.DeepEqual(output,expectedoutput){
		t.Errorf("Expected %v but got %v",expectedoutput,output)
	}
}


func Test_MainService(t *testing.T){
	newservice:=NewService()
	f := func(c rune) bool {
		//assuming some package names have numbers as well.
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	input:=[]string{"DEPEND TELNET TCPIP NETCARD","DEPEND TCPIP NETCARD",
		"DEPEND DNS TCPIP NETCARD",
		"DEPEND BROWSER TCPIP HTML",
		"INSTALL NETCARD",
		"INSTALL TELNET",
		"INSTALL foo",
		"REMOVE NETCARD",
		"INSTALL BROWSER",
		"INSTALL DNS",
		"LIST",
		"REMOVE TELNET",
		"REMOVE NETCARD",
		"REMOVE DNS",
		"REMOVE NETCARD",
		"INSTALL NETCARD",
		"REMOVE TCPIP",
		"REMOVE BROWSER",
		"REMOVE TCPIP",
		"LIST",
		"END",
	}
	for _,val:=range input{
		fmt.Println(val)

		//vals is string slice which contains all the words in the line.
		vals:=strings.FieldsFunc(val,f)


		//get the length of the vals slice
		valLen:=len(vals)

		//if its zero just continue to next line.
		if valLen==0{
			continue
		}
		//The first value in vals slice determines the command to be run. so thats why the switch case. Based on the value, appropriate functions will run.

		switch vals[0]{

		case "DEPEND":
			if valLen<2{
				continue
			}else{
				newservice.ProcessDependency(vals[1:valLen],valLen-1)
			}
			break


		case "INSTALL":
			if valLen==2{
				_,err:=newservice.Install(vals[1],true)
				if err!=nil{
					log.Println("Error while working on INSTALL",err)
					log.Println("Please Issue right command")
				}
			}
			break


		case "REMOVE":

			if valLen==2 {
				_,err:=newservice.Remove(vals[1],true)
				if err!=nil{
					log.Println("Erorr in remove case",err)
					log.Println("Please Issue right command")
				}
			}
			break
		case "LIST":
			newservice.ListInstalledPackages()
			break
		case "END":
			//fmt.Println("END")
			return
		}
	}
}