package mocks

// type MockLogger struct {
// 	mock.Mock
// }

// func (m *MockLogger) GetLogger() *zerolog.Logger {
// 	args := m.Called()
// 	return args.Get(0).(*zerolog.Logger)
// }

// func (m *MockLogger) WithFields(fields map[string]any) logger.LoggerInterface {
// 	args := m.Called(fields)
// 	return args.Get(0).(logger.LoggerInterface)
// }

// func (m *MockLogger) WithContext(requestID, method string) logger.LoggerInterface {
// 	args := m.Called(requestID, method)
// 	return args.Get(0).(logger.LoggerInterface)
// }

// func (m *MockLogger) AddFields(key, val string) logger.LoggerInterface {
// 	args := m.Called(key, val)
// 	return args.Get(0).(logger.LoggerInterface)
// }

// func (m *MockLogger) Info(msg string, fields map[string]any) {
// 	m.Called(msg, fields)
// }

// func (m *MockLogger) Error(err error, msg string, fields map[string]any) {
// 	m.Called(err, msg, fields)
// }

// func (m *MockLogger) Fatal(err error, msg string, fields map[string]any) {
// 	m.Called(err, msg, fields)
// }

// func (m *MockLogger) Warn(msg string, fields map[string]any, err error) {
// 	m.Called(msg, fields, err)
// }

// func (m *MockLogger) Debug(msg string, fields map[string]any) {
// 	m.Called(msg, fields)
// }

// func (m *MockLogger) WrapError(message string, err error) error {
// 	args := m.Called(message, err)
// 	if args.Error(0) != nil {
// 		return args.Error(0)
// 	}
// 	return fmt.Errorf("%s", message)
// }

// func (m *MockLogger) WrapDetailError(message string, err error, messDetail any) error {
// 	args := m.Called(message, err, messDetail)
// 	if args.Error(0) != nil {
// 		return args.Error(0)
// 	}
// 	return fmt.Errorf("%s", message)
// }
