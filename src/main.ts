import { mount } from "svelte";
import "./app.css";

import HouseCost from "./lib/HouseCost.svelte";
const app = mount(HouseCost, {
  target: document.getElementById("app")!,
});

export default app;
