// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/devtron-labs/image-scanner/api"
	"github.com/devtron-labs/image-scanner/client"
	"github.com/devtron-labs/image-scanner/internal/logger"
	"github.com/devtron-labs/image-scanner/internal/sql"
	"github.com/devtron-labs/image-scanner/internal/sql/repository"
	"github.com/devtron-labs/image-scanner/pkg/grafeasService"
	"github.com/devtron-labs/image-scanner/pkg/klarService"
	"github.com/devtron-labs/image-scanner/pkg/security"
	"github.com/devtron-labs/image-scanner/pkg/user"
	"github.com/devtron-labs/image-scanner/pubsub"
)

// Injectors from Wire.go:

func InitializeApp() (*App, error) {
	sugaredLogger := logger.NewSugardLogger()
	klarConfig, err := klarService.GetKlarConfig()
	if err != nil {
		return nil, err
	}
	apiClient := grafeasService.GetGrafeasClient()
	httpClient := logger.NewHttpClient()
	grafeasServiceImpl := grafeasService.NewKlarServiceImpl(sugaredLogger, apiClient, httpClient)
	config, err := sql.GetConfig()
	if err != nil {
		return nil, err
	}
	db, err := sql.NewDbConnection(config, sugaredLogger)
	if err != nil {
		return nil, err
	}
	userRepositoryImpl := repository.NewUserRepositoryImpl(db)
	imageScanHistoryRepositoryImpl := repository.NewImageScanHistoryRepositoryImpl(db, sugaredLogger)
	imageScanResultRepositoryImpl := repository.NewImageScanResultRepositoryImpl(db, sugaredLogger)
	imageScanObjectMetaRepositoryImpl := repository.NewImageScanObjectMetaRepositoryImpl(db, sugaredLogger)
	cveStoreRepositoryImpl := repository.NewCveStoreRepositoryImpl(db, sugaredLogger)
	imageScanDeployInfoRepositoryImpl := repository.NewImageScanDeployInfoRepositoryImpl(db, sugaredLogger)
	ciArtifactRepositoryImpl := repository.NewCiArtifactRepositoryImpl(db, sugaredLogger)
	imageScanServiceImpl := security.NewImageScanServiceImpl(sugaredLogger, imageScanHistoryRepositoryImpl, imageScanResultRepositoryImpl, imageScanObjectMetaRepositoryImpl, cveStoreRepositoryImpl, imageScanDeployInfoRepositoryImpl, ciArtifactRepositoryImpl)
	klarServiceImpl := klarService.NewKlarServiceImpl(sugaredLogger, klarConfig, grafeasServiceImpl, userRepositoryImpl, imageScanServiceImpl)
	pubSubClient, err := client.NewPubSubClient(sugaredLogger)
	if err != nil {
		return nil, err
	}
	testPublishImpl := pubsub.NewTestPublishImpl(pubSubClient, sugaredLogger, klarServiceImpl)
	userServiceImpl := user.NewUserServiceImpl(sugaredLogger, userRepositoryImpl)
	restHandlerImpl := api.NewRestHandlerImpl(sugaredLogger, klarServiceImpl, testPublishImpl, grafeasServiceImpl, userServiceImpl, imageScanServiceImpl)
	muxRouter := api.NewMuxRouter(sugaredLogger, restHandlerImpl)
	natSubscriptionImpl, err := pubsub.NewNatSubscription(pubSubClient, sugaredLogger, klarServiceImpl)
	if err != nil {
		return nil, err
	}
	app := NewApp(muxRouter, sugaredLogger, db, natSubscriptionImpl)
	return app, nil
}