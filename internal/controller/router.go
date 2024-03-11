package controller

func (s *Server) SetupRoutes() {
	v1 := s.App.Group("")
	v1.POST("/generate/:uuid", s.handler.Token.Generate)
	v1.POST("/refresh", s.handler.Token.Refresh)
}
