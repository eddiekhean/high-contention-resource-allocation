import { Router } from "express";
import { MazeGenerator } from "../core/MazeGenerator";

const router = Router();
router.get("/health", (_, res) => {
    res.json({ status: "ok" });
});

router.post("/generate", (req, res) => {
    const { rows, cols, loop_ratio = 0, seed } = req.body;
    const gen = new MazeGenerator(rows, cols, seed);
    res.json(gen.generate(loop_ratio));
});

export default router;
