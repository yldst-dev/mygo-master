package view

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"starter/internal/domain"
	"starter/internal/usecase"
)

const repeatedClick = 500 * time.Millisecond

type App struct {
	tasks     *usecase.Tasks
	setDark   func(bool)
	dark      bool
	draft     string
	problem   string
	input     ui.Handle
	removedAt time.Time
}

func LoadFailure(err error) func(*ui.Context) {
	return func(c *ui.Context) {
		ui.Column(c).Fill().Padding(24).Gap(8).Children(func() {
			ui.Text(c, "The saved tasks could not be read").FontSize(18).Bold().Role(ui.RoleHeading)
			ui.Text(c, "The file is untouched. Fix or move it, then open the app again.")
			ui.Text(c, err.Error()).FontSize(13).TextColor(c.Theme().TextMuted).Selectable()
		})
	}
}

func New(tasks *usecase.Tasks, dark bool, setDark func(bool)) *App {
	return &App{tasks: tasks, dark: dark, setDark: setDark}
}

func (a *App) View(c *ui.Context) {
	darkness := progress(c.Root(), "theme", a.dark, themeMotion)
	c.SetTheme(fadedTheme(darkness))
	ui.Column(c).Fill().Padding(24).Gap(16).Children(func() {
		a.header(c)
		a.composer(c)
		steadyScroll(c, func() {
			for _, task := range a.tasks.List() {
				a.row(c, task)
			}
		})
	})
}

func (a *App) header(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).Gap(12).Children(func() {
		ui.Text(c, "Tasks").FontSize(22).Bold().Role(ui.RoleHeading).Grow(1)
		ui.Box(c).ClipY().Children(func() {
			ui.Text(c, fmt.Sprintf("%d left", a.tasks.Remaining())).TextColor(t.TextMuted).NoWrap().ContentTransition(rollDigits)
		})
		toggle(c, &a.dark, "Dark mode", func() { a.setDark(a.dark) })
	})
}

func (a *App) composer(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).Gap(8).Children(func() {
		input := ui.TextInput(c.Key("draft"), &a.draft).Bind(&a.input).Label("New task").Placeholder("New task").Height(40).Grow(1)
		if primaryButton(c, "Add").Clicked() || input.Submitted() {
			a.add()
		}
	})
	if a.problem != "" {
		ui.Text(c.Key("problem"), a.problem).TextColor(t.Danger).Transition(popMotion)
	}
}

func (a *App) add() {
	err := a.tasks.Add(a.draft)
	a.report(err)
	if err == nil {
		a.draft = ""
	}
}

func (a *App) report(err error) {
	a.problem = ""
	if err != nil {
		a.problem = err.Error()
	}
}

func (a *App) row(c *ui.Context, task domain.Task) {
	t := c.Theme()
	row := ui.Row(c.Key(task.ID)).Height(44).PaddingX(8).Gap(8).Radius(10).ClipY().Transition(rowMotion)
	row.Background(t.SurfaceHover.Alpha(progress(row, "hover", row.Hovered(), fadeMotion))).Children(func() {
		done := task.Done
		ui.Checkbox(c, &done, task.Title).Label(task.Title).Grow(1).MinWidth(0).FillHeight().OnChange(func() { a.report(a.tasks.Toggle(task.ID)) })
		remove := softButton(c, "Remove").Label("Remove " + task.Title).KeepFocus().Shrink(0)
		byKeyboard := remove.Focused()
		remove.OnClick(func() {
			a.remove(c, task.ID)
			if byKeyboard {
				a.input.Focus()
			}
		})
	})
}

func (a *App) remove(c *ui.Context, id int) {
	if c.Now().Sub(a.removedAt) < repeatedClick {
		return
	}
	a.removedAt = c.Now()
	a.report(a.tasks.Remove(id))
}
