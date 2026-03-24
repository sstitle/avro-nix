from .model import Item
from .repository import NotFound


class MemoryRepository:
    """In-process adapter backed by a dict. Suitable for tests and local DI."""

    def __init__(self) -> None:
        self._store: dict[str, Item] = {}

    def get(self, id: str) -> Item:
        if id not in self._store:
            raise NotFound(id)
        return self._store[id]

    def list(self) -> list[Item]:
        return list(self._store.values())

    def save(self, item: Item) -> None:
        self._store[item.id] = item

    def delete(self, id: str) -> None:
        self._store.pop(id, None)
