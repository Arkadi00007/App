package usecase
//
//import (
//	"Test_App/internal/domain"
//	"Test_App/pkg/code"
//	"Test_App/pkg/email"
//	"Test_App/pkg/jwt"
//	"context"
//	"errors"
//
//	"golang.org/x/crypto/bcrypt"
//	"time"
//)
//
//type userUseCase struct {
//	userRepo   domain.UserRepository
//	verifyRepo domain.VerificationCodeRepository
//	mailer     *email.Sender
//	jwtSecret string
//}
//
//func NewUserUseCase(
//	userRepo domain.UserRepository,
//	codeRepo domain.VerificationCodeRepository,
//	mailer *email.Sender,
//	jwtSecret string,
//) domain.UserUseCase {
//	return &userUseCase{
//		userRepo:   userRepo,
//		verifyRepo: codeRepo,
//		mailer:     mailer,
//		jwtSecret:  jwtSecret,
//	}
//}
//
//
////func (uc *userUseCase) sendCode(ctx context.Context, user *domain.User, codeType string) error {
////	// 1. проверка cooldown
////	latest, err := uc.codeRepo.GetLatestCode(ctx, user.ID, codeType)
////	if err != nil {
////		return err
////	}
////
////	if latest != nil {
////		elapsed := time.Since(latest.CreatedAt)
////		if elapsed < domain.ResendCooldown {
////			remaining := domain.ResendCooldown - elapsed
////			return fmt.Errorf("подождите %d секунд перед повторной отправкой", int(remaining.Seconds()))
////		}
////	}
////
////	// 2. создаём новый код
////	newCode := code.Generate()
////	vc := &domain.VerificationCode{
////		UserID:    user.ID,
////		Code:      newCode,
////		Type:      codeType,
////		ExpiresAt: time.Now().Add(15 * time.Minute),
////	}
////	if err := uc.codeRepo.Create(ctx, vc); err != nil {
////		return err
////	}
////
////	// 3. отправляем письмо
////	if codeType == domain.CodeTypeEmailVerify {
////		return uc.mailer.SendVerificationCode(user.Email, newCode)
////	}
////	return uc.mailer.SendResetPasswordCode(user.Email, newCode)
////}
//
//func (uc *userUseCase) Register(ctx context.Context, emailAddr, password, name string) error {
//	// 1. проверяем что email не занят
//	existing, err := uc.userRepo.GetUserByEmail(ctx, emailAddr)
//	if err != nil {
//		return err
//	}
//	if existing != nil {
//		return errors.New("email уже занят")
//	}
//
//	// 2. хэшируем пароль
//	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//	if err != nil {
//		return err
//	}
//
//	// 3. создаём юзера
//	user := &domain.User{
//		Email:        emailAddr,
//		PasswordHash: string(hash),
//		Name:         name,
//	}
//	if err := uc.userRepo.CreateUser(ctx, user); err != nil {
//		return err
//	}
//
//	// 4. создаём код верификации
//	verificationCode := code.Generate()
//	vc := &domain.VerificationCode{
//		UserID:    user.ID,
//		Code:      verificationCode,
//		Type:      domain.CodeTypeEmailVerify,
//		ExpiresAt: time.Now().Add(15 * time.Minute),
//	}
//	if err := uc.verifyRepo.CreateCode(ctx, vc); err != nil {
//		return err
//	}
//
//	// 5. отправляем код на email
//	return uc.mailer.SendVerificationCode(emailAddr, verificationCode)
//}
//
//func (uc *userUseCase) VerifyEmail(ctx context.Context, emailAddr, inputCode string) (*domain.User, string, error) {
//	// 1. находим юзера
//	user, err := uc.userRepo.GetUserByEmail(ctx,emailAddr)
//	if err != nil {
//		return nil, "", err
//	}
//	if user == nil {
//		return nil, "", errors.New("пользователь не найден")
//	}
//	if user.IsVerified {
//		return nil, "", errors.New("email уже подтверждён")
//	}
//
//	// 2. находим активный код
//	vc, err := uc.verifyRepo.GetActiveCode(ctx, user.ID, domain.CodeTypeEmailVerify)
//	if err != nil {
//		return nil, "", err
//	}
//	if vc == nil {
//		return nil, "", errors.New("код не найден или истёк")
//	}
//
//	// 3. проверяем код
//	if vc.Code != inputCode {
//		return nil, "", errors.New("неверный код")
//	}
//
//	// 4. помечаем код как использованный
//	if err := uc.verifyRepo.MarkAsUsedCode(ctx, vc.ID); err != nil {
//		return nil, "", err
//	}
//
//	// 5. помечаем юзера как верифицированного
//	user.IsVerified = true
//	if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
//		return nil, "", err
//	}
//
//	// 6. генерируем токен
//	token, err := jwt.Generate(user.ID, user.Role, uc.jwtSecret)
//	if err != nil {
//		return nil, "", err
//	}
//
//	return user, token, nil
//}
//
//func (uc *userUseCase) Login(ctx context.Context, emailAddr, password string) (*domain.User, string, error) {
//	// 1. находим юзера
//	user, err := uc.userRepo.GetUserByEmail(ctx, emailAddr)
//	if err != nil {
//		return nil, "", err
//	}
//	if user == nil {
//		return nil, "", errors.New("неверный email или пароль")
//	}
//
//	// 2. проверяем верификацию
//	if !user.IsVerified {
//		return nil, "", errors.New("email не подтверждён")
//	}
//
//	// 3. проверяем пароль
//	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
//		return nil, "", errors.New("неверный email или пароль")
//	}
//
//	// 4. генерируем токен
//	token, err := jwt.Generate(user.ID, user.Role, uc.jwtSecret)
//	if err != nil {
//		return nil, "", err
//	}
//
//	return user, token, nil
//}
//
//
//
//``
//func (uc *userUseCase) ForgotPassword(ctx context.Context, mailAddr string) error {
//	// 1. находим юзера
//	user, err := uc.userRepo.GetUserByEmail(ctx, mailAddr)
//	if err != nil {
//		return err
//	}
//	// намеренно не говорим что юзер не найден — безопасность
//	if user == nil {
//		return nil
//	}
//
//	// 2. генерируем код
//	resetCode := code.Generate()
//	vc := &domain.VerificationCode{
//		UserID:    user.ID,
//		Code:      resetCode,
//		Type:      domain.CodeTypeResetPassword,
//		ExpiresAt: time.Now().Add(15 * time.Minute),
//	}
//	if err := uc.verifyRepo.CreateCode(ctx, vc); err != nil {
//		return err
//	}
//
//	// 3. отправляем код
//	return uc.mailer.SendResetPasswordCode(mailAddr, resetCode)
//}
//
//func (uc *userUseCase) ResetPassword(ctx context.Context, emailAddr, inputCode, newPassword string) error {
//	// 1. находим юзера
//	user, err := uc.userRepo.GetUserByEmail(ctx, emailAddr)
//	if err != nil {
//		return err
//	}
//	if user == nil {
//		return errors.New("пользователь не найден")
//	}
//
//	// 2. находим активный код
//	vc, err := uc.verifyRepo.GetActiveCode(ctx, user.ID, domain.CodeTypeResetPassword)
//	if err != nil {
//		return err
//	}
//	if vc == nil {
//		return errors.New("код не найден или истёк")
//	}
//
//	// 3. проверяем код
//	if vc.Code != inputCode {
//		return errors.New("неверный код")
//	}
//
//	// 4. помечаем код как использованный
//	if err := uc.verifyRepo.MarkAsUsedCode(ctx, vc.ID); err != nil {
//		return err
//	}
//
//	// 5. хэшируем новый пароль и обновляем
//	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
//	if err != nil {
//		return err
//	}
//	user.PasswordHash = string(hash)
//	return uc.userRepo.UpdateUser(ctx, user)
//}
