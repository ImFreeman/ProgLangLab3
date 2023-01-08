# ProgLangLab2
## Start server
go run main.go -numOfNodes=5

## Send message (with Windows PowerShell)
curl -Method POST -Uri http://localhost:3000 -Body '{"data":"someData", "recv":3, "ttl":3}' -UseBasicParsing
