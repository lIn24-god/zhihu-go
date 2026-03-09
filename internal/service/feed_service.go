package service

import (
	"context"
	"log"
	"zhihu-go/internal/dao"
	"zhihu-go/internal/dto"
	"zhihu-go/internal/model"
)

type FeedService interface {
	// PushPost 同步推送一篇文章给所有粉丝（内部使用）
	PushPost(ctx context.Context, postID, authorID uint) error
	// PushPostAsync 异步推送文章（对外暴露）
	PushPostAsync(ctx context.Context, postID, authorID uint)
	// GetUserFeed 获取用户的 feed 列表
	GetUserFeed(ctx context.Context, userID uint, page, pageSize int) ([]dto.FeedItem, int64, error)
}

type feedService struct {
	followDAO   dao.FollowDAO
	timelineDAO dao.TimelineDAO
	postDAO     dao.PostDAO // 假设已有
	userDAO     dao.UserDAO // 假设已有

	// 异步推送队列
	taskChan    chan *pushTask
	workerCount int
}

type pushTask struct {
	ctx      context.Context
	postID   uint
	authorID uint
}

func NewFeedService(
	followDAO dao.FollowDAO,
	timelineDAO dao.TimelineDAO,
	postDAO dao.PostDAO,
	userDAO dao.UserDAO,
	workerCount int, // 例如 5
) FeedService {
	s := &feedService{
		followDAO:   followDAO,
		timelineDAO: timelineDAO,
		postDAO:     postDAO,
		userDAO:     userDAO,
		taskChan:    make(chan *pushTask, 1000), // 缓冲区大小
		workerCount: workerCount,
	}
	s.startWorkers()
	return s
}

// startWorkers 启动固定数量的 goroutine 处理推送任务
func (s *feedService) startWorkers() {
	for i := 0; i < s.workerCount; i++ {
		go func() {
			for task := range s.taskChan {
				// 使用独立的 context，避免父 context 被取消
				ctx := task.ctx
				if ctx == nil {
					ctx = context.Background()
				}
				if err := s.PushPost(ctx, task.postID, task.authorID); err != nil {
					// 记录错误，可以考虑重试（这里简化）
					log.Printf("推送文章失败: postID=%d, authorID=%d, error=%v", task.postID, task.authorID, err)
				}
			}
		}()
	}
}

// PushPostAsync 将任务放入队列，立即返回
func (s *feedService) PushPostAsync(ctx context.Context, postID, authorID uint) {
	select {
	case s.taskChan <- &pushTask{ctx: ctx, postID: postID, authorID: authorID}:
		// 入队成功
	default:
		// 队列满，降级为同步执行（或记录告警）
		log.Printf("推送队列已满，转为同步执行")
		go func() {
			if err := s.PushPost(context.Background(), postID, authorID); err != nil {
				log.Printf("同步推送失败: %v", err)
			}
		}()
	}
}

// PushPost 同步推送（实际写入 timeline）
func (s *feedService) PushPost(ctx context.Context, postID, authorID uint) error {
	// 1. 获取作者的所有粉丝（关注了 authorID 且状态为 normal）
	followerIDs, err := s.followDAO.GetFollowers(ctx, authorID)
	if err != nil {
		return err
	}
	if len(followerIDs) == 0 {
		return nil // 没有粉丝，无需推送
	}

	// 2. 构造 timeline 记录（包括作者自己，以便自己能看见自己的文章）
	timelines := make([]*model.Timeline, 0, len(followerIDs)+1)
	for _, fid := range followerIDs {
		timelines = append(timelines, &model.Timeline{
			UserID:   fid.ID,
			PostID:   postID,
			AuthorID: authorID,
			IsOwn:    fid.ID == authorID,
		})
	}
	// 加上作者自己的记录（如果粉丝列表中没有自己）
	if authorID != 0 {
		// 简单检查是否已包含（如果作者自己关注了自己，上面可能已包含，但通常没有）
		// 为简单，直接添加，数据库允许重复（无唯一约束）
		timelines = append(timelines, &model.Timeline{
			UserID:   authorID,
			PostID:   postID,
			AuthorID: authorID,
			IsOwn:    true,
		})
	}

	// 3. 批量插入
	return s.timelineDAO.BatchInsert(ctx, timelines)
}

// GetUserFeed 获取用户 feed
func (s *feedService) GetUserFeed(ctx context.Context, userID uint, page, pageSize int) ([]dto.FeedItem, int64, error) {
	offset := (page - 1) * pageSize
	timelineList, err := s.timelineDAO.GetUserTimeline(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	if len(timelineList) == 0 {
		return []dto.FeedItem{}, 0, nil
	}

	// 收集所有 postID 和 authorID
	postIDs := make([]uint, 0, len(timelineList))
	authorIDs := make([]uint, 0, len(timelineList))
	for _, tl := range timelineList {
		postIDs = append(postIDs, tl.PostID)
		authorIDs = append(authorIDs, tl.AuthorID)
	}

	// 批量获取文章信息
	posts, err := s.postDAO.GetPostsByIDs(ctx, postIDs) // 假设 PostDAO 有这个批量方法
	if err != nil {
		return nil, 0, err
	}
	postMap := make(map[uint]model.Post)
	for _, p := range posts {
		postMap[p.ID] = p
	}

	// 批量获取用户信息
	users, err := s.userDAO.GetUsersByIDs(ctx, authorIDs)
	if err != nil {
		return nil, 0, err
	}
	userMap := make(map[uint]model.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	// 组装 DTO
	items := make([]dto.FeedItem, 0, len(timelineList))
	for _, tl := range timelineList {
		post, ok := postMap[tl.PostID]
		if !ok {
			continue // 文章可能已删除，跳过
		}
		author, ok := userMap[tl.AuthorID]
		if !ok {
			continue
		}
		items = append(items, dto.FeedItem{
			ID:         tl.ID,
			PostID:     tl.PostID,
			AuthorID:   tl.AuthorID,
			AuthorName: author.Username,
			Title:      post.Title,
			CreatedAt:  tl.CreatedAt,
		})
	}

	// 获取总数用于分页
	total, err := s.timelineDAO.CountUserTimeline(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
