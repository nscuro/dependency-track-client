package dtrack

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

func (c Client) UploadBOM(request BOMUploadRequest) (string, error) {
	res, err := c.restClient.R().
		SetBody(request).
		SetHeader("Content-Type", "application/json").
		SetResult(&bomSubmitResponse{}).
		Put("/api/v1/bom")
	if err != nil {
		return "", err
	}

	if err = c.checkResponseStatus(res, 200); err != nil {
		return "", err
	}

	return res.Result().(*bomSubmitResponse).Token, nil
}

func (c Client) IsTokenBeingProcessed(uploadToken string) (bool, error) {
	res, err := c.restClient.R().
		SetPathParams(map[string]string{
			"token": uploadToken,
		}).
		SetResult(&tokenProcessingResponse{}).
		Get("/api/v1/bom/token/{token}")
	if err != nil {
		return false, err
	}

	if err = c.checkResponseStatus(res, 200); err != nil {
		return false, err
	}

	return res.Result().(*tokenProcessingResponse).Processing, nil
}

func (c Client) ExportProjectAsCycloneDX(uuid string) (string, error) {
	res, err := c.restClient.R().
		SetHeader("Accept", "application/xml").
		SetPathParams(map[string]string{
			"uuid": uuid,
		}).
		Get("/api/v1/bom/cyclonedx/project/{uuid}")
	if err != nil {
		return "", err
	}

	if err = c.checkResponseStatus(res, 200); err != nil {
		return "", err
	}

	return res.String(), nil
}
