# Todo

- [ ] Deleting a tracker causes problems for new trackers?
- [ ] Add Ctrl+Enter to submit from a "textarea"

# Ideas

```go
package main

// tracks a named data point. database table
type _Tracker struct {
	tracker_id          string
	friendly_name       string
	related_tracker_ids []string
	note                string // markdown
	tags                []string

	value_type           string // "Number" | "String" | "List"
	value_decimal_places int    // -2 for money, 0 for integers, 2 for hundreds
	value_max            int    // Length for strings
	value_min            int    // Length for strings
}

// point in time
type _Event[T int | string | []string] struct {
	event_id   string
	tracker_id string // parent,
	timestamp  string // data time
	note       string // markdown

	target_id string

	value                T
	value_type           string // "Number" | "String" | "List"
	value_decimal_places int    // -2 for money, 0 for integers, 2 for hundreds

	goal           bool
	goal_value     T
	goal_type      string // "Equal" | "Greater Than" | "Less Than" | "Contains"
	goal_completed bool
}

type _Target[T int | string | []string] struct {
	target_id     string
	tracker_id    string
	friendly_name string
	note          string // markdown
	tags          []string
	date          string // data time

	reminder bool

	goal       bool
	goal_type  string // "Equal" | "Greater Than" | "Less Than" | "Contains"
	goal_value T

	repeat                         bool
	repeat_type                    string // "Daily" | "Weekly" | "Monthly" | "Yearly"
	repeat_every_value             int    // "Daily "Musically (BPM)" |" Hourly
	repeat_every_monday            bool   // Week
	repeat_every_tuesday           bool   // Week
	repeat_every_wednesday         bool   // Week
	repeat_every_thursday          bool   // Week
	repeat_every_friday            bool   // Week
	repeat_every_saturday          bool   // Week
	repeat_every_sunday            bool   // Week
	repeat_monthly_on_this_weekday bool   // Month
}

type _Value_Update_Event[T int | string | []string] struct {
	event_id   string
	tracker_id string
	timestamp  string // data time
	note       string // markdown

	value                T
	value_type           string // "Number" | "String" | "List"
	value_decimal_places int    // -2 for money, 0 for integers, 2 for hundreds
}

type _Target_Event[T int | string | []string] struct {
	event_id   string
	target_id  string
	tracker_id string
	timestamp  string // data time
	note       string // markdown

	value                T
	value_type           string // "Number" | "String" | "List"
	value_decimal_places int    // -2 for money, 0 for integers, 2 for hundreds

	goal           bool
	goal_value     T
	goal_type      string // "Equal" | "Greater Than" | "Less Than" | "Contains"
	goal_completed bool
}

```