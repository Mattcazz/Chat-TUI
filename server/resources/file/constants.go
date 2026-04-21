package file

type serverPath string

const (
	TmpUploadsPath   serverPath = "./uploads/tmp/"
	FinalUploadsPath serverPath = "./uploads/assembled/"
)

type FileStatus string

const (
	FileStatusUploading FileStatus = "uploading"
	FileStatusReady     FileStatus = "ready"
	FileStatusExpired   FileStatus = "expired"
)

type UploadSessionStatus string

const (
	FileSessionStatusUploading  UploadSessionStatus = "uploading"
	FileSessionStatusAssembling UploadSessionStatus = "assembling"
	FileSessionStatusCompleted  UploadSessionStatus = "completed"
	FileSessionStatusCanceled   UploadSessionStatus = "canceled"
	FileSessionStatusPending    UploadSessionStatus = "pending"
)

const TimeToExpireUploadSession int32 = 24 * 60 * 60 // 24 hours in seconds
