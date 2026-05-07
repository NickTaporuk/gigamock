package scenarioType

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/NickTaporuk/gigamock/src/scenarios"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type s3Metrics struct {
	mu      sync.RWMutex
	Buckets map[string]*s3Metric `json:"buckets"`
}

type s3Metric struct {
	Calls   int64 `json:"calls"`
	Gets    int64 `json:"gets"`
	Puts    int64 `json:"puts"`
	Deletes int64 `json:"deletes"`
	Lists   int64 `json:"lists"`
	Errors  int64 `json:"errors"`
	DryRuns int64 `json:"dryRuns"`
}

type s3Object struct {
	Key          string
	Body         []byte
	ContentType  string
	ETag         string
	Metadata     map[string]string
	LastModified time.Time
}

type s3ObjectStore struct {
	mu      sync.RWMutex
	Buckets map[string]map[string]s3Object
}

var (
	s3RuntimeMetrics = &s3Metrics{Buckets: map[string]*s3Metric{}}
	s3RuntimeStore   = &s3ObjectStore{Buckets: map[string]map[string]s3Object{}}
)

// S3Provider serves a small S3-compatible path-style API backed by memory.
type S3Provider struct {
	w         http.ResponseWriter
	req       *http.Request
	scenarios scenarios.S3Scenarios
	logger    *logrus.Entry
}

func NewS3Provider(w http.ResponseWriter, req *http.Request, lgr *logrus.Entry) *S3Provider {
	return &S3Provider{w: w, req: req, logger: lgr}
}

func (s *S3Provider) Unmarshal(rawScenarios []map[string]interface{}) error {
	return mapstructure.Decode(rawScenarios, &s.scenarios)
}

func (s S3Provider) Validate() error {
	if len(s.scenarios) == 0 {
		return fmt.Errorf("s3 scenarios are required")
	}
	for index, scenario := range s.scenarios {
		if err := scenario.Validate(); err != nil {
			return fmt.Errorf("s3 scenario %d is invalid: %w", index, err)
		}
	}
	return nil
}

func (s *S3Provider) Retrieve(scenarioNumber int) {
	scenario, err := s.scenarioByNumber(scenarioNumber)
	if err != nil {
		s.recordS3Metric("unknown", s3Metric{Calls: 1, Errors: 1})
		s.writeJSONError(http.StatusInternalServerError, "s3 scenario failed")
		return
	}

	bucket, key := s.bucketAndKey(scenario)
	if bucket == "" {
		s.recordS3Metric("unknown", s3Metric{Calls: 1, Errors: 1})
		s.writeJSONError(http.StatusBadRequest, "s3 bucket is required")
		return
	}

	s.recordS3Metric(bucket, s3Metric{Calls: 1})
	if scenario.DryRun {
		s.recordS3Metric(bucket, s3Metric{DryRuns: 1})
		s.writeDryRunResponse(bucket, key)
		return
	}

	switch s.req.Method {
	case http.MethodPut:
		s.putObject(bucket, key, scenario)
	case http.MethodGet:
		if key == "" {
			s.listBucket(bucket)
			return
		}
		s.getObject(bucket, key, scenario)
	case http.MethodHead:
		s.headObject(bucket, key, scenario)
	case http.MethodDelete:
		s.deleteObject(bucket, key)
	default:
		s.recordS3Metric(bucket, s3Metric{Errors: 1})
		s.writeJSONError(http.StatusMethodNotAllowed, "s3 method is not supported")
	}
}

func (s *S3Provider) putObject(bucket string, key string, scenario scenarios.S3Scenario) {
	if key == "" {
		s.recordS3Metric(bucket, s3Metric{Errors: 1})
		s.writeJSONError(http.StatusBadRequest, "s3 object key is required")
		return
	}

	body, err := io.ReadAll(s.req.Body)
	if err != nil {
		s.recordS3Metric(bucket, s3Metric{Errors: 1})
		s.writeJSONError(http.StatusBadRequest, "failed to read object body")
		return
	}
	contentType := s.req.Header.Get("Content-Type")
	if contentType == "" {
		contentType = scenario.ContentType
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	object := s3Object{
		Key:          key,
		Body:         body,
		ContentType:  contentType,
		ETag:         s3ETag(body),
		Metadata:     s3MetadataFromHeaders(s.req.Header, scenario.Metadata),
		LastModified: time.Now().UTC(),
	}
	s3RuntimeStore.put(bucket, object)
	s.recordS3Metric(bucket, s3Metric{Puts: 1})

	s.w.Header().Set("ETag", object.ETag)
	s.w.WriteHeader(http.StatusOK)
}

func (s *S3Provider) getObject(bucket string, key string, scenario scenarios.S3Scenario) {
	object, ok := s3RuntimeStore.get(bucket, key)
	if !ok && scenario.Body != "" {
		object = s3Object{
			Key:          key,
			Body:         []byte(scenario.Body),
			ContentType:  defaultString(scenario.ContentType, "application/octet-stream"),
			ETag:         s3ETag([]byte(scenario.Body)),
			Metadata:     scenario.Metadata,
			LastModified: time.Now().UTC(),
		}
		ok = true
	}
	if !ok {
		s.recordS3Metric(bucket, s3Metric{Errors: 1})
		s.writeS3Error(http.StatusNotFound, "NoSuchKey", "The specified key does not exist.")
		return
	}

	s.writeObjectHeaders(object, scenario.Headers)
	s.recordS3Metric(bucket, s3Metric{Gets: 1})
	s.w.WriteHeader(http.StatusOK)
	s.w.Write(object.Body)
}

func (s *S3Provider) headObject(bucket string, key string, scenario scenarios.S3Scenario) {
	object, ok := s3RuntimeStore.get(bucket, key)
	if !ok && scenario.Body != "" {
		object = s3Object{
			Key:          key,
			Body:         []byte(scenario.Body),
			ContentType:  defaultString(scenario.ContentType, "application/octet-stream"),
			ETag:         s3ETag([]byte(scenario.Body)),
			Metadata:     scenario.Metadata,
			LastModified: time.Now().UTC(),
		}
		ok = true
	}
	if !ok {
		s.recordS3Metric(bucket, s3Metric{Errors: 1})
		s.writeS3Error(http.StatusNotFound, "NoSuchKey", "The specified key does not exist.")
		return
	}

	s.writeObjectHeaders(object, scenario.Headers)
	s.recordS3Metric(bucket, s3Metric{Gets: 1})
	s.w.WriteHeader(http.StatusOK)
}

func (s *S3Provider) deleteObject(bucket string, key string) {
	if key == "" {
		s.recordS3Metric(bucket, s3Metric{Errors: 1})
		s.writeJSONError(http.StatusBadRequest, "s3 object key is required")
		return
	}
	s3RuntimeStore.delete(bucket, key)
	s.recordS3Metric(bucket, s3Metric{Deletes: 1})
	s.w.WriteHeader(http.StatusNoContent)
}

func (s *S3Provider) listBucket(bucket string) {
	keys := s3RuntimeStore.keys(bucket)
	response := s3ListBucketResult{
		XMLName:  xml.Name{Local: "ListBucketResult"},
		Name:     bucket,
		KeyCount: len(keys),
		MaxKeys:  len(keys),
		Contents: make([]s3ListBucketContent, 0, len(keys)),
	}
	for _, key := range keys {
		object, _ := s3RuntimeStore.get(bucket, key)
		response.Contents = append(response.Contents, s3ListBucketContent{
			Key:          key,
			LastModified: object.LastModified.Format(time.RFC3339),
			ETag:         object.ETag,
			Size:         len(object.Body),
		})
	}

	data, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		s.recordS3Metric(bucket, s3Metric{Errors: 1})
		s.writeJSONError(http.StatusInternalServerError, "failed to render bucket list")
		return
	}
	s.recordS3Metric(bucket, s3Metric{Lists: 1})
	s.w.Header().Set("Content-Type", "application/xml")
	s.w.WriteHeader(http.StatusOK)
	s.w.Write([]byte(xml.Header))
	s.w.Write(data)
}

func (s *S3Provider) writeObjectHeaders(object s3Object, headers map[string]string) {
	s.w.Header().Set("Content-Type", object.ContentType)
	s.w.Header().Set("Content-Length", fmt.Sprint(len(object.Body)))
	s.w.Header().Set("ETag", object.ETag)
	s.w.Header().Set("Last-Modified", object.LastModified.Format(http.TimeFormat))
	for key, value := range object.Metadata {
		s.w.Header().Set("x-amz-meta-"+key, value)
	}
	for key, value := range headers {
		s.w.Header().Set(key, value)
	}
}

func (s S3Provider) scenarioByNumber(scenarioNumber int) (scenarios.S3Scenario, error) {
	if len(s.scenarios) == 0 {
		return scenarios.S3Scenario{}, fmt.Errorf("s3 scenarios are empty")
	}
	if scenarioNumber < 0 || scenarioNumber >= len(s.scenarios) {
		return s.scenarios[0], nil
	}
	return s.scenarios[scenarioNumber], nil
}

func (s S3Provider) bucketAndKey(scenario scenarios.S3Scenario) (string, string) {
	bucket := scenario.Bucket
	key := scenario.Key
	if s.req == nil {
		return bucket, key
	}

	parts := strings.Split(strings.TrimPrefix(s.req.URL.Path, "/"), "/")
	if len(parts) >= 2 && (parts[0] == "s3" || parts[0] == "s3-dry-run") {
		bucket = defaultString(bucket, parts[1])
	}
	if len(parts) >= 3 {
		key = defaultString(key, strings.Join(parts[2:], "/"))
	}
	return bucket, key
}

func (s *S3Provider) writeDryRunResponse(bucket string, key string) {
	s.w.Header().Set("Content-Type", "application/json")
	s.w.WriteHeader(http.StatusOK)
	json.NewEncoder(s.w).Encode(map[string]interface{}{
		"s3":     true,
		"dryRun": true,
		"bucket": bucket,
		"key":    key,
	})
}

func (s *S3Provider) writeJSONError(statusCode int, message string) {
	s.w.Header().Set("Content-Type", "application/json")
	s.w.WriteHeader(statusCode)
	json.NewEncoder(s.w).Encode(map[string]string{"error": message})
}

func (s *S3Provider) writeS3Error(statusCode int, code string, message string) {
	data, _ := xml.Marshal(s3Error{Code: code, Message: message})
	s.w.Header().Set("Content-Type", "application/xml")
	s.w.WriteHeader(statusCode)
	s.w.Write([]byte(xml.Header))
	s.w.Write(data)
}

func (s S3Provider) recordS3Metric(bucket string, delta s3Metric) {
	recordS3Metric(bucket, delta)
}

func recordS3Metric(bucket string, delta s3Metric) {
	if bucket == "" {
		bucket = "unknown"
	}
	s3RuntimeMetrics.mu.Lock()
	defer s3RuntimeMetrics.mu.Unlock()

	metric, ok := s3RuntimeMetrics.Buckets[bucket]
	if !ok {
		metric = &s3Metric{}
		s3RuntimeMetrics.Buckets[bucket] = metric
	}
	metric.Calls += delta.Calls
	metric.Gets += delta.Gets
	metric.Puts += delta.Puts
	metric.Deletes += delta.Deletes
	metric.Lists += delta.Lists
	metric.Errors += delta.Errors
	metric.DryRuns += delta.DryRuns
}

// S3MetricsSnapshot returns a copy of S3 runtime metrics.
func S3MetricsSnapshot() map[string]s3Metric {
	s3RuntimeMetrics.mu.RLock()
	defer s3RuntimeMetrics.mu.RUnlock()

	out := make(map[string]s3Metric, len(s3RuntimeMetrics.Buckets))
	for bucket, metric := range s3RuntimeMetrics.Buckets {
		out[bucket] = *metric
	}
	return out
}

func (s *s3ObjectStore) put(bucket string, object s3Object) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Buckets[bucket]; !ok {
		s.Buckets[bucket] = map[string]s3Object{}
	}
	s.Buckets[bucket][object.Key] = object
}

func (s *s3ObjectStore) get(bucket string, key string) (s3Object, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	objects, ok := s.Buckets[bucket]
	if !ok {
		return s3Object{}, false
	}
	object, ok := objects[key]
	return object, ok
}

func (s *s3ObjectStore) delete(bucket string, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if objects, ok := s.Buckets[bucket]; ok {
		delete(objects, key)
	}
}

func (s *s3ObjectStore) keys(bucket string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	objects := s.Buckets[bucket]
	keys := make([]string, 0, len(objects))
	for key := range objects {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func s3ETag(body []byte) string {
	sum := md5.Sum(body)
	return `"` + hex.EncodeToString(sum[:]) + `"`
}

func s3MetadataFromHeaders(headers http.Header, defaults map[string]string) map[string]string {
	metadata := map[string]string{}
	for key, value := range defaults {
		metadata[key] = value
	}
	for key, values := range headers {
		lower := strings.ToLower(key)
		if strings.HasPrefix(lower, "x-amz-meta-") && len(values) > 0 {
			metadata[strings.TrimPrefix(lower, "x-amz-meta-")] = values[0]
		}
	}
	return metadata
}

func defaultString(value string, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

type s3ListBucketResult struct {
	XMLName  xml.Name              `xml:"ListBucketResult"`
	Name     string                `xml:"Name"`
	KeyCount int                   `xml:"KeyCount"`
	MaxKeys  int                   `xml:"MaxKeys"`
	Contents []s3ListBucketContent `xml:"Contents"`
}

type s3ListBucketContent struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	ETag         string `xml:"ETag"`
	Size         int    `xml:"Size"`
}

type s3Error struct {
	XMLName xml.Name `xml:"Error"`
	Code    string   `xml:"Code"`
	Message string   `xml:"Message"`
}
