package repository

import (
	"encoding/json"

	"../driver"
	"../model"
	jsoniter "github.com/json-iterator/go"
)

type LoadDataJobRepository interface {
	FetchBasicInformation() []model.LoadDataJob
	Create(loadDataJob *model.LoadDataJob) string
	FindById(loadDataJob string) *model.LoadDataJob
	Delete(loadDataJob *model.LoadDataJob) *model.LoadDataJob
	Update(uid string, loadDataJob *model.LoadDataJob) string
}

type LoadDataJobRepositoryDGraph struct {
}

func (b LoadDataJobRepositoryDGraph) FetchBasicInformation() ([]model.LoadDataJob, error) {
	query :=
		`
		{
			findAllLoadDataJobs(func: has(load_data_job_id)) {
				uid
				load_data_job_id
				load_data_job_date
				load_data_job_email
			}
		}	
	`
	res, err := driver.RunQuery(query)
	if err != nil {
		return nil, err
	}
	var loadDataJobsFound []model.LoadDataJob
	var objmap map[string]json.RawMessage
	err = json.Unmarshal(res, &objmap)
	if err != nil {
		return nil, err
	}
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["findAllLoadDataJobs"], &loadDataJobsFound)
	return loadDataJobsFound, nil
}

func (b LoadDataJobRepositoryDGraph) Create(loadDataJob *model.LoadDataJob) error {
	loadDataJob.UID = "_:" + loadDataJob.ID
	err := driver.RunMutation(loadDataJob)
	if err != nil {
		return err
	}
	loadDataJob, err = b.FindById(loadDataJob.ID)
	if err != nil {
		return err
	}
	return nil
}

func (b LoadDataJobRepositoryDGraph) Delete(loadDataJob *model.LoadDataJob) error {
	err := driver.RunMutationForDelete(loadDataJob)
	if err != nil {
		return err
	}
	return nil
}

func (b LoadDataJobRepositoryDGraph) Update(uid string, loadDataJob *model.LoadDataJob) error {
	loadDataJob.UID = uid
	err := driver.RunMutation(loadDataJob)
	if err != nil {
		return err
	}
	return nil
}

func (b LoadDataJobRepositoryDGraph) FindById(loadDataJob_id string) (*model.LoadDataJob, error) {
	query :=
		`
		{
			findLoadDataJobById(func: eq(load_data_job_id, "` + loadDataJob_id + `"), first: 1) {
				uid
				load_data_job_id
				load_data_job_date
				load_data_job_email
			}
		}	
	`
	res, err := driver.RunQuery(query)
	if err != nil {
		return nil, err
	}
	var loadDataJobsFound []model.LoadDataJob
	var objmap map[string]json.RawMessage
	err = json.Unmarshal(res, &objmap)
	if err != nil {
		return nil, err
	}
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["findLoadDataJobById"], &loadDataJobsFound)
	if err != nil {
		return nil, err
	}
	if len(loadDataJobsFound) > 0 {
		return &loadDataJobsFound[0], nil
	}
	return nil, nil
}
