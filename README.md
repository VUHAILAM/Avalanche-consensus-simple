# Avalanche-consensus-simple

Simple implementation of Avalanche Snowball consensus
## How to run:
```shell
docker compose up
```
## API:
- GET    /api/v1/health: Get status of all peer in discovery
- POST   /api/v1/create: Create a node

## ToDo:
- Add more tests
- Implement chain using DAGs
- Implement hash data
- Add transactions

## References
- https://docs.avax.network/overview/getting-started/avalanche-consensus
- https://github.com/ava-labs/mastering-avalanche/blob/main/chapter_09.md
- https://github.com/Jeiwan/blockchain_go