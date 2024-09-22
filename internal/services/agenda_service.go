package services

type AgendaService interface {
	StartBroker() error
	// CreateAppointment cria um horario para um paciente com um dado medico
	// Retorna um erro caso ocorra algum problema durante a criacao do horario
	// Retorna nil caso a criacao ocorra com sucesso
	AgendaTask() error
}
