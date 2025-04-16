import { mount } from "svelte";
import "./app.css";
import HouseCost from "./lib/HouseCost.svelte";
import Test from "./lib/Test.svelte";

const app = mount(HouseCost, {
  target: document.getElementById("app")!,
});

export default app;
