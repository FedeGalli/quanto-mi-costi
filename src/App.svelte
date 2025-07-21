<script>
    import Router from "svelte-spa-router";
    import HouseCost from "./lib/HouseCost.svelte";
    import SignIn from "./lib/SignIn.svelte";
    import { onMount, onDestroy } from "svelte";
    import { initAuthStore, destroyAuthStore } from "./lib/auth/auth-store";

    onMount(() => {
        // Initialize the auth store when the app starts
        initAuthStore();
    });

    onDestroy(() => {
        // Clean up the auth store when the app is destroyed
        destroyAuthStore();
    });

    // Define your routes
    const routes = {
        "/": HouseCost,
        "/signin": SignIn,
    };
</script>

<Router {routes} />
