import { writable, type Writable } from "svelte/store";

export const collapsedDetailTiles: Writable<Record<string, boolean>> = writable(
  {},
);
