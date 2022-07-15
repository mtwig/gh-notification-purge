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
	Commits      int8   `json:"commits"`
	Comments     int8   `json:"comments"`
	Additions    int8   `json:"additions"`
	Deletions    int8   `json:"deletions"`
	ChangedFiles int8   `json:"changed_files"`
}
