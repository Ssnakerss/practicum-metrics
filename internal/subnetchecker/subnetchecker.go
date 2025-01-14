package subnetchecker

import (
	"fmt"
	"net"
	"net/http"

	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"go.uber.org/zap"
)

type SubNetChecker struct {
	trustedSubnet *net.IPNet
}

func NewSubNetChecker(subnet string) (*SubNetChecker, error) {
	_, s, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, fmt.Errorf("subnet %s parse CIDR error: %v", subnet, err)
	}
	return &SubNetChecker{s}, nil
}

func (s *SubNetChecker) IsTrusted(agentIP string) (bool, error) {
	ipToCheck := net.ParseIP(agentIP)
	if ipToCheck == nil {
		return false, fmt.Errorf("wrong real ip from agent: %s", agentIP)
	}
	return s.trustedSubnet.Contains(ipToCheck), nil
}

// middleware to use with router for ip filering
func (s *SubNetChecker) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		realIP := r.Header.Get("X-Real-IP")
		if realIP == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		logger.Log.Debug("subnetchecker middleware", zap.String("real ip", realIP))
		res, err := s.IsTrusted(realIP)
		if err != nil {
			logger.SLog.Warnf("subnetchecker middleware", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if res {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	})
}

func GetLocalIP() (string, error) {
	// Получаею список всех сетевых интерфейсов
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("interfaces getting error: %w", err)
	}

	for _, iface := range interfaces {
		// Получаю адрес каждого интерфейса
		addrs, err := iface.Addrs()
		if err != nil {

			continue
		}
		// Перебираю адреса и вывожу их
		for _, addr := range addrs {
			// проверяю, является ли это IP-адресом IPv4
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("ip address of host is empty")
}
