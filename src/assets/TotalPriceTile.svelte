<script lang="ts">
    import NumberFlow, { NumberFlowGroup } from "@number-flow/svelte";
    import { slide } from "svelte/transition";
    export let number: number;
    export let dif: number;
</script>

<div
    class="border p-2 rounded-lg shadow text-center"
    transition:slide={{ duration: 500 }}
>
    <h1 class="text-xl font-bold mb-3">💸 Final Total Price</h1>
    <NumberFlowGroup
        style="--number-flow-char-height: 0.85em"
        class="flex items-center gap-4 font-semibold"
    >
        <NumberFlow
            value={number}
            format={{
                style: "currency",
                currency: "EUR",
                maximumFractionDigits: 0,
            }}
            locales={"it-IT"}
            class="text-3xl"
        />
        <NumberFlow
            value={dif}
            format={{
                style: "currency",
                currency: "EUR",
                signDisplay: "always",
                maximumFractionDigits: 0,
            }}
            class="text-lg transition-colors duration-300 text-red-500"
        />
    </NumberFlowGroup>
</div>
