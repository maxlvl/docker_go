package main

import (
  "io"
	"os"
	"os/exec"
  "path/filepath"
  "syscall"
  "fmt"
)

func copyFile(src, dst string) error {
  fmt.Println(src)

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return dstFile.Sync()
}

type nullReader struct{}
type nullWriter struct{}
func (nullReader) Read(p []byte) (n int, err error)  { return len(p), nil }
func (nullWriter) Write(p []byte) (n int, err error) { return len(p), nil }

func main() {
	command := os.Args[3]
	args := os.Args[4:len(os.Args)]
	cmd := exec.Command(command, args...)

  dir := "chroot_dir"
  err := os.Mkdir(dir, 0755)
  if err != nil {
    fmt.Println("oops!")
    return
  }

  err = syscall.Chroot(dir)
  if err != nil {
    fmt.Println("oopser!")
    return 
  }


  dockerExplorerPath := "/usr/local/bin/docker-explorer"
  dockerExplorerDestinationPath := filepath.Join(dir, "docker-explorer")

  err = copyFile(dockerExplorerPath, dockerExplorerDestinationPath)
	if err != nil {
 	  fmt.Println("Error copying docker-explorer path:", err)
	 	return
	}

	commandPath, err := exec.LookPath(command)
	if err != nil {
		fmt.Println("Error finding command path:", err)
		return
	}

  fmt.Println(commandPath)
	dstCommandPath := filepath.Join(dir, filepath.Base(commandPath))
	err = copyFile(commandPath, dstCommandPath)
	if err != nil {
		fmt.Println("Error copying command:", err)
		return
	}

  err = os.Chdir("/")
  if err != nil {
    fmt.Println("oopsest!")
    return 
  } 
  cmd.Stdin = nullReader{}
 	cmd.Stdout = nullWriter{}
 	cmd.Stderr = nullWriter{}

  cmd.Path = dstCommandPath
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
