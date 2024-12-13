package model

type News struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TopicName   string `json:"topic_name"`
}

type UserWithNews struct {
	Email    string `json:"email"`
	NewsList []News `json:"news_list,omitempty"`
}
