#!/usr/bin/env python3
# /// script
# dependencies = []
# ///
"""
HTTP+JSON server exposing a Repository adapter.
Routes mirror the messages in ports/item/repository.avpr.

Usage:
  uv run adapters/python/server.py [--addr host:port] [--adapter memory]
"""

import argparse
import json
import sys
from http.server import BaseHTTPRequestHandler, HTTPServer
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))

from item.memory import MemoryRepository
from item.model import Item
from item.repository import NotFound


def make_handler(repo):
    class Handler(BaseHTTPRequestHandler):
        def log_message(self, format, *args):  # noqa: A002
            pass  # suppress default per-request logging

        def respond(self, status: int, body) -> None:
            data = json.dumps(body).encode()
            self.send_response(status)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(data)

        def read_body(self) -> dict:
            length = int(self.headers.get("Content-Length", 0))
            return json.loads(self.rfile.read(length)) if length else {}

        def do_POST(self):
            body = self.read_body()
            try:
                if self.path == "/item/get":
                    it = repo.get(body["id"])
                    self.respond(200, {"id": it.id, "name": it.name})
                elif self.path == "/item/list":
                    items = repo.list()
                    self.respond(200, [{"id": i.id, "name": i.name} for i in items])
                elif self.path == "/item/save":
                    repo.save(Item(id=body["id"], name=body["name"]))
                    self.respond(200, None)
                elif self.path == "/item/delete":
                    repo.delete(body["id"])
                    self.respond(200, None)
                else:
                    self.respond(404, {"error": "unknown route"})
            except NotFound as e:
                self.respond(404, {"error": str(e), "id": e.id})

    return Handler


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--addr", default="localhost:8080")
    parser.add_argument("--adapter", default="memory", choices=["memory"])
    args = parser.parse_args()

    host, port = args.addr.rsplit(":", 1)
    repo = MemoryRepository()

    print(f"ItemRepository server ({args.adapter} adapter) → http://{args.addr}")
    HTTPServer((host, int(port)), make_handler(repo)).serve_forever()


if __name__ == "__main__":
    main()
