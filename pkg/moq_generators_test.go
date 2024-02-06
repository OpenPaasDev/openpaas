package pkg

//go:generate moq -pkg provider -stub -out ./provider/moq_ansible_client_test.go ./ansible Client:MockAnsibleClient
