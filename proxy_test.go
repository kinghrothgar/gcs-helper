package main

import (
	"net/http"
	"testing"
)

func TestServerProxyOnly(t *testing.T) {
	addr, cleanup := startServer(t, Config{BucketName: "my-bucket"})
	defer cleanup()
	var tests = []serverTest{
		{
			"healthcheck through the proxy",
			"GET",
			addr,
			nil,
			http.StatusOK,
			nil,
			"",
		},
		{
			"download file",
			"GET",
			addr + "/musics/music/music1.txt",
			nil,
			http.StatusOK,
			http.Header{
				"Accept-Ranges":  []string{"bytes"},
				"Content-Length": []string{"15"},
			},
			"some nice music",
		},
		{
			"download file - range",
			"GET",
			addr + "/musics/music/music2.txt",
			http.Header{
				"Range": []string{"bytes=2-10"},
			},
			http.StatusPartialContent,
			http.Header{
				"Accept-Ranges":  []string{"bytes"},
				"Content-Length": []string{"8"},
				"Content-Range":  []string{"bytes 2-10/16"},
			},
			"me nicer",
		},
		{
			"file attrs",
			"HEAD",
			addr + "/musics/music/music2.txt",
			nil,
			http.StatusOK,
			http.Header{
				"Accept-Ranges":  []string{"bytes"},
				"Content-Length": []string{"16"},
			},
			"",
		},
		{
			"download file - object not found",
			"GET",
			addr + "/musics/music/some-music.txt",
			nil,
			http.StatusNotFound,
			nil,
			"storage: object doesn't exist\n",
		},
		{
			"file attrs - object not found",
			"HEAD",
			addr + "/musics/music/some-music.txt",
			nil,
			http.StatusNotFound,
			nil,
			"",
		},
		{
			"method not allowed - POST",
			"POST",
			addr + "/whatever",
			nil,
			http.StatusMethodNotAllowed,
			nil,
			"method not allowed\n",
		},
		{
			"method not allowed - PUT",
			"PUT",
			addr + "/whatever",
			nil,
			http.StatusMethodNotAllowed,
			nil,
			"method not allowed\n",
		},
	}
	for _, test := range tests {
		t.Run(test.testCase, test.run)
	}
}

func TestServerProxyHandlerBucketNotFound(t *testing.T) {
	addr, cleanup := startServer(t, Config{BucketName: "some-bucket"})
	defer cleanup()
	req, _ := http.NewRequest("HEAD", addr+"/whatever", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("wrong status code\nwant %d\ngot  %d", http.StatusNotFound, resp.StatusCode)
	}
}
