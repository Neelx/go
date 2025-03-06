package no

import (
	"bufio"
	"fmt"
	"os"
)

type Task struct {
	ID       int
	Title    string
	Complete bool
}

func Todo() {
	var tasks []Task
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n--- TODO List Menu ---")
		fmt.Println("1. Add Task")
		fmt.Println("2. List Tasks")
		fmt.Println("3. Mark Task as Complete")
		fmt.Println("4. Exit")
		fmt.Print("Choose an option: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Print("Enter task title: ")
			scanner.Scan()
			title := scanner.Text()
			task := Task{
				ID:       len(tasks) + 1,
				Title:    title,
				Complete: false,
			}
			tasks = append(tasks, task)
			fmt.Println("Task added successfully!")

		case 2:
			if len(tasks) == 0 {
				fmt.Println("No tasks found!")
			} else {
				fmt.Println("\nYour Tasks:")
				for _, task := range tasks {
					status := " "
					if task.Complete {
						status = "✓"
					}
					fmt.Printf("[%s] %d. %s\n", status, task.ID, task.Title)
				}
			}

		case 3:
			fmt.Print("Enter task ID to mark as complete: ")
			var id int
			fmt.Scanln(&id)
			found := false
			for i := range tasks {
				if tasks[i].ID == id {
					tasks[i].Complete = true
					found = true
					fmt.Println("Task marked as complete!")
					break
				}
			}
			if !found {
				fmt.Println("Task not found!")
			}

		case 4:
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid option! Please try again.")
		}
	}
}
func main() {
	Todo()
}
