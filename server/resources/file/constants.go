package file

type serverPath string

const (
	tmpUploads   serverPath = "./uploads/tmp/"
	finalUploads serverPath = "./uploads/assembled/"
)
