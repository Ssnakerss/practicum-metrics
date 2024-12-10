cls
go build -o ./cmd/server/  ./cmd/server/server.go
go build -o ./cmd/agent/  ./cmd/agent/agent.go

metricstest -test.v -test.run=^TestIteration1$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=.\ -server-port=8080 -file-storage-path=d:\temp\filest.txt
metricstest -test.v -test.run=^TestIteration2[AB]$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=.\ -server-port=8080 -file-storage-path=d:\temp\filest.txt
metricstest -test.v -test.run=^TestIteration3[AB]$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=.\ -server-port=8080 -file-storage-path=d:\temp\filest.txt
metricstest -test.v -test.run=^TestIteration4$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=.\ -server-port=8080 -file-storage-path=d:\temp\filest.txt
metricstest -test.v -test.run=^TestIteration5$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=.\ -server-port=8080 -file-storage-path=d:\temp\filest.txt
metricstest -test.v -test.run=^TestIteration6$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=.\ -server-port=8080 -file-storage-path=d:\temp\filest.txt
metricstest -test.v -test.run=^TestIteration7$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=.\ -server-port=8080 -file-storage-path=d:\temp\filest.txt
metricstest -test.v -test.run=^TestIteration8$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=.\ -server-port=8080 -file-storage-path=d:\temp\filest.txt
metricstest -test.v -test.run=^TestIteration9$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=.\ -server-port=8080 -file-storage-path=d:\temp\filest.txt
metricstest -test.v -test.run=^TestIteration10[AB]$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=.\ -server-port=8080 -file-storage-path=d:\temp\filest.txt -database-dsn=postgres://orchestra:orchestra12qwaszx@pg-ext.os.serk.lan:5103/orchestra?sslmode=disable
metricstest -test.v -test.run=^TestIteration11$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=.\ -server-port=8080 -file-storage-path=d:\temp\filest.txt -database-dsn=postgres://orchestra:orchestra12qwaszx@pg-ext.os.serk.lan:5103/orchestra?sslmode=disable
metricstest -test.v -test.run=^TestIteration12$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=.\ -server-port=8080 -file-storage-path=d:\temp\filest.txt -database-dsn=postgres://orchestra:orchestra12qwaszx@pg-ext.os.serk.lan:5103/orchestra?sslmode=disable
metricstest -test.v -test.run=^TestIteration13$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=.\ -server-port=8080 -file-storage-path=d:\temp\filest.txt -database-dsn=postgres://orchestra:orchestra12qwaszx@pg-ext.os.serk.lan:5103/orchestra?sslmode=disable
metricstest -test.v -test.run=^TestIteration14$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=.\ -server-port=8080 -file-storage-path=d:\temp\filest.txt -database-dsn=postgres://orchestra:orchestra12qwaszx@pg-ext.os.serk.lan:5103/orchestra?sslmode=disable -key=12qwaszx