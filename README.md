# Completa AI Backend

Backend em Go responsável pelos endpoints do app. Agora é possível executar toda a stack (API + PostgreSQL) usando Docker e Docker Compose.

## Pré-requisitos

- Docker e Docker Compose instalados.
- Arquivo `.env` (opcional) para sobrescrever variáveis de ambiente. Copie de `.env.example` e ajuste os valores reais de Supabase.

## Subindo o ambiente

```bash
cd completa_ai_backend
docker compose up --build
```

O primeiro `up` vai:

1. Criar a imagem da API usando o `Dockerfile`.
2. Inicializar o PostgreSQL (porta host padrão `54322`) e aplicar os arquivos em `migrations/`.
3. Iniciar o container `backend` escutando em `http://localhost:8080`.

As variáveis usadas pela API agora ficam centralizadas no arquivo `.env` (lido automaticamente pelo Docker Compose). Para alterar algum valor basta editar esse arquivo ou exportar a variável no shell antes de rodar o `compose`. Principais chaves:

- `PORT` / `DB_PORT` — portas expostas para a API e para o banco.
- `DATABASE_URL` — string de conexão que a API usa (aponta para o serviço `db`).
- `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD` — credenciais do container Postgres.
- `SUPABASE_URL`, `SUPABASE_ANON_KEY`, `SUPABASE_JWT_SECRET` — credenciais usadas no fluxo de autenticação.

## Comandos úteis

- `docker compose logs -f backend` — acompanha os logs da API.
- `docker compose exec db psql -U postgres` — abre o psql dentro do banco.
- `docker compose down` — derruba os serviços mantendo os dados.
- `docker compose down -v` — derruba tudo e apaga o volume do banco (útil para reset).

## Saúde do serviço

- `GET http://localhost:8080/health` retorna `{"status":"ok"}` quando a API estiver pronta.

## Testando endpoints protegidos

1. Crie um usuário de teste no banco (o UUID precisa casar com o que você vai colocar no token):

   ```bash
   docker compose exec db psql -U postgres -c "INSERT INTO auth.users (id, email) VALUES ('00000000-0000-0000-0000-000000000000', 'local@test.com') ON CONFLICT (id) DO NOTHING;"
   ```

2. Gere um JWT HS256 usando o mesmo segredo configurado em `SUPABASE_JWT_SECRET` (exemplo com Python):

   ```bash
   python3 - <<'PY'
   import json, base64, hmac, hashlib
   user_id = "00000000-0000-0000-0000-000000000000"
   header = {"alg": "HS256", "typ": "JWT"}
   payload = {"sub": user_id}
   key = b"local-jwt-secret"
   def b64(obj): return base64.urlsafe_b64encode(json.dumps(obj, separators=(',',':')).encode()).rstrip(b'=')
   signing_input = b'.'.join([b64(header), b64(payload)])
   signature = base64.urlsafe_b64encode(hmac.new(key, signing_input, hashlib.sha256).digest()).rstrip(b'=')
   print((signing_input + b'.' + signature).decode())
   PY
   ```

3. Use esse token para chamar os endpoints:

   ```bash
   curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/collection
   curl -H "Authorization: Bearer <TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"stickers":{"BRA_01":2},"client_time":"2026-01-01T12:00:00Z"}' \
     http://localhost:8080/collection/sync
   ```

Com isso o backend fica pronto para ser executado localmente ou em produção apenas com Docker.
