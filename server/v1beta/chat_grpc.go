package v1beta

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	chats "github.com/choral-io/gommerce-protobuf-go/chats/v1beta"
	gender_v1 "github.com/choral-io/gommerce-protobuf-go/types/v1/gender"
	sqlpb "github.com/choral-io/gommerce-protobuf-go/types/v1/sqlpb"
	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-core/data"
	"github.com/choral-io/gommerce-server-core/logging"
	"github.com/choral-io/gommerce-server-core/secure"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nats-io/nats.go"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type chatsServiceServer struct {
	chats.UnimplementedChatsServiceServer

	bdb    bun.IDB
	nsc    *nats.Conn
	logger logging.Logger
}

func NewChatsServiceServer(bdb bun.IDB, nsc *nats.Conn, logger logging.Logger) (chats.ChatsServiceServer, error) {
	s := &chatsServiceServer{
		bdb:    bdb,
		nsc:    nsc,
		logger: logger,
	}
	if _, err := nsc.Subscribe("chat.records.s.*", func(msg *nats.Msg) {
		sid := strings.Split(msg.Subject, ".")[3]
		cms := make([]string, 0)
		if err := bdb.NewSelect().Model((*models.ChatMember)(nil)).
			Column("user_id").Where("session_id = ?", sid).
			Scan(context.Background(), &cms); err != nil {
			s.logger.Error(context.Background(), "failed to fetch chat members", "error", err)
			return
		}
		for _, cm := range cms {
			if err := nsc.Publish(fmt.Sprintf("chat.records.u.%s", cm), msg.Data); err != nil {
				s.logger.Error(context.Background(), "failed to publish chat record", "error", err)
			}
		}
	}); err != nil {
		return nil, fmt.Errorf("failed to subscribe to chat.records.s.*: %w", err)
	}
	return s, nil
}

func (s *chatsServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&chats.ChatsService_ServiceDesc, s)
}

func (s *chatsServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return chats.RegisterChatsServiceHandler(ctx, mux, conn)
}

func (s *chatsServiceServer) Authorize(ctx context.Context, procedure string) error {
	return secure.Authorize(ctx, secure.AuthFuncAuthenticated, secure.AuthFuncRequireSchema(secure.AUTH_SCHEMA_BEARER))
}

func (s *chatsServiceServer) ListSessions(ctx context.Context, req *chats.ListSessionsRequest) (*chats.ListSessionsResponse, error) {
	user := secure.IdentityFromContext(ctx)
	members := []models.ChatMember{}
	total, err := s.bdb.NewSelect().Model(&members).
		Relation("Session").Relation("Session.Members").Relation("Session.Members.Profile").
		Where("user_id = ?", user.Token().Subject()).Apply(data.WithPaging(req)).ScanAndCount(ctx)
	if err != nil {
		return nil, err
	}
	res := &chats.ListSessionsResponse{
		Page:  req.Page,
		Size:  req.Size,
		Total: int64(total),
		Items: make([]*chats.Session, len(members)),
	}
	sids := make([]string, len(members))
	for i, m := range members {
		sids[i] = m.SessionId
	}
	records := make([]models.ChatRecord, 0, len(members))
	rids := s.bdb.NewSelect().Model((*models.ChatRecord)(nil)).ColumnExpr("MAX(id) AS id").Where("session_id IN (?)", bun.In(sids))
	if err := s.bdb.NewSelect().Model(&records).Where("id IN (?)", rids).Scan(ctx); err != nil {
		return nil, err
	}
	for i, m := range members {
		res.Items[i] = &chats.Session{
			Id:           m.SessionId,
			Readonly:     m.Session.Readonly,
			CreatedAt:    timestamppb.New(m.Session.CreatedAt),
			UpdatedAt:    sqlpb.FromNullTime(m.Session.UpdatedAt),
			DeletedAt:    sqlpb.FromNullTime(m.Session.DeletedAt),
			IconUrl:      sqlpb.FromNullString(m.Session.IconUrl),
			Title:        m.Session.Title,
			Introduction: sqlpb.FromNullString(m.Session.Introduction),
			Members:      make([]*chats.Member, len(m.Session.Members)),
			ReadCursor:   sqlpb.FromNullString(m.ReadCursor),
		}
		for j, m := range m.Session.Members {
			res.Items[i].Members[j] = &chats.Member{
				UserId:      m.UserId,
				SessionId:   m.SessionId,
				CreatedAt:   timestamppb.New(m.CreatedAt),
				UpdatedAt:   sqlpb.FromNullTime(m.UpdatedAt),
				Permission:  m.Permission,
				Gender:      gender_v1.FromSqlNullString(m.Profile.Gender),
				DisplayName: m.DisplayName.String,
				AvatarUrl:   sqlpb.FromNullString(m.Profile.AvatarUrl),
			}
			if !m.DisplayName.Valid && m.Profile != nil {
				res.Items[i].Members[j].DisplayName = m.Profile.DisplayName
			}
		}
		for _, r := range records {
			if r.SessionId == m.SessionId {
				res.Items[i].LastRecord = &chats.Record{
					Id:        r.Id,
					SessionId: r.SessionId,
					CreatorId: r.CreatorId,
					CreatedAt: timestamppb.New(r.CreatedAt),
					UpdatedAt: sqlpb.FromNullTime(r.UpdatedAt),
					DeletedAt: sqlpb.FromNullTime(r.DeletedAt),
					Version:   r.Version,
					Headers:   r.Headers,
					Content:   r.Content,
				}
				break
			}
		}
	}
	return res, nil
}

func (s *chatsServiceServer) DescribeSession(ctx context.Context, req *chats.DescribeSessionRequest) (*chats.DescribeSessionResponse, error) {
	user := secure.IdentityFromContext(ctx)
	member := models.ChatMember{
		SessionId: req.SessionId,
		UserId:    user.Token().Subject(),
	}
	if err := s.bdb.NewSelect().Model(&member).
		Relation("Session").Relation("Session.Members").Relation("Session.Members.Profile").
		WherePK("session_id", "user_id").Scan(ctx); errors.Is(err, sql.ErrNoRows) {
		return &chats.DescribeSessionResponse{Item: nil}, nil
	} else if err != nil {
		return nil, err
	}
	session := &chats.Session{
		Id:           member.SessionId,
		Readonly:     member.Session.Readonly,
		CreatedAt:    timestamppb.New(member.Session.CreatedAt),
		UpdatedAt:    sqlpb.FromNullTime(member.Session.UpdatedAt),
		DeletedAt:    sqlpb.FromNullTime(member.Session.DeletedAt),
		IconUrl:      sqlpb.FromNullString(member.Session.IconUrl),
		Title:        member.Session.Title,
		Introduction: sqlpb.FromNullString(member.Session.Introduction),
		Members:      make([]*chats.Member, len(member.Session.Members)),
		ReadCursor:   sqlpb.FromNullString(member.ReadCursor),
	}
	records := make([]models.ChatRecord, 0, 1)
	rids := s.bdb.NewSelect().Model((*models.ChatRecord)(nil)).ColumnExpr("MAX(id) AS id").Where("session_id = ?", session.Id)
	if err := s.bdb.NewSelect().Model(&records).Where("id IN (?)", rids).Scan(ctx); err != nil {
		return nil, err
	}
	for j, m := range member.Session.Members {
		session.Members[j] = &chats.Member{
			UserId:      m.UserId,
			SessionId:   m.SessionId,
			CreatedAt:   timestamppb.New(m.CreatedAt),
			UpdatedAt:   sqlpb.FromNullTime(m.UpdatedAt),
			Permission:  m.Permission,
			Gender:      gender_v1.FromSqlNullString(m.Profile.Gender),
			DisplayName: m.DisplayName.String,
			AvatarUrl:   sqlpb.FromNullString(m.Profile.AvatarUrl),
		}
		if !m.DisplayName.Valid && m.Profile != nil {
			session.Members[j].DisplayName = m.Profile.DisplayName
		}
	}
	if len(records) > 0 {
		r := records[0]
		session.LastRecord = &chats.Record{
			Id:        r.Id,
			SessionId: r.SessionId,
			CreatorId: r.CreatorId,
			CreatedAt: timestamppb.New(r.CreatedAt),
			UpdatedAt: sqlpb.FromNullTime(r.UpdatedAt),
			DeletedAt: sqlpb.FromNullTime(r.DeletedAt),
			Version:   r.Version,
			Headers:   r.Headers,
			Content:   r.Content,
		}
	}
	return &chats.DescribeSessionResponse{Item: session}, nil
}

func (s *chatsServiceServer) ReadSession(ctx context.Context, req *chats.ReadSessionRequest) (*chats.ReadSessionResponse, error) {
	userId := secure.IdentityFromContext(ctx).Token().Subject()
	if _, err := s.bdb.NewUpdate().
		Model((*models.ChatMember)(nil)).
		Set("read_cursor = ?", sqlpb.ToNullString(req.ReadCursor)).
		Where("session_id = ? AND user_id = ?", req.SessionId, userId).
		Exec(ctx); err != nil {
		return nil, err
	}
	return &chats.ReadSessionResponse{}, nil
}

func (s *chatsServiceServer) SendRecord(ctx context.Context, req *chats.SendRecordRequest) (*chats.SendRecordResponse, error) {
	user := secure.IdentityFromContext(ctx)
	record := models.ChatRecord{
		SessionId: req.SessionId,
		CreatorId: user.Token().Subject(),
		Version:   "0.0.1", // version of the record is used to determine the format of the record
		Headers:   req.Headers,
		Content:   req.Content,
	}
	if _, err := s.bdb.NewInsert().Model(&record).Exec(ctx); err != nil {
		return nil, err
	}
	event := &chats.Record{
		Id:        record.Id,
		SessionId: record.SessionId,
		CreatorId: record.CreatorId,
		CreatedAt: timestamppb.New(record.CreatedAt),
		Version:   record.Version,
		Headers:   record.Headers,
		Content:   record.Content,
	}
	data, _ := proto.Marshal(event)
	if err := s.nsc.Publish(fmt.Sprintf("chat.records.s.%s", req.SessionId), data); err != nil {
		s.logger.Error(ctx, "failed to publish chat record", "error", err)
	}
	return &chats.SendRecordResponse{
		Id: event.Id,
	}, nil
}

func (s *chatsServiceServer) WatchRecords(_ *chats.WatchRecordsRequest, srv chats.ChatsService_WatchRecordsServer) error {
	user := secure.IdentityFromContext(srv.Context())
	csc := make(chan *nats.Msg)
	cls, err := s.nsc.ChanSubscribe(fmt.Sprintf("chat.records.u.%s", user.Token().Subject()), csc)
	if err != nil {
		return err
	}
	defer func() {
		if err := cls.Unsubscribe(); err != nil {
			s.logger.Error(srv.Context(), "failed to unsubscribe from chat.records", "error", err)
		} else {
			s.logger.Info(srv.Context(), "unsubscribed from chat.records")
		}
	}()
	for {
		select {
		case msg := <-csc:
			record := &chats.Record{}
			proto.Unmarshal(msg.Data, record)
			srv.Send(&chats.WatchRecordsResponse{
				Items: []*chats.Record{record},
			})
		case <-srv.Context().Done():
			if err := srv.Context().Err(); err == context.Canceled {
				return nil
			} else {
				return err
			}
		}
	}
}
