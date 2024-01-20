package types

type Result struct {
	File_id        string `json:"file_id"`
	File_unique_id string `json:"file_unique_id"`
	File_size      int    `json:"file_size"`
	File_path      string `json:"file_path"`
}

type Data struct {
	Ok     bool   `json:"ok"`
	Result Result `json:"result"`
}
type Image struct {
	Src string `json:"src"`
}
