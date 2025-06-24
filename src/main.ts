import { mount } from "svelte";
import "./app.css";
import HouseCost from "./lib/HouseCost.svelte";
import ChartTest from "./lib/ChartTest.svelte";

const app = mount(HouseCost, {
  target: document.getElementById("app")!,
});

export default app;
