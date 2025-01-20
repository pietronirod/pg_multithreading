package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	BrasilAPIURL string
	ViaCEPURL    string
	Timeout      time.Duration
}

// Estrutura para a resposta do BrasilAPI
type BrasilAPIResponse struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

// Estrutura para a resposta do ViaCEP
type ViaCEPResponse struct {
	CEP        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
}

// Estrutura comum para uso no código
type Address struct {
	CEP        string
	Logradouro string
	Bairro     string
	Cidade     string
	UF         string
}

// Estrutura para a resposta da API junto com a fonte
type APIResponse struct {
	Result Address
	Source string
}

// Estrutura para erros detalhados
type DetailedError struct {
	API      string
	Message  string
	Duration time.Duration
}

func (e *DetailedError) Error() string {
	return fmt.Sprintf("Erro na API %s: %s (durou %v)", e.API, e.Message, e.Duration)
}

var (
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:      10,
			IdleConnTimeout:   30 * time.Second,
			DisableKeepAlives: false,
		},
	}
)

func loadConfig() Config {
	brasilAPIURL := os.Getenv("BRASIL_API_URL")
	if brasilAPIURL == "" {
		brasilAPIURL = "https://brasilapi.com.br/api/cep/v1/"
	}

	viaCEPURL := os.Getenv("VIACEP_URL")
	if viaCEPURL == "" {
		viaCEPURL = "https://viacep.com.br/ws/"
	}

	timeoutStr := os.Getenv("API_TIMEOUT")
	timeout := 1 * time.Second
	if timeoutStr != "" {
		if t, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = t
		}
	}

	return Config{
		BrasilAPIURL: brasilAPIURL,
		ViaCEPURL:    viaCEPURL,
		Timeout:      timeout,
	}
}

func fetchAPI(ctx context.Context, url, source string) (Address, error) {
	start := time.Now()
	log.Printf("Iniciando requisição para %s (%s)", source, url)

	var address Address
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Address{}, &DetailedError{
			API:      source,
			Message:  err.Error(),
			Duration: time.Since(start),
		}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return Address{}, &DetailedError{
			API:      source,
			Message:  err.Error(),
			Duration: time.Since(start),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Address{}, &DetailedError{
			API:      source,
			Message:  err.Error(),
			Duration: time.Since(start),
		}
	}

	switch source {
	case "BrasilAPI":
		var brasilResponse BrasilAPIResponse
		err = json.Unmarshal(body, &brasilResponse)
		if err != nil {
			return Address{}, &DetailedError{
				API:      source,
				Message:  err.Error(),
				Duration: time.Since(start),
			}
		}
		address = Address{
			CEP:        brasilResponse.CEP,
			Logradouro: brasilResponse.Street,
			Bairro:     brasilResponse.Neighborhood,
			Cidade:     brasilResponse.City,
			UF:         brasilResponse.State,
		}
	case "ViaCEP":
		var viaCEPResponse ViaCEPResponse
		err = json.Unmarshal(body, &viaCEPResponse)
		if err != nil {
			return Address{}, &DetailedError{
				API:      source,
				Message:  err.Error(),
				Duration: time.Since(start),
			}
		}
		address = Address{
			CEP:        viaCEPResponse.CEP,
			Logradouro: viaCEPResponse.Logradouro,
			Bairro:     viaCEPResponse.Bairro,
			Cidade:     viaCEPResponse.Localidade,
			UF:         viaCEPResponse.UF,
		}
	default:
		return Address{}, &DetailedError{
			API:      source,
			Message:  "API desconhecida",
			Duration: time.Since(start),
		}
	}

	log.Printf("Requisição para %s completada em %v", source, time.Since(start))
	return address, nil
}

func fetchAPIWithRetry(ctx context.Context, url, source string, retries int) (Address, error) {
	var address Address
	var err error

	for i := 0; i < retries; i++ {
		address, err = fetchAPI(ctx, url, source)
		if err == nil {
			return address, nil
		}
		time.Sleep(time.Duration(i) * 100 * time.Millisecond)
	}

	return Address{}, err
}

func FetchFastestAPI(ctx context.Context, cep string, config Config) (Address, string, error) {
	ctx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	result := make(chan APIResponse, 1)
	errChan := make(chan error, 2)

	apis := map[string]string{
		"BrasilAPI": config.BrasilAPIURL + cep,
		"ViaCEP":    config.ViaCEPURL + cep + "/json",
	}

	for source, url := range apis {
		go func(ctx context.Context, url, source string) {
			select {
			case <-ctx.Done():
				return
			default:
				address, err := fetchAPIWithRetry(ctx, url, source, 3)
				if err != nil {
					errChan <- err
					return
				}
				select {
				case result <- APIResponse{Result: address, Source: source}:
					cancel()
				case <-ctx.Done():
				}
			}
		}(ctx, url, source)
	}

	select {
	case res := <-result:
		return res.Result, res.Source, nil
	case <-ctx.Done():
		return Address{}, "", errors.New("timeout")
	case err := <-errChan:
		return Address{}, "", err
	}
}

func main() {
	config := loadConfig()
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	cep := "01153000"
	result, source, err := FetchFastestAPI(ctx, cep, config)
	if err != nil {
		log.Println("Erro:", err)
	} else {
		log.Printf("Resultado da API %s: %+v\n", source, result)
	}
}
