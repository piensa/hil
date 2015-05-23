/*
 * Minimal object storage library (C) 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package objectstorage

// bucketHandler is an http.Handler that verifies bucket responses and validates incoming requests
import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"
)

type bucketHandler struct {
	resource string
}

func (h bucketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		switch {
		case r.URL.Path == "/":
			response := []byte("<ListAllMyBucketsResult xmlns=\"http://doc.s3.amazonaws.com/2006-03-01\"><Buckets><Bucket><Name>bucket</Name><CreationDate>2015-05-20T23:05:09.230Z</CreationDate></Bucket></Buckets><Owner><ID>minio</ID><DisplayName>minio</DisplayName></Owner></ListAllMyBucketsResult>")
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			w.Write(response)
		case r.URL.Path == "/bucket":
			response := []byte("<ListBucketResult xmlns=\"http://doc.s3.amazonaws.com/2006-03-01\"><Contents><ETag>259d04a13802ae09c7e41be50ccc6baa</ETag><Key>object</Key><LastModified>2015-05-21T18:24:21.097Z</LastModified><Size>22061</Size><Owner><ID>minio</ID><DisplayName>minio</DisplayName></Owner><StorageClass>STANDARD</StorageClass></Contents><Delimiter></Delimiter><EncodingType></EncodingType><IsTruncated>false</IsTruncated><Marker></Marker><MaxKeys>1000</MaxKeys><Name>testbucket</Name><NextMarker></NextMarker><Prefix></Prefix></ListBucketResult>")
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			w.Write(response)
		}
	case r.Method == "PUT":
		switch {
		case r.URL.Path == h.resource:
			_, ok := r.URL.Query()["acl"]
			if ok {
				if r.Header.Get("x-amz-acl") != "public-read-write" {
					w.WriteHeader(http.StatusNotImplemented)
				}
			}
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	case r.Method == "HEAD":
		switch {
		case r.URL.Path == h.resource:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusForbidden)
		}
	case r.Method == "DELETE":
		switch {
		case r.URL.Path != h.resource:
			w.WriteHeader(http.StatusNotFound)
		default:
			h.resource = ""
			w.WriteHeader(http.StatusOK)
		}
	}
}

// objectHandler is an http.Handler that verifies object responses and validates incoming requests
type objectHandler struct {
	resource string
	data     []byte
}

func (h objectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "PUT":
		length, err := strconv.Atoi(r.Header.Get("Content-Length"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var buffer bytes.Buffer
		_, err = io.CopyN(&buffer, r.Body, int64(length))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !bytes.Equal(h.data, buffer.Bytes()) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("ETag", "9af2f8218b150c351ad802c6f3d66abe")
		w.WriteHeader(http.StatusOK)
	case r.Method == "HEAD":
		if r.URL.Path != h.resource {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Length", strconv.Itoa(len(h.data)))
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		w.Header().Set("ETag", "9af2f8218b150c351ad802c6f3d66abe")
		w.WriteHeader(http.StatusOK)
	case r.Method == "POST":
	case r.Method == "GET":
		if r.URL.Path != h.resource {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Length", strconv.Itoa(len(h.data)))
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		w.Header().Set("ETag", "9af2f8218b150c351ad802c6f3d66abe")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, bytes.NewReader(h.data))
	case r.Method == "DELETE":
		if r.URL.Path != h.resource {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.resource = ""
		h.data = nil
		w.WriteHeader(http.StatusOK)
	}
}