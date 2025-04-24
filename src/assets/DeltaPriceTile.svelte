<script lang="ts">
    import NumberFlow, { NumberFlowGroup } from "@number-flow/svelte";
    import { slide } from "svelte/transition";
    export let showVal: boolean;
    export let number: number;
    export let delta: number;
    export let name: string;
</script>

<div
    class="border p-2 rounded-lg shadow text-center"
    transition:slide={{ duration: 500 }}
>
    <h1 class="text-lg font-semibold mb-2">{name}</h1>
    <NumberFlowGroup
        style="--number-flow-char-height: 0.85em"
        class="flex items-center gap-4 font-semibold"
    >
        <NumberFlow
            value={number ? number : 0}
            format={{
                style: "currency",
                currency: "EUR",
                maximumFractionDigits: 0,
            }}
            locales={"it-IT"}
            class="text-3xl"
        />
        <NumberFlow
            value={showVal ? delta : 0}
            format={{
                style: "percent",
                maximumFractionDigits: 2,
                minimumFractionDigits: 2,
                signDisplay: "always",
            }}
            class="text-lg transition-colors duration-300 text-red-500"
        />
    </NumberFlowGroup>
</div>
