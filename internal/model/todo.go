package model

// Todo is standard model
type Todo struct {
	ID        int    `db:"id" json:"id"`
	Title     string `db:"title" json:"title"`
	IsDone    int    `db:"is_done" json:"isDone"`
	Detail    string `db:"detail" json:"detail"`
	CreatedAt string `db:"created_at" json:"createdAt"`
	UpdatedAt string `db:"updated_at" json:"updatedAt"`
}
