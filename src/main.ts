import { mount } from "svelte";
import "./app.css";
import Chart from "chart.js/auto";

import HouseCost from "./lib/HouseCost.svelte";
const app = mount(HouseCost, {
  target: document.getElementById("app")!,
});

export default app;
