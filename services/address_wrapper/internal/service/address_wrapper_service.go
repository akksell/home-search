package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"homesearch.axel.to/address_wrapper/config"
	"homesearch.axel.to/address_wrapper/internal/util"

	base "homesearch.axel.to/base/types"
	pb "homesearch.axel.to/services/address_wrapper/api"
)

type googleGeocodeAPIResponse struct {
	Results []struct {
		PlaceId string `json:"placeId"`
	}
}

type AddressWrapperService struct {
	Config *config.AppConfig
	pb.UnimplementedAddressWrapperServiceServer
}

func (aws *AddressWrapperService) GetPlaceId(ctx context.Context, req *pb.PlaceIdRequest) (*pb.PlaceIdResponse, error) {
	// TODO: validate address here

	structuredAddressParams, err := getStructuredAddress(req.GetAddress())
	if err != nil {
		// TODO: error responses
		return &pb.PlaceIdResponse{}, err
	}

	apiKeyQueryParam := "key=" + aws.Config.GoogleGeocoderAPIKey
	geocodeAPI := aws.Config.GoogleGeocoderAPIHost + "/geocode/address"
	requestURL := geocodeAPI + "?" + structuredAddressParams + handleAmpersand(structuredAddressParams, apiKeyQueryParam)
	geocodeRequest, err := http.NewRequest(http.MethodGet, requestURL, nil)
	// TODO: use error + log after creating request

	geocodeRequest.Header.Set("X-Goog-FieldMask", "results.placeId")
	httpClient := &http.Client{}
	response, err := httpClient.Do(geocodeRequest)
	// TODO: log and return clean/obfuscated error
	if err != nil {
		return &pb.PlaceIdResponse{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	var bodyAsJson googleGeocodeAPIResponse
	err = json.Unmarshal(body, &bodyAsJson)
	// TODO: log error and return clean/obfuscated error
	if err != nil {
		return &pb.PlaceIdResponse{}, err
	}

	if len(bodyAsJson.Results) == 0 {
		return &pb.PlaceIdResponse{}, fmt.Errorf("No results for address")
	}

	placeId := bodyAsJson.Results[0].PlaceId
	// TODO: store place_id + address in firestore db so it can queried from there instead of
	//		 potentially hitting quota from requests for common/reused addresses
	return &pb.PlaceIdResponse{
		PlaceId: placeId,
	}, nil
}

// based off of google geocode v4 api
func getStructuredAddress(address *base.Address) (string, error) {
	components := ""
	if address.GetStreet() != "" {
		component := "address.addressLines=" + util.URLEncode(address.GetStreet())
		components += handleAmpersand(components, component)
	}
	if address.GetCity() != "" {
		component := "address.locality=" + util.URLEncode(address.GetCity())
		components += handleAmpersand(components, component)
	}
	if address.GetPostalCode() != "" {
		component := "address.postalCode=" + util.URLEncode(address.GetPostalCode())
		components += handleAmpersand(components, component)
	}
	if address.GetStateProvinceCode() != "" {
		component := "address.administrativeArea=" + util.URLEncode(address.GetStateProvinceCode())
		components += handleAmpersand(components, component)
	}
	if address.GetCountryCode() != "" {
		component := "address.regionCode=" + util.URLEncode(address.GetCountryCode())
		components += handleAmpersand(components, component)
	}
	return components, nil
}

func handleAmpersand(existingComponents, component string) string {
	if len(existingComponents) != 0 {
		return "&" + component
	}
	return component
}
