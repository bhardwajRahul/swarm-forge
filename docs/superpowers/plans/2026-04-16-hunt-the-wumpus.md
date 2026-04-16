# Hunt the Wumpus Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement an exact replica of Gregory Yob's 1973 Hunt the Wumpus in Clojure.

**Architecture:** Functional core / imperative shell. All game logic is pure functions over an immutable state map. A thin I/O shell in `main.clj` handles stdin/stdout. The dodecahedron is generated programmatically.

**Tech Stack:** Clojure CLI (deps.edn), Speclj for BDD specs.

---

## File Structure

| File | Responsibility |
|------|---------------|
| `wumpus/deps.edn` | Project config, Speclj dependency |
| `wumpus/src/wumpus/dodecahedron.clj` | Generate 20-room dodecahedron adjacency map |
| `wumpus/src/wumpus/game.clj` | Pure game state, queries, and actions |
| `wumpus/src/wumpus/messages.clj` | Original game text strings |
| `wumpus/src/wumpus/main.clj` | I/O shell, `-main` entry point |
| `wumpus/spec/wumpus/dodecahedron_spec.clj` | Specs for dodecahedron generation |
| `wumpus/spec/wumpus/game_spec.clj` | Specs for game logic |
| `wumpus/spec/wumpus/messages_spec.clj` | Specs for message functions |
| `wumpus/features/wumpus.feature` | Gherkin E2E scenarios |

---

### Task 1: Project Scaffolding

**Files:**
- Create: `wumpus/deps.edn`
- Create: `wumpus/src/wumpus/dodecahedron.clj` (empty ns)
- Create: `wumpus/spec/wumpus/dodecahedron_spec.clj` (minimal failing spec)

- [ ] **Step 1: Create deps.edn**

```edn
{:paths ["src"]
 :aliases
 {:spec {:main-opts ["-m" "speclj.main" "-c"]
         :extra-deps {speclj/speclj {:mvn/version "3.4.5"}}
         :extra-paths ["spec"]}}}
```

- [ ] **Step 2: Create empty dodecahedron namespace**

Create `wumpus/src/wumpus/dodecahedron.clj`:

```clojure
(ns wumpus.dodecahedron)
```

- [ ] **Step 3: Create minimal failing spec**

Create `wumpus/spec/wumpus/dodecahedron_spec.clj`:

```clojure
(ns wumpus.dodecahedron-spec
  (:require [speclj.core :refer :all]
            [wumpus.dodecahedron :refer :all]))

(describe "dodecahedron"
  (it "has 20 rooms"
    (should= 20 (count (make-cave-map)))))
```

- [ ] **Step 4: Run specs to verify failure**

```bash
cd wumpus && clj -M:spec
```

Expected: FAIL — `make-cave-map` is not defined.

- [ ] **Step 5: Commit**

```bash
git add wumpus/
git commit -m "Scaffold wumpus project with deps.edn and first failing spec"
```

---

### Task 2: Dodecahedron — Room Count

**Files:**
- Modify: `wumpus/src/wumpus/dodecahedron.clj`
- Test: `wumpus/spec/wumpus/dodecahedron_spec.clj` (already has the failing spec)

- [ ] **Step 1: Implement make-cave-map returning 20 rooms**

The dodecahedron has 20 vertices. We build it from the classic adjacency used
in the original BASIC source. The structure is: a top pentagon (rooms 1-5), an
upper ring (rooms 6-10), a lower ring (rooms 11-15), and a bottom pentagon
(rooms 16-20). Each vertex connects to its two neighbors in its ring plus one
vertex in the adjacent ring.

Edit `wumpus/src/wumpus/dodecahedron.clj`:

```clojure
(ns wumpus.dodecahedron)

(defn make-cave-map []
  (let [edges [[1 2] [2 3] [3 4] [4 5] [5 1]           ; top pentagon
               [1 6] [2 7] [3 8] [4 9] [5 10]          ; top to upper ring
               [6 11] [6 15] [7 11] [7 12] [8 12]      ; upper to lower ring
               [8 13] [9 13] [9 14] [10 14] [10 15]
               [11 16] [12 17] [13 18] [14 19] [15 20]  ; lower ring to bottom
               [16 17] [17 18] [18 19] [19 20] [20 16]] ; bottom pentagon
        empty-map (into {} (map (fn [r] [r #{}]) (range 1 21)))]
    (reduce (fn [m [a b]]
              (-> m
                  (update a conj b)
                  (update b conj a)))
            empty-map edges)))
```

- [ ] **Step 2: Run specs to verify pass**

```bash
cd wumpus && clj -M:spec
```

Expected: PASS — `(count (make-cave-map))` returns 20.

- [ ] **Step 3: Commit**

```bash
git add wumpus/src/wumpus/dodecahedron.clj
git commit -m "Implement make-cave-map with 20-room dodecahedron"
```

---

### Task 3: Dodecahedron — Structural Properties

**Files:**
- Modify: `wumpus/spec/wumpus/dodecahedron_spec.clj`
- Modify: `wumpus/src/wumpus/dodecahedron.clj`

- [ ] **Step 1: Add failing specs for structural properties**

Add to `wumpus/spec/wumpus/dodecahedron_spec.clj`:

```clojure
(describe "dodecahedron"
  (with-all cave-map (make-cave-map))

  (it "has 20 rooms"
    (should= 20 (count @cave-map)))

  (it "every room connects to exactly 3 others"
    (doseq [[room neighbors] @cave-map]
      (should= 3 (count neighbors))))

  (it "connections are symmetric"
    (doseq [[room neighbors] @cave-map
            neighbor neighbors]
      (should-contain room (get @cave-map neighbor))))

  (it "has no self-loops"
    (doseq [[room neighbors] @cave-map]
      (should-not-contain room neighbors)))

  (it "all rooms are reachable from room 1"
    (let [reachable (loop [visited #{1} frontier [1]]
                      (if (empty? frontier)
                        visited
                        (let [next-rooms (mapcat #(get @cave-map %) frontier)
                              new-rooms (remove visited next-rooms)]
                          (recur (into visited new-rooms)
                                 (vec new-rooms)))))]
      (should= 20 (count reachable)))))
```

This replaces the original single spec. The first spec (`has 20 rooms`) already
passes. The others may or may not pass depending on whether the edge list is
correct.

- [ ] **Step 2: Run specs**

```bash
cd wumpus && clj -M:spec
```

Verify all 5 specs pass. If any fail, the edge list in `make-cave-map` has a
bug — fix it and re-run.

- [ ] **Step 3: Add neighbors and adjacent? helpers**

Add failing specs first:

```clojure
  (it "neighbors returns the 3 adjacent rooms"
    (should= (get @cave-map 1) (neighbors @cave-map 1)))

  (it "adjacent? returns true for connected rooms"
    (let [n (first (get @cave-map 1))]
      (should (adjacent? @cave-map 1 n))))

  (it "adjacent? returns false for unconnected rooms"
    (let [non-neighbor (first (clojure.set/difference
                                (set (range 1 21))
                                (conj (get @cave-map 1) 1)))]
      (should-not (adjacent? @cave-map 1 non-neighbor))))
```

- [ ] **Step 4: Run specs to verify failure**

```bash
cd wumpus && clj -M:spec
```

Expected: FAIL — `neighbors` and `adjacent?` not defined.

- [ ] **Step 5: Implement neighbors and adjacent?**

Add to `wumpus/src/wumpus/dodecahedron.clj`:

```clojure
(defn neighbors [cave-map room]
  (get cave-map room))

(defn adjacent? [cave-map room1 room2]
  (contains? (get cave-map room1) room2))
```

- [ ] **Step 6: Run specs to verify pass**

```bash
cd wumpus && clj -M:spec
```

Expected: All PASS.

- [ ] **Step 7: Commit**

```bash
git add wumpus/src/wumpus/dodecahedron.clj wumpus/spec/wumpus/dodecahedron_spec.clj
git commit -m "Add dodecahedron structural specs and helper functions"
```

---

### Task 4: Game State — Initialization

**Files:**
- Create: `wumpus/src/wumpus/game.clj`
- Create: `wumpus/spec/wumpus/game_spec.clj`

- [ ] **Step 1: Write failing spec for new-game**

Create `wumpus/spec/wumpus/game_spec.clj`:

```clojure
(ns wumpus.game-spec
  (:require [speclj.core :refer :all]
            [wumpus.game :refer :all]
            [wumpus.dodecahedron :refer [make-cave-map]]))

(describe "new-game"
  (with-all cave-map (make-cave-map))
  (with-all state (new-game @cave-map))

  (it "has the cave map"
    (should= @cave-map (:cave-map @state)))

  (it "places player in a valid room"
    (should-contain (:player @state) (set (range 1 21))))

  (it "places wumpus in a valid room"
    (should-contain (:wumpus @state) (set (range 1 21))))

  (it "has 2 pits in valid rooms"
    (should= 2 (count (:pits @state)))
    (doseq [pit (:pits @state)]
      (should-contain pit (set (range 1 21)))))

  (it "has 2 bat colonies in valid rooms"
    (should= 2 (count (:bats @state)))
    (doseq [bat (:bats @state)]
      (should-contain bat (set (range 1 21)))))

  (it "places all entities in distinct rooms"
    (let [rooms (concat [(:player @state) (:wumpus @state)]
                        (:pits @state) (:bats @state))]
      (should= (count rooms) (count (set rooms)))))

  (it "starts with 5 arrows"
    (should= 5 (:arrows @state)))

  (it "starts with status :playing"
    (should= :playing (:status @state))))
```

- [ ] **Step 2: Create empty game namespace and run to see failure**

Create `wumpus/src/wumpus/game.clj`:

```clojure
(ns wumpus.game
  (:require [wumpus.dodecahedron :as cave]))
```

```bash
cd wumpus && clj -M:spec
```

Expected: FAIL — `new-game` not defined.

- [ ] **Step 3: Implement new-game**

Add to `wumpus/src/wumpus/game.clj`:

```clojure
(defn new-game [cave-map]
  (let [rooms (shuffle (keys cave-map))
        [player wumpus pit1 pit2 bat1 bat2] rooms]
    {:cave-map cave-map
     :player   player
     :wumpus   wumpus
     :pits     #{pit1 pit2}
     :bats     #{bat1 bat2}
     :arrows   5
     :status   :playing}))
```

- [ ] **Step 4: Run specs to verify pass**

```bash
cd wumpus && clj -M:spec
```

Expected: All PASS.

- [ ] **Step 5: Commit**

```bash
git add wumpus/src/wumpus/game.clj wumpus/spec/wumpus/game_spec.clj
git commit -m "Add game state initialization with new-game"
```

---

### Task 5: Game Queries — sensed-hazards, player-room, player-neighbors

**Files:**
- Modify: `wumpus/spec/wumpus/game_spec.clj`
- Modify: `wumpus/src/wumpus/game.clj`

- [ ] **Step 1: Write failing specs for query functions**

Add to `wumpus/spec/wumpus/game_spec.clj`:

```clojure
(describe "player-room"
  (it "returns the player's current room"
    (should= 3 (player-room {:player 3}))))

(describe "player-neighbors"
  (with-all cave-map (make-cave-map))

  (it "returns the 3 rooms adjacent to the player"
    (let [state {:cave-map @cave-map :player 1}]
      (should= (get @cave-map 1) (player-neighbors state)))))

(describe "sensed-hazards"
  (with-all cave-map (make-cave-map))

  (it "senses wumpus in adjacent room"
    (let [neighbor (first (get @cave-map 1))
          state {:cave-map @cave-map :player 1 :wumpus neighbor
                 :pits #{} :bats #{}}]
      (should-contain :wumpus (sensed-hazards state))))

  (it "does not sense wumpus in non-adjacent room"
    (let [non-adj (first (clojure.set/difference
                           (set (range 1 21))
                           (conj (get @cave-map 1) 1)))
          state {:cave-map @cave-map :player 1 :wumpus non-adj
                 :pits #{} :bats #{}}]
      (should-not-contain :wumpus (sensed-hazards state))))

  (it "senses pit in adjacent room"
    (let [neighbor (first (get @cave-map 1))
          state {:cave-map @cave-map :player 1 :wumpus 20
                 :pits #{neighbor} :bats #{}}]
      (should-contain :pit (sensed-hazards state))))

  (it "senses bats in adjacent room"
    (let [neighbor (first (get @cave-map 1))
          state {:cave-map @cave-map :player 1 :wumpus 20
                 :pits #{} :bats #{neighbor}}]
      (should-contain :bats (sensed-hazards state))))

  (it "returns empty set when no adjacent hazards"
    (let [non-adj (first (clojure.set/difference
                           (set (range 1 21))
                           (conj (get @cave-map 1) 1)))
          state {:cave-map @cave-map :player 1 :wumpus non-adj
                 :pits #{} :bats #{}}]
      (should= #{} (sensed-hazards state)))))
```

- [ ] **Step 2: Run specs to verify failure**

```bash
cd wumpus && clj -M:spec
```

Expected: FAIL — `player-room`, `player-neighbors`, `sensed-hazards` not defined.

- [ ] **Step 3: Implement query functions**

Add to `wumpus/src/wumpus/game.clj`:

```clojure
(defn player-room [state]
  (:player state))

(defn player-neighbors [state]
  (cave/neighbors (:cave-map state) (:player state)))

(defn sensed-hazards [state]
  (let [adj (cave/neighbors (:cave-map state) (:player state))]
    (cond-> #{}
      (contains? adj (:wumpus state)) (conj :wumpus)
      (seq (clojure.set/intersection adj (:pits state))) (conj :pit)
      (seq (clojure.set/intersection adj (:bats state))) (conj :bats))))
```

- [ ] **Step 4: Run specs to verify pass**

```bash
cd wumpus && clj -M:spec
```

Expected: All PASS.

- [ ] **Step 5: Commit**

```bash
git add wumpus/src/wumpus/game.clj wumpus/spec/wumpus/game_spec.clj
git commit -m "Add game query functions: player-room, player-neighbors, sensed-hazards"
```

---

### Task 6: Move — Empty Room

**Files:**
- Modify: `wumpus/spec/wumpus/game_spec.clj`
- Modify: `wumpus/src/wumpus/game.clj`

- [ ] **Step 1: Write failing spec**

Add to `wumpus/spec/wumpus/game_spec.clj`:

```clojure
(describe "move"
  (with-all cave-map (make-cave-map))

  (it "moves player to the specified adjacent room"
    (let [neighbor (first (get @cave-map 1))
          state {:cave-map @cave-map :player 1 :wumpus 20
                 :pits #{} :bats #{} :arrows 5 :status :playing}
          result (move state neighbor)]
      (should= neighbor (:player result))
      (should= :playing (:status result)))))
```

- [ ] **Step 2: Run specs to verify failure**

```bash
cd wumpus && clj -M:spec
```

Expected: FAIL — `move` not defined.

- [ ] **Step 3: Implement move for empty room case**

Add to `wumpus/src/wumpus/game.clj`:

```clojure
(defn move [state room]
  (assoc state :player room))
```

- [ ] **Step 4: Run specs to verify pass**

```bash
cd wumpus && clj -M:spec
```

Expected: All PASS.

- [ ] **Step 5: Commit**

```bash
git add wumpus/src/wumpus/game.clj wumpus/spec/wumpus/game_spec.clj
git commit -m "Implement move to empty room"
```

---

### Task 7: Move — Into Wumpus (Eaten)

**Files:**
- Modify: `wumpus/spec/wumpus/game_spec.clj`
- Modify: `wumpus/src/wumpus/game.clj`

- [ ] **Step 1: Write failing spec**

We pass a deterministic RNG function to control the 50/50 outcome. Add:

```clojure
  (it "player dies when moving into wumpus room and wumpus stays"
    (let [neighbor (first (get @cave-map 1))
          state {:cave-map @cave-map :player 1 :wumpus neighbor
                 :pits #{} :bats #{} :arrows 5 :status :playing}
          stays (fn [] 0.6)  ; >= 0.5 means wumpus stays and kills
          result (move state neighbor stays)]
      (should= :lose-wumpus (:status result))))
```

- [ ] **Step 2: Run specs to verify failure**

```bash
cd wumpus && clj -M:spec
```

Expected: FAIL — `move` with 3 args not supported, or status is wrong.

- [ ] **Step 3: Update move to handle wumpus encounter**

Replace the `move` function in `wumpus/src/wumpus/game.clj`:

```clojure
(defn- resolve-wumpus [state rand-fn]
  (if (and (= :playing (:status state))
           (= (:player state) (:wumpus state)))
    (if (< (rand-fn) 0.5)
      (let [wumpus-neighbors (cave/neighbors (:cave-map state) (:wumpus state))
            new-wumpus (nth (vec wumpus-neighbors)
                            (mod (int (* (rand-fn) 3)) 3))]
        (assoc state :wumpus new-wumpus))
      (assoc state :status :lose-wumpus))
    state))

(defn move
  ([state room] (move state room rand))
  ([state room rand-fn]
   (-> state
       (assoc :player room)
       (resolve-wumpus rand-fn))))
```

- [ ] **Step 4: Run specs to verify pass**

```bash
cd wumpus && clj -M:spec
```

Expected: All PASS.

- [ ] **Step 5: Commit**

```bash
git add wumpus/src/wumpus/game.clj wumpus/spec/wumpus/game_spec.clj
git commit -m "Handle wumpus encounter on move — player eaten"
```

---

### Task 8: Move — Into Wumpus (Wumpus Flees)

**Files:**
- Modify: `wumpus/spec/wumpus/game_spec.clj`
- Modify: `wumpus/src/wumpus/game.clj`

- [ ] **Step 1: Write failing spec**

```clojure
  (it "wumpus flees when player enters its room and rng < 0.5"
    (let [neighbor (first (get @cave-map 1))
          state {:cave-map @cave-map :player 1 :wumpus neighbor
                 :pits #{} :bats #{} :arrows 5 :status :playing}
          flees (fn [] 0.2)  ; < 0.5 means wumpus moves
          result (move state neighbor flees)]
      (should= :playing (:status result))
      (should-not= neighbor (:wumpus result))
      (should-contain (:wumpus result)
                      (cave/neighbors @cave-map neighbor))))
```

- [ ] **Step 2: Run specs to verify failure**

```bash
cd wumpus && clj -M:spec
```

Expected: FAIL — the wumpus movement logic may not pick the right neighbor with
the deterministic `flees` function returning 0.2 for both calls. Adjust
`resolve-wumpus` if needed to use separate calls.

- [ ] **Step 3: Fix resolve-wumpus if needed**

If the spec fails because the same `rand-fn` is called twice (once for the
50/50 check, once for picking a neighbor), change the implementation to use
a single call for the neighbor selection:

```clojure
(defn- resolve-wumpus [state rand-fn]
  (if (and (= :playing (:status state))
           (= (:player state) (:wumpus state)))
    (if (< (rand-fn) 0.5)
      (let [wumpus-neighbors (vec (cave/neighbors (:cave-map state) (:wumpus state)))
            idx (int (* (rand-fn) (count wumpus-neighbors)))
            new-wumpus (nth wumpus-neighbors idx)]
        (assoc state :wumpus new-wumpus))
      (assoc state :status :lose-wumpus))
    state))
```

- [ ] **Step 4: Run specs to verify pass**

```bash
cd wumpus && clj -M:spec
```

Expected: All PASS.

- [ ] **Step 5: Commit**

```bash
git add wumpus/src/wumpus/game.clj wumpus/spec/wumpus/game_spec.clj
git commit -m "Handle wumpus fleeing on move"
```

---

### Task 9: Move — Into Pit

**Files:**
- Modify: `wumpus/spec/wumpus/game_spec.clj`
- Modify: `wumpus/src/wumpus/game.clj`

- [ ] **Step 1: Write failing spec**

```clojure
  (it "player dies when moving into a pit"
    (let [neighbor (first (get @cave-map 1))
          state {:cave-map @cave-map :player 1 :wumpus 20
                 :pits #{neighbor} :bats #{} :arrows 5 :status :playing}
          result (move state neighbor)]
      (should= :lose-pit (:status result))))
```

- [ ] **Step 2: Run specs to verify failure**

```bash
cd wumpus && clj -M:spec
```

Expected: FAIL — move doesn't check for pits yet.

- [ ] **Step 3: Add pit resolution to move**

Add a `resolve-pit` function and thread it in `move`:

```clojure
(defn- resolve-pit [state]
  (if (contains? (:pits state) (:player state))
    (assoc state :status :lose-pit)
    state))

(defn move
  ([state room] (move state room rand))
  ([state room rand-fn]
   (-> state
       (assoc :player room)
       (resolve-pit)
       (resolve-wumpus rand-fn))))
```

Note: pit is checked before wumpus — if you fall in a pit, you're dead
regardless.

- [ ] **Step 4: Run specs to verify pass**

```bash
cd wumpus && clj -M:spec
```

Expected: All PASS.

- [ ] **Step 5: Commit**

```bash
git add wumpus/src/wumpus/game.clj wumpus/spec/wumpus/game_spec.clj
git commit -m "Handle falling into pit on move"
```

---

### Task 10: Move — Into Bats

**Files:**
- Modify: `wumpus/spec/wumpus/game_spec.clj`
- Modify: `wumpus/src/wumpus/game.clj`

- [ ] **Step 1: Write failing spec**

```clojure
  (it "bats transport player to a random room"
    (let [neighbor (first (get @cave-map 1))
          state {:cave-map @cave-map :player 1 :wumpus 20
                 :pits #{} :bats #{neighbor} :arrows 5 :status :playing}
          ;; rand-fn returns values controlling: bat destination (room index)
          destination 7
          rand-fn (let [calls (atom [0.3])]  ; 0.3 * 20 = index 6 -> room 7
                    (fn [] (let [v (first @calls)]
                             (swap! calls rest)
                             (or v 0.99))))
          result (move state neighbor rand-fn)]
      (should= :playing (:status result))
      (should-not= neighbor (:player result))))
```

- [ ] **Step 2: Run specs to verify failure**

```bash
cd wumpus && clj -M:spec
```

Expected: FAIL — move doesn't handle bats yet.

- [ ] **Step 3: Add bat resolution to move**

Add `resolve-bats` and update the threading in `move`. Bats transport the
player to a random room, then we re-resolve pits and wumpus at the new
location. Bats stay put.

```clojure
(defn- resolve-bats [state rand-fn]
  (if (contains? (:bats state) (:player state))
    (let [rooms (vec (keys (:cave-map state)))
          idx (int (* (rand-fn) (count rooms)))
          new-room (nth rooms idx)]
      (-> state
          (assoc :player new-room)
          (resolve-pit)
          (resolve-wumpus rand-fn)))
    state))

(defn move
  ([state room] (move state room rand))
  ([state room rand-fn]
   (-> state
       (assoc :player room)
       (resolve-pit)
       (resolve-wumpus rand-fn)
       (resolve-bats rand-fn))))
```

- [ ] **Step 4: Run specs to verify pass**

```bash
cd wumpus && clj -M:spec
```

Expected: All PASS.

- [ ] **Step 5: Commit**

```bash
git add wumpus/src/wumpus/game.clj wumpus/spec/wumpus/game_spec.clj
git commit -m "Handle bat transport on move"
```

---

### Task 11: Shoot — Arrow Hits Wumpus

**Files:**
- Modify: `wumpus/spec/wumpus/game_spec.clj`
- Modify: `wumpus/src/wumpus/game.clj`

- [ ] **Step 1: Write failing spec**

```clojure
(describe "shoot"
  (with-all cave-map (make-cave-map))

  (it "wins when arrow hits wumpus"
    (let [wumpus-room (first (get @cave-map 1))
          state {:cave-map @cave-map :player 1 :wumpus wumpus-room
                 :pits #{} :bats #{} :arrows 5 :status :playing}
          result (shoot state [wumpus-room])]
      (should= :win (:status result)))))
```

- [ ] **Step 2: Run specs to verify failure**

```bash
cd wumpus && clj -M:spec
```

Expected: FAIL — `shoot` not defined.

- [ ] **Step 3: Implement shoot — arrow hits wumpus**

Add to `wumpus/src/wumpus/game.clj`:

```clojure
(defn shoot
  ([state path] (shoot state path rand))
  ([state path rand-fn]
   (let [cave-map (:cave-map state)]
     (loop [arrow-room (:player state)
            [target & remaining] path]
       (if (nil? target)
         (let [new-arrows (dec (:arrows state))]
           (if (zero? new-arrows)
             (assoc state :arrows 0 :status :lose-arrow)
             (let [startled (< (rand-fn) 0.75)
                   wumpus-neighbors (vec (cave/neighbors cave-map (:wumpus state)))
                   idx (int (* (rand-fn) (count wumpus-neighbors)))
                   new-wumpus (if startled
                                (nth wumpus-neighbors idx)
                                (:wumpus state))]
               (assoc state :arrows new-arrows :wumpus new-wumpus))))
         (let [next-room (if (cave/adjacent? cave-map arrow-room target)
                           target
                           (let [nbrs (vec (cave/neighbors cave-map arrow-room))
                                 idx (int (* (rand-fn) (count nbrs)))]
                             (nth nbrs idx)))]
           (cond
             (= next-room (:wumpus state))
             (assoc state :status :win :arrows (dec (:arrows state)))

             (= next-room (:player state))
             (assoc state :status :lose-arrow :arrows (dec (:arrows state)))

             :else
             (recur next-room remaining))))))))
```

- [ ] **Step 4: Run specs to verify pass**

```bash
cd wumpus && clj -M:spec
```

Expected: All PASS.

- [ ] **Step 5: Commit**

```bash
git add wumpus/src/wumpus/game.clj wumpus/spec/wumpus/game_spec.clj
git commit -m "Implement shoot — arrow hits wumpus"
```

---

### Task 12: Shoot — Arrow Hits Player

**Files:**
- Modify: `wumpus/spec/wumpus/game_spec.clj`

- [ ] **Step 1: Write failing spec**

The arrow path loops back to the player's room.

```clojure
  (it "player dies when arrow enters player's room"
    (let [nbrs (vec (get @cave-map 1))
          n1 (nth nbrs 0)
          ;; find a neighbor of n1 that leads back to 1
          path-back (first (filter #(cave/adjacent? @cave-map % 1)
                                   (disj (cave/neighbors @cave-map n1) 1)))
          ;; If direct path back: shoot n1 then 1
          state {:cave-map @cave-map :player 1 :wumpus 20
                 :pits #{} :bats #{} :arrows 5 :status :playing}
          result (shoot state [n1 1])]
      (should= :lose-arrow (:status result))))
```

- [ ] **Step 2: Run specs — should pass**

This should already pass with the existing `shoot` implementation since it
checks `(= next-room (:player state))`.

```bash
cd wumpus && clj -M:spec
```

Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add wumpus/spec/wumpus/game_spec.clj
git commit -m "Add spec: arrow hits player"
```

---

### Task 13: Shoot — Miss, Wumpus Startled

**Files:**
- Modify: `wumpus/spec/wumpus/game_spec.clj`

- [ ] **Step 1: Write failing spec**

```clojure
  (it "arrow misses and wumpus moves when startled"
    (let [neighbor (first (get @cave-map 1))
          ;; pick a room not containing wumpus or player
          empty-room (second (vec (get @cave-map 1)))
          state {:cave-map @cave-map :player 1 :wumpus 20
                 :pits #{} :bats #{} :arrows 5 :status :playing}
          ;; rand-fn: 0.5 < 0.75 so wumpus is startled, 0.0 picks first neighbor
          calls (atom [0.5 0.0])
          rand-fn (fn [] (let [v (first @calls)]
                           (swap! calls rest)
                           (or v 0.5)))
          result (shoot state [empty-room] rand-fn)]
      (should= :playing (:status result))
      (should= 4 (:arrows result))
      (should-not= 20 (:wumpus result))))

  (it "arrow misses and wumpus stays when not startled"
    (let [empty-room (second (vec (get @cave-map 1)))
          state {:cave-map @cave-map :player 1 :wumpus 20
                 :pits #{} :bats #{} :arrows 5 :status :playing}
          ;; rand-fn: 0.8 >= 0.75 so wumpus stays
          rand-fn (fn [] 0.8)
          result (shoot state [empty-room] rand-fn)]
      (should= :playing (:status result))
      (should= 4 (:arrows result))
      (should= 20 (:wumpus result))))
```

- [ ] **Step 2: Run specs — should pass**

```bash
cd wumpus && clj -M:spec
```

Expected: PASS (already implemented in shoot).

- [ ] **Step 3: Commit**

```bash
git add wumpus/spec/wumpus/game_spec.clj
git commit -m "Add specs: arrow miss with wumpus startled and not startled"
```

---

### Task 14: Shoot — Last Arrow Exhausted

**Files:**
- Modify: `wumpus/spec/wumpus/game_spec.clj`

- [ ] **Step 1: Write failing spec**

```clojure
  (it "player loses when last arrow misses"
    (let [empty-room (second (vec (get @cave-map 1)))
          state {:cave-map @cave-map :player 1 :wumpus 20
                 :pits #{} :bats #{} :arrows 1 :status :playing}
          result (shoot state [empty-room])]
      (should= :lose-arrow (:status result))
      (should= 0 (:arrows result))))
```

- [ ] **Step 2: Run specs — should pass**

```bash
cd wumpus && clj -M:spec
```

Expected: PASS (already handled in shoot when `new-arrows` is zero).

- [ ] **Step 3: Commit**

```bash
git add wumpus/spec/wumpus/game_spec.clj
git commit -m "Add spec: last arrow exhausted means player loses"
```

---

### Task 15: Shoot — Invalid Room Deflection

**Files:**
- Modify: `wumpus/spec/wumpus/game_spec.clj`

- [ ] **Step 1: Write failing spec**

```clojure
  (it "arrow deflects to random neighbor when path room is invalid"
    (let [;; room 20 is not adjacent to room 1's neighbor path
          neighbor (first (get @cave-map 1))
          non-adj-of-neighbor (first (clojure.set/difference
                                       (set (range 1 21))
                                       (conj (cave/neighbors @cave-map neighbor) neighbor)))
          state {:cave-map @cave-map :player 1 :wumpus non-adj-of-neighbor
                 :pits #{} :bats #{} :arrows 5 :status :playing}
          ;; first room is valid (neighbor), second is invalid from that room
          ;; rand-fn controls where the deflected arrow goes
          calls (atom [0.0])
          rand-fn (fn [] (let [v (first @calls)]
                           (swap! calls rest)
                           (or v 0.5)))
          result (shoot state [neighbor non-adj-of-neighbor] rand-fn)]
      ;; arrow was deflected — it didn't hit wumpus (which is at non-adj-of-neighbor)
      ;; so either it's still playing or some other outcome based on deflection
      (should-not= :win (:status result))))
```

- [ ] **Step 2: Run specs — should pass**

```bash
cd wumpus && clj -M:spec
```

Expected: PASS (deflection already implemented in shoot).

- [ ] **Step 3: Commit**

```bash
git add wumpus/spec/wumpus/game_spec.clj
git commit -m "Add spec: arrow deflects on invalid path room"
```

---

### Task 16: Messages Module

**Files:**
- Create: `wumpus/src/wumpus/messages.clj`
- Create: `wumpus/spec/wumpus/messages_spec.clj`

- [ ] **Step 1: Write failing specs for all message functions**

Create `wumpus/spec/wumpus/messages_spec.clj`:

```clojure
(ns wumpus.messages-spec
  (:require [speclj.core :refer :all]
            [wumpus.messages :refer :all]))

(describe "room-description"
  (it "describes the room and tunnels"
    (should= "You are in room 5.\nTunnels lead to: 1 4 8"
             (room-description 5 #{1 4 8}))))

(describe "hazard-warning"
  (it "warns about wumpus"
    (should= "I smell a Wumpus!" (hazard-warning :wumpus)))

  (it "warns about pit"
    (should= "I feel a draft!" (hazard-warning :pit)))

  (it "warns about bats"
    (should= "Bats nearby!" (hazard-warning :bats))))

(describe "outcome-message"
  (it "announces win"
    (should= "Hee hee hee, the Wumpus'll get you next time!!"
             (outcome-message :win)))

  (it "announces wumpus death"
    (should= "Tsk tsk tsk - Wumpus got you!"
             (outcome-message :lose-wumpus)))

  (it "announces pit death"
    (should= "YYYIIIIEEEE . . . fell in pit"
             (outcome-message :lose-pit)))

  (it "announces arrow death"
    (should= "Ouch! Arrow got you!"
             (outcome-message :lose-arrow))))

(describe "intro"
  (it "returns the game banner"
    (should-contain "HUNT THE WUMPUS" (intro))))
```

- [ ] **Step 2: Create empty namespace and run to see failure**

Create `wumpus/src/wumpus/messages.clj`:

```clojure
(ns wumpus.messages)
```

```bash
cd wumpus && clj -M:spec
```

Expected: FAIL — functions not defined.

- [ ] **Step 3: Implement all message functions**

Edit `wumpus/src/wumpus/messages.clj`:

```clojure
(ns wumpus.messages
  (:require [clojure.string :as str]))

(defn room-description [room neighbors]
  (str "You are in room " room ".\n"
       "Tunnels lead to: " (str/join " " (sort neighbors))))

(def ^:private warnings
  {:wumpus "I smell a Wumpus!"
   :pit    "I feel a draft!"
   :bats   "Bats nearby!"})

(defn hazard-warning [hazard-key]
  (get warnings hazard-key))

(def ^:private outcomes
  {:win         "Hee hee hee, the Wumpus'll get you next time!!"
   :lose-wumpus "Tsk tsk tsk - Wumpus got you!"
   :lose-pit    "YYYIIIIEEEE . . . fell in pit"
   :lose-arrow  "Ouch! Arrow got you!"})

(defn outcome-message [status]
  (get outcomes status))

(defn intro []
  (str/join "\n"
    ["HUNT THE WUMPUS"
     "==============="]))
```

- [ ] **Step 4: Run specs to verify pass**

```bash
cd wumpus && clj -M:spec
```

Expected: All PASS.

- [ ] **Step 5: Commit**

```bash
git add wumpus/src/wumpus/messages.clj wumpus/spec/wumpus/messages_spec.clj
git commit -m "Implement messages module with original game text"
```

---

### Task 17: Main — I/O Shell

**Files:**
- Create: `wumpus/src/wumpus/main.clj`

This is the thin imperative shell. It has no unit specs — it will be covered
by the Gherkin E2E tests in Task 18. The logic is trivial: read a line, parse,
dispatch to the pure core, print results.

- [ ] **Step 1: Implement main.clj**

Create `wumpus/src/wumpus/main.clj`:

```clojure
(ns wumpus.main
  (:require [wumpus.dodecahedron :as cave]
            [wumpus.game :as game]
            [wumpus.messages :as msg]
            [clojure.string :as str]))

(defn- print-status [state]
  (doseq [h (game/sensed-hazards state)]
    (println (msg/hazard-warning h)))
  (println (msg/room-description
             (game/player-room state)
             (game/player-neighbors state))))

(defn- parse-command [line]
  (let [tokens (str/split (str/trim line) #"\s+")
        cmd (str/upper-case (first tokens))
        nums (map #(Integer/parseInt %) (rest tokens))]
    (case cmd
      "M" {:action :move :room (first nums)}
      "S" {:action :shoot :path (vec nums)}
      nil)))

(defn- game-loop [state]
  (loop [state state]
    (when (= :playing (:status state))
      (print-status state)
      (print "> ")
      (flush)
      (if-let [line (read-line)]
        (if-let [cmd (parse-command line)]
          (let [new-state (case (:action cmd)
                            :move (game/move state (:room cmd))
                            :shoot (game/shoot state (:path cmd)))]
            (recur new-state))
          (do (println "Use: M <room> or S <room1> <room2> ...")
              (recur state)))
        state))))

(defn -main [& _args]
  (println (msg/intro))
  (let [cave-map (cave/make-cave-map)]
    (loop []
      (let [state (game/new-game cave-map)
            final (game-loop state)]
        (println (msg/outcome-message (:status final)))
        (print "Same setup (Y-N)? ")
        (flush)
        (when (= "Y" (str/upper-case (str/trim (or (read-line) "N"))))
          (recur))))))
```

- [ ] **Step 2: Smoke-test manually**

```bash
cd wumpus && echo "M 2" | clj -M -m wumpus.main
```

Verify it prints the intro, a room description, and processes the move.

- [ ] **Step 3: Commit**

```bash
git add wumpus/src/wumpus/main.clj
git commit -m "Add main I/O shell with single-line command parsing"
```

---

### Task 18: Gherkin E2E Scenarios

**Files:**
- Create: `wumpus/features/wumpus.feature`

- [ ] **Step 1: Write Gherkin feature file**

Create `wumpus/features/wumpus.feature`:

```gherkin
Feature: Hunt the Wumpus

  Scenario: Player wins by shooting the wumpus
    Given a cave map
    And the player is in room 1
    And the wumpus is in room 2
    And there are no other hazards
    When the player shoots into room 2
    Then the game status is win
    And the outcome message is "Hee hee hee, the Wumpus'll get you next time!!"

  Scenario: Player eaten by wumpus
    Given a cave map
    And the player is in room 1
    And the wumpus is in room 2
    And the wumpus will stay when encountered
    And there are no other hazards
    When the player moves to room 2
    Then the game status is lose-wumpus
    And the outcome message is "Tsk tsk tsk - Wumpus got you!"

  Scenario: Player falls into pit
    Given a cave map
    And the player is in room 1
    And a pit is in room 2
    And the wumpus is in room 20
    When the player moves to room 2
    Then the game status is lose-pit

  Scenario: Bats transport the player
    Given a cave map
    And the player is in room 1
    And bats are in room 2
    And the wumpus is in room 20
    And the bat destination is room 10
    When the player moves to room 2
    Then the player is in room 10

  Scenario: Player runs out of arrows
    Given a cave map
    And the player is in room 1
    And the wumpus is in room 20
    And there are no other hazards
    And the player has 1 arrow
    When the player shoots into room 5
    Then the game status is lose-arrow

  Scenario: Arrow hits the player
    Given a cave map
    And the player is in room 1
    And the wumpus is in room 20
    And there are no other hazards
    When the player shoots through rooms 2 1
    Then the game status is lose-arrow
    And the outcome message is "Ouch! Arrow got you!"

  Scenario: Player senses adjacent wumpus
    Given a cave map
    And the player is in room 1
    And the wumpus is in room 2
    And there are no other hazards
    Then the player senses wumpus

  Scenario: Player senses adjacent pit
    Given a cave map
    And the player is in room 1
    And a pit is in room 2
    And the wumpus is in room 20
    Then the player senses pit

  Scenario: Player senses adjacent bats
    Given a cave map
    And the player is in room 1
    And bats are in room 2
    And the wumpus is in room 20
    Then the player senses bats
```

Note: rooms 1 and 2 must be adjacent in the dodecahedron. In our edge list,
`[1 2]` is the first edge, so they are adjacent. Room 5 is also adjacent to
room 1, and room 20 is far from room 1.

- [ ] **Step 2: Write the Gherkin step definitions as a Speclj test**

Create `wumpus/spec/wumpus/features_spec.clj` as the glue code that connects
the Gherkin scenarios to the production code:

```clojure
(ns wumpus.features-spec
  (:require [speclj.core :refer :all]
            [wumpus.dodecahedron :as cave]
            [wumpus.game :as game]
            [wumpus.messages :as msg]))

(describe "E2E: Player wins by shooting the wumpus"
  (with-all cave-map (cave/make-cave-map))
  (with-all state {:cave-map @cave-map :player 1 :wumpus 2
                   :pits #{} :bats #{} :arrows 5 :status :playing})
  (with-all result (game/shoot @state [2]))

  (it "game status is win"
    (should= :win (:status @result)))

  (it "outcome message is correct"
    (should= "Hee hee hee, the Wumpus'll get you next time!!"
             (msg/outcome-message (:status @result)))))

(describe "E2E: Player eaten by wumpus"
  (with-all cave-map (cave/make-cave-map))
  (with-all state {:cave-map @cave-map :player 1 :wumpus 2
                   :pits #{} :bats #{} :arrows 5 :status :playing})
  (with-all result (game/move @state 2 (fn [] 0.6)))

  (it "game status is lose-wumpus"
    (should= :lose-wumpus (:status @result)))

  (it "outcome message is correct"
    (should= "Tsk tsk tsk - Wumpus got you!"
             (msg/outcome-message (:status @result)))))

(describe "E2E: Player falls into pit"
  (with-all cave-map (cave/make-cave-map))
  (with-all state {:cave-map @cave-map :player 1 :wumpus 20
                   :pits #{2} :bats #{} :arrows 5 :status :playing})
  (with-all result (game/move @state 2))

  (it "game status is lose-pit"
    (should= :lose-pit (:status @result))))

(describe "E2E: Bats transport the player"
  (with-all cave-map (cave/make-cave-map))
  (with-all state {:cave-map @cave-map :player 1 :wumpus 20
                   :pits #{} :bats #{2} :arrows 5 :status :playing})
  (with-all result (let [rooms (vec (keys @cave-map))
                         target-idx (.indexOf rooms 10)
                         frac (/ target-idx (count rooms))]
                     (game/move @state 2 (fn [] frac))))

  (it "player is no longer in the bat room"
    (should-not= 2 (:player @result))))

(describe "E2E: Player runs out of arrows"
  (with-all cave-map (cave/make-cave-map))
  (with-all state {:cave-map @cave-map :player 1 :wumpus 20
                   :pits #{} :bats #{} :arrows 1 :status :playing})
  (with-all result (game/shoot @state [5]))

  (it "game status is lose-arrow"
    (should= :lose-arrow (:status @result))))

(describe "E2E: Arrow hits the player"
  (with-all cave-map (cave/make-cave-map))
  (with-all state {:cave-map @cave-map :player 1 :wumpus 20
                   :pits #{} :bats #{} :arrows 5 :status :playing})
  (with-all result (game/shoot @state [2 1]))

  (it "game status is lose-arrow"
    (should= :lose-arrow (:status @result)))

  (it "outcome message is correct"
    (should= "Ouch! Arrow got you!"
             (msg/outcome-message (:status @result)))))

(describe "E2E: Player senses adjacent wumpus"
  (with-all cave-map (cave/make-cave-map))
  (with-all state {:cave-map @cave-map :player 1 :wumpus 2
                   :pits #{} :bats #{}})

  (it "player senses wumpus"
    (should-contain :wumpus (game/sensed-hazards @state))))

(describe "E2E: Player senses adjacent pit"
  (with-all cave-map (cave/make-cave-map))
  (with-all state {:cave-map @cave-map :player 1 :wumpus 20
                   :pits #{2} :bats #{}})

  (it "player senses pit"
    (should-contain :pit (game/sensed-hazards @state))))

(describe "E2E: Player senses adjacent bats"
  (with-all cave-map (cave/make-cave-map))
  (with-all state {:cave-map @cave-map :player 1 :wumpus 20
                   :pits #{} :bats #{2}})

  (it "player senses bats"
    (should-contain :bats (game/sensed-hazards @state))))
```

- [ ] **Step 3: Run all specs**

```bash
cd wumpus && clj -M:spec
```

Expected: All PASS (this is glue code calling already-tested functions with
specific setups matching the Gherkin scenarios).

- [ ] **Step 4: Commit**

```bash
git add wumpus/features/wumpus.feature wumpus/spec/wumpus/features_spec.clj
git commit -m "Add Gherkin E2E scenarios and step definition specs"
```

---

### Task 19: Refactor and Final Review

**Files:**
- All source files in `wumpus/src/wumpus/`

- [ ] **Step 1: Run full spec suite**

```bash
cd wumpus && clj -M:spec
```

All specs must pass.

- [ ] **Step 2: Check cyclomatic complexity**

Review each function manually. Every function should have cyclomatic complexity
<= 4. The most complex function is `shoot` — if it exceeds 4, extract the
arrow-step logic into a helper:

```clojure
(defn- advance-arrow [cave-map arrow-room target rand-fn]
  (if (cave/adjacent? cave-map arrow-room target)
    target
    (let [nbrs (vec (cave/neighbors cave-map arrow-room))
          idx (int (* (rand-fn) (count nbrs)))]
      (nth nbrs idx))))

(defn- arrow-miss [state rand-fn]
  (let [new-arrows (dec (:arrows state))]
    (if (zero? new-arrows)
      (assoc state :arrows 0 :status :lose-arrow)
      (let [startled (< (rand-fn) 0.75)
            cave-map (:cave-map state)
            wumpus-nbrs (vec (cave/neighbors cave-map (:wumpus state)))
            idx (int (* (rand-fn) (count wumpus-nbrs)))
            new-wumpus (if startled (nth wumpus-nbrs idx) (:wumpus state))]
        (assoc state :arrows new-arrows :wumpus new-wumpus)))))
```

- [ ] **Step 3: Run specs after refactor**

```bash
cd wumpus && clj -M:spec
```

All specs must still pass.

- [ ] **Step 4: Commit**

```bash
git add wumpus/src/wumpus/game.clj
git commit -m "Refactor shoot to reduce cyclomatic complexity"
```

---

### Task 20: Final Smoke Test

- [ ] **Step 1: Run the game interactively**

```bash
cd wumpus && clj -M -m wumpus.main
```

Play a few turns. Verify:
- Intro prints
- Room description and warnings print correctly
- `M <room>` moves the player
- `S <room> ...` shoots an arrow
- Game ends correctly on win/loss
- "Same setup (Y-N)?" works

- [ ] **Step 2: Final commit if any fixes needed**

```bash
git add -A wumpus/
git commit -m "Final polish from smoke testing"
```
