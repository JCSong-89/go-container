package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		panic("what?")
	}
}

func run() {
	fmt.Printf("Running: %v\v\n", os.Args[2:])
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	/*
		Cloneflags는 프로세스 생성을 제어하는데 사용.
		&syscall.SysProcAttr{
		 Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
		 //Cloneflags is linux only
		}
	*/

	err := cmd.Run()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
		os.Exit(1)
	}
}

func child() {
	//os.Getpid()는 현재 프로세스의 PID를 반환
	fmt.Printf("Running %v as PID %d\n", os.Args[2:], os.Getpid())
	//syscall.Sethostname([]byte("container")) mac os는 할수가 없다.
	err := syscall.Chroot("/home/ubuntu/Downloads/rootfs")
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
		os.Exit(1)
	}

	err = syscall.Chdir("/")
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
		os.Exit(1)
	}

	// syscall.Mount("proc", "proc", "proc", 0, "") mac os 미지원
	defer func() {
		err := syscall.Unmount("proc", 0)
		if err != nil {

		}
	}()

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
		os.Exit(1)
	}
}

func cg() {
	cgroups := "/sys/fs/cgroup/"
	pids := path.Join(cgroups, "pids")
	err := os.Mkdir(path.Join(pids, "test"), 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	err = os.WriteFile(path.Join(pids, "test/pids.max"), []byte("20"), 0700)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(path.Join(pids, "test/notify_on_release"), []byte("1"), 0700)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(path.Join(pids, "test/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
	if err != nil {
		panic(err)
	}
}
