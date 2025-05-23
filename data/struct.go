package data

type Post struct {
  PostID     int
  Title      string
  Content    string
  Categories []string
  Comments   []CommentWithLike
}

type PostWithLike struct {
  Post
  IsLike       int
  LikeCount    int
  DislikeCount int
}

type Comment struct {
  CommentID int
  Content   string
}

type CommentWithLike struct {
  Comment
  IsLike       int
  LikeCount    int
  DislikeCount int
}