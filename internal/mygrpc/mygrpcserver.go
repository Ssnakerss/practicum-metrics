package mygrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Ssnakerss/practicum-metrics/internal/app"
	"github.com/Ssnakerss/practicum-metrics/internal/compression"
	"github.com/Ssnakerss/practicum-metrics/internal/dtadapter"
	"github.com/Ssnakerss/practicum-metrics/internal/encrypt"
	"github.com/Ssnakerss/practicum-metrics/internal/hash"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/proto"
	"go.uber.org/zap"
)

type Server struct {
	//grpc
	proto.UnimplementedMetricsServer

	Da *dtadapter.Adapter
	C  *app.ServerConfig
	E  *encrypt.Coder
	L  *zap.SugaredLogger
}

// grpc function
func (s *Server) SaveJSONMetrics(ctx context.Context, in *proto.JSONSaveRequest) (*proto.JSONSaveResponse, error) {
	response := proto.JSONSaveResponse{}

	mcsj := make([]metric.MetricJSON, 0)
	mcs := make([]metric.Metric, 0)

	//in.JSONMetrics - merics  json array
	//workflow:  check hash -> decrypt -> unzip -> unmarhal
	body := in.JSONMetrics

	//checking hash
	calcHash, err := hash.MakeSHA256(body, s.C.Key)
	if err != nil {
		s.L.Warnw("grpc savejsonmetrics", "cannot calculate hash", err)
		response.Error = fmt.Sprintf("hash calculation error : %s", err.Error())
		return &response, err
	}
	if calcHash != in.Hash {
		s.L.Warn("grpc savejsonmetrics", "has mismatch")
		response.Error = "hash check error : hash mismatch"
		return &response, err
	}

	//decrypting
	if s.E != nil {
		body, err = s.E.Decrypt(body)
		if err != nil {
			s.L.Warnw("grpc savejsonmetrics", "cannot decrypt message body", err)
			response.Error = fmt.Sprintf("metrics decrypt error : %s", err.Error())
			return &response, err
		}
	}

	//unzipping
	body, err = compression.Decompress(body)
	if err != nil {
		s.L.Warnw("grpc savejsonmetrics", "cannot unzip message body", err)
		response.Error = fmt.Sprintf("metrics unzip error : %s", err.Error())
		return &response, err
	}

	//unmarshal
	//tyrying to convert json to []metric
	if err := json.Unmarshal(body, &mcsj); err != nil {
		s.L.Warnw("grpc savejsonmetrics", "cannot convert json to []metric", err)
		response.Error = fmt.Sprintf("metrics convertion error : %s", err)
		return &response, err
	}

	//converting from transport to store format
	for _, m := range mcsj {
		mcs = append(mcs, *metric.ConvertMetricI2S(&m))
	}

	// //Записываем получившийся массив в хранилище
	err = s.Da.WriteAll(&mcs)
	if err != nil {
		s.L.Warnw("grpc savejsonmetrics", "data save error", err)
		response.Error = fmt.Sprintf("error saving to storage : %s", err)
		return &response, err
	}

	s.L.Infow("grpc savejsonmetrics", "saved items", len(mcs))
	response.Message = fmt.Sprintf("saved %d items", len(mcs))
	return &response, nil
}
