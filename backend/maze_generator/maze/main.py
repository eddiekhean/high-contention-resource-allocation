import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from .api.routes import router as maze_router

app = FastAPI(
    title="Maze Generator Service",
    description="Backend service for generating perfect and braided mazes.",
    version="1.0.0"
)

# Configure CORS to allow requests from any origin (for dev/demo purposes)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(maze_router, prefix="/maze", tags=["maze"])

if __name__ == "__main__":
    uvicorn.run("maze.main:app", host="0.0.0.0", port=8000, reload=True)
