# Todo API

Une API REST simple pour gérer des tâches (todos), conçue pour être déployée sur Kubernetes.

GO Code: IA
Docerfile: IA
github_action : kaez
k8s: kaez

## Fonctionnalités

- CRUD complet pour les todos
- Base de données SQLite
- Health checks (liveness & readiness)
- Métriques Prometheus
- Logs structurés
- Configuration via variables d'environnement

## API Endpoints

### Todos

- `GET /api/todos` - Récupérer toutes les tâches
- `GET /api/todos/{id}` - Récupérer une tâche par ID
- `POST /api/todos` - Créer une nouvelle tâche
- `PUT /api/todos/{id}` - Mettre à jour une tâche
- `DELETE /api/todos/{id}` - Supprimer une tâche

### Health Checks

- `GET /health/live` - Liveness probe (pour k8s)
- `GET /health/ready` - Readiness probe (pour k8s)

### Métriques

- `GET /metrics` - Métriques au format Prometheus

## Exemples d'utilisation

### Créer un todo

```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Apprendre Kubernetes",
    "description": "Déployer mon application Go sur k8s"
  }'
```

### Lister tous les todos

```bash
curl http://localhost:8080/api/todos
```

### Récupérer un todo spécifique

```bash
curl http://localhost:8080/api/todos/1
```

### Mettre à jour un todo

```bash
curl -X PUT http://localhost:8080/api/todos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "completed": true
  }'
```

### Supprimer un todo

```bash
curl -X DELETE http://localhost:8080/api/todos/1
```

## Exécution locale

### Avec Go

```bash
go run cmd/api/main.go
```

### Avec Docker

```bash
docker build -t todo-api .
docker run -p 8080:8080 todo-api
```

### Avec volume pour la persistance

```bash
docker run -p 8080:8080 -v $(pwd)/data:/data todo-api
```

## Variables d'environnement

- `PORT` - Port d'écoute (défaut: 8080)
- `DB_PATH` - Chemin vers la base de données SQLite (défaut: ./todos.db)

## Tests des health checks

```bash
# Liveness
curl http://localhost:8080/health/live

# Readiness
curl http://localhost:8080/health/ready
```

## Métriques Prometheus

```bash
curl http://localhost:8080/metrics
```

Métriques disponibles:

- `http_requests_total` - Nombre total de requêtes HTTP
- `todos_total` - Nombre total de todos
- `go_memstats_alloc_bytes` - Mémoire allouée
- `go_goroutines` - Nombre de goroutines actives

## Déploiement sur Kubernetes

Ce projet est prêt pour Kubernetes avec:

- Health checks pour liveness et readiness probes
- Métriques Prometheus pour le monitoring
- Configuration via variables d'environnement
- Support des volumes pour la persistance
- Shutdown gracieux

Exemple de configuration k8s minimale:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: todo-api
spec:
  replicas: 3
  template:
    spec:
      containers:
        - name: todo-api
          image: todo-api:latest
          ports:
            - containerPort: 8080
          env:
            - name: PORT
              value: "8080"
          livenessProbe:
            httpGet:
              path: /health/live
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health/ready
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
```

## Structure du projet

```
.
├── cmd/
│   └── api/
│       └── main.go              # Point d'entrée de l'application
├── internal/
│   ├── handlers/
│   │   ├── health_handler.go   # Health checks
│   │   ├── metrics_handler.go  # Métriques Prometheus
│   │   └── todo_handler.go     # Handlers CRUD todos
│   ├── middleware/
│   │   └── logging.go           # Middleware de logging
│   ├── models/
│   │   └── todo.go              # Modèles de données
│   └── repository/
│       └── todo_repository.go   # Accès à la base de données
├── Dockerfile
├── go.mod
└── README.md
```

## Technologies utilisées

- Go 1.23
- SQLite3
- HTTP standard library
- Docker multi-stage builds
