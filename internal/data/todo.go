package data

import (
	"context"
	"time"

	"shunfeng-miniprogram/internal/biz"
)

// todoRepo 待办事项仓储的 GORM 实现。
type todoRepo struct{}

// NewTodoRepo 创建待办事项仓储。
func NewTodoRepo() biz.TodoRepo {
	return &todoRepo{}
}

func (r *todoRepo) FindByID(ctx context.Context, id int64) (*biz.Todo, error) {
	var po Todo
	if err := DB.WithContext(ctx).First(&po, id).Error; err != nil {
		return nil, biz.ErrTodoNotFound
	}
	return toTodoBiz(&po), nil
}

func (r *todoRepo) ListTodos(ctx context.Context, opts ...biz.ListOption) ([]*biz.Todo, error) {
	options := biz.ListOptions{Limit: 20}
	for _, opt := range opts {
		opt(&options)
	}
	if options.Offset < 0 || options.Limit <= 0 {
		return nil, biz.ErrTodoInvalidArgument
	}

	var pos []Todo
	tx := DB.WithContext(ctx).Order("id ASC")
	if options.Limit > 0 {
		tx = tx.Limit(options.Limit)
	}
	if options.Offset > 0 {
		tx = tx.Offset(options.Offset)
	}
	if err := tx.Find(&pos).Error; err != nil {
		return nil, err
	}

	todos := make([]*biz.Todo, len(pos))
	for i, po := range pos {
		todos[i] = toTodoBiz(&po)
	}
	return todos, nil
}

func (r *todoRepo) CreateTodo(ctx context.Context, todo *biz.Todo) (*biz.Todo, error) {
	po := newTodo(todo)
	po.CreateTime = time.Now()
	po.UpdateTime = time.Now()

	if err := DB.WithContext(ctx).Create(po).Error; err != nil {
		return nil, err
	}
	return toTodoBiz(po), nil
}

func (r *todoRepo) UpdateTodo(ctx context.Context, todo *biz.Todo) (*biz.Todo, error) {
	var po Todo
	if err := DB.WithContext(ctx).First(&po, todo.ID).Error; err != nil {
		return nil, biz.ErrTodoNotFound
	}

	updated := newTodo(todo)
	updated.ID = po.ID
	updated.CreateTime = po.CreateTime
	updated.UpdateTime = time.Now()

	if err := DB.WithContext(ctx).Model(&po).Updates(updated).Error; err != nil {
		return nil, err
	}
	return toTodoBiz(updated), nil
}

func (r *todoRepo) DeleteTodo(ctx context.Context, id int64) error {
	result := DB.WithContext(ctx).Delete(&Todo{}, id)
	if result.RowsAffected == 0 {
		return biz.ErrTodoNotFound
	}
	return result.Error
}

// newTodo 将 biz.Todo 转换为 Todo（DO → PO）。
func newTodo(in *biz.Todo) *Todo {
	if in == nil {
		return nil
	}
	return &Todo{
		Title:     in.Title,
		Content:   in.Content,
		Completed: in.Completed,
	}
}

// toTodoBiz 将 Todo 转换为 biz.Todo（PO → DO）。
func toTodoBiz(in *Todo) *biz.Todo {
	if in == nil {
		return nil
	}
	return &biz.Todo{
		ID:         in.ID,
		Title:      in.Title,
		Content:    in.Content,
		Completed:  in.Completed,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}
}
