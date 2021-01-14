package model

//LoadDataJob type details
type LoadDataJob struct {
	UID   string `predicate:"uid,omitempty"`
	ID    string `predicate:"load_data_job_id,omitempty"`
	Date  int64  `predicate:"load_data_job_date,omitempty"`
	Email string `predicate:"load_data_job_email,omitempty"`
}
