package dtrack

import "context"

type BOMUploadRequest struct {
	ProjectUUID    string `json:"project,omitempty"`
	ProjectName    string `json:"projectName,omitempty"`
	ProjectVersion string `json:"projectVersion,omitempty"`
	AutoCreate     bool   `json:"autoCreate"`
	BOM            string `json:"bom"`
}

type bomSubmitResponse struct {
	Token string `json:"token"`
}

type tokenProcessingResponse struct {
	Processing bool `json:"processing"`
}

type BOMService interface {
	ExportProjectAsCycloneDX(ctx context.Context, projectUUID string) (string, error)
	IsBeingProcessed(ctx context.Context, uploadToken string) (bool, error)
	Upload(ctx context.Context, req BOMUploadRequest) (string, error)
}

type bomServiceImpl struct {
	client *Client
}

func (b bomServiceImpl) ExportProjectAsCycloneDX(ctx context.Context, projectUUID string) (string, error) {
	res, err := b.client.restClient.R().
		SetContext(ctx).
		SetHeader("Accept", "application/xml").
		SetPathParam("uuid", projectUUID).
		Get("/api/v1/bom/cyclonedx/project/{uuid}")
	if err != nil {
		return "", err
	}

	if err = b.client.checkResponseStatus(res, 200); err != nil {
		return "", err
	}

	return res.String(), nil
}

func (b bomServiceImpl) IsBeingProcessed(ctx context.Context, uploadToken string) (bool, error) {
	res, err := b.client.restClient.R().
		SetContext(ctx).
		SetPathParam("token", uploadToken).
		SetResult(&tokenProcessingResponse{}).
		Get("/api/v1/bom/token/{token}")
	if err != nil {
		return false, err
	}

	if err = b.client.checkResponseStatus(res, 200); err != nil {
		return false, err
	}

	return res.Result().(*tokenProcessingResponse).Processing, nil
}

func (b bomServiceImpl) Upload(ctx context.Context, req BOMUploadRequest) (string, error) {
	res, err := b.client.restClient.R().
		SetContext(ctx).
		SetBody(req).
		SetHeader("Content-Type", "application/json").
		SetResult(&bomSubmitResponse{}).
		Put("/api/v1/bom")
	if err != nil {
		return "", err
	}

	if err = b.client.checkResponseStatus(res, 200); err != nil {
		return "", err
	}

	return res.Result().(*bomSubmitResponse).Token, nil
}
