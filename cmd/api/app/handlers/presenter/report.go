package presenter

// import "users_example/internal/supervisor/report"

// type jsonReport struct {
// 	Id             string       `json:"id"`
// 	CreatedAt      string       `json:"created_at"`
// 	LazyDevelopers int          `json:"lazy_developers"`
// 	Details        []jsonDetail `json:"details"`
// }

// type jsonDetail struct {
// 	DeveloperName  string `json:"developer_name"`
// 	Seniority      string `json:"seniority"`
// 	TaskStatus     string `json:"task_status"`
// 	LastActivityAt string `json:"last_activity_at"`
// }

// func Report(r *report.Report) jsonReport {
// 	toRet := jsonReport{
// 		Id:             r.Id,
// 		CreatedAt:      r.CreatedAt.String(),
// 		LazyDevelopers: r.LazyDevelopers,
// 		Details:        []jsonDetail{},
// 	}

// 	for _, detail := range r.Details {
// 		toRet.Details = append(toRet.Details, jsonDetail{
// 			DeveloperName:  detail.DeveloperName,
// 			Seniority:      detail.Seniority,
// 			TaskStatus:     detail.TaskStatus,
// 			LastActivityAt: detail.LastActivityAt.String(),
// 		})
// 	}

// 	return toRet
// }
