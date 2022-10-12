package app

type NotificationDetail struct {
	NotificationUrl    string `json:"url"`
	NotificationReason string `json:"reason""`
	NotificationType   string `json:"type"`
	Repository         string `json:"repository"`
	NotificationTitle  string `json:"title""`
	Pull               string `json:"pull""`
}

type Pull struct {
	State        string `json:"state"`
	Merged       bool   `json:"merged"`
	User         string `json:"user"`
	Commits      int    `json:"commits"`
	Comments     int    `json:"comments"`
	Additions    int    `json:"additions"`
	Deletions    int    `json:"deletions"`
	ChangedFiles int    `json:"changed_files"`
}
