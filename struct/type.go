package types

// Result 表示文件的详细信息
type Result struct {
	FileID       string `json:"file_id"`        // 文件 ID
	FileUniqueID string `json:"file_unique_id"` // 文件唯一 ID
	FileSize     int    `json:"file_size"`      // 文件大小
	FilePath     string `json:"file_path"`      // 文件路径
}

// Data 表示 API 返回的数据结构
type Data struct {
	Ok     bool   `json:"ok"`     // 请求是否成功
	Result Result `json:"result"` // 文件结果信息
}

// Image 表示图片的源信息
type Image struct {
	Src string `json:"src"` // 图片源 URL
}
