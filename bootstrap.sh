# Create Service
curl -X POST http://localhost:8001/services/ -H 'Content-Type: application/json' -d '{ "name": "auth-v1", "url": "http://auth.local:9000" }'

# Add Auth Service routes
curl -X POST http://localhost:8001/services/auth-v1/routes/ -H 'Content-Type: application/json' -d '{ "paths": ["/auth"] }'

# Confirm Setup complete
curl http://localhost:8000/auth/healthcheck

# Enable the JWT Plugin
curl -X POST http://localhost:8001/services/auth-v1/plugins -H 'Content-Type: application/json' -d '{ "name": "jwt" }'
