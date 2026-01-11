from collections import deque
from typing import List, Tuple, Set, Dict

def get_neighbors(r: int, c: int, rows: int, cols: int) -> List[Tuple[int, int]]:
    """
    Returns valid neighbors for a given cell (r, c) within the grid boundaries.
    Used for initial graph traversal where walls don't strict movement yet (if needed),
    but primarily we need to know neighbors *respecting walls* for the BFS.
    """
    moves = [(-1, 0), (1, 0), (0, -1), (0, 1)] # Top, Bottom, Left, Right
    neighbors = []
    for dr, dc in moves:
        nr, nc, = r + dr, c + dc
        if 0 <= nr < rows and 0 <= nc < cols:
            neighbors.append((nr, nc))
    return neighbors

def bfs_dist_map(start_node: Tuple[int, int], grid: List[List[Dict[str, bool]]], rows: int, cols: int) -> Dict[Tuple[int, int], int]:
    """
    Performs BFS from start_node to find distances to all other reachable nodes.
    grid is expected to be a 2D list of cells, where each cell is a dict/object having 'walls'.
    Internal representation: grid[r][c]['top'] is True if there is a wall.
    Returns a dict of {node: distance}.
    """
    distances = {start_node: 0}
    queue = deque([start_node])
    
    while queue:
        r, c = queue.popleft()
        dist = distances[(r, c)]
        
        # Check all 4 directions. Movement is allowed only if there is NO wall.
        # Top
        if r > 0 and not grid[r][c]['top']:
            nr, nc = r - 1, c
            if (nr, nc) not in distances:
                distances[(nr, nc)] = dist + 1
                queue.append((nr, nc))
        # Bottom
        if r < rows - 1 and not grid[r][c]['bottom']:
            nr, nc = r + 1, c
            if (nr, nc) not in distances:
                distances[(nr, nc)] = dist + 1
                queue.append((nr, nc))
        # Left
        if c > 0 and not grid[r][c]['left']:
            nr, nc = r, c - 1
            if (nr, nc) not in distances:
                distances[(nr, nc)] = dist + 1
                queue.append((nr, nc))
        # Right
        if c < cols - 1 and not grid[r][c]['right']:
            nr, nc = r, c + 1
            if (nr, nc) not in distances:
                distances[(nr, nc)] = dist + 1
                queue.append((nr, nc))
                
    return distances

def find_longest_path_boundary_points(grid: List[List[Dict[str, bool]]], rows: int, cols: int) -> Tuple[Tuple[int, int], Tuple[int, int]]:
    """
    Finds two points on the boundary that maximize the shortest-path distance between them.
    Strategy:
    1. Identify all boundary cells.
    2. Pick an arbitrary boundary cell A.
    3. BFS from A to find the furthest boundary cell B.
    4. BFS from B to find the furthest boundary cell C.
    5. Return B and C. (Approximate diameter of the graph restricted to boundary nodes).
    
    Refinement for "Maximized shortest-path":
    Running BFS from ALL boundary nodes might be too expensive O(Perimeter * Cells).
    Double BFS is a good heuristic for tree-like mazes (perfect mazes).
    For braided mazes, it's also a decent heuristic.
    """
    
    boundary_cells = []
    for r in range(rows):
        for c in range(cols):
            if r == 0 or r == rows - 1 or c == 0 or c == cols - 1:
                boundary_cells.append((r, c))
    
    if not boundary_cells:
        return (0, 0), (0, 0) # Should not happen for rows, cols >= 2

    # Heuristic: Start from a corner (0,0) - likely to be on boundary.
    start_node = (0, 0) 
    
    # 1. BFS from start_node to find distances to all cells
    dists_from_start = bfs_dist_map(start_node, grid, rows, cols)
    
    # 2. Find the furthest boundary cell from start_node
    furthest_cell_1 = start_node
    max_dist = 0
    
    # Iterate only over accessible boundary cells (to ensure connectivity)
    for cell in boundary_cells:
        if cell in dists_from_start:
            if dists_from_start[cell] > max_dist:
                max_dist = dists_from_start[cell]
                furthest_cell_1 = cell
    
    # 3. BFS from furthest_cell_1 to find the absolute furthest boundary cell (Diameter end)
    dists_from_f1 = bfs_dist_map(furthest_cell_1, grid, rows, cols)
    
    furthest_cell_2 = furthest_cell_1
    max_dist_2 = 0
    
    for cell in boundary_cells:
        if cell in dists_from_f1:
            if dists_from_f1[cell] > max_dist_2:
                max_dist_2 = dists_from_f1[cell]
                furthest_cell_2 = cell
                
    return furthest_cell_1, furthest_cell_2
