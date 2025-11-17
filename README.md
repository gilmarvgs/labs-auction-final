# labs-auction-goexpert

Projeto de exemplo para gerenciamento de leilões em Go.

Principais pontos implementados
- Cálculo do tempo do leilão via variável de ambiente `AUCTION_INTERVAL` (função `getAuctionInterval`).
- Goroutine em `CreateAuction` que fecha (atualiza status para `Completed`) o leilão quando o tempo expira.
- Teste unitário que valida o fechamento automático: `internal/infra/database/auction/create_auction_test.go`.

Como executar localmente (PowerShell)

- Compilar o projeto:
```powershell
go build ./...
```

- Rodar todos os testes:
```powershell
go test ./...
```

- Rodar apenas o teste do fechamento do leilão:
```powershell
go test ./internal/infra/database/auction -run TestCreateAuction_automaticallyClosesAfterInterval -v
```

Variáveis de ambiente úteis
- `AUCTION_INTERVAL` — duração do leilão (ex.: `20s`, `1m`). Se não setada, o sistema usa valores padrão definidos nos repositórios.

Como publicar suas mudanças (PowerShell)

1) Adicionar e commitar o `README.md` localmente:
```powershell
git add README.md
git commit -m "docs: add README with run/test instructions"
```

2) Enviar para o remote `origin` (branch `master`):
```powershell
git push origin master
```

Observações
- Se você estiver em outro branch, substitua `master` pelo nome do branch atual.
- Se preferir usar SSH para o remote, configure o remote para `git@github.com:gilmarvgs/labs-auction-final.git` antes do push.

Se quiser, eu posso também criar um `README` mais detalhado ou um arquivo `CONTRIBUTING.md`.
