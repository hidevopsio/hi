package at

// EnableScheduling
type EnableScheduling struct {
	Annotation

	BaseAnnotation
}

// Scheduler
type Scheduled struct {
	Annotation

	BaseAnnotation

	// limit times
	AtLimit *int `at:"limit" json:"-"`

	// standard cron expressions
	AtCron *string `at:"cron" json:"-"`

	// number
	AtEvery *int `at:"every" json:"-"`
	// valid units are: milliseconds, seconds, minutes, hours, days, weeks, months
	AtUnit *string `at:"unit" json:"-"`
	// at
	AtTime *string `at:"time" json:"-"`

	// tag
	AtTag *string `at:"tag" json:"-"`

	// delay
	AtDelay *int64 `at:"delay" json:"-"`
}

