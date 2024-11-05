package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/5aradise/link-forge/internal/database"
	"github.com/5aradise/link-forge/internal/handlers/mocks"
	"github.com/5aradise/link-forge/internal/types"
	"github.com/5aradise/link-forge/pkg/api"
	"github.com/5aradise/link-forge/pkg/logger"
)

func TestURLHandlers(t *testing.T) {
	lMock := logger.NewMock()
	sMock := mocks.NewURLStorage(t)

	s := NewURLService(lMock, sMock)

	r := http.NewServeMux()
	r.HandleFunc(http.MethodPost+" /", s.CreateURL)
	r.HandleFunc(http.MethodGet+" /", s.ListURLs)
	r.HandleFunc(http.MethodGet+" /{alias}", s.RedirectURL)
	r.HandleFunc(http.MethodDelete+" /{alias}", s.DeleteURL)

	t.Run("Create", func(t *testing.T) {
		cases := map[string]struct {
			req  CreateURLRequest
			res  CreateURLResponse
			code int
		}{
			"Normal": {
				req: CreateURLRequest{
					URL:   "http://test.com",
					Alias: "test",
				},
				res: CreateURLResponse{
					Response: api.ResOK(),
					Alias:    "test",
				},
				code: http.StatusCreated,
			},
			"Same_alias": {
				req: CreateURLRequest{
					URL:   "http://test.com",
					Alias: "same",
				},
				res: CreateURLResponse{
					Response: api.ResError("alias already exists"),
				},
				code: http.StatusBadRequest,
			},
			"Empty_url": {
				req: CreateURLRequest{
					Alias: "test",
				},
				res: CreateURLResponse{
					Response: api.ResError("empty url field"),
				},
				code: http.StatusBadRequest,
			},
			"Invalid_url": {
				req: CreateURLRequest{
					URL:   "test.com",
					Alias: "test",
				},
				res: CreateURLResponse{
					Response: api.ResError("invalid url"),
				},
				code: http.StatusBadRequest,
			},
			// "Empty alias": {
			// 	req: CreateURLRequest{
			// 		URL:   "http://test.com",
			// 		Alias: "",
			// 	},
			// 	res: CreateURLResponse{
			// 		Response: api.ResOK(),
			// 		Alias:    "a",
			// 	},
			// },
		}

		sMock.On("CreateURL", context.Background(), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(func(ctx context.Context, alias string, url string) (types.URL, error) {
				if alias == "same" {
					return types.URL{}, database.ErrAliasExists
				}
				return types.URL{Id: 1, Alias: alias, Url: url}, nil
			})

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				reqBody, err := json.Marshal(tc.req)
				require.NoError(err)

				code, body, _, err := serveHTTP(r, http.MethodPost, "", reqBody)
				require.NoError(err)

				assert.Equal(tc.code, code)

				var res CreateURLResponse
				require.NoError(json.Unmarshal(body, &res))

				require.Equal(tc.res, res)
			})
		}
	})

	t.Run("List", func(t *testing.T) {
		cases := map[string]struct {
			res  ListURLsResponse
			code int
		}{
			"Normal": {
				res: ListURLsResponse{
					Response: api.ResOK(),
					URLs: []types.URL{
						{Id: 1, Alias: "a", Url: "http://test1.com"},
						{Id: 2, Alias: "b", Url: "http://test2.com"},
						{Id: 3, Alias: "c", Url: "http://test3.com"}},
				},
				code: http.StatusOK,
			},
		}

		sMock.On("ListURLs", context.Background()).
			Return([]types.URL{
				{Id: 1, Alias: "a", Url: "http://test1.com"},
				{Id: 2, Alias: "b", Url: "http://test2.com"},
				{Id: 3, Alias: "c", Url: "http://test3.com"}}, nil)

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				code, body, _, err := serveHTTP(r, http.MethodGet, "", []byte{})
				require.NoError(err)

				assert.Equal(tc.code, code)

				var res ListURLsResponse
				require.NoError(json.Unmarshal(body, &res))

				require.Equal(tc.res, res)
			})
		}
	})

	t.Run("Redirect", func(t *testing.T) {
		cases := map[string]struct {
			path string
			code int
			url  string
		}{
			"Normal": {
				path: "alias",
				url:  "http://test.com/",
				code: http.StatusFound,
			},
			"Wrong_alias": {
				path: "wrong",
				url:  "",
				code: http.StatusNotFound,
			},
		}

		sMock.On("GetURLByAlias", context.Background(), "alias").
			Return(types.URL{Id: 1, Alias: "alias", Url: "http://test.com/"}, nil)
		sMock.On("GetURLByAlias", context.Background(), "wrong").
			Return(types.URL{}, database.ErrURLUnfound)

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				assert := assert.New(t)

				code, _, head, _ := serveHTTP(r, http.MethodGet, tc.path, []byte{})

				assert.Equal(tc.code, code)

				assert.Equal(tc.url, head.Get("Location"))
			})
		}
	})

	t.Run("Delete", func(t *testing.T) {
		cases := map[string]struct {
			path string
			code int
			res  any
		}{
			"Normal": {
				path: "alias",
				code: http.StatusOK,
				res:  api.ResOK(),
			},
			"Wrong_alias": {
				path: "unfound",
				code: http.StatusBadRequest,
				res:  api.ResError("url with this alias unfound"),
			},
		}

		sMock.On("DeleteURLByAlias", context.Background(), "alias").
			Return(types.URL{Id: 1, Alias: "alias", Url: ""}, nil)
		sMock.On("DeleteURLByAlias", context.Background(), "unfound").
			Return(types.URL{}, database.ErrURLUnfound)

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				code, body, _, _ := serveHTTP(r, http.MethodDelete, tc.path, []byte{})

				assert.Equal(tc.code, code)

				var res api.Response
				require.NoError(json.Unmarshal(body, &res))

				require.Equal(tc.res, res)
			})
		}
	})
}

func serveHTTP(r http.Handler, method, path string, reqBody []byte) (code int, body []byte, header http.Header, err error) {
	req := httptest.NewRequest(method, "/"+path, bytes.NewReader(reqBody))
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	body, err = io.ReadAll(res.Body)

	return res.Code, body, res.Header(), err
}
