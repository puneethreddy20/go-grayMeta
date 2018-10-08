/*
Created & Written by: Puneeth Reddy
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)


type packageinfo struct {
	//number of currently installed packages which need this package as dependency.
	n int

	//Its dependencies
	Required []string

	//if installed true, else false
	Installed bool


	//check if its Explicitly installed or implicitly installed. This should only matter when Installed is true.
	ExplicitlyInstalled bool

}


type Service struct {

	PackageInfo map[string]packageinfo
}


func (p *packageinfo) IsInstalled() bool{
	return p.Installed
}


// To Printout/ return all the installed packages.
func (s *Service) ListInstalledPackages()[]string {
	var output []string
	if s.PackageInfo != nil {
		for key, value := range s.PackageInfo {
			if value.IsInstalled() {
				output=append(output,key)
				fmt.Println(key)
			}
		}
	}
	return output
}

//To remove the installed packages
func (s *Service) Remove(input string,ExplicitCall bool)(bool,error){

	val,ok:=s.PackageInfo[input]

	if ok{

		//check if its installed or not, if installed then remove it.
		if val.IsInstalled(){
			//checking if there any other packages that depend on this package
			if val.n>0{
				if ExplicitCall {
					val.ExplicitlyInstalled=false
					s.PackageInfo[input]=val
					fmt.Println(input, "is still needed")
				}
				return false,nil
			}else{

				//remove the package.. if its implicit call(not explicit)
				if ExplicitCall==false{
					fmt.Println(input,"is no longer needed")
				}
				val.Installed=false
				val.ExplicitlyInstalled=false
				fmt.Println(input,"successfully removed")

				s.PackageInfo[input]=val

				//Update its dependencies
				if val.Required!=nil {
					for _, eachpackage := range val.Required {
						eachpackageInfo := s.PackageInfo[eachpackage]

						//decrement the value, because of its package which need it, is removed.
						eachpackageInfo.n--

						//update it in map.
						s.PackageInfo[eachpackage] = eachpackageInfo

						//After decrementing if the value is zero, then remove it as well.
						if eachpackageInfo.n == 0 {
							//if the package is explicitly installed then dont remove it.
							if eachpackageInfo.ExplicitlyInstalled{
								continue
							}
							_, err := s.Remove(eachpackage, false)
							if err != nil {
								log.Println("Error occured while removing", eachpackage)
								return false, err
							}

						}
					}
				}

			}
		}else{
			fmt.Println(input,"is not installed")
		}
	}else{
		fmt.Println(input,"is not installed")
	}
	return true,nil
}



//To Install the packages
func (s *Service) Install(input string,ExplicitCall bool)(bool,error){
	val,ok:=s.PackageInfo[input]
	if ok{
		//check if the package is already installed, if installed return true.
		if val.IsInstalled(){
			if ExplicitCall{
				fmt.Println(input,"is already installed")
				//if its explicitly called then make it true.
				val.ExplicitlyInstalled=true
				s.PackageInfo[input]=val
			}
			return true,nil

		}else{
			if ExplicitCall{
				val.ExplicitlyInstalled=true
			}else{
				val.ExplicitlyInstalled=false
			}

			//if not installed go through the package dependencies and installed them.

			//if it doesn't have any dependencies and is not installed.
			if val.Required==nil{
				val.Installed=true
				s.PackageInfo[input]=val
				fmt.Println(input,"successfully installed")
				return true,nil
			}

			for _,eachpackage:=range val.Required{

				eachpackageInfo,ok:=s.PackageInfo[eachpackage]

				if !ok{
					_, err := s.Install(eachpackage,false)
					if err != nil {
						log.Println("Error while installing", eachpackage)
						return false, err
					}
				}else{

					if !(eachpackageInfo.IsInstalled()){
						_, err := s.Install(eachpackage,false)
						if err != nil {
							log.Println("Error while installing", eachpackage)
							return false, err
						}
					}

				}

				UpdatedInfo,ok:=s.PackageInfo[eachpackage]

				if !(UpdatedInfo.IsInstalled()){
					//update the installed to true
					UpdatedInfo.Installed = true
					fmt.Println(eachpackage,"successfully installed")
				}
				//increment the which number of packages would require this package.
				UpdatedInfo.n++
				s.PackageInfo[eachpackage] = UpdatedInfo

			}

			fmt.Println(input,"successfully installed")

			//update the installed value.
			val.Installed=true

			s.PackageInfo[input]=val
		}
	}else{
		var newpackageinfo packageinfo

		if ExplicitCall{
			newpackageinfo = packageinfo{0, nil, true,true}

		}else {
			newpackageinfo = packageinfo{0, nil, true,false}
		}

		s.PackageInfo[input]=newpackageinfo
		fmt.Println(input,"successfully installed")
	}
	return true,nil
}

//To process dependencies and insert in the map
func (s *Service) ProcessDependency(input []string,len int){
	if len==1{
		s.PackageInfo[input[0]]=packageinfo{0,nil,false,false}
		return
	}
	//Create the entry in the PackageInfo map.
	s.PackageInfo[input[0]]=packageinfo{0,input[1:len],false,false}
}


//Create a NewService
func NewService()*Service{
	return &Service{make(map[string]packageinfo)}
}



func main() {

	//scanner to get input from Stdin
	scanner:=bufio.NewScanner(os.Stdin)

	//create a newservice instance.
	newservice:=NewService()


	f := func(c rune) bool {
		//assuming some package names have numbers as well.
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	for scanner.Scan(){

		//gets the first line into input variable
		input:=scanner.Text()

		fmt.Println(input)

		//vals is string slice which contains all the words in the line.
		vals:=strings.FieldsFunc(input,f)

		//fmt.Println(vals)

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
/*
input:

DEPEND TELNET TCPIP NETCARD
DEPEND TCPIP NETCARD
DEPEND DNS TCPIP NETCARD
DEPEND BROWSER TCPIP HTML
INSTALL NETCARD
INSTALL TELNET
INSTALL foo
REMOVE NETCARD
INSTALL BROWSER
INSTALL DNS
LIST
REMOVE TELNET
REMOVE NETCARD
REMOVE DNS
REMOVE NETCARD
INSTALL NETCARD
REMOVE TCPIP
REMOVE BROWSER
REMOVE TCPIP
LIST
END

output:

DEPEND TELNET TCPIP NETCARD
DEPEND TCPIP NETCARD
DEPEND DNS TCPIP NETCARD
DEPEND BROWSER TCPIP HTML
INSTALL NETCARD
NETCARD successfully installed
INSTALL TELNET
TCPIP successfully installed
TELNET successfully installed
INSTALL foo
foo successfully installed
REMOVE NETCARD
NETCARD is still needed
INSTALL BROWSER
HTML successfully installed
BROWSER successfully installed
INSTALL DNS
DNS successfully installed
LIST
HTML
BROWSER
DNS
NETCARD
foo
TCPIP
TELNET
REMOVE TELNET
TELNET successfully removed
REMOVE NETCARD
NETCARD is still needed
REMOVE DNS
DNS successfully removed
REMOVE NETCARD
NETCARD is still needed
INSTALL NETCARD
NETCARD is already installed
REMOVE TCPIP
TCPIP is still needed
REMOVE BROWSER
BROWSER successfully removed
REMOVE TCPIP
TCPIP successfully removed
LIST
NETCARD
foo
END


 */