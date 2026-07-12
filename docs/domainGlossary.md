# Domain Glossary

| Term | Meaning |
|---|---|
| Match | One authoritative 6–10 player game and its durable history. |
| Chapter | One government, Quest, and possible Trial cycle; at most five. |
| Phase | A normative uppercase engine state. |
| Scenario | Versioned script containing crises, quests, endings, and compatible roles. |
| Role / Faction / Objective | A private identity, its broad allegiance, and its personal success condition. |
| Command | Idempotent player or system intention validated against a revision. |
| Domain Event | Canonical fact produced by accepted engine execution. |
| Public / Private Event | Recipient-safe presentation of a canonical fact. |
| Projection | Explicit client view derived from internal state. |
| Snapshot | Planned compressed recovery state; the required future cadence is every 20 events and phase boundary. |
| Revision / Sequence | Monotonic state version / strictly ordered event position. |
| Quest | Party action resolved from owner-hidden Sigils. |
| Trial | Post-Quest accusation and banishment vote. |
| Ending | One canonical outcome set selected by precedence. |
| Fallen Player | Banished participant retaining voice and one Last Vote. |
| Sigil | Secret Quest submission such as Aid or Betray. |
| Regent / Pathfinder | Rotating nominator / elected party leader. |
| Fracture / Ruin | Failed-government pressure / catastrophe track. |
| Usurper / Wanderer | Hidden Shadow leader / independent oath-driven role. |

## GDD consistency note

GDD Figure 3 labels the post-`FINAL_RECKONING` state `EPILOGUE_REMATCH`, while the immediately following normative allowed-command table defines `EPILOGUE`. The revised architecture brief also declares `EPILOGUE` canonical. Therefore the engine, contracts, tests, persistence, and telemetry use only `EPILOGUE`. `ROOM_LOBBY`, `SCRIPT_SELECTION`, and `ROOM_DEPARTURE` describe outgoing product destinations, not replacement match-engine phases. The source discrepancy remains documented rather than silently rewritten.
