package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/mail"
	"strings"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/app/email"
	"github.com/Kushian01100111/Tickermaster/internal/app/otpChallenge"
	"github.com/Kushian01100111/Tickermaster/internal/app/user"
	"github.com/Kushian01100111/Tickermaster/internal/domain/session"
	u "github.com/Kushian01100111/Tickermaster/internal/domain/user"
	"github.com/Kushian01100111/Tickermaster/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrRefreshInvalidString = errors.New("invalid refresh string")
	ErrRefreshInvalid       = errors.New("invalid refresh token")
	ErrRefreshSession       = errors.New("there is no session related to this token")
	ErrRefreshRevoked       = errors.New("this token has been already revoke")
	ErrRefreshExpired       = errors.New("this token has already expired")
	ErrEmailInvalid         = errors.New("invalid mail")
	ErrInvalidCreadentials  = errors.New("invalid creadentials")
	ErrMethodNotAllowed     = errors.New("this method of sign in is no available for this user")
	ErrHashRequired         = errors.New("hash is required")
	ErrUserRequired         = errors.New("user is required")
	ErrUserAlreadyExists    = errors.New("user is already created")
	ErrCreatedAtRequired    = errors.New("createdAt date required")
	ErrExpiredAtRequired    = errors.New("expiredAt date required")

	ErrInvalidOTP     = errors.New("invalid OTP")
	ErrEmptyOTP       = errors.New("empty otp Challenge")
	ErrOTPNotFound    = errors.New("otp Challenge not found")
	ErrOTPExpired     = errors.New("otp Challenge is expired")
	ErrOTPInvalidCode = errors.New("invalid otp Challenge code")

	ErrPasswordRequired = errors.New("password is required")
	ErrRole             = errors.New("invalid role type")
	ErrAuthMethod       = errors.New("invalid authmethod")
	ErrAuthMethodLen    = errors.New("invalid amount of authmethods")
)

type UserParams struct {
	Email      string
	Role       string
	Password   string
	AuthMethod string
}

type LoginParams struct {
	Email    string
	Password string
}

type VerifyParams struct {
	Email string
	Code  string
}

type Session struct {
	User         u.User
	AccessToken  string
	RefreshToken string
	ExpiresInSec int64
}

type AuthService interface {
	Login(ctx context.Context, params LoginParams) (*Session, error)
	Refresh(ctx context.Context, refresh string) (*Session, error)
	Logout(ctx context.Context, refresh string) error
	LogoutAll(ctx context.Context, refresh string) error
	CreateUser(ctx context.Context, params UserParams) (*Session, error)
	SignupRequest(ctx context.Context, email string) error
	SignupVerify(ctx context.Context, param VerifyParams) (*Session, error)
	LoginRequest(ctx context.Context, email string) error
	LoginVerify(ctx context.Context, param VerifyParams) (*Session, error)
}

type authService struct {
	otpSrv     otpChallenge.OTPService
	authRepo   repository.AuthRepository
	userSrv    user.UserService
	mailer     email.EmailSender
	hasher     session.PasswordHasher
	jwt        *session.JWTManager
	otpTTL     time.Duration
	refreshTTL time.Duration
}

type AuthConfig struct {
	OTPTTL     time.Duration
	RefreshTTL time.Duration
}

func NewAuthService(
	otpRepo otpChallenge.OTPService,
	authRepo repository.AuthRepository,
	userRepo user.UserService,
	mailer email.EmailSender,
	hasher session.PasswordHasher,
	jwt *session.JWTManager,
	config AuthConfig) AuthService {

	otpTTL := config.OTPTTL
	if otpTTL <= 0 {
		otpTTL = 10 * time.Minute
	}

	refreshTTL := config.RefreshTTL
	if refreshTTL <= 0 {
		refreshTTL = 30 * 24 * time.Hour
	}
	return &authService{
		otpSrv:     otpRepo,
		authRepo:   authRepo,
		userSrv:    userRepo,
		mailer:     mailer,
		hasher:     hasher,
		jwt:        jwt,
		otpTTL:     otpTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *authService) Login(ctx context.Context, params LoginParams) (*Session, error) {
	email := normalizeEmail(params.Email)
	pass := strings.TrimSpace(params.Password)

	if email == "" || pass == "" {
		return nil, ErrInvalidCreadentials
	}

	user, err := s.userSrv.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, err
	}

	if err := s.challengePass(*user, pass); err != nil {
		_ = s.userSrv.FailedLogin(ctx, user)
		return nil, err
	}
	_ = s.userSrv.ResetFailedLogin(ctx, user)

	return s.issueSession(ctx, user)
}

func (s *authService) Refresh(ctx context.Context, refresh string) (*Session, error) {
	refreshToken := strings.TrimSpace(refresh)
	if refreshToken == "" {
		return nil, ErrRefreshInvalidString
	}

	hash := sha256Hex(refresh)
	sess, err := s.authRepo.GetByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, ErrRefreshSession
	}
	if sess.RevokedAt != nil {
		return nil, ErrRefreshRevoked
	}
	if time.Now().After(sess.ExpiresAt) {
		return nil, ErrRefreshExpired
	}

	user, err := s.userSrv.GetUser(ctx, sess.UserID.Hex())
	if err != nil {
		return nil, ErrRefreshInvalid
	}

	err = s.authRepo.RevokeRefreshToken(ctx, sess.Hash)
	if err != nil {
		return nil, err
	}

	return s.issueSession(ctx, user)
}

func (s *authService) Logout(ctx context.Context, refresh string) error {
	refreshToken := strings.TrimSpace(refresh)
	if refresh == "" {
		return nil
	}
	hash := sha256Hex(refreshToken)
	return s.authRepo.RevokeRefreshToken(ctx, hash)
}

func (s *authService) LogoutAll(ctx context.Context, idHex string) error {
	oid, err := bson.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	return s.authRepo.RevokeAllByUserID(ctx, oid)
}

func (s *authService) CreateUser(ctx context.Context, params UserParams) (*Session, error) {
	if err := validateParam(params); err != nil {
		return nil, err
	}

	email := normalizeEmail(params.Email)
	if u, err := s.userSrv.GetByEmail(ctx, email); err == nil || u != nil {
		return nil, ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	temp := string(hash)

	u, err := s.userSrv.CreateUserB(ctx, &u.User{
		Email:        email,
		Role:         params.Role,
		PasswordHash: &temp,
	})
	if err != nil {
		return nil, err
	}

	return s.issueSession(ctx, u)
}

func (s *authService) SignupRequest(ctx context.Context, email string) error {
	mail := normalizeEmail(email)
	if mail == "" {
		return ErrEmailInvalid
	}

	if u, err := s.userSrv.GetByEmail(ctx, mail); err == nil || u != nil {
		return ErrUserAlreadyExists
	}

	code, err := new6DigitCode()
	if err != nil {
		return err
	}

	ch := otpChallenge.OTPParams{
		Email:     mail,
		CodeHash:  sha256Hex(code),
		ExpiresAt: time.Now().Add(s.otpTTL),
		Purpuse:   "signUp",
		Attempts:  0,
		CreatedAt: time.Now(),
	}
	if err := s.otpSrv.CreateOrReplace(ctx, ch); err != nil {
		return err
	}

	return s.mailer.SendSignUpCode(ctx, mail, code)
}

func (s *authService) SignupVerify(ctx context.Context, params VerifyParams) (*Session, error) {
	email := normalizeEmail(params.Email)
	code := strings.TrimSpace(params.Code)

	if email == "" || code == "" {
		return nil, ErrInvalidOTP
	}

	if u, err := s.userSrv.GetByEmail(ctx, email); err == nil || u != nil {
		return nil, ErrUserAlreadyExists
	}

	if err := s.verifyAndConsumeOTP(ctx, email, code, "signUp"); err != nil {
		return nil, err
	}

	u, err := s.userSrv.CreateUser(ctx, user.UserParams{
		Email:      email,
		Role:       "costumer",
		AuthMethod: "email_otp",
	})
	if err != nil {
		return nil, err
	}

	return s.issueSession(ctx, u)
}

func (s *authService) LoginRequest(ctx context.Context, email string) error {
	mail := normalizeEmail(email)
	if mail == "" {
		return ErrEmailInvalid
	}

	if u, err := s.userSrv.GetByEmail(ctx, mail); err != nil || u == nil {
		return err
	}

	code, err := new6DigitCode()
	if err != nil {
		return err
	}

	ch := otpChallenge.OTPParams{
		Email:     mail,
		CodeHash:  sha256Hex(code),
		ExpiresAt: time.Now().Add(s.otpTTL),
		Purpuse:   "login",
		Attempts:  0,
		CreatedAt: time.Now(),
	}
	if err := s.otpSrv.CreateOrReplace(ctx, ch); err != nil {
		return err
	}

	return s.mailer.SendLoginCode(ctx, mail, code)
}

func (s *authService) LoginVerify(ctx context.Context, params VerifyParams) (*Session, error) {
	email := normalizeEmail(params.Email)
	code := strings.TrimSpace(params.Code)
	if email == "" || code == "" {
		return nil, ErrEmptyOTP
	}

	if err := s.verifyAndConsumeOTP(ctx, email, code, "login"); err != nil {
		return nil, err
	}

	u, err := s.userSrv.GetByEmail(ctx, email)
	if err != nil || u == nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			u, err = s.userSrv.CreateUser(ctx, user.UserParams{
				Email:      email,
				Role:       "costumer",
				AuthMethod: "email_otp",
			})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return s.issueSession(ctx, u)
}

func (s *authService) verifyAndConsumeOTP(ctx context.Context, email, code, purpuse string) error {
	ch, err := s.otpSrv.GetActiveByEmail(ctx, email, purpuse)
	if err != nil || ch == nil {
		return ErrOTPNotFound
	}

	if ch.ConsumedAt != nil || time.Now().After(ch.ExpiresAt) {
		return ErrOTPExpired
	}

	if sha256Hex(code) != ch.CodeHash {
		_ = s.otpSrv.IncAttempts(ctx, email)
		return ErrOTPInvalidCode
	}

	_ = s.otpSrv.Consume(ctx, email)
	return nil
}

func (s *authService) issueSession(ctx context.Context, user *u.User) (*Session, error) {
	access, exp, err := s.jwt.NewAccessToken(user.ID.Hex(), user.Role, nil)
	if err != nil {
		return nil, err
	}

	refresh, err := NewRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshHash := sha256Hex(refresh)
	now := time.Now()

	session := session.RefreshSession{
		UserID:    user.ID,
		Hash:      refreshHash,
		ExpiresAt: now.Add(s.refreshTTL),
		CreatedAt: now,
	}

	if err := validateRefreshSession(session); err != nil {
		return nil, err
	}

	if err := s.authRepo.CreateRefreshToken(ctx, session); err != nil {
		return nil, err
	}

	return &Session{
		User:         *user,
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresInSec: int64(time.Until(exp).Seconds()),
	}, nil
}

func validateRefreshSession(refresh session.RefreshSession) error {
	if strings.TrimSpace(refresh.Hash) == "" {
		return ErrHashRequired
	}

	if refresh.UserID.IsZero() {
		return ErrUserRequired
	}

	if refresh.CreatedAt.IsZero() {
		return ErrCreatedAtRequired
	}

	if refresh.ExpiresAt.IsZero() {
		return ErrExpiredAtRequired
	}

	return nil
}

func (s *authService) challengePass(user u.User, pass string) error {
	if user.PasswordHash == nil || *user.PasswordHash == "" {
		return ErrMethodNotAllowed
	}

	if err := s.hasher.Compare(*user.PasswordHash, pass); err != nil {
		return ErrInvalidCreadentials
	}

	return nil
}

///

func NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func new6DigitCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func sha256Hex(s string) string {
	s = strings.TrimSpace(s)
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

//

func validateParam(params UserParams) error {
	email := strings.TrimSpace(strings.ToLower(params.Email))
	if err := validateEmail(email); err != nil {
		return err
	}

	if err := validateRole(params.Role); err != nil {
		return err
	}

	if err := validateAuthMethod(params.AuthMethod); err != nil {
		return ErrAuthMethod
	}

	if params.AuthMethod == "password" {
		if strings.TrimSpace(params.Password) == "" {
			return ErrPasswordRequired
		}
	}

	return nil
}

func validateEmail(email string) error {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return ErrEmailInvalid
	}

	if addr.Address != email {
		return ErrEmailInvalid
	}

	return nil
}

func passwordMethodOk(methods []string) bool {
	for _, val := range methods {
		if val == "password" {
			return true
		}
	}
	return false
}

type Role string

const (
	RoleCostumer Role = "costumer"
	RoleEditor   Role = "editor"
	RoleAdmin    Role = "admin"
)

func validateRole(role string) error {
	switch Role(role) {
	case RoleCostumer, RoleEditor, RoleAdmin:
		return nil
	default:
		return ErrRole
	}
}

type AuthMethod string

const (
	AuthPass  AuthMethod = "password"
	AuthEmail AuthMethod = "email_otp"
)

func validateAuthMethod(method string) error {
	switch AuthMethod(method) {
	case AuthPass, AuthEmail:
		return nil
	default:
		return ErrAuthMethod
	}
}

func normalizeEmail(email string) string {
	s := strings.ToLower(strings.TrimSpace(email))
	if s == "" || !strings.Contains(email, "@") {
		return ""
	}
	return s
}
