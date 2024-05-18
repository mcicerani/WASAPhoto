package database

type UserDetails struct {
	Username string `json:"username"`
}

type Identifier struct {
	UserId    string `json:"user_id"`
	IsNewUser bool   `json:"is_new_user"`
}

type UserProfile struct {
	Username       string   `json:"username"`
	UserId         string   `json:"user_id"`
	FollowerCount  int      `json:"follower_count"`
	Followers      []string `json:"followers"`
	FollowingCount int      `json:"following_count"`
	Follows        []string `json:"follows"`
	PhotosCount    int      `json:"photos_count"`
	BannedUser     []string `json:"banned_user"`
}

type Photo struct {
	PhotosId    string   `json:"photos_id"`
	Url         string   `json:"url"`
	Timestamp   string   `json:"timestamp"`
	LikesNumber int      `json:"username"`
	Comments    []string `json:"comments"`
}

type Comment struct {
	UserId    string `json:"user_id"`
	PhotosId  string `json:"photos_id"`
	CommentId string `json:"comment_id"`
	Text      string `json:"text_comment"`
}
