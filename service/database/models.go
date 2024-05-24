package database

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Photo struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	URL       string `json:"url"`
	Timestamp string `json:"timestamp"`
}

type Like struct {
	ID      int `json:"id"`
	UserID  int `json:"user_id"`
	PhotoID int `json:"photo_id"`
}

type Comment struct {
	ID          int    `json:"id"`
	UserId      int    `json:"user_id"`
	PhotoId     int    `json:"photo_id"`
	CommentText string `json:"comment_text"`
	Timestamp   string `json:"timestamp"`
}

type Follower struct {
	ID         int `json:"id"`
	FollowerID int `json:"follower_id"`
	FollowedID int `json:"followed_id"`
}

type Ban struct {
	ID       int `json:"id"`
	UserID   int `json:"user_id"`
	BannedID int `json:"banned_id"`
}
