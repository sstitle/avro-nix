"""RPC adapter: implements Repository over HTTP+JSON.

Delegates all calls to a remote server, making local↔remote swap transparent
to the application core and service layer. The wire protocol mirrors the
messages in ports/item/repository.avpr.
"""

import json
import urllib.error
import urllib.request

from .model import Item
from .repository import NotFound


class RpcClient:
    """Adapter that satisfies Repository by calling a remote ItemRepository server."""

    def __init__(self, addr: str) -> None:
        self._base = f"http://{addr}"

    def _post(self, path: str, body: dict) -> dict | list | None:
        data = json.dumps(body).encode()
        req = urllib.request.Request(
            self._base + path,
            data=data,
            headers={"Content-Type": "application/json"},
        )
        try:
            with urllib.request.urlopen(req) as resp:
                return json.loads(resp.read())
        except urllib.error.HTTPError as e:
            if e.code == 404:
                payload = json.loads(e.read())
                raise NotFound(payload.get("id", "?"))
            raise

    def get(self, id: str) -> Item:
        data = self._post("/item/get", {"id": id})
        return Item(**data)

    def list(self) -> list[Item]:
        data = self._post("/item/list", {})
        return [Item(**d) for d in (data or [])]

    def save(self, item: Item) -> None:
        self._post("/item/save", {"id": item.id, "name": item.name})

    def delete(self, id: str) -> None:
        self._post("/item/delete", {"id": id})
