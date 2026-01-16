package output

import (
	"fmt"
	"io"

	"github.com/whoaa512/asana-cli/internal/models"
)

type Formatter interface {
	Print(v any) error
	PrintError(err error) error
	PrintTasks(tasks []models.Task) error
	PrintTaskList(list *models.ListResponse[models.Task]) error
}

func NewFormatter(format string, w io.Writer) Formatter {
	if format == "brief" {
		return &Brief{w: w}
	}
	return NewJSON(w)
}

type Brief struct {
	w io.Writer
}

func (b *Brief) Print(v any) error {
	if task, ok := v.(models.Task); ok {
		return b.printTask(task)
	}
	if task, ok := v.(*models.Task); ok && task != nil {
		return b.printTask(*task)
	}
	json := NewJSON(b.w)
	return json.Print(v)
}

func (b *Brief) printTask(t models.Task) error {
	if t.DueOn != "" {
		_, err := fmt.Fprintf(b.w, "%s  %s  (due %s)\n", t.GID, t.Name, t.DueOn)
		return err
	}
	_, err := fmt.Fprintf(b.w, "%s  %s\n", t.GID, t.Name)
	return err
}

func (b *Brief) PrintError(err error) error {
	json := NewJSON(b.w)
	return json.PrintError(err)
}

func (b *Brief) PrintTasks(tasks []models.Task) error {
	for _, t := range tasks {
		var err error
		if t.DueOn != "" {
			_, err = fmt.Fprintf(b.w, "%s  %s  (due %s)\n", t.GID, t.Name, t.DueOn)
		} else {
			_, err = fmt.Fprintf(b.w, "%s  %s\n", t.GID, t.Name)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (j *JSON) PrintTasks(tasks []models.Task) error {
	return j.Print(map[string]any{"data": tasks})
}

func (b *Brief) PrintTaskList(list *models.ListResponse[models.Task]) error {
	return b.PrintTasks(list.Data)
}

func (j *JSON) PrintTaskList(list *models.ListResponse[models.Task]) error {
	return j.Print(list)
}
