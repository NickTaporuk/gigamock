package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"gopkg.in/yaml.v3"

	"github.com/NickTaporuk/gigamock/src/common"
	"github.com/NickTaporuk/gigamock/src/fileWalkers"
)

type grpcMockService interface {
	mustEmbedGRPCMockService()
}

type grpcRuntime struct {
	store   *map[string]fileWalkers.IndexedData
	logger  *logrus.Entry
	methods map[string]*grpcMethodRuntime
}

func (g *grpcRuntime) mustEmbedGRPCMockService() {}

type grpcMockFile struct {
	Path      string             `json:"path" yaml:"path"`
	Method    string             `json:"method" yaml:"method"`
	Type      string             `json:"type" yaml:"type"`
	Proto     grpcProtoConfig    `json:"proto" yaml:"proto"`
	Scenarios []grpcMockScenario `json:"scenarios" yaml:"scenarios"`
}

type grpcProtoConfig struct {
	File          string   `json:"file" yaml:"file"`
	DescriptorSet string   `json:"descriptorSet" yaml:"descriptorSet"`
	ImportPaths   []string `json:"importPaths" yaml:"importPaths"`
	Service       string   `json:"service" yaml:"service"`
	Method        string   `json:"method" yaml:"method"`
}

type grpcMockScenario struct {
	Name     string             `json:"name" yaml:"name"`
	Request  grpcRequestConfig  `json:"request" yaml:"request"`
	Response grpcResponseConfig `json:"response" yaml:"response"`
	Stream   grpcStreamConfig   `json:"stream" yaml:"stream"`
}

type grpcRequestConfig struct {
	Match map[string]interface{} `json:"match" yaml:"match"`
}

type grpcResponseConfig struct {
	Code     string                 `json:"code" yaml:"code"`
	Message  string                 `json:"message" yaml:"message"`
	Metadata map[string]string      `json:"metadata" yaml:"metadata"`
	Body     map[string]interface{} `json:"body" yaml:"body"`
}

type grpcStreamConfig struct {
	SendOnConnect []map[string]interface{} `json:"sendOnConnect" yaml:"sendOnConnect"`
	Steps         []grpcStreamStep         `json:"steps" yaml:"steps"`
	OnReceive     []grpcStreamRule         `json:"onReceive" yaml:"onReceive"`
}

type grpcStreamStep struct {
	Receive map[string]interface{} `json:"receive" yaml:"receive"`
	Send    map[string]interface{} `json:"send" yaml:"send"`
	Close   *grpcCloseConfig       `json:"close" yaml:"close"`
}

type grpcStreamRule struct {
	Match map[string]interface{} `json:"match" yaml:"match"`
	Send  map[string]interface{} `json:"send" yaml:"send"`
	Close *grpcCloseConfig       `json:"close" yaml:"close"`
}

type grpcCloseConfig struct {
	Code    string `json:"code" yaml:"code"`
	Message string `json:"message" yaml:"message"`
}

type grpcMethodRuntime struct {
	key      string
	filePath string
	config   grpcMockFile
	method   protoreflect.MethodDescriptor
	store    *map[string]fileWalkers.IndexedData
}

func (di *Dispatcher) startGRPCServer(addr string) (*grpc.Server, net.Listener, error) {
	runtime, files, services, err := di.buildGRPCRuntime()
	if err != nil {
		return nil, nil, err
	}
	if len(runtime.methods) == 0 {
		di.logger.Info("gRPC mock server skipped: no type=grpc scenarios found")
		return nil, nil, nil
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	grpcServer := grpc.NewServer()
	for serviceName, service := range services {
		serviceDesc := grpc.ServiceDesc{
			ServiceName: serviceName,
			HandlerType: (*grpcMockService)(nil),
			Methods:     service.unaryMethods(),
			Streams:     service.streamMethods(),
			Metadata:    service.metadata,
		}
		grpcServer.RegisterService(&serviceDesc, runtime)
	}

	rpb.RegisterServerReflectionServer(grpcServer, reflection.NewServer(reflection.ServerOptions{
		Services:           grpcServer,
		DescriptorResolver: files.AsResolver(),
	}))

	go func() {
		di.logger.Infof("Ready to accept gRPC connections on %s", addr)
		if err := grpcServer.Serve(lis); err != nil {
			di.logger.WithError(err).Error("gRPC server retrieved an error")
		}
	}()

	return grpcServer, lis, nil
}

type grpcServiceRuntime struct {
	metadata string
	methods  []*grpcMethodRuntime
}

func (s *grpcServiceRuntime) unaryMethods() []grpc.MethodDesc {
	methods := make([]grpc.MethodDesc, 0, len(s.methods))
	for _, method := range s.methods {
		if method.method.IsStreamingClient() || method.method.IsStreamingServer() {
			continue
		}
		current := method
		methods = append(methods, grpc.MethodDesc{
			MethodName: string(current.method.Name()),
			Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				runtime := srv.(*grpcRuntime)
				req := dynamicpb.NewMessage(current.method.Input())
				if err := dec(req); err != nil {
					return nil, err
				}
				handler := func(ctx context.Context, request interface{}) (interface{}, error) {
					return runtime.invokeUnary(ctx, current, request.(proto.Message))
				}
				if interceptor == nil {
					return handler(ctx, req)
				}
				return interceptor(ctx, req, &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/" + string(current.method.Parent().FullName()) + "/" + string(current.method.Name()),
				}, handler)
			},
		})
	}
	return methods
}

func (s *grpcServiceRuntime) streamMethods() []grpc.StreamDesc {
	streams := make([]grpc.StreamDesc, 0, len(s.methods))
	for _, method := range s.methods {
		if !method.method.IsStreamingClient() && !method.method.IsStreamingServer() {
			continue
		}
		current := method
		streams = append(streams, grpc.StreamDesc{
			StreamName: string(current.method.Name()),
			Handler: func(srv interface{}, stream grpc.ServerStream) error {
				return srv.(*grpcRuntime).invokeStream(current, stream)
			},
			ServerStreams: current.method.IsStreamingServer(),
			ClientStreams: current.method.IsStreamingClient(),
		})
	}
	return streams
}

func (di *Dispatcher) buildGRPCRuntime() (*grpcRuntime, linker.Files, map[string]*grpcServiceRuntime, error) {
	runtime := &grpcRuntime{
		store:   &di.indexedFiles,
		logger:  di.logger,
		methods: map[string]*grpcMethodRuntime{},
	}
	services := map[string]*grpcServiceRuntime{}
	var allFiles linker.Files
	compiledFiles := map[string]linker.Files{}

	for key, indexedData := range di.indexedFiles {
		config, err := readGRPCMockFile(indexedData.FilePath)
		if err != nil {
			return nil, nil, nil, err
		}
		if config.Type != common.GRPCScenarioType {
			continue
		}

		serviceName, methodName, err := grpcServiceAndMethod(config)
		if err != nil {
			return nil, nil, nil, err
		}

		protoFile := config.Proto.File
		if protoFile == "" {
			return nil, nil, nil, fmt.Errorf("grpc scenario %s must define proto.file", indexedData.FilePath)
		}

		files, ok := compiledFiles[protoFile]
		if !ok {
			files, err = compileProto(di.ctx, protoFile, config.Proto.ImportPaths)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("compile proto %s: %w", protoFile, err)
			}
			compiledFiles[protoFile] = files
			allFiles = append(allFiles, files...)
		}

		desc, err := files.AsResolver().FindDescriptorByName(protoreflect.FullName(serviceName))
		if err != nil {
			return nil, nil, nil, err
		}
		serviceDesc, ok := desc.(protoreflect.ServiceDescriptor)
		if !ok {
			return nil, nil, nil, fmt.Errorf("%s is not a gRPC service", serviceName)
		}
		methodDesc := serviceDesc.Methods().ByName(protoreflect.Name(methodName))
		if methodDesc == nil {
			return nil, nil, nil, fmt.Errorf("method %s/%s is not found in %s", serviceName, methodName, protoFile)
		}

		fullMethod := "/" + serviceName + "/" + methodName
		methodRuntime := &grpcMethodRuntime{
			key:      key,
			filePath: indexedData.FilePath,
			config:   config,
			method:   methodDesc,
			store:    &di.indexedFiles,
		}
		runtime.methods[fullMethod] = methodRuntime
		if _, ok := services[serviceName]; !ok {
			services[serviceName] = &grpcServiceRuntime{metadata: protoFile}
		}
		services[serviceName].methods = append(services[serviceName].methods, methodRuntime)
	}

	return runtime, allFiles, services, nil
}

func readGRPCMockFile(filePath string) (grpcMockFile, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return grpcMockFile{}, err
	}

	var config grpcMockFile
	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".json":
		err = json.Unmarshal(raw, &config)
	default:
		err = yaml.Unmarshal(raw, &config)
	}
	return config, err
}

func compileProto(ctx context.Context, protoFile string, importPaths []string) (linker.Files, error) {
	target := filepath.ToSlash(protoFile)
	if filepath.IsAbs(protoFile) {
		target = filepath.Base(protoFile)
		if len(importPaths) == 0 {
			importPaths = []string{filepath.Dir(protoFile)}
		}
	}
	if len(importPaths) == 0 {
		importPaths = []string{"."}
	}
	for i := range importPaths {
		importPaths[i] = filepath.Clean(importPaths[i])
	}

	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{
			ImportPaths: importPaths,
		}),
	}
	return compiler.Compile(ctx, target)
}

func grpcServiceAndMethod(config grpcMockFile) (string, string, error) {
	serviceName := config.Proto.Service
	methodName := config.Proto.Method
	if serviceName != "" && methodName != "" {
		return serviceName, methodName, nil
	}

	parts := strings.Split(strings.TrimPrefix(config.Path, "/"), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("grpc path must use /package.Service/Method format: %s", config.Path)
	}
	if serviceName == "" {
		serviceName = parts[0]
	}
	if methodName == "" {
		methodName = parts[1]
	}
	return serviceName, methodName, nil
}

func (g *grpcRuntime) invokeUnary(ctx context.Context, method *grpcMethodRuntime, req proto.Message) (interface{}, error) {
	scenario := method.activeScenario(req)
	if len(scenario.Response.Metadata) > 0 {
		pairs := make([]string, 0, len(scenario.Response.Metadata)*2)
		for key, value := range scenario.Response.Metadata {
			pairs = append(pairs, key, value)
		}
		grpc.SetHeader(ctx, metadata.Pairs(pairs...))
	}

	code := grpcStatusCode(scenario.Response.Code)
	if code != codes.OK {
		return nil, status.Error(code, scenario.Response.Message)
	}

	return buildDynamicMessage(method.method.Output(), scenario.Response.Body)
}

func (g *grpcRuntime) invokeStream(method *grpcMethodRuntime, stream grpc.ServerStream) error {
	scenario := method.activeScenario(nil)
	for _, payload := range scenario.Stream.SendOnConnect {
		if err := sendDynamicMessage(stream, method.method.Output(), payload); err != nil {
			return err
		}
	}

	for _, step := range scenario.Stream.Steps {
		if step.Receive != nil {
			msg := dynamicpb.NewMessage(method.method.Input())
			if err := stream.RecvMsg(msg); err != nil {
				return err
			}
			if !matchesMessage(msg, step.Receive) {
				return status.Error(codes.InvalidArgument, "stream receive payload does not match scenario")
			}
		}
		if step.Send != nil {
			if err := sendDynamicMessage(stream, method.method.Output(), step.Send); err != nil {
				return err
			}
		}
		if step.Close != nil {
			return grpcCloseError(step.Close)
		}
	}

	if len(scenario.Stream.OnReceive) == 0 {
		return nil
	}

	for {
		msg := dynamicpb.NewMessage(method.method.Input())
		if err := stream.RecvMsg(msg); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		for _, rule := range scenario.Stream.OnReceive {
			if !matchesMessage(msg, rule.Match) {
				continue
			}
			if rule.Send != nil {
				if err := sendDynamicMessage(stream, method.method.Output(), rule.Send); err != nil {
					return err
				}
			}
			if rule.Close != nil {
				return grpcCloseError(rule.Close)
			}
		}
	}
}

func (m *grpcMethodRuntime) activeScenario(req proto.Message) grpcMockScenario {
	scenarioNumber := 0
	if m.store != nil {
		if indexedData, ok := (*m.store)[m.key]; ok {
			scenarioNumber = indexedData.ScenarioNumber
		}
	}
	if scenarioNumber >= 0 && scenarioNumber < len(m.config.Scenarios) {
		selected := m.config.Scenarios[scenarioNumber]
		if scenarioNumber != 0 || req == nil || len(selected.Request.Match) == 0 || matchesMessage(req, selected.Request.Match) {
			return selected
		}
	}
	if req != nil {
		for _, scenario := range m.config.Scenarios {
			if matchesMessage(req, scenario.Request.Match) {
				return scenario
			}
		}
	}
	if len(m.config.Scenarios) == 0 {
		return grpcMockScenario{Response: grpcResponseConfig{Code: codes.NotFound.String(), Message: "scenario is not found"}}
	}
	return m.config.Scenarios[0]
}

func matchesMessage(msg proto.Message, match map[string]interface{}) bool {
	if len(match) == 0 {
		return true
	}
	raw, err := protojson.MarshalOptions{UseProtoNames: false}.Marshal(msg)
	if err != nil {
		return false
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return false
	}
	return matchesMap(payload, normalizeJSONValue(match).(map[string]interface{}))
}

func matchesMap(payload map[string]interface{}, expected map[string]interface{}) bool {
	for key, expectedValue := range expected {
		actualValue, ok := payload[key]
		if !ok {
			return false
		}
		expectedMap, expectedIsMap := expectedValue.(map[string]interface{})
		actualMap, actualIsMap := actualValue.(map[string]interface{})
		if expectedIsMap {
			if !actualIsMap || !matchesMap(actualMap, expectedMap) {
				return false
			}
			continue
		}
		if fmt.Sprint(actualValue) != fmt.Sprint(expectedValue) {
			return false
		}
	}
	return true
}

func buildDynamicMessage(desc protoreflect.MessageDescriptor, body map[string]interface{}) (proto.Message, error) {
	msg := dynamicpb.NewMessage(desc)
	if len(body) == 0 {
		return msg, nil
	}
	raw, err := json.Marshal(normalizeJSONValue(body))
	if err != nil {
		return nil, err
	}
	if err := protojson.Unmarshal(raw, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func sendDynamicMessage(stream grpc.ServerStream, desc protoreflect.MessageDescriptor, body map[string]interface{}) error {
	msg, err := buildDynamicMessage(desc, body)
	if err != nil {
		return err
	}
	return stream.SendMsg(msg)
}

func normalizeJSONValue(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(typed))
		for key, value := range typed {
			out[key] = normalizeJSONValue(value)
		}
		return out
	case map[interface{}]interface{}:
		out := make(map[string]interface{}, len(typed))
		for key, value := range typed {
			out[fmt.Sprint(key)] = normalizeJSONValue(value)
		}
		return out
	case []interface{}:
		out := make([]interface{}, 0, len(typed))
		for _, value := range typed {
			out = append(out, normalizeJSONValue(value))
		}
		return out
	default:
		return typed
	}
}

func grpcCloseError(closeConfig *grpcCloseConfig) error {
	code := grpcStatusCode(closeConfig.Code)
	if code == codes.OK {
		return nil
	}
	return status.Error(code, closeConfig.Message)
}

func grpcStatusCode(raw string) codes.Code {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "", "OK":
		return codes.OK
	case "CANCELLED":
		return codes.Canceled
	case "UNKNOWN":
		return codes.Unknown
	case "INVALID_ARGUMENT":
		return codes.InvalidArgument
	case "DEADLINE_EXCEEDED":
		return codes.DeadlineExceeded
	case "NOT_FOUND":
		return codes.NotFound
	case "ALREADY_EXISTS":
		return codes.AlreadyExists
	case "PERMISSION_DENIED":
		return codes.PermissionDenied
	case "RESOURCE_EXHAUSTED":
		return codes.ResourceExhausted
	case "FAILED_PRECONDITION":
		return codes.FailedPrecondition
	case "ABORTED":
		return codes.Aborted
	case "OUT_OF_RANGE":
		return codes.OutOfRange
	case "UNIMPLEMENTED":
		return codes.Unimplemented
	case "INTERNAL":
		return codes.Internal
	case "UNAVAILABLE":
		return codes.Unavailable
	case "DATA_LOSS":
		return codes.DataLoss
	case "UNAUTHENTICATED":
		return codes.Unauthenticated
	default:
		return codes.Unknown
	}
}
