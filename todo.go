package main

import (
	"errors"
	table2 "github.com/aquasecurity/table"
	"os"
	"strconv"
	"strings"
	"time"
)

type Todo struct {
	Title       string
	Completed   bool
	DueDate     *time.Time
	CreatedAt   time.Time
	CompletedAt *time.Time
}

type Todos []Todo

func parseNaturalLanguage(input string) (*time.Time, string) {
	now := time.Now()
	words := strings.Fields(input)
	var dueDate *time.Time
	var titleWords []string

	// Loop through the words to identify date/time-related words
	for i, word := range words {
		lowerWord := strings.ToLower(word)

		switch lowerWord {
		case "today":
			// Set due date to today at 11:59 PM
			d := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 0, 0, now.Location())
			dueDate = &d
		case "tomorrow":
			// Set due date to tomorrow at 11:59 PM
			tomorrow := now.Add(24 * time.Hour)
			d := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 23, 59, 0, 0, now.Location())
			dueDate = &d
		case "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday":
			// Find the next occurrence of the specified day
			weekdayMap := map[string]time.Weekday{
				"monday":    time.Monday,
				"tuesday":   time.Tuesday,
				"wednesday": time.Wednesday,
				"thursday":  time.Thursday,
				"friday":    time.Friday,
				"saturday":  time.Saturday,
				"sunday":    time.Sunday,
			}
			targetDay := weekdayMap[lowerWord]
			offset := (int(targetDay) - int(now.Weekday()) + 7) % 7
			if offset == 0 {
				offset = 7 // Next week's same day
			}
			nextDay := now.AddDate(0, 0, offset)
			d := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 23, 59, 0, 0, nextDay.Location())
			dueDate = &d
		case "at":
			// Check if the next word specifies a time (e.g., "10am", "3:30pm")
			if i+1 < len(words) {
				timePart := words[i+1]
				parsedTime, err := time.Parse("3:04pm", timePart)
				if err == nil {
					hour, m, _ := parsedTime.Clock()
					if dueDate == nil {
						// If no date is specified, assume today at the given time
						today := time.Date(now.Year(), now.Month(), now.Day(), hour, m, 0, 0, now.Location())
						dueDate = &today
					} else {
						// Add the time to the existing date
						*dueDate = time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), hour, m, 0, 0, dueDate.Location())
					}
					break
				}
			}
		default:
			titleWords = append(titleWords, word)
		}
	}

	// Default time to 11:59 PM if a date is found but no specific time is set
	if dueDate != nil && dueDate.Hour() == 0 && dueDate.Minute() == 0 && dueDate.Second() == 0 {
		*dueDate = time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 23, 59, 0, 0, dueDate.Location())
	}

	// Join remaining words as the title
	title := strings.Join(titleWords, " ")
	return dueDate, title
}

func (todos *Todos) add(input string) {
	dueDate, title := parseNaturalLanguage(input)

	todo := Todo{
		Title:       title,
		Completed:   false,
		DueDate:     dueDate,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}

	*todos = append(*todos, todo)
}

func (todos *Todos) validateIndex(index int) error {
	if index < 0 || index >= len(*todos) {
		return errors.New("index out of range")
	}
	return nil
}

func (todos *Todos) delete(index int) error {
	t := *todos
	if err := t.validateIndex(index); err != nil {
		return err
	}
	*todos = append(t[:index], t[index+1:]...)

	return nil
}

func (todos *Todos) toggle(index int) error {
	t := *todos
	if err := t.validateIndex(index); err != nil {
		return err
	}
	isCompleted := t[index].Completed
	if !isCompleted {
		completionTime := time.Now()
		t[index].CompletedAt = &completionTime
	}
	t[index].Completed = !isCompleted
	return nil
}

func (todos *Todos) edit(index int, title string) error {
	t := *todos
	if err := t.validateIndex(index); err != nil {
		return err
	}
	t[index].Title = title
	return nil
}

func (todos *Todos) print() {
	// Initialize the table
	table := table2.New(os.Stdout)
	table.SetRowLines(false)
	table.SetHeaders("Index", "Title", "Completed", "Created At", "Due At", "Completed At")

	// Iterate through todos and populate the table
	for i, t := range *todos {
		completed := "❌"
		completedAt := "Not Set"
		dueDate := "Not Set"

		// Update display values based on the todo's state
		if t.Completed {
			completed = "✅"
			if t.CompletedAt != nil {
				completedAt = t.CompletedAt.Format(time.RFC1123)
			}
		}
		if t.DueDate != nil {
			dueDate = t.DueDate.Format(time.RFC1123)
		}

		// Add the row to the table
		table.AddRow(
			strconv.Itoa(i),                  // Index
			t.Title,                          // Title
			completed,                        // Completed status
			t.CreatedAt.Format(time.RFC1123), // Created At
			dueDate,                          // Due At
			completedAt,                      // Completed At
		)
	}

	// Render the table
	table.Render()
}
