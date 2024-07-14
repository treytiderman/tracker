# Todo

- [ ] Add Tailwind css as a static file
    - https://github.com/tailwindlabs/tailwindcss-from-zero-to-production/tree/main/01-setting-up-tailwindcss
- [ ] Fix hover sticking on mobile
	- https://dev.to/truongductri01/solving-the-sticky-hover-effect-on-mobile-with-tailwindcss-i5p
- [ ] Deleting a tracker causes problems for new trackers?

# Ideas

```go
func Tracker_Table_Create(db *sql.DB) {
	// lock the tracker from adding changing fields
}

func Tracker_Table_Migrate(db *sql.DB) {
	// unlock the tracker from adding changing fields
	// update fields
	// migrate data to new schema
	// re-lock the tracker
}
```

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