package main

import (
  "os"
	"os/exec"
  "syscall"
  "io/ioutil"
  "fmt"
  "log"
  "path"
)



func main() {
	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

  tmpDir, err := ioutil.TempDir("", "")
  if err != nil {
    log.Fatal("tmpdir: ", err)
  }

  commandPath := path.Join(tmpDir, command)

  if err := exec.Command("mkdir", "-p", path.Dir(commandPath)).Run(); err != nil {
    log.Fatal("mkdir: ", err)
  }
  

  if err := exec.Command("cp", "-f", command, commandPath).Run(); err != nil {
    log.Fatal("copy: ", err)
  }


  err = syscall.Chroot(tmpDir)
  if err != nil {
    fmt.Println("oopser!")
    return 
  }

  err = os.Chdir("/")
  if err != nil {
    fmt.Println("oopsest!")
    return 
  } 

	cmd := exec.Command(command, args...)

  cmd.Stdin = os.Stdin
 	cmd.Stdout = os.Stdout
 	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
    fmt.Println(err)
		if exitError, ok := err.(*exec.ExitError); ok {
      fmt.Println(exitError)
			exitCode := exitError.ProcessState.ExitCode()
			os.Exit(exitCode)
		}
	}

}
