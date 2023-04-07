# mcsdeploy

This util helps to configure local deploy considering

- configs
  - service1-v1
    - config.yaml 
  - service2-v1
      - config.yaml
  - service3-v1
      - config.yaml
  - ...
- envoy
  - Dockerfile
  - envoy.yaml
- docker-compose.yml
- Makefile

You have to set file `services.yaml` like
```yaml
services:
- name: service1
  version: v1
- name: service2
  version: v1
- name: service3
  version: v1
- ...
```

