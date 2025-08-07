<script lang="ts">
    import NumberFlow, { NumberFlowGroup } from "@number-flow/svelte";
    import { slide, fade } from "svelte/transition";
    export let showVal: boolean;
    export let number: number;
    export let delta: number;
    export let name: string;
    export let duration: number;
</script>

<div class="flex-1 min-w-[160px]">
    <div
        class="bg-neutral-900 border border-purple-400 rounded-xl p-5 flex-1 text-center"
        transition:slide={{ duration: 500 }}
    >
        <h1 class="text-lg font-bold mb-1">
            {name}
            <NumberFlow value={duration} class="text-lg font-bold" />
            {" anni"}
        </h1>
        <NumberFlowGroup
            style="--number-flow-char-height: 0.85em"
            class="flex items-center font-semibold"
        >
            <NumberFlow
                value={number ? number : 0}
                format={{
                    style: "currency",
                    currency: "EUR",
                    maximumFractionDigits: 0,
                }}
                locales={"it-IT"}
                class="text-[1.5rem] font-semibold"
            />
            <NumberFlow
                value={showVal ? delta : 0}
                format={{
                    style: "percent",
                    maximumFractionDigits: 2,
                    minimumFractionDigits: 2,
                    signDisplay: "always",
                }}
                class="text-lg text-red-500 font-semibold"
            /> <span class="text-xs"> sul mutuo</span>
        </NumberFlowGroup>
    </div>
</div>
