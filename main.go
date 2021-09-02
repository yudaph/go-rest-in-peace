package gorip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/valyala/fasthttp"
)

type Config fasthttp.Client
var (
	a = fasthttp.Client{}
)

func New(config ...Config) {
	if len(config) >0 {
		a = fasthttp.Client(config[0])
	}
}

func Get(url string, target interface{}, headerOptional ...map[string]string) error {

	head := map[string]string{}
	if len(headerOptional) > 0 {
		head = headerOptional[0]
	}

	return request(&url, fasthttp.MethodGet, &[]byte{}, &target, &head)
}

func GetGoroutine(wg *sync.WaitGroup, url string, target interface{}, ErrChanel chan error, headerOptional ...map[string]string) {
	defer wg.Done()

	head := map[string]string{}
	if len(headerOptional) > 0 {
		head = headerOptional[0]
	}

	err:= request(&url, fasthttp.MethodGet, &[]byte{}, &target, &head)
	ErrChanel <- err
}

func Post(url string, body map[string]interface{}, target interface{}, headerOptional ...map[string]string) error {
	
	head := map[string]string{}
	if len(headerOptional) > 0 {
		head = headerOptional[0]
	}

	content, err := json.Marshal(body)
	if err!=nil{
		return err
	}
	
	return request(&url, fasthttp.MethodPost, &content, &target, &head)
}

func PostGoroutine(wg *sync.WaitGroup, url string, body map[string]interface{}, target interface{}, ErrChanel chan error, headerOptional ...map[string]string){
	defer wg.Done()

	head := map[string]string{}
	if len(headerOptional) > 0 {
		head = headerOptional[0]
	}

	content, err := json.Marshal(body)
	if err != nil{	
		ErrChanel <- err
	}

	err = request(&url, fasthttp.MethodPost, &content, &target, &head)
	ErrChanel <- err
}

func Put(url string, body map[string]interface{}, target interface{}, headerOptional ...map[string]string) error {
	
	head := map[string]string{}
	if len(headerOptional) > 0 {
		head = headerOptional[0]
	}

	content, err := json.Marshal(body)
	if err!=nil{
		return err
	}
	
	return request(&url, fasthttp.MethodPut, &content, &target, &head)
}

func PutGoroutine(wg *sync.WaitGroup, url string, body map[string]interface{}, target interface{}, ErrChanel chan error, headerOptional ...map[string]string){
	defer wg.Done()

	head := map[string]string{}
	if len(headerOptional) > 0 {
		head = headerOptional[0]
	}

	content, err := json.Marshal(body)
	if err != nil{	
		ErrChanel <- err
	}

	err = request(&url, fasthttp.MethodPut, &content, &target, &head)
	ErrChanel <- err
}

func Patch(url string, body map[string]interface{}, target interface{}, headerOptional ...map[string]string) error {
	
	head := map[string]string{}
	if len(headerOptional) > 0 {
		head = headerOptional[0]
	}

	content, err := json.Marshal(body)
	if err!=nil{
		return err
	}
	
	return request(&url, fasthttp.MethodPatch, &content, &target, &head)
}

func PatchGoroutine(wg *sync.WaitGroup, url string, body map[string]interface{}, target interface{}, ErrChanel chan error, headerOptional ...map[string]string){
	defer wg.Done()

	head := map[string]string{}
	if len(headerOptional) > 0 {
		head = headerOptional[0]
	}

	content, err := json.Marshal(body)
	if err != nil{	
		ErrChanel <- err
	}

	err = request(&url, fasthttp.MethodPatch, &content, &target, &head)
	ErrChanel <- err
}

func Delete(url string, target interface{}, headerOptional ...map[string]string) error {

	head := map[string]string{}
	if len(headerOptional) > 0 {
		head = headerOptional[0]
	}

	return request(&url, fasthttp.MethodDelete, &[]byte{}, &target, &head)
}

func DeleteGoroutine(wg *sync.WaitGroup, url string, target interface{}, ErrChanel chan error, headerOptional ...map[string]string) {
	defer wg.Done()

	head := map[string]string{}
	if len(headerOptional) > 0 {
		head = headerOptional[0]
	}

	err:= request(&url, fasthttp.MethodDelete, &[]byte{}, &target, &head)
	ErrChanel <- err
}

func request(url *string, method string, content *[]byte, target *interface{}, header *map[string]string) error {

	//fmt.Println("Isi dari a = ",a)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	if len(*header) > 0 {
		for i,v := range *header{
			fmt.Println("key = ", i)
			fmt.Println("value = ", v)
		}
	}

	req.AppendBody(*content)
	req.SetRequestURI(*url)
	req.URI().DisablePathNormalizing = true
	req.Header.SetMethod(method)

	// fasthttp does not automatically request a gzipped response.
	// We must explicitly ask for it.
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// Perform the request
	err := a.Do(req, resp)
	if err != nil {
		fmt.Printf("Client get failed: %s\n", err)
		return err
	}
	if resp.StatusCode() < 200 || resp.StatusCode() > 210 {
		//fmt.Printf("Expected status code %d but got %d\n", fasthttp.StatusOK, resp.StatusCode())
		contentEncoding := resp.Header.Peek("Content-Encoding")
		var errBody []byte
		if bytes.EqualFold(contentEncoding, []byte("gzip")) {
			fmt.Println("Unzipping...")
			errBody, _ = resp.BodyGunzip()
		} else {
			errBody = resp.Body()
		}
		//fmt.Println(string(errBody))
		return fmt.Errorf("status code %d, "+string(errBody), resp.StatusCode())
	}

	// Verify the content type
	contentType := resp.Header.Peek("Content-Type")
	if bytes.Index(contentType, []byte("application/json")) != 0 {
		fmt.Printf("Expected content type application/json but got %s\n", contentType)
		return err
	}

	// Do we need to decompress the response?
	contentEncoding := resp.Header.Peek("Content-Encoding")
	var body []byte
	if bytes.EqualFold(contentEncoding, []byte("gzip")) {
		fmt.Println("Unzipping...")
		body, _ = resp.BodyGunzip()
	} else {
		body = resp.Body()
	}
	
	err = json.Unmarshal(body, target)
	if err!=nil {
		return err
	}
	return nil
}