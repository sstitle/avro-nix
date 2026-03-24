"""
Port: Repository interface for Item entities.

Derived from ports/item/repository.avpr. Any class implementing this Protocol
is a valid adapter — memory, file-backed, or RPC client. The application core
and service layer import only from this module; they never reference a concrete
adapter.
"""

from typing import Protocol, runtime_checkable

from .model import Item


class NotFound(Exception):
    """Raised when no Item exists for the requested ID.

    Derived from com.example.errors.NotFound in ports/item/repository.avpr.
    """

    def __init__(self, id: str) -> None:
        super().__init__(f"item not found: {id}")
        self.id = id


@runtime_checkable
class Repository(Protocol):
    """The port. Application core and service layer depend only on this."""

    def get(self, id: str) -> Item: ...
    def list(self) -> list[Item]: ...
    def save(self, item: Item) -> None: ...
    def delete(self, id: str) -> None: ...
