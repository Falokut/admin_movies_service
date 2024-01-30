package service

import (
	"context"
	"sync"

	"github.com/Falokut/admin_movies_service/internal/events"
	image_processing_service "github.com/Falokut/image_processing_service/pkg/image_processing_service/v1/protos"
	images_storage_service "github.com/Falokut/images_storage_service/pkg/images_storage_service/v1/protos"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/sirupsen/logrus"
)

type imagesService struct {
	logger               *logrus.Logger
	errorHandler         errorHandler
	processingClient     image_processing_service.ImageProcessingServiceV1Client
	storage              images_storage_service.ImagesStorageServiceV1Client
	storageConn          *grpc.ClientConn
	processingClientConn *grpc.ClientConn
	imagesMQ             events.ImagesEventsMQ
	basePhotoUrl         string
}

func (s *imagesService) Shutdown() {
	s.storageConn.Close()
	s.processingClientConn.Close()
}

func NewImageService(logger *logrus.Logger,
	basePhotoUrl, storageAddr, imageProcessingAddr string,
	imagesMQ events.ImagesEventsMQ) (*imagesService, error) {
	errorHandler := newErrorHandler(logger)
	processingClientConn, err := getGrpcConnection(imageProcessingAddr)
	if err != nil {
		return nil, err
	}
	processingClient := image_processing_service.NewImageProcessingServiceV1Client(processingClientConn)

	storageConn, err := getGrpcConnection(storageAddr)
	if err != nil {
		return nil, err
	}
	storage := images_storage_service.NewImagesStorageServiceV1Client(storageConn)

	return &imagesService{
		logger:               logger,
		errorHandler:         errorHandler,
		basePhotoUrl:         basePhotoUrl,
		processingClient:     processingClient,
		processingClientConn: processingClientConn,
		storage:              storage,
		imagesMQ:             imagesMQ,
		storageConn:          storageConn,
	}, nil
}

type CheckImageConfig struct {
	MaxImageWidth  int32
	MaxImageHeight int32
	MinImageWidth  int32
	MinImageHeight int32
	AllowedTypes   []string
}

// Returns ref to the int var if var>0, othervise returns nil.
func getIntRef(i int32) *int32 {
	if i <= 0 {
		return nil
	}
	return &i
}

func (c CheckImageConfig) getMaxImageHeight() *int32 {
	return getIntRef(c.MaxImageHeight)
}
func (c CheckImageConfig) getMaxImageWidth() *int32 {
	return getIntRef(c.MaxImageHeight)
}

func (c CheckImageConfig) getMinImageHeight() *int32 {
	return getIntRef(c.MinImageHeight)
}
func (c CheckImageConfig) getMinImageWidth() *int32 {
	return getIntRef(c.MinImageWidth)
}

type UploadPictureParam struct {
	Image                []byte
	ValidateImage        bool
	CheckImage           CheckImageConfig
	ResizeImage          bool
	ImageProcessingParam struct {
		ImageHeight       int32
		ImageWidth        int32
		ImageResizeMethod image_processing_service.ResampleFilter
	}
	Category string
}

type ReplacePicturesParam struct {
	Image                []byte
	ValidateImage        bool
	CheckImage           CheckImageConfig
	ResizeImage          bool
	ImageProcessingParam struct {
		ImageHeight       int32
		ImageWidth        int32
		ImageResizeMethod image_processing_service.ResampleFilter
	}
	Category  string
	ImageName string // if empty, will be created new image
}

func (s *imagesService) UploadPicture(ctx context.Context, picture UploadPictureParam) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "imagesService.UploadPicture")
	defer span.Finish()

	if picture.ValidateImage {
		if err := s.checkImage(ctx, picture.Image, picture.CheckImage); err != nil {
			span.SetTag("grpc.status", status.Code(err))
			ext.LogError(span, err)
			return "", err
		}
	}

	var img = picture.Image
	var err error
	if picture.ResizeImage {
		img, err = s.ResizeImage(ctx, img, picture.ImageProcessingParam.ImageHeight,
			picture.ImageProcessingParam.ImageWidth, picture.ImageProcessingParam.ImageResizeMethod)
		if err != nil {
			span.SetTag("grpc.status", status.Code(err))
			ext.LogError(span, err)
			return "", err
		}
	}

	s.logger.Info("Creating stream")

	res, err := s.storage.UploadImage(ctx, &images_storage_service.UploadImageRequest{
		Category: picture.Category,
		Image:    img,
	})
	if err != nil {
		return "", status.Errorf(status.Code(err), "can't upload image %s", err.Error())
	}

	return res.ImageId, nil
}

func (s *imagesService) ResizeImage(ctx context.Context, image []byte, height, width int32,
	ResizeMethod image_processing_service.ResampleFilter) ([]byte, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "imagesService.ResizeImage")
	defer span.Finish()

	resized, err := s.processingClient.Resize(ctx, &image_processing_service.ResizeRequest{
		Image:          &image_processing_service.Image{Image: image},
		ResampleFilter: ResizeMethod,
		Width:          width,
		Height:         height,
	})

	if err != nil {
		span.SetTag("grpc.status", status.Code(err))
		ext.LogError(span, err)
		return []byte{}, err
	}
	if resized == nil {
		return []byte{}, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, "can't resize image")
	}

	span.SetTag("grpc.status", codes.OK)
	return resized.Data, nil
}

func (s *imagesService) checkImage(ctx context.Context, image []byte, cfg CheckImageConfig) error {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"imagesService.checkImage")
	defer span.Finish()

	res, err := s.processingClient.Validate(ctx,
		&image_processing_service.ValidateRequest{
			Image:          &image_processing_service.Image{Image: image},
			SupportedTypes: cfg.AllowedTypes,
			MaxWidth:       cfg.getMaxImageWidth(),
			MaxHeight:      cfg.getMaxImageHeight(),
			MinHeight:      cfg.getMinImageHeight(),
			MinWidth:       cfg.getMinImageWidth(),
		})
	if status.Code(err) == codes.InvalidArgument {
		var msg string
		if res != nil {
			msg = res.GetDetails()
		}
		return s.errorHandler.createExtendedErrorResponceWithSpan(span, ErrInvalidImage, "", msg)
	} else if err != nil {
		var msg string
		if res != nil {
			msg = res.GetDetails()
		}
		return s.errorHandler.createExtendedErrorResponceWithSpan(span, err, "", msg)

	}

	span.SetTag("grpc.status", codes.OK)
	return nil
}

func (s *imagesService) DeletePicture(ctx context.Context, category, id string) error {
	_, err := s.storage.DeleteImage(ctx,
		&images_storage_service.ImageRequest{ImageId: id, Category: category})
	if status.Code(err) != codes.NotFound {
		return nil
	}
	if err != nil {
		return s.imagesMQ.DeleteImageRequest(ctx, id, category)
	}
	return nil
}

// Returns picture url for GET request
func (s *imagesService) GetPictureURL(pictureID, category string) string {
	if pictureID == "" {
		return ""
	}

	return s.basePhotoUrl + "/" + category + "/" + pictureID
}

func (s *imagesService) UploadPictures(ctx context.Context,
	pictures map[string]UploadPictureParam) (map[string]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "imagesService.UploadPictures")
	defer span.Finish()
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var resCh = make(chan struct {
		id  string
		key string
		err error
	})

	f := func(key string, picture UploadPictureParam) {
		defer wg.Done()
		var st = struct {
			id  string
			key string
			err error
		}{key: key}

		if len(picture.Image) == 0 {
			resCh <- st
			return
		}
		select {
		case <-ctx.Done():
			resCh <- st
			return
		default:
			id, err := s.UploadPicture(ctx, picture)
			st.id = id
			st.err = err
			resCh <- st
		}
	}

	wg.Add(len(pictures))
	for key, value := range pictures {
		go f(key, value)
	}

	go func() {
		wg.Wait()
		s.logger.Debug("channel closed")
		close(resCh)
	}()

	var images = make(map[string]string, 3)
	for res := range resCh {
		if res.err != nil {
			return map[string]string{}, res.err
		}
		if res.id == "" {
			continue
		}
		images[res.key] = res.id
		s.logger.Debug("Uploading pictures loop")
	}

	s.logger.Info("Uploading pictures success")
	return images, nil
}

func getGrpcConnection(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(
			otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
}

func (s *imagesService) ReplacePicture(ctx context.Context, picture ReplacePicturesParam, createIfNotExist bool) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "imagesService.ReplacePicture")
	defer span.Finish()

	if picture.ValidateImage {
		if err := s.checkImage(ctx, picture.Image, picture.CheckImage); err != nil {
			span.SetTag("grpc.status", status.Code(err))
			ext.LogError(span, err)
			return "", err
		}
	}

	var img = picture.Image
	var err error
	if picture.ResizeImage {
		img, err = s.ResizeImage(ctx, img, picture.ImageProcessingParam.ImageHeight,
			picture.ImageProcessingParam.ImageWidth, picture.ImageProcessingParam.ImageResizeMethod)
		if err != nil {
			span.SetTag("grpc.status", status.Code(err))
			ext.LogError(span, err)
			return "", err
		}
	}

	s.logger.Info("Creating stream")
	res, err := s.storage.ReplaceImage(ctx, &images_storage_service.ReplaceImageRequest{
		Category:         picture.Category,
		CreateIfNotExist: createIfNotExist,
		ImageData:        img,
		ImageId:          picture.ImageName,
	})
	if err != nil {
		return "", s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error()+" error while sending close")
	}
	if res.ImageId == "" {
		res.ImageId = picture.ImageName
	}

	return res.ImageId, nil
}

func (s *imagesService) ReplacePictures(ctx context.Context,
	pictures map[string]ReplacePicturesParam) (map[string]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "imagesService.UploadPictures")
	defer span.Finish()
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var resCh = make(chan struct {
		id  string
		key string
		err error
	})

	f := func(key string, picture ReplacePicturesParam) {
		defer wg.Done()
		var st = struct {
			id  string
			key string
			err error
		}{key: key}

		if len(picture.Image) == 0 {
			resCh <- st
			return
		}
		select {
		case <-ctx.Done():
			resCh <- st
			return
		default:
			id, err := s.ReplacePicture(ctx, picture, true)
			st.id = id
			st.err = err
			resCh <- st
		}
	}

	wg.Add(len(pictures))
	for key, value := range pictures {
		go f(key, value)
	}

	go func() {
		wg.Wait()
		s.logger.Debug("channel closed")
		close(resCh)
	}()

	var images = make(map[string]string, 3)
	for res := range resCh {
		if res.err != nil {
			return images, res.err
		}
		if res.id == "" {
			continue
		}
		images[res.key] = res.id
		s.logger.Debug("Uploading pictures loop")
	}

	s.logger.Info("Uploading pictures success")
	return images, nil
}
