package session

import (
	v4 "github.com/ExerciseCoding/template/internal/web_server/v4"
	"github.com/google/uuid"
)

type Manager struct {
	Propagator
	Store
	CtxSessKey string
}

func (m *Manager) GetSession(ctx *v4.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 1)
	}
	// ctx.Req.Context().Value(m.CtxSessKey)
	val, ok := ctx.UserValues[m.CtxSessKey]
	if ok {
		return val.(Session), nil
	}
	// 尝试缓存住session
	sessId, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}

    sess, err := m.Get(ctx.Req.Context(), sessId)
	if err != nil {
		return nil, err
	}

	ctx.UserValues[m.CtxSessKey] = sess
	// ctx.Req = ctx.Req.WithContext(context.WithValue(ctx.Req.Context(), m.CtxSessKey, sess))
	return sess, err
}

func (m *Manager) InitSession(ctx *v4.Context) (Session, error) {
	id := uuid.New().String()
	sess, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}

	// 注入进去HTTP响应里面
	err = m.Inject(id, ctx.Resp)
	return sess, err
}


func (m *Manager) RefreshSession(ctx *v4.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	return m.Refresh(ctx.Req.Context(), sess.ID())
}

func (m *Manager) RemoveSession(ctx *v4.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Store.Remove(ctx.Req.Context(), sess.ID())
	if err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Resp)
}