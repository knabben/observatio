package llm

import "time"

type Agent struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Status       string    `json:"status"`
	Activity     string    `json:"activity"`
	LastUpdate   time.Time `json:"last_update"`
	Capabilities []string  `json:"capabilities"`
}

func (s *ObservationService) initializeAgents() {
	agents := []*Agent{
		{
			ID:           "cluster-agent",
			Name:         "Cluster Agent",
			Type:         "analysis",
			Status:       "active",
			Activity:     "Monitoring cluster health and analyzing issues",
			Capabilities: []string{"cluster-analysis", "resource-monitoring", "pattern-detection"},
			LastUpdate:   time.Now(),
		},
	}
	for _, agent := range agents {
		s.agents[agent.ID] = agent
	}
}
