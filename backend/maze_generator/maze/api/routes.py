from fastapi import APIRouter, HTTPException
from ..core.models import MazeGenerationRequest, MazeResponse
from ..core.generator import MazeGenerator

router = APIRouter()

@router.post("/generate", response_model=MazeResponse)
async def generate_maze(request: MazeGenerationRequest):
    """
    Generate a maze based on the provided specifications.
    Algorithm: Recursive Backtracker (DFS)
    """
    try:
        generator = MazeGenerator(
            rows=request.rows,
            cols=request.cols,
            seed=request.seed
        )
        maze = generator.generate(loop_ratio=request.loop_ratio)
        return maze
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
