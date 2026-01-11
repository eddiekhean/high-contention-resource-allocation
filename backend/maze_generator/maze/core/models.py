from typing import List, Optional
from pydantic import BaseModel, Field, field_validator

class MazeGenerationRequest(BaseModel):
    rows: int = Field(..., ge=2, description="Number of rows")
    cols: int = Field(..., ge=2, description="Number of columns")
    loop_ratio: float = Field(0.0, ge=0.0, le=0.5, description="Probability of removing a wall to create a loop (0.0 to 0.5)")
    seed: Optional[int] = Field(None, description="Random seed for deterministic generation")

class Walls(BaseModel):
    top: bool
    right: bool
    bottom: bool
    left: bool

class Cell(BaseModel):
    x: int
    y: int
    walls: Walls

class Point(BaseModel):
    x: int
    y: int

class MazeResponse(BaseModel):
    rows: int
    cols: int
    start: Point
    end: Point
    cells: List[Cell]
