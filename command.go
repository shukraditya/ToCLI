package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type cmdFlags struct {
	Add    string
	Del    int
	Edit   string
	Toggle int
	List   bool
}

func NewCommandFlags() *cmdFlags {
	cf := cmdFlags{}
	flag.StringVar(&cf.Add, "add", "", "Add a new task")
	flag.StringVar(&cf.Edit, "edit", "", "Edit Title by index value in the format {id}:{new title}")
	flag.IntVar(&cf.Del, "del", -1, "Specify a todo by index to delete")
	flag.IntVar(&cf.Toggle, "toggle", -1, "Toggle by index")
	flag.BoolVar(&cf.List, "list", false, "List all todos")

	flag.Parse()
	return &cf
}

func (cf *cmdFlags) execute(todos *Todos) {
	switch {
	case cf.List:
		todos.print()
	case cf.Add != "":
		todos.add(cf.Add)
	case cf.Edit != "":
		parts := strings.SplitN(cf.Edit, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Use the format {id}:{new title}")
			os.Exit(1)
		}
		index, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("Invalid index")
			os.Exit(1)
		}
		err = todos.edit(index, parts[1])
		if err != nil {
			return
		}
	case cf.Toggle != -1:
		err := todos.toggle(cf.Toggle)
		if err != nil {
			return
		}
	case cf.Del != -1:
		err := todos.delete(cf.Del)
		if err != nil {
			return
		}
	default:
		todos.print()
		fmt.Println("You did not pass a flag.")
	}

}
