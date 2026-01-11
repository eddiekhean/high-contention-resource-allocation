import sys
import os
import random
from fastapi.testclient import TestClient

# sys.path hack removed, use -m flag to run


# Run this as a module from parent directory: python -m maze.verify

from .main import app
from .core.generator import MazeGenerator

client = TestClient(app)

def test_maze_generation_logic():
    print("Testing MazeGenerator logic...")
    rows, cols = 5, 5
    generator = MazeGenerator(rows=rows, cols=cols, seed=42)
    maze = generator.generate(loop_ratio=0.0)
    
    assert maze.rows == rows
    assert maze.cols == cols
    assert len(maze.cells) == rows * cols
    
    # Check start and end are on boundary
    assert maze.start.x == 0 or maze.start.x == cols - 1 or maze.start.y == 0 or maze.start.y == rows - 1
    assert maze.end.x == 0 or maze.end.x == cols - 1 or maze.end.y == 0 or maze.end.y == rows - 1
    
    print("Logic test passed!")

def test_determinism():
    print("Testing determinism with seed...")
    g1 = MazeGenerator(5, 5, seed=123)
    m1 = g1.generate()
    
    g2 = MazeGenerator(5, 5, seed=123)
    m2 = g2.generate()
    
    assert m1.dict() == m2.dict()
    print("Determinism test passed!")

def test_api_endpoint():
    print("Testing POST /maze/generate...")
    payload = {
        "rows": 4,
        "cols": 4,
        "loop_ratio": 0.2,
        "seed": 999
    }
    response = client.post("/maze/generate", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert data["rows"] == 4
    assert data["cols"] == 4
    assert "start" in data
    assert "end" in data
    assert len(data["cells"]) == 16
    print("API test passed!")

def test_invalid_input():
    print("Testing invalid input...")
    # rows < 2
    response = client.post("/maze/generate", json={"rows": 1, "cols": 5})
    assert response.status_code == 422
    
    # loop_ratio > 0.5
    response = client.post("/maze/generate", json={"rows": 5, "cols": 5, "loop_ratio": 0.6})
    assert response.status_code == 422
    print("Invalid input test passed!")

if __name__ == "__main__":
    try:
        test_maze_generation_logic()
        test_determinism()
        test_api_endpoint()
        test_invalid_input()
        print("\nAll tests passed successfully!")
    except Exception as e:
        print(f"\nTests failed: {e}")
        sys.exit(1)
