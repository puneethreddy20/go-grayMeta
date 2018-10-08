# go-grayMeta

##### How to build and run the application:


###### Using Go compiler on host machine:


  ```
  $ cd go-grayMeta

  $ go get -d -v ./...

  $ go test -v ./...

  //if there is no binary in the folder. Please run the below command. else ignore and run the binary
  $ go build


    //run the binary
  $ ./go-grayMeta
  ```




##### Implementation:

  The Main functionalities for this application are Creating Dependency Map for packages, Install, Remove, ListAll Packages.

###### Creating Dependency map:
  ```
  DEPEND package1 package2 package3
  ```
  In the above command, package2 and package3 are the dependencies for package1. The Information
  will be stored in a map as package1 as 'key' and value as packageinfo struct which has string slice which stores the dependencies.

######  Install:
  ```
  INSTALL package1
  ```


  With the above command, we explicitly install package1 for which we need to install its dependencies first. The dependencies are installed implicitly. If package2,3 has dependencies we would know from the map, if not
  present in map then we create an entry with string slice as nil, install it. There are two important flags in packageinfo struct which are Installed and ExplicitlyInstalled, which would be helpful to identify if the package is
  installed or not and if installed, whether it is explicitly installed or implicitly installed. So package2 and package3 flags are set as true for Installed and false for ExplicitlyInstalled. The package info struct also an varibale which a
  value which determines number of currently installed packages which need this package as dependency.(n is incremented from 0 to 1) as package1 depends on package2



  ```
  INSTALL package2
  ```

  After the above command, As package2 is already installed, then we check set the explicitlyInstalled flag as true.


###### Remove:
  ```
    REMOVE package1
  ```


  With the above command, we first check the 'n' value from which we get to know if there are any packages which depend on package1. If there are any packages that depend on, it returns as its still needed.
  Otherwise, its removed and then we into its dependencies which are package2,package3 and check if those packages are needed by other packages and if check if those are explicitly installed or implicitly installed.
  If other packages depend on them or if they are explicitly installed then we don't remove them. If there are no other packages that depend on them and if they are implicitly installed then we remove them.

###### ListAll:
  ```
  LIST
  ```

  It Prints out all the currently installed packages to the console.



###### END:
 ```
 END
 ```
 END command, returns the main function, which means termination.


##### Further TODO:

Need to go through some edge test cases as well. works for most of them, need to think more to find out if there are any cases that fail tests.


