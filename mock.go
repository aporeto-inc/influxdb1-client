package client

import (
	"context"
	"sync"
	"testing"
	"time"
)

// A TestClient is the interface of mockable test influxdb Client.
type TestClient interface {
	Client
	MockPing(t *testing.T, impl func(timeout time.Duration) (time.Duration, string, error))
	MockWrite(t *testing.T, impl func(bp BatchPoints) error)
	MockQuery(t *testing.T, impl func(q Query) (*Response, error))
	MockQueryWithContext(t *testing.T, impl func(ctx context.Context, q Query) (*Response, error))
	MockQueryAsChunk(t *testing.T, impl func(q Query) (*ChunkedResponse, error))
	MockQueryAsChunkWithContext(t *testing.T, impl func(ctx context.Context, q Query) (*ChunkedResponse, error))
	MockClose(t *testing.T, impl func() error)
}

type mockedClient struct {
	pingMock                    func(timeout time.Duration) (time.Duration, string, error)
	writeMock                   func(bp BatchPoints) error
	queryMock                   func(q Query) (*Response, error)
	queryWithContextMock        func(ctx context.Context, q Query) (*Response, error)
	queryAsChunkMock            func(q Query) (*ChunkedResponse, error)
	queryAsChunkWithContextMock func(ctx context.Context, q Query) (*ChunkedResponse, error)
	closeMock                   func() error
}

type testClient struct {
	mocks       map[*testing.T]*mockedClient
	lock        *sync.Mutex
	currentTest *testing.T
}

// NewTestClient returns a new TestTokenManager.
func NewTestClient() TestClient {
	return &testClient{
		lock:  &sync.Mutex{},
		mocks: map[*testing.T]*mockedClient{},
	}
}

func (m *testClient) MockPing(t *testing.T, impl func(timeout time.Duration) (time.Duration, string, error)) {

	m.lock.Lock()
	defer m.lock.Unlock()

	m.currentMocks(t).pingMock = impl
}

func (m *testClient) MockWrite(t *testing.T, impl func(bp BatchPoints) error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	m.currentMocks(t).writeMock = impl
}

func (m *testClient) MockQuery(t *testing.T, impl func(q Query) (*Response, error)) {

	m.lock.Lock()
	defer m.lock.Unlock()

	m.currentMocks(t).queryMock = impl
}

func (m *testClient) MockQueryWithContext(t *testing.T, impl func(ctx context.Context, q Query) (*Response, error)) {

	m.lock.Lock()
	defer m.lock.Unlock()

	m.currentMocks(t).queryWithContextMock = impl
}

func (m *testClient) MockQueryAsChunk(t *testing.T, impl func(q Query) (*ChunkedResponse, error)) {

	m.lock.Lock()
	defer m.lock.Unlock()

	m.currentMocks(t).queryAsChunkMock = impl
}

func (m *testClient) MockQueryAsChunkWithContext(t *testing.T, impl func(ctx context.Context, q Query) (*ChunkedResponse, error)) {

	m.lock.Lock()
	defer m.lock.Unlock()

	m.currentMocks(t).queryAsChunkWithContextMock = impl
}

func (m *testClient) MockClose(t *testing.T, impl func() error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	m.currentMocks(t).closeMock = impl
}

func (m *testClient) Ping(timeout time.Duration) (time.Duration, string, error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	if mock := m.currentMocks(m.currentTest); mock != nil && mock.pingMock != nil {
		return mock.pingMock(timeout)
	}

	return 0, "", nil
}

func (m *testClient) Write(bp BatchPoints) error {

	m.lock.Lock()
	defer m.lock.Unlock()

	if mock := m.currentMocks(m.currentTest); mock != nil && mock.writeMock != nil {
		return mock.writeMock(bp)
	}

	return nil
}

func (m *testClient) Query(q Query) (*Response, error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	if mock := m.currentMocks(m.currentTest); mock != nil && mock.queryMock != nil {
		return mock.queryMock(q)
	}

	return nil, nil
}

func (m *testClient) QueryWithContext(ctx context.Context, q Query) (*Response, error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	if mock := m.currentMocks(m.currentTest); mock != nil && mock.queryWithContextMock != nil {
		return mock.queryWithContextMock(ctx, q)
	}

	return nil, nil
}

func (m *testClient) QueryAsChunk(q Query) (*ChunkedResponse, error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	if mock := m.currentMocks(m.currentTest); mock != nil && mock.queryAsChunkMock != nil {
		return mock.queryAsChunkMock(q)
	}

	return nil, nil
}

func (m *testClient) QueryAsChunkWithContext(ctx context.Context, q Query) (*ChunkedResponse, error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	if mock := m.currentMocks(m.currentTest); mock != nil && mock.queryAsChunkWithContextMock != nil {
		return mock.queryAsChunkWithContextMock(ctx, q)
	}

	return nil, nil
}

func (m *testClient) Close() error {

	m.lock.Lock()
	defer m.lock.Unlock()

	if mock := m.currentMocks(m.currentTest); mock != nil && mock.closeMock != nil {
		return mock.closeMock()
	}

	return nil
}

func (m *testClient) currentMocks(t *testing.T) *mockedClient {

	mocks := m.mocks[t]

	if mocks == nil {
		mocks = &mockedClient{}
		m.mocks[t] = mocks
	}

	m.currentTest = t

	return mocks
}
