package dto

type Calendar struct {
	ID        uint    `json:"id"         form:"id"`
	StartTime string  `json:"start_time" form:"start_time" binding:"required,timeFormat"`
	EndTime   string  `json:"end_time"   form:"end_time"   binding:"required,timeFormat"`
	Price     float64 `json:"price"      form:"price"`
	Monday    bool    `json:"monday"     form:"monday"   `
	Tuesday   bool    `json:"tuesday"    form:"tuesday"  `
	Wednesday bool    `json:"wednesday"  form:"wednesday"`
	Thursday  bool    `json:"thursday"   form:"thursday" `
	Friday    bool    `json:"friday"     form:"friday"   `
	Saturday  bool    `json:"saturday"   form:"saturday" `
	Sunday    bool    `json:"sunday"     form:"sunday"   `
}

type UpdateCalendar struct {
	Times []Calendar `json:"times"  form:"times" binding:"required,dive"`
}
