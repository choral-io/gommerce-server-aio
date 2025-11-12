package chats_v1beta

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nats-io/nats.go"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/choral-io/gommerce-server-core/data"
	"github.com/choral-io/gommerce-server-core/logging"
	"github.com/choral-io/gommerce-server-core/secure"

	"github.com/choral-io/gommerce-server-aio/data/models"

	chats_pb "github.com/choral-io/gommerce-protobuf-go/chats/v1beta"
	gender "github.com/choral-io/gommerce-protobuf-go/types/v1/gender"
	sqlpb "github.com/choral-io/gommerce-protobuf-go/types/v1/sqlpb"
)

type ChatsServiceServer struct {
	chats_pb.UnimplementedChatsServiceServer

	bdb    bun.IDB
	nsc    *nats.Conn
	logger logging.Logger
}

func NewChatsServiceServer(bdb bun.IDB, nsc *nats.Conn, logger logging.Logger) (chats_pb.ChatsServiceServer, error) {
	s := &ChatsServiceServer{
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

func (s *ChatsServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&chats_pb.ChatsService_ServiceDesc, s)
}

func (s *ChatsServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return chats_pb.RegisterChatsServiceHandler(ctx, mux, conn)
}

func (s *ChatsServiceServer) Authorize(ctx context.Context, _ string) error {
	return secure.Authorize(ctx, secure.AuthFuncAuthenticated, secure.AuthFuncRequireSchema(secure.AuthSchemaBearer))
}

func (s *ChatsServiceServer) ListSessions(ctx context.Context, req *chats_pb.ListSessionsRequest) (*chats_pb.ListSessionsResponse, error) {
	user := secure.IdentityFromContext(ctx)
	var members []models.ChatMember
	total, err := s.bdb.NewSelect().Model(&members).
		Relation("Session").Relation("Session.Members").Relation("Session.Members.Profile").
		Where("user_id = ?", user.Token().Subject()).Apply(data.WithPaging(req)).ScanAndCount(ctx)
	if err != nil {
		return nil, err
	}
	res := &chats_pb.ListSessionsResponse{
		Page:  req.Page,
		Size:  req.Size,
		Total: int64(total),
		Items: make([]*chats_pb.Session, len(members)),
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
		res.Items[i] = &chats_pb.Session{
			Id:           m.SessionId,
			Readonly:     m.Session.Readonly,
			CreatedAt:    timestamppb.New(m.Session.CreatedAt),
			UpdatedAt:    sqlpb.FromNullTime(m.Session.UpdatedAt),
			DeletedAt:    sqlpb.FromNullTime(m.Session.DeletedAt),
			IconUrl:      sqlpb.FromNullString(m.Session.IconUrl),
			Title:        m.Session.Title,
			Introduction: sqlpb.FromNullString(m.Session.Introduction),
			Members:      make([]*chats_pb.Member, len(m.Session.Members)),
			ReadCursor:   sqlpb.FromNullString(m.ReadCursor),
		}
		for j, m := range m.Session.Members {
			res.Items[i].Members[j] = &chats_pb.Member{
				UserId:      m.UserId,
				SessionId:   m.SessionId,
				CreatedAt:   timestamppb.New(m.CreatedAt),
				UpdatedAt:   sqlpb.FromNullTime(m.UpdatedAt),
				Permission:  m.Permission,
				Gender:      sqlpb.EnumFromNullName[gender.Gender](m.Profile.Gender),
				DisplayName: m.DisplayName.String,
				AvatarUrl:   sqlpb.FromNullString(m.Profile.AvatarUrl),
			}
			if !m.DisplayName.Valid && m.Profile != nil {
				res.Items[i].Members[j].DisplayName = m.Profile.DisplayName
			}
		}
		for _, r := range records {
			if r.SessionId == m.SessionId {
				res.Items[i].LastRecord = &chats_pb.Record{
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

func (s *ChatsServiceServer) DescribeSession(ctx context.Context, req *chats_pb.DescribeSessionRequest) (*chats_pb.DescribeSessionResponse, error) {
	user := secure.IdentityFromContext(ctx)
	member := models.ChatMember{
		SessionId: req.SessionId,
		UserId:    user.Token().Subject(),
	}
	if err := s.bdb.NewSelect().Model(&member).
		Relation("Session").Relation("Session.Members").Relation("Session.Members.Profile").
		WherePK("session_id", "user_id").Scan(ctx); errors.Is(err, sql.ErrNoRows) {
		return &chats_pb.DescribeSessionResponse{Item: nil}, nil
	} else if err != nil {
		return nil, err
	}
	session := &chats_pb.Session{
		Id:           member.SessionId,
		Readonly:     member.Session.Readonly,
		CreatedAt:    timestamppb.New(member.Session.CreatedAt),
		UpdatedAt:    sqlpb.FromNullTime(member.Session.UpdatedAt),
		DeletedAt:    sqlpb.FromNullTime(member.Session.DeletedAt),
		IconUrl:      sqlpb.FromNullString(member.Session.IconUrl),
		Title:        member.Session.Title,
		Introduction: sqlpb.FromNullString(member.Session.Introduction),
		Members:      make([]*chats_pb.Member, len(member.Session.Members)),
		ReadCursor:   sqlpb.FromNullString(member.ReadCursor),
	}
	records := make([]models.ChatRecord, 0, 1)
	rids := s.bdb.NewSelect().Model((*models.ChatRecord)(nil)).ColumnExpr("MAX(id) AS id").Where("session_id = ?", session.Id)
	if err := s.bdb.NewSelect().Model(&records).Where("id IN (?)", rids).Scan(ctx); err != nil {
		return nil, err
	}
	for j, m := range member.Session.Members {
		session.Members[j] = &chats_pb.Member{
			UserId:      m.UserId,
			SessionId:   m.SessionId,
			CreatedAt:   timestamppb.New(m.CreatedAt),
			UpdatedAt:   sqlpb.FromNullTime(m.UpdatedAt),
			Permission:  m.Permission,
			Gender:      sqlpb.EnumFromNullName[gender.Gender](m.Profile.Gender),
			DisplayName: m.DisplayName.String,
			AvatarUrl:   sqlpb.FromNullString(m.Profile.AvatarUrl),
		}
		if !m.DisplayName.Valid && m.Profile != nil {
			session.Members[j].DisplayName = m.Profile.DisplayName
		}
	}
	if len(records) > 0 {
		r := records[0]
		session.LastRecord = &chats_pb.Record{
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
	return &chats_pb.DescribeSessionResponse{Item: session}, nil
}

func (s *ChatsServiceServer) ReadSession(ctx context.Context, req *chats_pb.ReadSessionRequest) (*chats_pb.ReadSessionResponse, error) {
	userId := secure.IdentityFromContext(ctx).Token().Subject()
	if _, err := s.bdb.NewUpdate().
		Model((*models.ChatMember)(nil)).
		Set("read_cursor = ?", sqlpb.ToNullString(req.Cursor)).
		Where("session_id = ? AND user_id = ?", req.SessionId, userId).
		Exec(ctx); err != nil {
		return nil, err
	}
	return &chats_pb.ReadSessionResponse{}, nil
}

func (s *ChatsServiceServer) SendRecord(ctx context.Context, req *chats_pb.SendRecordRequest) (*chats_pb.SendRecordResponse, error) {
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
	event := &chats_pb.Record{
		Id:        record.Id,
		SessionId: record.SessionId,
		CreatorId: record.CreatorId,
		CreatedAt: timestamppb.New(record.CreatedAt),
		Version:   record.Version,
		Headers:   record.Headers,
		Content:   record.Content,
	}
	bytes, _ := proto.Marshal(event)
	if err := s.nsc.Publish(fmt.Sprintf("chat.records.s.%s", req.SessionId), bytes); err != nil {
		s.logger.Error(ctx, "failed to publish chat record", "error", err)
	}
	return &chats_pb.SendRecordResponse{
		Id: event.Id,
	}, nil
}

func (s *ChatsServiceServer) ListRecords(ctx context.Context, req *chats_pb.ListRecordsRequest) (*chats_pb.ListRecordsResponse, error) {
	records := make([]models.ChatRecord, 0)
	query := s.bdb.NewSelect().Model(&records)
	if req.SessionId.GetValue() != "" {
		query.Where("session_id = ?", req.SessionId.GetValue())
	}
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	res := &chats_pb.ListRecordsResponse{
		Items: make([]*chats_pb.Record, len(records)),
	}
	for i, r := range records {
		res.Items[i] = &chats_pb.Record{
			Id:        r.Id,
			SessionId: r.SessionId,
			CreatorId: r.CreatorId,
			CreatedAt: timestamppb.New(r.CreatedAt),
			UpdatedAt: sqlpb.FromNullTime(r.UpdatedAt),
			Version:   r.Version,
			Headers:   r.Headers,
			Content:   r.Content,
		}
	}
	return res, nil
}

func (s *ChatsServiceServer) WatchRecords(_ *chats_pb.WatchRecordsRequest, srv chats_pb.ChatsService_WatchRecordsServer) error {
	user := secure.IdentityFromContext(srv.Context())
	csc := make(chan *nats.Msg, 64) // buffered to reduce NATS backpressure
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
	keepaliveTimer := time.NewTimer(25 * time.Second)
	defer func() {
		if !keepaliveTimer.Stop() {
			select {
			case <-keepaliveTimer.C:
			default:
			}
		}
	}()
	for {
		select {
		case msg := <-csc:
			record := &chats_pb.Record{}
			if err := proto.Unmarshal(msg.Data, record); err != nil {
				s.logger.Error(srv.Context(), "failed to unmarshal chat record", "error", err)
			} else if err := srv.Send(&chats_pb.WatchRecordsResponse{Items: []*chats_pb.Record{record}}); err != nil {
				s.logger.Error(srv.Context(), "failed to send chat records", "error", err)
				return err
			}
			if !keepaliveTimer.Stop() {
				select {
				case <-keepaliveTimer.C:
				default:
				}
			}
			keepaliveTimer.Reset(25 * time.Second)
		case <-keepaliveTimer.C: // send empty message to keep the stream alive
			if err := srv.Send(&chats_pb.WatchRecordsResponse{Items: []*chats_pb.Record{}}); err != nil {
				s.logger.Error(srv.Context(), "failed to send chat records", "error", err)
				return err
			}
			keepaliveTimer.Reset(25 * time.Second)
		case <-srv.Context().Done():
			if err := srv.Context().Err(); errors.Is(err, context.Canceled) {
				return nil
			} else {
				return err
			}
		}
	}
}
