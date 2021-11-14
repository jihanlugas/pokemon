package controller

import (
	"encoding/binary"
	"errors"
	"gopokemon/constant"
	"gopokemon/cryption"
	"gopokemon/log"
	"gopokemon/response"
	"gopokemon/validator"
	"sync/atomic"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

var json jsoniter.API
var Validate *validator.CustomValidator
var ErrNilLastLogin error
var CacheUserAuth map[int64]atomic.Value

type UserAuthToken struct {
	PassVersion int
	LastLogin   int64
}

func init() {
	json = jsoniter.ConfigFastest
	Validate = validator.NewValidator()
	ErrNilLastLogin = errors.New("Last login can not be nil")
	CacheUserAuth = make(map[int64]atomic.Value)
}

type UserLogin struct {
	UserID      int64
}

func getLoginToken(userID int64, expiredAt time.Time) ([]byte, error) {
	expiredUnix := expiredAt.Unix()

	tokenPayload := make([]byte, constant.TokenPayloadLen)
	binary.BigEndian.PutUint64(tokenPayload, uint64(expiredUnix)) // Expired date
	binary.BigEndian.PutUint64(tokenPayload[8:], uint64(userID))

	return cryption.EncryptAES64(tokenPayload)
}

func getUserLoginInfo(c echo.Context) (UserLogin, error) {
	if u, ok := c.Get(constant.TokenUserContext).(UserLogin); ok {
		return u, nil
	} else {
		return UserLogin{}, response.ErrorForce("Akses tidak diterima", response.Payload{})
	}
}

func Ping(c echo.Context) error {
	return response.Success("Hallo　世界", response.Payload{}).SendJSON(c)
}

func errorInternal(c echo.Context, err error) {
	log.System.Error().Err(err).Str("Host", c.Request().Host).Str("Path", c.Path()).Send()
	panic(err)
}

//func upload(ctx context.Context, tx pgx.Tx, file *multipart.FileHeader, docType constant.DocType, userID int64, c echo.Context) (error, int64) {
//	// asumsi penggunaan fungsi ini adalah check file (berapa mb? extension bener/tidak? dll) sudah dilakukan sebelumnya.
//	var err error
//	var photo model.PublicPhoto
//
//	var newPhotoID int64
//	err = tx.QueryRow(ctx, `select public.next_id() as tid;`).Scan(&newPhotoID)
//	if err != nil {
//		return err, 0
//	}
//
//	photo.PhotoID = newPhotoID
//	photo.ClientName = file.Filename
//	photo.Ext = filepath.Ext(file.Filename)
//	photo.ServerName = strconv.FormatInt(int64(newPhotoID), 10) + photo.Ext
//	photo.PhotoPath = constant.FolderForFile[docType]
//	photo.PhotoSize = int64(file.Size)
//	photo.PhotoWidth = 0
//	photo.PhotoHeigth = 0
//	photo.CreateBy = userID
//	err = photo.Insert(ctx, tx)
//	if err != nil {
//		return err, 0
//	}
//
//	src, err := file.Open()
//	if err != nil {
//		return err, 0
//	}
//	defer src.Close()
//
//	// Destination
//	dst, err := os.Create(config.UploadFilePath + "/" + constant.FolderForFile[docType] + "/" + photo.ServerName)
//	if err != nil {
//		return err, 0
//	}
//	defer dst.Close()
//
//	// Copy
//	if _, err = io.Copy(dst, src); err != nil {
//		return err, 0
//	}
//	return nil, newPhotoID
//}

//func deletePhoto(ctx context.Context, tx pgx.Tx, photoID int64) error {
//	var qPhoto model.PublicPhoto
//	var err error
//	qPhoto.PhotoID = photoID
//
//	err = pgxscan.Get(ctx, tx, &qPhoto, "SELECT photo_id, client_name, server_name, ext, photo_path, photo_size, photo_width, photo_heigth FROM public.photo WHERE photo_id = $1", qPhoto.PhotoID)
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			return nil
//		}
//		return err
//	}
//
//	// remove photo from hdd
//	if err = os.Remove(config.UploadFilePath + "/" + qPhoto.PhotoPath + "/" + qPhoto.ServerName); err != nil {
//		return err
//	}
//
//	_, err = tx.Exec(ctx, "DELETE FROM public.photo WHERE photo_id = $1", qPhoto.PhotoID)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
