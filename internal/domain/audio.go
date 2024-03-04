package domain

type Audio struct {
	Id      int    `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Idx     string `db:"idx" json:"idx"`
	VideoId int    `db:"video_id" json:"video_id"`
}
