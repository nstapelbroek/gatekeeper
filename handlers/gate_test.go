package handlers

//func TestGateHandler_PostOpen(t *testing.T) {
//	adapterInstance := dummy.NewDummyAdapter()
//	gateHandler := NewGateHandler(adapterInstance, 2, "TCP:22")
//	app := setupGin(gateHandler.PostOpen)
//	request, _ := http.NewRequest("GET", "/", nil)
//	response := performRequest(app, request)
//
//	assert.Equal(t, http.StatusNotFound, response.Code)
//	assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("content-type"))
//	assert.Contains(t, string(response.Body.Bytes()), "Page not found")
//}

//
//func TestGateOpenHandler(t *testing.T) {
//	req := prepareRequest(t)
//	adapterInstance := dummy.NewDummyAdapter()
//	gateHandler := NewGateHandler(adapterInstance, 2, "TCP:22")
//
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(gateHandler.PostOpen)
//	handler.ServeHTTP(rr, req)
//
//	if status := rr.Code; status != http.StatusCreated {
//		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
//	}
//
//	expected := `127.0.0.1 has been whitelisted for 2 seconds`
//	if rr.Body.String() != expected {
//		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
//	}
//}
