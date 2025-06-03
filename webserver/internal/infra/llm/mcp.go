package llm

type MCPServer struct {
	service *ObservationService
	port    string
}

func NewMCPServer(service *ObservationService, port string) *MCPServer {
	return &MCPServer{
		service: service,
		port:    port,
	}
}

//
//func (mcp *MCPServer) listTools(c *gin.Context) {
//	tools := mcp.service.tools
//
//	response := map[string]interface{}{
//		"tools": tools,
//	}
//
//	c.JSON(http.StatusOK, response)
//}
//
//func (mcp *MCPServer) callTool(c *gin.Context) {
//	var request struct {
//		Name      string                 `json:"name"`
//		Arguments map[string]interface{} `json:"arguments"`
//	}
//
//	if err := c.ShouldBindJSON(&request); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	result, err := mcp.service.executeTool(c.Request.Context(), request.Name, request.Arguments)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"content": []map[string]interface{}{
//			{
//				"type": "text",
//				"text": result,
//			},
//		},
//	})
//}
//
//func (mcp *MCPServer) getResources(c *gin.Context) {
//	resources := []map[string]interface{}{
//		{
//			"uri":         "cluster://current/health",
//			"name":        "Cluster Health",
//			"description": "Current cluster health status and metrics",
//			"mimeType":    "application/json",
//		},
//		{
//			"uri":         "cluster://current/issues",
//			"name":        "Active Issues",
//			"description": "Currently active cluster issues and alerts",
//			"mimeType":    "application/json",
//		},
//		{
//			"uri":         "cluster://current/workflows",
//			"name":        "Active Workflows",
//			"description": "Currently running remediation workflows",
//			"mimeType":    "application/json",
//		},
//	}
//
//	c.JSON(http.StatusOK, gin.H{"resources": resources})
//}
//
//func (mcp *MCPServer) getPrompts(c *gin.Context) {
//	prompts := []map[string]interface{}{
//		{
//			"name":        "troubleshoot_issue",
//			"description": "Generate troubleshooting steps for a specific cluster issue",
//			"arguments": []map[string]interface{}{
//				{
//					"name":        "issue_description",
//					"description": "Description of the cluster issue",
//					"required":    true,
//				},
//				{
//					"name":        "severity",
//					"description": "Issue severity level",
//					"required":    false,
//				},
//			},
//		},
//		{
//			"name":        "optimize_resources",
//			"description": "Analyze and suggest resource optimization strategies",
//			"arguments": []map[string]interface{}{
//				{
//					"name":        "resource_type",
//					"description": "Type of resource to optimize (cpu, memory, storage)",
//					"required":    true,
//				},
//				{
//					"name":        "namespace",
//					"description": "Kubernetes namespace to focus on",
//					"required":    false,
//				},
//			},
//		},
//	}
//
//	c.JSON(http.StatusOK, gin.H{"prompts": prompts})
//}
//
//func (mcp *MCPServer) Start() error {
//	router := gin.Default()
//
//	// MCP-specific endpoints
//	router.POST("/mcp/tools/list", mcp.listTools)
//	router.POST("/mcp/tools/call", mcp.callTool)
//	router.GET("/mcp/resources", mcp.getResources)
//	router.POST("/mcp/prompts", mcp.getPrompts)
//
//	log.Printf("Starting MCP server on port %s", mcp.port)
//	return router.Run(":" + mcp.port)
//}
