package model

type ScheduleTo struct {
	ID                 uint64   `json:"shedule_id"`
	NameMedication     string   `json:"name_medication"`
	MedicationPerDay   int      `json:"medication_per_day"`
	ScheduleMedication []string `json:"schedule"`
}
